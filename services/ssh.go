package services

import "fmt"

type sshlib interface {
	Execute(hoststring, uersname, password, cmd string) (string, error)
}

type SSHConfig struct {
	Address  string
	Port     string
	Username string
	Password string
	Key      string
}

type SSHService struct {
	l sshlib
}

func NewSSHService(l sshlib) *SSHService {
	return &SSHService{l}
}

func (s *SSHService) RunCommand(cfg SSHConfig, cmd string) (string, error) {
	fmt.Println(cmd)
	return s.l.Execute(cfg.Address+":"+cfg.Port, cfg.Username, cfg.Password, cmd)
}
