package services

import (
	"sisyphos/models"
)

type sshlib interface {
	Execute(m models.SSHConfig, cmd string) (string, error)
}

type SSHService struct {
	l sshlib
}

func NewSSHService(l sshlib) *SSHService {
	return &SSHService{l}
}

func (s *SSHService) RunCommand(cfg models.SSHConfig, cmd string) (string, error) {
	return s.l.Execute(cfg, cmd)
}
