package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	Logger "./logger"
)

func main() {
	// parse command line args
	// TODO: later these flags should be read from environment variables
	host, port, version, message, logpath := parseFlags()
	fmt.Println(host, port, version, message, logpath)

	// start server in terminals
	// TODO: later start the server in service/deamon mode
	startServer(host, port, version, message, logpath)
	ll := &Logger.Info{Type: Logger.WARN, Message: "Error ocurred"}
	fmt.Println(ll)
	fmt.Println(time.Now())

	infoCh := make(chan Logger.Info, 1000)
	shutdownCh := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(1)
	go Logger.Logger("txt1.txt", infoCh, shutdownCh, &wg)
	//infoCh <- *ll
	reader := bufio.NewReader(os.Stdin)
	for {
		txt, _ := reader.ReadString('\n')
		txt = strings.TrimSpace(txt)
		if strings.ToLower(txt) == "shutdown" {
			shutdownCh <- true
			break
		}
		infoCh <- Logger.Info{LogTime: time.Now(), Type: Logger.INFO,
			Package: "main", Method: "main()", ErrorCode: "0", Message: txt, Error: errors.New("no-error")}
		//fmt.Println(txt)
	}

}

// parse command line args
func parseFlags() (host, port, version, message, logpath string) {
	Host := flag.String("h", "127.0.0.1", "host on which this application will listen")
	Port := flag.String("p", "8008", "port on which this application will listen")
	Version := flag.String("v", "1.0.1", "version of the application [test features]")
	Message := flag.String("m", "Hello World from Endurance-Microservice-Gateway", "custom message for the application [test features]")
	Logpath := flag.String("l", "./log", "log file path")

	flag.Parse()
	host = *Host
	port = ":" + *Port // :8008
	version = *Version
	message = *Message
	logpath = *Logpath
	return host, port, version, message, logpath
}

func startServer(host, port, version, message, logpath string) {

}
