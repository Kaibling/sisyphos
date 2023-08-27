package cron

import "github.com/adhocore/gronx"

func Validate(s string) bool {
	gron := gronx.New()
	return gron.IsValid(s)
}
