package common

import (
	"math/rand"
	"os"
	"path/filepath"
	"time"
	"strings"
	"fmt"
	"strconv"
	"net"
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

func GenerateID(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GetIpPort(ipPort string) (ip string, port int, err error) {
	index := strings.LastIndex(ipPort,":")
	if ipPort == "" {
		err = fmt.Errorf("empty ip and port\n")
		return
	}
	if index > -1 && index < len(ipPort) {
		ip = ipPort[0:index]
		port,err = strconv.Atoi(ipPort[index+1:])
		if err != nil {
			ip = ipPort
			port = 0
			err = nil
		}
	} else {
		ip = ipPort
		port = 0
		err = nil
	}
	if ip != "" {
		ipaddr := net.ParseIP(ip)
		if ipaddr == nil {
			err = fmt.Errorf("error ip and port\n")
		}
	}
	return
}

