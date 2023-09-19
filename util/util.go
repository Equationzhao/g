package util

import (
	"cmp"
	"path/filepath"
	"strings"
)

func Max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Min[T cmp.Ordered](a, b T) T {
	if a > b {
		return b
	}
	return a
}

func RemoveSep(s string) string {
	return strings.TrimRight(s, string(filepath.Separator))
}
