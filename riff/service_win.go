// +build windows

package riff

import "syscall"

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
