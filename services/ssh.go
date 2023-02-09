package services

import (
	"sisyphos/models"
)

type sshlib interface {
	Execute(m models.SSHConfig, cmd string) (string, error)
	ReadHostKey(host string, port int) (string, error)
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

func (s *SSHService) ReadHostKey(host string, port int) (string, error) {
	return s.l.ReadHostKey(host, port)
}
