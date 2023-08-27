package models

import (
	"errors"
	"time"

	"sisyphos/lib/utils"
)

type Action struct {
	DBInfo
	Name         *string                `json:"name"`
	Groups       []string               `json:"groups"`
	ScheduleExpr *string                `json:"schedule_expression"`
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
	if a.Variables == nil {
		a.Variables = map[string]any{}
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = time.Now()
	}
	if a.UpdatedAt.IsZero() {
		a.UpdatedAt = time.Now()
	}
}

type OrderdAction struct {
	Name  string `json:"name"`
	Order int    `json:"order"`
}
