package models

import "time"

type Log struct {
	DBInfo
	Url       string    `json:"name"`
	Body      string    `json:"body"`
	Method    string    `json:"metoh"`
	User      string    `json:"user"`
	RequestID string    `json:"request_id"`
	CreatedAt time.Time `gorm:"created_at"`
}
