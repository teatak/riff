//go:build windows || solaris
// +build windows solaris

package riff

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/teatak/riff/common"
)

func (s *Service) Restart() error {
	pid := s.GetPid()
	if pid != 0 {
		err := s.Stop()
		if err != nil {
			return err
		} else {
			err = s.Start()
			if err != nil {
				return err
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
		out := common.MakeFile(common.BinDir + "/logs/" + s.Name + "/stderr.log")
		cmd.Stderr = out
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
			err := p.Kill()
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
