package models

import (
	"sisyphos/lib/utils"
	"time"
)

type Run struct {
	ID        string    `json:"id"`
	Action    string    `json:"action"`
	Host      *string   `json:"host"`
	RequestID string    `json:"request_id"`
	ParentID  string    `json:"parent_id,omitempty"`
	Childs    []*Run    `json:"childs,omitempty"`
	User      string    `json:"user"`
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
		ID:        utils.NewULID().String(),
		User:      username,
		RequestID: reqID,
		ParentID:  parentID,
	}
}

func (r *Run) SetEndTime() {
	r.EndTime = time.Now()
	r.Duration = r.EndTime.Sub(r.StartTime).String()
}
