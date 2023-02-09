package models

import (
	"encoding/json"
	"fmt"
	"sisyphos/lib/utils"
)

type Host struct {
	Name     *string  `json:"name"`
	Username *string  `json:"username"`
	Password *string  `json:"password"`
	SSHKey   *string  `json:"ssh_key"`
	KnownKey *string  `json:"known_key"`
	Address  *string  `json:"address"`
	Port     *int     `json:"port"`
	Tags     []string `json:"tags"`
}

func (h *Host) ToSSHConfig() SSHConfig {
	return SSHConfig{
		Address:    utils.PtrRead(h.Address),
		Port:       utils.PtrRead(h.Port),
		Username:   utils.PtrRead(h.Username),
		Password:   utils.PtrRead(h.Password),
		PrivateKey: utils.PtrRead(h.SSHKey),
		KnownKey:   utils.PtrRead(h.KnownKey),
	}
}

type Connection struct {
	Host
	Order int `json:"order"`
}

type Service struct {
	HostName string `json:"name"`
	Order    int    `json:"order,omitempty"`
}

func (s *Service) FromJson(d map[string]interface{}) {
	b, err := json.Marshal(d)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = json.Unmarshal(b, s)
	if err != nil {
		fmt.Println(err.Error())
	}
}
