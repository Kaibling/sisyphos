package models

import (
	"time"

	"sisyphos/lib/utils"
)

type Run struct {
	// Action     ActionExt `json:"action,omitempty"`
	Action    string    `json:"action"`
	Host      *string   `json:"host"`
	RequestID string    `json:"request_id"`
	ParentID  string    `json:"parent_id"`
	User      string    `json:"user"`
	RunID     string    `json:"run_id"`
	StartTime time.Time `json:"start_date"`
	EndTime   time.Time `json:"end_date"`
	Duration  string    `json:"duration"`
	Output    string    `json:"output"`
	Error     string    `json:"error"`
	Status    string    `json:"status"`
}

func NewRun(actionName, username, reqID, parentID string) *Run {
	return &Run{
		Action:    actionName,
		StartTime: time.Now(),
		RunID:     utils.NewULID().String(),
		User:      username,
		RequestID: reqID,
		ParentID:  parentID,
	}
}

func (r *Run) SetEndTime() {
	r.EndTime = time.Now()
	r.Duration = r.EndTime.Sub(r.StartTime).String()
}
