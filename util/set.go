package util

import (
	"unsafe"

	"github.com/alphadose/haxmap"
	"golang.org/x/exp/constraints"
)

type hashable interface {
	constraints.Integer | constraints.Float | constraints.Complex | ~string | uintptr | ~unsafe.Pointer
}

type SafeSet[T hashable] struct {
	internal *haxmap.Map[T, struct{}]
}

func (s *SafeSet[T]) Add(k T) {
	s.internal.Set(k, struct{}{})
}

func (s *SafeSet[T]) Contains(k T) bool {
	_, t := s.internal.Get(k)
	return t
}

func NewSet[T hashable]() *SafeSet[T] {
	return &SafeSet[T]{
		internal: haxmap.New[T, struct{}](10),
	}
}
