package models

import "time"

type Token struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires_at"`
}
