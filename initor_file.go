package log

import (
	"errors"
	"fmt"
	"os"
	"sync"
)

func fileAppendix(fileName string, num int) string {
	if num == 0 {
		return fileName
	}

	return fileName + "." + fmt.Sprint(num)
}

// Write logs to file
func NewLoggerFile(fileName string) LogCB {
	return func() (logger DoLog, closer Closer) {
		adder := 1
		
		for {
			_, err := os.Stat(fileAppendix(fileName, adder))
			if errors.Is(err, os.ErrNotExist) {
				break
			}

			adder++
		}

		file, err := os.Create(fileAppendix(fileName, adder))

		if err != nil {
			panic(err)
		}

		fileLock := &sync.Mutex{}

		logger = func (_ LogLevel, prefix, msg string) {
			fileLock.Lock()
			file.Write([]byte(prefix + " " + msg + "\n"))
			fileLock.Unlock()
		}
		
		closer = func() {
			fileLock.Lock()
			file.Close()
			fileLock.Unlock()
		}

		return
	}
}
