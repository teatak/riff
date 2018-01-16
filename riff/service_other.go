// +build !windows

package riff

import "syscall"

func (s *Service) Restart() error {
	pid := s.GetPid()
	if pid != 0 {
		if s.Config.Grace {
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
