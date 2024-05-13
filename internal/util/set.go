package util

import (
	constval "github.com/Equationzhao/g/internal/global"
	"github.com/alphadose/haxmap"
)

type SafeSet[T constval.Hashable] struct {
	internal *haxmap.Map[T, struct{}]
}

func (s *SafeSet[T]) Add(k T) {
	s.internal.Set(k, struct{}{})
}

func (s *SafeSet[T]) Contains(k T) bool {
	_, t := s.internal.Get(k)
	return t
}

func NewSet[T constval.Hashable]() *SafeSet[T] {
	return &SafeSet[T]{
		internal: haxmap.New[T, struct{}](10),
	}
}
