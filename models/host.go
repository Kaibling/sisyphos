package models

import (
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

type OrderedHost struct {
	Name  string `json:"name"`
	Order int    `json:"order"`
}
