package osbased

import (
	"os"
	"testing"
)

func TestOwnerID(t *testing.T) {
	s, _ := os.Stat("time_linux.go")

	got := OwnerID(s)
	println(got)
}
