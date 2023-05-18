package log

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func fileNameGen(fileName string, num int) string {
	return fileName + "." + fmt.Sprint(num)
}

type FileInitorMode int

const (
	// This mode will overwrite the previous file that was written
	// Under this mode, the name will always be consistent
	FILE_OVERWRITE FileInitorMode = iota
	// This mode will append to the file names. The files will be named as 'fileName.0', 'fileName.1' etc. The most recent file will be the one with the highest extension. 
	FILE_ASCENDING
	// This mode will 'shift' the files up by 1. The files will be named as 'fileName.0', 'fileName.1' etc. The most recent file will be the .0 one.
	FILE_DESCENDING
)

func theNames(dir, name string) []int {
	files, _ := os.ReadDir(dir)

	names := []int{}
	
	for _, f := range files {
		cut, found := strings.CutPrefix(f.Name(), name)
		if !found {
			continue
		}
		// the `.NNNNN` is not present
		if len(cut) < 2 {
			continue
		}

		cut = cut[1:]

		n, err := strconv.Atoi(cut)
		if err != nil || n < 0 {
			continue
		}

		names = append(names, n)
	}

	sort.IntSlice(names).Sort()

	return names
}

// maxNumber must be either >= 0. 
// maxNumber indicates the max amount of files that can be stored. 0 indicates that an infinite amount of files can be stored.
// if mode == FILE_OVERWRITE, maxNumber is ignored.
// Please note the negative numbers do not count towards the maxNumber count! These are only made if a user renames a file to a negative number.
func NewLoggerFileComplex(fileName string, mode FileInitorMode, maxNumber int) (LogCB) {
	var file *os.File
	var err error
	
	if fileName == "" {
		panic("fileName must be specified in the file logger!")
	}

	if fileName[len(fileName)-1] == '/' {
		panic("fileName must be a file for the file logger!")
	}

	info, err := os.Stat(fileName)

	if err == nil {
		if info.IsDir() {
			panic("fileName points to a path in the file logger!")
		}
	}

	if mode != FILE_OVERWRITE {
		if mode < 0 {
			panic("maxNumber must be >= 0 for the file logger!")
		}
	}
	
	switch mode {
	case FILE_OVERWRITE:
		file, err = os.Create(fileName)
	case FILE_ASCENDING:
		dir, name := path.Split(fileName)
		
		names := theNames(dir, name)
		
		if maxNumber != 0 && len(names) >= maxNumber {
			for _, n := range names[:maxNumber] {
				os.Remove(fileNameGen(fileName, n))
			}
			names = names[maxNumber:]

			for i, n := range names {
				os.Rename(fileNameGen(fileName, n), fileNameGen(fileName, n - 1))
				names[i]--
			}
		}

		newN := -1

		if len(names) != 0 {
			newN = names[len(names)-1]
		}

		file, err = os.Create(fileNameGen(fileName, newN + 1))
	case FILE_DESCENDING:
		dir, name := path.Split(fileName)
		names := theNames(dir, name)
		
		if maxNumber != 0 && len(names) >= maxNumber {
			for _, n := range names[maxNumber:] {
				os.Remove(fileNameGen(fileName, n))
			}
		}

		names = names[:maxNumber]

		for _, n := range names {
			os.Rename(fileNameGen(fileName, n), fileNameGen(fileName, n + 1))
		}

		file, err = os.Create(fileNameGen(fileName, 0))
	}
	
	if err != nil {
		panic(err)
	}
	
	fileLock := &sync.Mutex{}
	
	return func() (logger DoLog, closer Closer) {
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

// Write logs to file
// This is an alias for NewLoggerFileComplex(fileName, FILE_ASCENDING, 10)
func NewLoggerFile(fileName string) LogCB {
	return NewLoggerFileComplex(fileName, FILE_ASCENDING, 10)
}
