package cached

import (
	constval "github.com/Equationzhao/g/const"
	"github.com/Equationzhao/pathbeautify"
	"github.com/alphadose/haxmap"
)

func GetUserHomeDir() string {
	return pathbeautify.GetUserHomeDir()
}

type Map[k constval.Hashable, v any] struct {
	*haxmap.Map[k, v]
}

func NewCacheMap[k constval.Hashable, v any](len int) *Map[k, v] {
	return &Map[k, v]{
		haxmap.New[k, v](uintptr(len)),
	}
}

func (m Map[k, v]) Keys() []k {
	keys := make([]k, 0, m.Len())
	m.ForEach(func(k k, v v) bool {
		keys = append(keys, k)
		return true
	})
	return keys
}

func (m Map[k, v]) Values() []v {
	values := make([]v, 0, m.Len())
	m.ForEach(func(k k, v v) bool {
		values = append(values, v)
		return true
	})
	return values
}

func (m Map[k, v]) Pairs() []Pair[k, v] {
	pairs := make([]Pair[k, v], 0, m.Len())
	m.ForEach(func(key k, value v) bool {
		pairs = append(pairs, Pair[k, v]{
			First:  key,
			Second: value,
		})
		return true
	})
	return pairs

}

// Pair is a struct that contains two variables ptr
type Pair[T, U any] struct {
	First  T
	Second U
}

// MakePair return a new Pair
// receive two value
func MakePair[T, U any](first T, second U) Pair[T, U] {
	return Pair[T, U]{
		First:  first,
		Second: second,
	}
}

// Set the pair
// Copy the `first` and `second` to the pair
func (p *Pair[T, U]) Set(first T, second U) {
	p.First = first
	p.Second = second
}

func (p Pair[T, U]) Key() T {
	return p.First
}

func (p Pair[T, U]) Value() U {
	return p.Second
}
