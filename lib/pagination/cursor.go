package pagination

import (
	"fmt"
	"strings"
)

type Cursor struct {
	Before *string
	After  *string
}

func (c *Cursor) CreateAfter(field string, ids ...string) {
	if len(ids) == 1 {
		c.After = ToPointer(fmt.Sprintf("AFTER|%s|%s", ids[0], field))
	}
	if len(ids) == 2 {
		c.After = ToPointer(fmt.Sprintf("AFTER|%s;%s|%s", ids[0], ids[1], field))
	}
}
func (c *Cursor) CreateBefore(field string, ids ...string) {
	if len(ids) == 1 {
		c.Before = ToPointer(fmt.Sprintf("BEFORE|%s|%s", ids[0], field))
	}
	if len(ids) == 2 {
		c.Before = ToPointer(fmt.Sprintf("BEFORE|%s;%s|%s", ids[0], ids[1], field))
	}
}
func (c *Cursor) Finish() {
	// TODO convert after and before to base64
}

func ParseCursor(s string) CursorInfo {
	parts := strings.Split(s, "|")
	ci := CursorInfo{}
	// TODO verify
	ci.Direction = parts[0]
	ci.SortField = parts[2]
	ids := strings.Split(parts[1], ";")
	ci.PrimaryId = ids[0]
	if len(ids) > 1 {
		ci.SortId = ids[1]
	}
	return ci
}
