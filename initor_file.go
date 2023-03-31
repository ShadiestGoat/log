package log

import (
	"os"
	"sync"
)

// Write logs to file
func NewLoggerFile(fileName string) LogCB {
	return func() (logger DoLog, closer Closer) {
		file, err := os.Create(fileName)

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
