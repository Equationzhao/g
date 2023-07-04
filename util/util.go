package util

import (
	"github.com/Equationzhao/g/util/cmp"
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
