package riff

import (
	"fmt"
	"github.com/gimke/riff/common"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Services map[string]*Service

type Service struct {
	Name        string `yaml:"name,omitempty"`
	IP          string `yaml:"ip,omitempty"`
	Port        uint16 `yaml:"port,omitempty"`
	Version     uint64
	State       stateType // Current state
	StateChange time.Time // Time last state change happened

	Env        []string `yaml:"env,omitempty"`
	Command    []string `yaml:"command,omitempty"`
	PidFile    string   `yaml:"pid_file,omitempty"`
	StdOutFile string   `yaml:"std_out_file,omitempty"`
	StdErrFile string   `yaml:"std_err_file,omitempty"`
	Grace      bool     `yaml:"grace,omitempty"`
	RunAtLoad  bool     `yaml:"run_at_load,omitempty"`
	KeepAlive  bool     `yaml:"keep_alive,omitempty"`
	*Deploy    `yaml:"deploy,omitempty"`
}

type Deploy struct {
	Provider   string `yaml:"provider,omitempty"`
	Token      string `yaml:"token,omitempty"`
	Repository string `yaml:"repository,omitempty"`
	Version    string `yaml:"version,omitempty"`
	Payload    string `yaml:"payload,omitempty"`
}

func (s *Service) Address() string {
	return net.JoinHostPort(s.IP, strconv.Itoa(int(s.Port)))
}

func (s *Services) Keys() []string {
	var keys = make([]string, 0, 0)
	for key, _ := range *s {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (s *Service) AutoRun() {
	if pid := s.GetPid(); pid == 0 && s.RunAtLoad {
		err := s.Start()
		if err != nil {
			server.Logger.Printf(errorServicePrefix+"%s running error %v", s.Name, err)
		} else {
			server.Logger.Printf(infoServicePrefix+"%s running success", s.Name)
		}
	}
}

func (s *Service) Start() error {
	if s.GetPid() != 0 {
		return fmt.Errorf(errorServicePrefix+"%s is already running", s.Name)
	}
	command := s.resoveCommand()
	dir, _ := filepath.Abs(filepath.Dir(command))

	cmd := exec.Command(command, s.Command[1:]...)
	if len(s.Env) > 0 {
		cmd.Env = append(os.Environ(), s.Env...)
	}
	cmd.Dir = dir

	if s.StdOutFile != "" {
		out := common.MakeFile(s.resovePath(s.StdOutFile))
		cmd.Stdout = out
	} else {
		out := common.MakeFile(common.BinDir + "/logs/" + s.Name + "/stdout.log")
		cmd.Stdout = out
	}
	if s.StdErrFile != "" {
		err := common.MakeFile(s.resovePath(s.StdErrFile))
		cmd.Stderr = err
	} else {
		err := common.MakeFile(common.BinDir + "/logs/" + s.Name + "/stderr.log")
		cmd.Stderr = err
	}

	err := cmd.Start()
	if err != nil {
		return err
	} else {
		go func() {
			cmd.Wait()
		}()
		if s.PidFile == "" {
			s.SetPid(cmd.Process.Pid)
		}
	}
	return nil
}

func (s *Service) Stop() error {
	pid := s.GetPid()
	if pid == 0 {
		return fmt.Errorf(errorServicePrefix+"%s has already been stopped", s.Name)
	} else {
		if p, find := s.processExist(pid); find {
			err := p.Signal(syscall.SIGINT)
			if err != nil {
				return err
			}
			quitStop := make(chan bool)
			go func() {
				for {
					if pid := s.GetPid(); pid == 0 {
						quitStop <- true
						break
					}
					time.Sleep(1 * time.Second)
				}
			}()
			<-quitStop
			if s.PidFile == "" {
				s.RemovePid()
			}
		}
	}
	return nil
}

func (s *Service) SetPid(pid int) {
	pidString := []byte(strconv.Itoa(pid))
	os.MkdirAll(filepath.Dir(s.pidFile()), 0755)
	ioutil.WriteFile(s.pidFile(), pidString, 0666)
}

func (s *Service) RemovePid() error {
	return os.Remove(s.pidFile())
}

func (s *Service) GetPid() int {
	content, err := ioutil.ReadFile(s.pidFile())
	if err != nil {
		return 0
	} else {
		pid, _ := strconv.Atoi(strings.Trim(string(content), "\n"))
		if _, find := s.processExist(pid); find {
			return pid
		} else {
			return 0
		}
	}
}

func (s *Service) resovePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	} else {
		if strings.HasPrefix(path, "."+string(os.PathSeparator)) {
			return common.BinDir + path[1:]
		} else {
			return common.BinDir + "/" + path
		}
	}
}

func (s *Service) resoveCommand() string {
	path := s.Command[0]
	if filepath.IsAbs(path) {
		return path
	} else {
		if strings.HasPrefix(path, "."+string(os.PathSeparator)) {
			return common.BinDir + path[1:]
		} else {
			return path
		}
	}
}

func (s *Service) pidFile() string {
	if s != nil && s.PidFile != "" {
		pid := s.resovePath(s.PidFile)
		return pid
	} else {
		return common.BinDir + "/run/" + s.Name + ".pid"
	}
}

func (s *Service) processExist(pid int) (*os.Process, bool) {
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
