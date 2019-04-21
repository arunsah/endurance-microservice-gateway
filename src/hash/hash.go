package hash

import (
	"time"

	Logger "../logger"
	"golang.org/x/crypto/bcrypt"
)

// Hash holds information related to the registry
type Hash struct {
	LogInfoCh     chan Logger.Info
	LogShutdownCh chan bool
}

// HashAndSalt takes byte slice and return hashed password
func (hash *Hash) HashAndSalt(pwd []byte) string {
	cost := bcrypt.DefaultCost // 10
	hashed, err := bcrypt.GenerateFromPassword(pwd, cost)
	if err != nil {
		hash.LogInfoCh <- Logger.Info{LogTime: time.Now(), Type: Logger.ERROR | Logger.STDOUT, Package: "hash",
			Method: "HashAndSalt()", ErrorCode: -1, Message: "Request served to ", Error: err}

	}
	hash.LogInfoCh <- Logger.Info{LogTime: time.Now(), Type: Logger.ERROR | Logger.STDOUT, Package: "hash",
		Method: "HashAndSalt()", ErrorCode: -1, Message: "Hashed assword: " + string(hashed), Error: nil}
	return string(hashed)
}

// ComparePasswords takes a hashedPwd (from  DB) and a plainPwd as byte slice and compares
func (hash *Hash) ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		hash.LogInfoCh <- Logger.Info{LogTime: time.Now(), Type: Logger.ERROR | Logger.STDOUT, Package: "hash",
			Method: "ComparePasswords()", ErrorCode: 0, Message: "Password didnot matched", Error: err}
		return false
	}
	hash.LogInfoCh <- Logger.Info{LogTime: time.Now(), Type: Logger.INFO | Logger.STDOUT, Package: "hash",
		Method: "ComparePasswords()", ErrorCode: 0, Message: "Password matched", Error: nil}
	return true
}
