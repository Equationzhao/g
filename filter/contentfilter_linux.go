//go:build linux

package filter

import (
	"github.com/Equationzhao/g/render"
	"os"
	"os/user"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func EnableOwner(renderer *render.Renderer) ContentOption {
	m := sync.RWMutex{}
	longestOwner := 0
	cache := sync.Map{}
	return func(info os.FileInfo) string {
		uid := strconv.FormatInt(int64(info.Sys().(*syscall.Stat_t).Uid), 10)
		nameAny, _ := cache.Load(uid)
		var name string
		if nameAny == nil {
			u, err := user.LookupId(uid)
			if err != nil {
				name = "uid:" + uid
			} else {
				name = u.Username
			}
			cache.Store(uid, name)
		} else {
			name = nameAny.(string)
		}

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
	cache := sync.Map{}
	return func(info os.FileInfo) string {
		gid := strconv.FormatInt(int64(info.Sys().(*syscall.Stat_t).Gid), 10)
		nameAny, _ := cache.Load(gid)
		var name string
		if nameAny == nil {
			g, err := user.LookupGroupId(gid)
			if err != nil {
				name = "gid:" + gid
			} else {
				name = g.Name
			}
			cache.Store(gid, name)
		} else {
			name = nameAny.(string)
		}

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
