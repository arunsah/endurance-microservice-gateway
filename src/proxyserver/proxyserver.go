package proxyserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	Logger "../logger"
)

// ProxyServer holds information related to the proxy server
type ProxyServer struct {
	Host          string
	Port          string
	Version       string
	Message       string
	Logpath       string
	Logfile       string
	LogInfoCh     chan Logger.Info
	LogShutdownCh chan bool
}

// RunServer init server and start it
func (proxyServer *ProxyServer) RunServer() {
	proxyServer.LogInfoCh <- Logger.Info{LogTime: time.Now(), Type: Logger.INFO | Logger.STDOUT, Package: "proxyserver",
		Method: "RunServer()", ErrorCode: 0, Message: "Starting Proxy Server at " + proxyServer.Host + proxyServer.Port, Error: nil}
	server := &http.Server{
		Addr:           proxyServer.Host + proxyServer.Port,
		Handler:        proxyServer,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := server.ListenAndServe()
	proxyServer.LogInfoCh <- Logger.Info{LogTime: time.Now(), Type: Logger.FATAL | Logger.STDOUT, Package: "proxyserver",
		Method: "RunServer()", ErrorCode: 2, Message: "Stopping Proxy Server at " + proxyServer.Host + proxyServer.Port, Error: err}
	proxyServer.LogShutdownCh <- true
	time.Sleep(1 * time.Second)
}

// This is redundant, should directly use the map
func (proxyServer *ProxyServer) parseURL(query url.Values) (port, service, version, msg string) {
	port = ":" + strings.TrimSpace(query["port"][0])
	service = strings.TrimSpace(query["service"][0])
	version = strings.TrimSpace(query["version"][0])
	msg = strings.TrimSpace(query["msg"][0])
	// req.URL.Path >> /inventory
	// req.URL.String() >>
	// /inventory?port=9001&service=inventory&version=v1&version=v2&version=v3&msg=helloworld
	// {"msg":["helloworld"],"port":["9001"],"service":["inventory"],"version":["v1","v2","v3"]}

	return port, service, version, msg
}

func (proxyServer *ProxyServer) queryToJSON(urlValue url.Values) []byte {
	json, err := json.Marshal(urlValue)
	if err != nil {
		proxyServer.LogInfoCh <- Logger.Info{LogTime: time.Now(), Type: Logger.ERROR | Logger.STDOUT, Package: "proxyserver",
			Method: "queryToJSON()", ErrorCode: 0, Message: "Error json.Marshal", Error: err}
	}
	return json
}

// staticFileHandler
func staticFileHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../assets/favicon.ico")
}

// ServeHTTP server the web request
func (proxyServer *ProxyServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	jsonQueryString := string(proxyServer.queryToJSON(req.URL.Query()))

	proxyServer.LogInfoCh <- Logger.Info{LogTime: time.Now(), Type: Logger.INFO | Logger.STDOUT, Package: "proxyserver",
		Method: "ServeHTTP()", ErrorCode: 0, Message: "Client connected #" + req.RemoteAddr, Error: nil}
	proxyServer.LogInfoCh <- Logger.Info{LogTime: time.Now(), Type: Logger.INFO | Logger.STDOUT, Package: "proxyserver",
		Method: "ServeHTTP()", ErrorCode: 0, Message: jsonQueryString, Error: nil}

	// Client connected #  [::1]:15533 /inventory /inventory?port=9001&service=inventory&version=v1&version=v2&version=v3&msg=helloworld {"msg":["helloworld"],"port":["9001"],"service":["inventory"],"version":["v1","v2","v3"]}

	if req.URL.String() == "/favicon.ico" { //strings.Contains(req.URL.Path, ".") {
		log.Println(req.URL.Path, req.URL.String(), req.RemoteAddr, "serving /favicon.ico")
		staticFileHandler(res, req)
		return
	}

	port, service, version, msg := proxyServer.parseURL(req.URL.Query())
	txt := ""
	txt = txt + "port: " + port + "\n"
	txt = txt + "service: " + service + "\n"
	txt = txt + "version: " + version + "\n"
	txt = txt + "msg: " + msg + "\n"
	txt = txt + "Path: " + req.URL.Path + "\n"

	fmt.Fprintf(res, "%s\n", txt)

	result, err := dialServer("localhost", port, service, version)
	if err != nil {
		log.Fatalln("ServeHTTP:ERROR1: ", err)
		//log.Fatalln("Retruing :")
		//result, err = dialServer("localhost", ":9001", service, version)
	}
	fmt.Fprintf(res, "%s\n", result)

	log.Println("Request served :", req.RemoteAddr)
}

func dialServer(host, port, service, version string) (string, error) {

	addr := "http://" + host + port + "/" + service
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
		MaxIdleConns:        100,
		IdleConnTimeout:     30 * time.Second,
		DisableCompression:  true,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	var netClient = &http.Client{
		Timeout:   time.Second * 50,
		Transport: netTransport,
	}

	resp, err := netClient.Get(addr)
	if err != nil {
		fmt.Println("dialServer:ERROR1", err)
		//return "", err
		addr := "http://" + host + ":9001" + "/" + service
		log.Println("Retrying inside DIAL with ", addr)
		resp, err = netClient.Get(addr)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("dialServer:ERROR2", err)
		return "", err
	}
	//fmt.Println(string(body))
	result := string(body)
	return result, nil
}
