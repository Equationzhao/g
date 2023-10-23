package util

import (
	"path/filepath"
	"strings"
)

func RemoveSep(s string) string {
	return strings.TrimRight(s, string(filepath.Separator))
}
