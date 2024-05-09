package util

import (
	"sync"
	"testing"
)

func TestSlice_AppendTo(t *testing.T) {
	s := NewSlice[int](10)
	s.AppendTo(1)
	s.AppendTo(2)
	s.AppendTo(3)
	if s.Len() != 3 {
		t.Errorf("AppendTo failed")
	}
}

func TestSlice_At(t *testing.T) {
	s := Slice[int]{
		data: []int{1, 2, 3},
		m:    sync.RWMutex{},
	}
	if s.At(1) != 2 {
		t.Errorf("At failed")
	}
}

func TestSlice_Clear(t *testing.T) {
	s := Slice[int]{
		data: []int{1, 2, 3},
		m:    sync.RWMutex{},
	}
	s.Clear()
	if s.Len() != 0 {
		t.Errorf("Clear failed")
	}
}

func TestSlice_GetRaw(t *testing.T) {
	s := Slice[int]{
		data: []int{1, 2, 3},
		m:    sync.RWMutex{},
	}
	gotRaw := s.GetRaw()
	if len(*gotRaw) != 3 {
		t.Errorf("GetRaw failed")
	}
	if gotRaw != &s.data {
		t.Errorf("GetRaw failed")
	}
}

func TestSlice_GetCopy(t *testing.T) {
	s := Slice[int]{
		data: []int{1, 2, 3},
		m:    sync.RWMutex{},
	}
	copied := s.GetCopy()
	if len(copied) != 3 {
		t.Errorf("GetCopy failed")
	}
	if &copied == &s.data {
		t.Errorf("GetCopy failed")
	}
}

func TestSlice_Len(t *testing.T) {
	s := Slice[int]{
		data: []int{1, 2, 3},
		m:    sync.RWMutex{},
	}
	if s.Len() != 3 {
		t.Errorf("Len failed")
	}
}

func TestSlice_Set(t *testing.T) {
	s := Slice[int]{
		data: []int{1, 2, 3},
		m:    sync.RWMutex{},
	}
	s.Set(1, 4)
	if s.At(1) != 4 {
		t.Errorf("Set failed")
	}
}
