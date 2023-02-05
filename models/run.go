package models

import (
	"time"

	"sisyphos/lib/utils"
)

type Run struct {
	// Action     ActionExtendedv3 `json:"action,omitempty"`
	Action    string    `json:"action"`
	RequestID string    `json:"request_id"`
	User      string    `json:"user"`
	RunID     string    `json:"run_id"`
	StartTime time.Time `json:"start_date"`
	EndTime   time.Time `json:"end_date"`
	Duration  string    `json:"duration"`
	Output    string    `json:"output"`
	Error     string    `json:"error"`
}

func NewRun(actionName string, username string, reqID string) *Run {
	return &Run{
		Action:    actionName,
		StartTime: time.Now(),
		RunID:     utils.NewULID().String(),
		User:      username,
		RequestID: reqID,
	}
}

func (r *Run) SetEndTime() {
	r.EndTime = time.Now()
	r.Duration = r.EndTime.Sub(r.StartTime).String()
}
