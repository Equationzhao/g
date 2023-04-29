package cached

import (
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

type CacheMap[k comparable, v any] struct {
	m     map[k]*v
	mutex *sync.RWMutex
	New   func(k) *v
}

func (p *CacheMap[k, v]) Set(key k, value v) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.m[key] = &value
}

func (p *CacheMap[k, v]) Get(key k) (valuePtr *v, ok bool) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	valuePtr, ok = p.m[key]
	return
}

func (p *CacheMap[k, v]) GetNew(key k) (valuePtr *v, ok bool) {
	valuePtr = p.New(key)
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.m[key] = valuePtr
	return
}

func (p *CacheMap[k, v]) GetCopy(key k) (value v, ok bool) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	valuePtr, ok := p.m[key]
	return *valuePtr, ok
}

// GetOrInit return the value itself(pointer)
// if already exists, return value and true
// if not, use init func and store the new value, return value and false
func (p *CacheMap[k, v]) GetOrInit(key k, init func() *v) (valuePtr *v, notNew bool) {
	p.mutex.RLock()
	valuePtr, notNew = p.m[key]
	p.mutex.RUnlock()
	if !notNew {
		p.mutex.Lock()
		valuePtr, notNew = p.m[key]
		if !notNew {
			valuePtr = init()
			p.m[key] = valuePtr
		}
		p.mutex.Unlock()
	}
	return
}

// GetCopyOrInit return the value copy
// if already exists, return value copy and true
// if not, use init func and store the new value, return value copy and false
func (p *CacheMap[k, v]) GetCopyOrInit(key k, init func() *v) (value v, notNew bool) {
	p.mutex.RLock()
	valuePtr, notNew := p.m[key]
	p.mutex.RUnlock()
	if !notNew {
		p.mutex.Lock()
		valuePtr, notNew = p.m[key]
		if !notNew {
			p.m[key] = init()
		}
		p.mutex.Unlock()
	}
	value = *valuePtr
	return
}

func NewCacheMap[k comparable, v any]() *CacheMap[k, v] {
	return &CacheMap[k, v]{
		m:     make(map[k]*v),
		mutex: new(sync.RWMutex),
	}
}
