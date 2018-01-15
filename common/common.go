package common

import (
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var BinDir string

func init() {
	rand.Seed(time.Now().UnixNano())
	bin, _ := os.Executable()
	realPath, err := os.Readlink(bin)
	if err == nil {
		bin = realPath
	}
	if filepath.Base(bin) == Name {
		BinDir = filepath.Dir(bin)
	} else {
		BinDir, _ = os.Getwd()
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func MakeFile(path string) *os.File {
	dir := filepath.Dir(path)
	if !IsExist(dir) {
		os.MkdirAll(dir, 0755)
	}
	file, _ := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	return file
}

func IsExist(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

func GenerateID(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandomNumber(n int) int {
	if n == 0 {
		return 0
	}
	return int(rand.Uint32() % uint32(n))
}

func GetIpPort(ipPort string) (ip string, port int) {
	var err error
	index := strings.LastIndex(ipPort, ":")
	if index > -1 && index < len(ipPort) {
		ip = ipPort[0:index]
		port, err = strconv.Atoi(ipPort[index+1:])
		if err != nil {
			ip = ipPort
			port = 0
		}
	} else {
		ip = ipPort
		port = 0
	}
	return
}
