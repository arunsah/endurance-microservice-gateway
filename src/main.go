package main

import (
	"flag"
	"sync"

	Logger "./logger"
	ProxyServer "./proxyserver"
)

var infoCh chan Logger.Info
var shutdownCh chan bool
var wg sync.WaitGroup

func initLogger(logpath, logfile string) {
	if infoCh == nil {
		infoCh = make(chan Logger.Info, 1000)
	}
	if shutdownCh == nil {
		shutdownCh = make(chan bool)
	}

	wg.Add(1)
	go Logger.Logger(logpath, logfile, infoCh, shutdownCh, &wg)

}

func main() {
	// parse command line args
	// TODO: later these flags should be read from environment variables
	host, port, version, message, logpath, logfile := parseFlags()

	// init logger
	initLogger(logpath, logfile)

	// start server in terminals
	// TODO: later start the server in service/deamon mode
	proxyServerInfo := ProxyServer.ProxyServer{Host: host, Port: port, Version: version, Message: message,
		Logpath: logpath, Logfile: logfile, LogInfoCh: infoCh, LogShutdownCh: shutdownCh}
	proxyServerInfo.RunServer()

}

// parse command line args
func parseFlags() (host, port, version, message, logpath, logfile string) {
	Host := flag.String("h", "127.0.0.1", "host on which this application will listen")
	Port := flag.String("p", "8008", "port on which this application will listen")
	Version := flag.String("v", "1.0.1", "version of the application [test features]")
	Message := flag.String("m", "Hello World from Endurance-Microservice-Gateway", "custom message for the application [test features]")
	Logpath := flag.String("logdir", "./log", "log file path")
	Logfile := flag.String("logfile", "webserver", "log file name prefix")

	flag.Parse()
	host = *Host
	port = ":" + *Port // :8008
	version = *Version
	message = *Message
	logpath = *Logpath
	logfile = *Logfile
	return host, port, version, message, logpath, logfile
}
