package logger

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

// Log level : DEBUG; INFO, INFORMATION; WARN, WARNING; ERROR, FAIL, FAILURE
// number will help to combine multiple levels
const (
	// STDOUT signify that the message will be printed on stdout
	STDOUT = 1
	// DEBUG signify normal debugging message
	DEBUG = 2
	// INFO signify normal info message
	INFO = 4
	// WARN signify warning message when some assumptions breaks but not harmful
	WARN = 8
	// ERROR signify error and will handled message
	ERROR = 16
	// FATAL signify fatal error, will exit from program
	FATAL = 32
)

// Info will use by different goroutines while logging to file or to screen
type Info struct {
	LogTime   time.Time
	Type      int // WARN | INFO | ERROR | DEBUG
	Package   string
	Method    string
	ErrorCode int
	Message   string
	Error     error
}

func (info *Info) String() string {
	txt := ""
	txt += info.LogTime.Format("2006-01-02T15:04:05.000000-07:30") + " : "
	txt += strconv.Itoa(info.Type) + " : "
	txt += info.Package + "."
	txt += info.Method + " ["
	txt += strconv.Itoa(info.ErrorCode) + "] "
	txt += info.Message
	if info.Error != nil {
		txt += " : " + info.Error.Error()
	}
	txt += ";\n"
	return txt
}

// Logger logs
func Logger(filepath, filename string, infoCh <-chan Info, shutdownCh <-chan bool, wg *sync.WaitGroup) {
	// create dir
	_, err := os.Stat(filepath)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(filepath, 0755)
		if errDir != nil {
			log.Fatal(err)
		}

	}
	fileTimeFormat := "-2006-01-02-15" // TODO: restore in production "-2006-01-02-15-04"
	fileFullName := filepath + "/" + filename + time.Now().Format(fileTimeFormat) + ".log"
	file, err := os.OpenFile(fileFullName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed opening log file: %s", err)
	}
	defer file.Close()

loop:
	for {
		select {
		case info := <-infoCh:
			file.WriteString(info.String())
			if info.Type&STDOUT == 1 {
				fmt.Println(info.String())
			}
		case shutdown := <-shutdownCh:
			log.Println("closing logger", shutdown)
			break loop
		}
	}
	wg.Done()
}
