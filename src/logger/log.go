package logger

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

const ( // iota is reset to 0
	// DEBUG signify normal debugging message
	DEBUG = iota // 0
	// INFO signify normal info message
	INFO = iota // 1
	// WARN signify warning message when some assumptions breaks but not harmful
	WARN = iota // 2
	// ERROR signify error and will handled message
	ERROR = iota // 3
	// FATAL signify fatal error, will exit from program
	FATAL = iota // 4
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
	txt += info.Method + " # "
	txt += strconv.Itoa(info.ErrorCode) + " : "
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
	fileFullName := filepath + "/" + filename + time.Now().Format("-2006-01-02-15-04") + ".log"
	file, err := os.OpenFile(fileFullName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()

loop:
	for {
		select {
		case info := <-infoCh:
			file.WriteString(info.String())
		case shutdown := <-shutdownCh:
			log.Println("closing logger", shutdown)
			break loop
		}
	}
	wg.Done()
}
