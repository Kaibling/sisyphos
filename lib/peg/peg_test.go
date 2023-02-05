package peg

import (
	"testing"
)

func TestInvalidCases(t *testing.T) {
	tc := "tag:all tag:ssh"

	got, err := Parse("", []byte(tc), GlobalStore("join", "joins"))
	if err != nil {
		t.Errorf("%q: want no error, got %v", tc, err)
	}

	t.Error(got)
}
