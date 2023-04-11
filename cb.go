package log

// Represents a callback to logging a message
type DoLog func (c LogLevel, prefix, msg string)

// Represents a callback to closing a logger
type Closer func ()

// Function that is used in Init() - adds another callback to the caller.
// closer can be nil here, in which case it will not be called 
type LogCB func () (logger DoLog, closer Closer)

var loggers = []DoLog{}
var closers = []func (){}

var ready bool

func Init(callbacks ...LogCB) {
	if len(callbacks) == 0 {
		panic("No callbacks for logger")
	}

	for _, cb := range callbacks {
		l, c := cb()
		loggers = append(loggers, l)
		
		if c != nil {
			closers = append(closers, c)
		}
	}

	ready = true
}