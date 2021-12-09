//go:build !windows
// +build !windows

package riff

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/teatak/riff/common"
)

func (s *Service) Restart() error {
	pid := s.GetPid()
	if pid != 0 {
		if s.Grace {
			if p, find := s.processExist(pid); find {
				err := p.Signal(syscall.SIGUSR2)
				if err != nil {
					return err
				}
			}
		} else {
			err := s.Stop()
			if err != nil {
				return err
			} else {
				err = s.Start()
				if err != nil {
					return err
				}
			}
		}
	} else {
		err := s.Start()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) Start() error {
	if s.GetPid() != 0 {
		return fmt.Errorf(errorServicePrefix+"%s is already running", s.Name)
	}
	command := s.resoveCommand()

	dir, _ := filepath.Abs(filepath.Dir(command))

	cmd := exec.Command(command, s.Command[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if len(s.Env) > 0 {
		cmd.Env = append(os.Environ(), s.Env...)
	}
	if s.Dir != "" {
		dir = s.Dir
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
			s.StartTime = time.Now()
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
			pgid, err := syscall.Getpgid(pid)
			if err == nil {
				if pgid == pid {
					//if pid == pgid
					err := syscall.Kill(-pgid, syscall.SIGKILL)
					if err != nil {
						return err
					}
				} else {
					err := p.Kill()
					if err != nil {
						return err
					}
				}
			} else {
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
