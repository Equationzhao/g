//go:build linux

package cached

import (
	"os/user"
	"sync"
)

type Uid = string
type Username = string

// usernameMap is a map from Uid to Username
// current not contained because it is cached in user.Current()
type usernameMap struct {
	m     map[Uid]Username
	mutex sync.RWMutex
}

func NewUsernameMap() *usernameMap {
	return &usernameMap{
		m: make(map[Uid]Username),
	}
}

func (m *usernameMap) Get(u Uid) Username {

	if c, _ := user.Current(); c.Uid == u {
		return c.Username
	}

	m.mutex.RLock()
	if username, ok := m.m[u]; ok {
		m.mutex.RUnlock()
		return username
	} else {
		m.mutex.RUnlock()
		m.mutex.Lock()
		if username, ok := m.m[u]; ok {
			m.mutex.Unlock()
			return username
		} else {
			targetUser, err := user.LookupId(u)
			if err != nil {
				targetUser.Username = "uid:" + u
			}
			m.m[u] = targetUser.Username
			m.mutex.Unlock()
			return targetUser.Username
		}
	}

}

type Gid = string
type Groupname = string

// groupnameMap is a map from Gid to Groupname
type groupnameMap struct {
	m     map[Gid]Groupname
	mutex sync.RWMutex
}

func NewGroupnameMap() *groupnameMap {
	return &groupnameMap{
		m: make(map[Gid]Groupname),
	}
}

func (m *groupnameMap) Get(g Gid) Groupname {

	m.mutex.RLock()
	if username, ok := m.m[g]; ok {
		m.mutex.RUnlock()
		return username
	} else {
		m.mutex.RUnlock()
		m.mutex.Lock()
		if username, ok = m.m[g]; ok {
			m.mutex.Unlock()
			return username
		} else {
			targetUser, err := user.LookupGroupId(g)
			if err != nil {
				targetUser.Name = "gid:" + g
			}
			m.m[g] = targetUser.Name
			m.mutex.Unlock()
			return targetUser.Name
		}
	}

}
