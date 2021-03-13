package cli

import (
	"github.com/teatak/riff/common"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"
)

func SetPid(pid int) {
	pidString := []byte(strconv.Itoa(pid))
	_ = os.MkdirAll(common.BinDir+"/run", 0755)
	_ = ioutil.WriteFile(common.BinDir+"/run/riff.pid", pidString, 0666)
}

func GetPid() int {
	content, err := ioutil.ReadFile(common.BinDir + "/run/riff.pid")
	if err != nil {
		return 0
	} else {
		pid, _ := strconv.Atoi(strings.Trim(string(content), "\n"))
		if _, find := ProcessExist(pid); find {
			return pid
		} else {
			return 0
		}
	}
}

func ProcessExist(pid int) (*os.Process, bool) {
	process, err := os.FindProcess(pid)
	if err != nil {
		return nil, false
	} else {
		err := process.Signal(syscall.Signal(0))
		if err != nil {
			return nil, false
		}
	}
	return process, true
}
