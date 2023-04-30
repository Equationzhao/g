package cached

import (
	"fmt"
	"os"
	"sync"
)

var (
	syncHomedir sync.Once
	userHomeDir string
)

func GetUserHomeDir() string {
	syncHomedir.Do(func() {
		userHomeDir, _ = os.UserHomeDir()
	})
	return userHomeDir
}

type Shard[k comparable, v any] struct {
	Lock        sync.RWMutex
	InternalMap map[k]*v
}

type Map[k comparable, v any] struct {
	m      []*Shard[k, v]
	shards int
	New    func(k) *v
	Hash   func(k) uint32
}

func (p *Map[k, v]) Load(key k) (value v, ok bool) {
	return p.GetCopy(key)
}

func (p *Map[k, v]) Store(key k, value v) {
	p.SetCopy(key, value)
}

func DefaultHash[k comparable](i k) uint32 {
	var hash uint32 = 5381 // magic constant, apparently this hash fewest collisions possible.
	data := fmt.Sprint(i)
	for _, c := range data {
		hash = ((hash << 5) + hash) + uint32(c)
	}
	return hash
}

const limit = 200

func (p *Map[k, v]) Free() {
	if p.shards > limit {
		wg := sync.WaitGroup{}
		wg.Add(p.shards)
		for _, m := range p.m {
			m := m
			go func() {
				m.Lock.Lock()
				m.InternalMap = make(map[k]*v)
				m.Lock.Unlock()
				wg.Done()
			}()
		}
		wg.Wait()
	} else {
		for _, m := range p.m {
			m.Lock.Lock()
			m.InternalMap = make(map[k]*v)
			m.Lock.Unlock()
		}
	}
}

func (p *Map[k, v]) Keys() []k {
	mu := new(sync.Mutex)
	length := 0
	for i := range p.m {
		length += len(p.m[i].InternalMap)
	}
	keys := make([]k, 0, length)
	wg := new(sync.WaitGroup)
	wg.Add(p.shards)
	for _, m := range p.m {
		m := m
		go func() {
			m.Lock.RLock()
			defer wg.Done()
			keysi := make([]k, 0, len(m.InternalMap))
			for k := range m.InternalMap {
				keysi = append(keysi, k)
			}
			mu.Lock()
			keys = append(keys, keysi...)
			mu.Unlock()
			m.Lock.RUnlock()
		}()
	}
	wg.Wait()
	return keys
}

func (p *Map[k, v]) SetCopy(key k, value v) {
	i := p.Hash(key) & uint32(p.shards-1)
	p.m[i].Lock.Lock()
	defer p.m[i].Lock.Unlock()
	p.m[i].InternalMap[key] = &value
}

func (p *Map[k, v]) Set(key k, valuePtr *v) {
	i := p.Hash(key) & uint32(p.shards-1)
	p.m[i].Lock.Lock()
	defer p.m[i].Lock.Unlock()
	p.m[i].InternalMap[key] = valuePtr
}

func (p *Map[k, v]) Get(key k) (valuePtr *v, ok bool) {
	i := p.Hash(key) & uint32(p.shards-1)
	p.m[i].Lock.RLock()
	defer p.m[i].Lock.RUnlock()
	valuePtr, ok = p.m[i].InternalMap[key]
	return
}

func (p *Map[k, v]) GetNew(key k) (valuePtr *v, ok bool) {
	valuePtr = p.New(key)
	i := p.Hash(key) & uint32(p.shards-1)
	p.m[i].Lock.Lock()
	defer p.m[i].Lock.Unlock()
	p.m[i].InternalMap[key] = valuePtr
	return
}

func (p *Map[k, v]) GetCopy(key k) (value v, ok bool) {
	i := p.Hash(key) & uint32(p.shards-1)
	p.m[i].Lock.RLock()
	defer p.m[i].Lock.RUnlock()
	valuePtr, ok := p.m[i].InternalMap[key]
	if ok {
		return *valuePtr, ok
	} else {
		return value, ok
	}
}

// GetOrInit return the value itself(pointer)
// if already exists, return value and false
// if not, use init func and store the new value, return value and true
func (p *Map[k, v]) GetOrInit(key k, init func() *v) (actual *v, initialized bool) {
	i := p.Hash(key) & uint32(p.shards-1)
	p.m[i].Lock.RLock()
	valuePtr, ok := p.m[i].InternalMap[key]
	p.m[i].Lock.RUnlock()
	if ok {
		// load
		return valuePtr, true
	}

	p.m[i].Lock.Lock()
	actual, ok = p.m[i].InternalMap[key]
	if !ok {
		// init
		actual = init()
		p.m[i].InternalMap[key] = actual
	}
	p.m[i].Lock.Unlock()
	return
}

// GetCopyOrInit return the value copy
// if already exists, return value copy and false
// if not, use init func and store the new value, return value copy and true
func (p *Map[k, v]) GetCopyOrInit(key k, init func() *v) (actual v, initialized bool) {
	i := p.Hash(key) & uint32(p.shards-1)
	p.m[i].Lock.RLock()
	actualPtr, ok := p.m[i].InternalMap[key]
	p.m[i].Lock.RUnlock()
	if ok {
		// load
		return *actualPtr, true
	}

	p.m[i].Lock.Lock()
	actualPtr, ok = p.m[i].InternalMap[key]
	if !ok {
		// init
		actualPtr = init()
		p.m[i].InternalMap[key] = actualPtr
	}
	p.m[i].Lock.Unlock()
	return *actualPtr, false
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (p *Map[k, v]) LoadOrStore(key k, value v) (actual v, loaded bool) {
	i := p.Hash(key) & uint32(p.shards-1)
	p.m[i].Lock.RLock()
	valuePtr, ok := p.m[i].InternalMap[key]
	p.m[i].Lock.RUnlock()
	if ok {
		// load
		return *valuePtr, true
	}

	p.m[i].Lock.Lock()
	valuePtr, ok = p.m[i].InternalMap[key]
	if ok {
		// load
		p.m[i].Lock.Unlock()
		return *valuePtr, true
	} else {
		// store
		p.m[i].InternalMap[key] = &value
		p.m[i].Lock.Unlock()
		return value, false
	}
}

func NewCacheMap[k comparable, v any](len int) *Map[k, v] {
	m := &Map[k, v]{
		m:      make([]*Shard[k, v], len),
		shards: len,
		New:    func(k) *v { return new(v) },
		Hash:   DefaultHash[k],
	}

	wg := sync.WaitGroup{}
	wg.Add(len)
	for i := 0; i < len; i++ {
		go func(i int) {
			m.m[i] = &Shard[k, v]{InternalMap: make(map[k]*v)}
			wg.Done()
		}(i)
	}
	wg.Wait()
	return m
}
