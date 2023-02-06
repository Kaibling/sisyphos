package models

import (
	"encoding/json"
	"fmt"
)

type Host struct {
	Name     *string  `json:"name"`
	Username *string  `json:"username"`
	Password *string  `json:"password"`
	SSHKey   *string  `json:"ssh_key"`
	Address  *string  `json:"address"`
	Port     *int     `json:"port"`
	Tags     []string `json:"tags"`
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
