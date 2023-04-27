//go:build linux

package filter

import (
	"github.com/Equationzhao/g/cached"
	"github.com/Equationzhao/g/render"
	"os"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func EnableOwner(renderer *render.Renderer) ContentOption {
	m := sync.RWMutex{}
	longestOwner := 0
	cache := cached.NewUsernameMap()
	return func(info os.FileInfo) string {
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

		time.Sleep(time.Microsecond * 5)

		return renderer.Owner(fillBlank(name, longestOwner))
	}
}

func EnableGroup(renderer *render.Renderer) ContentOption {
	m := sync.RWMutex{}
	longestGroup := 0
	cache := cached.NewGroupnameMap()
	return func(info os.FileInfo) string {
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

		time.Sleep(time.Microsecond * 5)

		return renderer.Owner(fillBlank(name, longestGroup))
	}
}
