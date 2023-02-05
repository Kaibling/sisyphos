package models

import (
	"encoding/json"
	"fmt"
)

type Host struct {
	Name    string   `json:"name"`
	SSHKey  string   `json:"ssh_key"`
	Address string   `json:"address"`
	Tags    []string `json:"tags"`
}
type Connection struct {
	Host
	Port  string `json:"port"`
	Order int    `json:"order"`
}

type Service struct {
	HostName string `json:"host_name"`
	Port     string `json:"port"`
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
