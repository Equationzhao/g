//go:build linux

package filter

import (
	"os"
	"strconv"
	"sync"
	"syscall"

	"github.com/Equationzhao/g/cached"
	"github.com/Equationzhao/g/render"
)

func (cf *ContentFilter) EnableOwner(renderer *render.Renderer) ContentOption {
	m := sync.RWMutex{}
	longestOwner := 0
	cache := cached.NewUsernameMap()
	return func(info os.FileInfo) string {
		wait := func(res string) string {
			cf.wgOwner.Wait()
			return renderer.Owner(fillBlank(res, longestOwner))
		}

		done := func(name string) {
			defer cf.wgOwner.Done()
			m.RLock()
			if len(name) > longestOwner {
				m.RUnlock()
				m.Lock()
				if len(name) > longestOwner {
					longestOwner = len(name)
				}
				m.Unlock()
			} else {
				m.RUnlock()
			}
		}
		uid := strconv.FormatInt(int64(info.Sys().(*syscall.Stat_t).Uid), 10)
		name := cache.Get(uid)

		m.RLock()
		if len(name) > longestOwner {
			m.RUnlock()
			m.Lock()
			if len(name) > longestOwner {
				longestOwner = len(name)
			}
			m.Unlock()
		} else {
			m.RUnlock()
		}
		done(name)
		return wait(name)
	}
}

func (cf *ContentFilter) EnableGroup(renderer *render.Renderer) ContentOption {
	m := sync.RWMutex{}
	longestGroup := 0
	cache := cached.NewGroupnameMap()
	return func(info os.FileInfo) string {
		wait := func(name string) string {
			cf.wgGroup.Wait()
			return renderer.Group(fillBlank(name, longestGroup))
		}

		done := func(name string) {
			defer cf.wgGroup.Done()
			m.RLock()
			if len(name) > longestGroup {
				m.RUnlock()
				m.Lock()
				if len(name) > longestGroup {
					longestGroup = len(name)
				}
				m.Unlock()
			} else {
				m.RUnlock()
			}
		}
		gid := strconv.FormatInt(int64(info.Sys().(*syscall.Stat_t).Gid), 10)
		name := cache.Get(gid)

		m.RLock()
		if len(name) > longestGroup {
			m.RUnlock()
			m.Lock()
			if len(name) > longestGroup {
				longestGroup = len(name)
			}
			m.Unlock()
		} else {
			m.RUnlock()
		}

		done(name)
		return wait(name)
	}
}
