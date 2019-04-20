package logger

import (
	"fmt"
	"log"
	SysLog "log"
	"os"
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
	LogTime     time.Time
	Type        int // WARN | INFO | ERROR | DEBUG
	FileName    string
	Method      string
	BlockNumber string
	Message     string
	Error       error
}

// Log is dummy struct
type Log struct {
	fileHandle os.File
}

var instance *Log
var logOnce sync.Once

// Logger log
func Logger() *Log {
	logOnce.Do(func() {
		instance = &Log{}
	})
	return instance
}

func AppendFile() {
	file, err := os.OpenFile("test.txt", os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()

	len, err := file.WriteString(" The Go language was conceived in September 2007 by Robert Griesemer, Rob Pike, and Ken Thompson at Google.")
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
	fmt.Printf("\nLength: %d bytes", len)
	fmt.Printf("\nFile Name: %s", file.Name())
}

func (log *Log) Error() {

}

// SetFile open a file in append mode and assign it to fileHandle
func (log *Log) SetFile(fileName string) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		SysLog.Fatalf("failed opening file: %s", err)
	}
	log.fileHandle = *file
}

// Close closes the log file
func (log *Log) Close() {
	log.fileHandle.Close()
}

// TestLogger test
func TestLogger() {
	//Logger().
}
