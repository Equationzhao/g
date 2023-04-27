package pathbeautify

import "testing"

func TestTransform(t *testing.T) {
	testcase := []string{"...", "...\\a\\b\\c", "...\\", ".", "~", "..", "~\\a"}
	for _, s := range testcase {
		Transform(&s)
		t.Log(s)
	}
}
