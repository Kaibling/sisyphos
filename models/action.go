package models

import (
	"errors"
	"sisyphos/lib/utils"
)

type Action struct {
	Name         *string                `json:"name"`
	Groups       []string               `json:"groups"`
	Script       *string                `json:"script"`
	Tags         []string               `json:"tags"`
	Actions      []OrderdAction         `json:"actions"`
	Hosts        []OrderedHost          `json:"hosts"`
	Variables    map[string]interface{} `json:"variables"`
	FailOnErrors *bool                  `json:"fail_on_errors"`
}

func (a *Action) Validate() error {
	if a.Name == nil || *a.Name == "" {
		return errors.New("action name missing")
	}
	return nil
}

func (a *Action) Default() {
	if a.FailOnErrors == nil {
		a.FailOnErrors = utils.ToPointer(true)
	}
}

type ActionExt struct {
	Name         *string
	Groups       []string `json:"groups"`
	Script       *string
	Triggers     []OrderdActionExt
	Tags         []string
	Hosts        []OrderedHost
	Variables    map[string]interface{}
	FailOnErrors *bool `json:"fail_on_errors"`
}

type OrderdAction struct {
	Action string `json:"action"`
	Order  int    `json:"order"`
}

type OrderdActionExt struct {
	ActionExt string `json:"action"`
	Order     int    `json:"order"`
}
