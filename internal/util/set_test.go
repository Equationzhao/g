package util

import (
	"testing"
)

func TestSafeSet(t *testing.T) {
	set := NewSet[string]()
	set.Add("name")
	if !set.Contains("name") {
		t.Errorf("Add failed")
	}
}
