package common

import (
	"math/rand"
	"os"
	"path/filepath"
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

func GenerateID(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
