package log

import (
	"fmt"
	"sync"
	"time"

	"github.com/fatih/color"
)

type logLevelInfo struct {
	Color  *color.Color
	Prefix string
}

type LogLevel int

const (
	LL_DEBUG LogLevel = iota
	LL_SUCCESS
	LL_WARN
	LL_ERROR
	LL_FATAL
)

func GetColor(l LogLevel) *color.Color {
	return levelInfo[l].Color
}

var levelInfo = map[LogLevel]logLevelInfo{
	LL_DEBUG: {
		Prefix: "DEBUG",
		Color:  color.New(color.FgCyan),
	},
	LL_SUCCESS: {
		Prefix: "SUCCESS",
		Color:  color.New(color.FgGreen),
	},
	LL_WARN: {
		Prefix: "WARNING",
		Color:  color.New(color.FgYellow),
	},
	LL_ERROR: {
		Prefix: "ERROR",
		Color:  color.New(color.FgRed),
	},
	LL_FATAL: {
		Prefix: "FATAL ERROR",
		Color:  color.New(color.FgWhite, color.BgRed),
	},
}

var gonnaPanic = false

func log(level LogLevel, msg string, args ...any) {
	levelInfo := levelInfo[level]

	msg = fmt.Sprintf(msg, args...)

	prefix := `[` + levelInfo.Prefix + `] [` + time.Now().Format(`02 Jan 2006 15:04:05`) + `]`

	wg := &sync.WaitGroup{}

	wg.Add(len(loggers))

	for _, l := range loggers {
		go func(l DoLog) {
			lvl := level

			for !ready {
				time.Sleep(200 * time.Millisecond)
			}

			l(lvl, prefix, msg)
			wg.Done()
		}(l)
	}

	if level == LL_FATAL {

		// Don't panic in case of race conditions!
		if gonnaPanic {
			level = LL_ERROR
		}

		gonnaPanic = true
	}

	go func() {
		wg.Wait()
		
		if level == LL_FATAL {
			panic(msg)
		}
	}()
}

func Debug(msg string, args ...any) {
	log(LL_DEBUG, msg, args...)
}

func Success(msg string, args ...any) {
	log(LL_SUCCESS, msg, args...)
}

func Warn(msg string, args ...any) {
	log(LL_WARN, msg, args...)
}

func Error(msg string, args ...any) {
	log(LL_ERROR, msg, args...)
}

// Warning! This function causes panic! 
func Fatal(msg string, args ...any) {
	log(LL_FATAL, msg, args...)
}

// "While {CONTEXT}: {ERROR}"
// This causes panic!
func FatalIfErr(err error, context string, args ...any) {
	if err != nil {
		context = fmt.Sprintf(context, args...)
		Fatal("Error while %s: %s", context, err.Error())
	}
}

// "While {CONTEXT}: {ERROR}"
// Returns true if err != nil, intended for inline usage:
// if ErrorIfErr(err, "fetching api info, status: %d", resp.StatusCode) { return }
func ErrorIfErr(err error, context string, args ...any) bool {
	if err != nil {
		context = fmt.Sprintf(context, args...)

		Error("Error while %s: %s", context, err.Error())
		return true
	}
	return false
}

// Close all the outputs.
// Note: The logger does not check if the output is close or not, so it is safe to call this only at the end of the application.
func Close() {
	for _, c := range closers {
		c()
	}
}
