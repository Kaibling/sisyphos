package models

import "errors"

type SSHConfig struct {
	Address    string
	Port       int
	Username   string
	Password   string
	PrivateKey string
	KnownKey   string
}

func (s *SSHConfig) Validate() error {
	if s.Address == "" {
		return errors.New("host address missing")
	}
	if s.Port == 0 {
		return errors.New("host port missing")
	}
	if s.Username == "" {
		return errors.New("host username missing")
	}
	return nil
}
