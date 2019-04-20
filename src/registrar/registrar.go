package registrar

import (
	"fmt"

	Logger "../logger"
)

// Registrar holds information related to the registry
type Registrar struct {
	Registry      map[string]map[uint]AppInfo
	Logpath       string
	Logfile       string
	LogInfoCh     chan Logger.Info
	LogShutdownCh chan bool
}

// AppInfo is primary data structure used by registrar and application to exchange information
//  about the service they provide.
type AppInfo struct {
	Host    string // "localhost"
	Port    string // "8888"
	Name    string // service name, all workers providing same service will have same Name
	Version string // service version, all workers providing same service will have same Name
	AppID   uint   // generated by registrar upon registration, unique for each workers and service
	Secret  string // secret session token will be used by service
	// we should not use Go dependent complex type such as time as we will require to receive these param from other systems
}

// Registry is map of map. Outer map will take Service Name as Key and inner map will take AppID as key.
// AppID will be generated by registrar

// func initMaps() {
// 	// if WorkerCount == nil {
// 	// 	WorkerCount = make(map[string]uint)
// 	// }
// 	if Registry == nil {
// 		Registry = make(map[string]map[uint]AppInfo)
// 	}
// }

func generateID() func() uint {
	id := uint(0)
	return func() uint {
		id++
		return id
	}
}

var nextAppID func() uint

// Register will add appData in registry and will return modified appData
func (registrar *Registrar) Register(appName string, appData AppInfo) AppInfo {
	if nextAppID == nil {
		nextAppID = generateID()
	}
	//initMaps()
	//WorkerCount[appName]++
	appData.AppID = nextAppID() // WorkerCount[appName]

	if registrar.Registry[appName] == nil {
		registrar.Registry[appName] = map[uint]AppInfo{}
	}
	registrar.Registry[appName][appData.AppID] = appData
	return appData
}

// UnRegister will un-register and remove the entry of appData and will return appData
func (registrar *Registrar) UnRegister(appName string, appID uint) AppInfo {
	//initMaps()
	appData, found := registrar.Registry[appName][appID]
	if !found {
		return AppInfo{}
	}
	delete(registrar.Registry[appName], appID)
	return appData
}

// HelloRegistrar says Hello
func (registrar *Registrar) HelloRegistrar() {
	fmt.Println("Hello from registrar")
}