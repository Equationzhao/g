package util

import "sync"

type Slice[T any] struct {
	data []T
	m    sync.RWMutex
}

func (s *Slice[T]) Clear() {
	s.m.Lock()
	s.data = s.data[:0]
	s.m.Unlock()
}

func (s *Slice[T]) AppendTo(d T) {
	s.m.Lock()
	s.data = append(s.data, d)
	s.m.Unlock()
}

func (s *Slice[T]) GetRaw() *[]T {
	return &s.data
}

func (s *Slice[T]) GetCopy() []T {
	s.m.RLock()
	defer s.m.RUnlock()
	copied := make([]T, len(s.data))
	copy(copied, s.data)
	return copied
}

func (s *Slice[T]) At(pos int) T {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.data[pos]
}

func (s *Slice[T]) Len() int {
	s.m.RLock()
	defer s.m.RUnlock()
	return len(s.data)
}

func (s *Slice[T]) Set(pos int, d T) {
	s.m.Lock()
	defer s.m.Unlock()
	s.data[pos] = d
}

func NewSlice[T any](size int) *Slice[T] {
	s := &Slice[T]{
		data: make([]T, 0, size),
	}
	return s
}
