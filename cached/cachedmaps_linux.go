//go:build linux

package cached

import (
	"os/user"

	"github.com/Equationzhao/tsmap"
)

type (
	Uid      = string
	Username = string
)

// usernameMap is a map from Uid to Username
// current not contained because it is cached in user.Current()
type usernameMap struct {
	m *tsmap.Map[Uid, Username]
}

func NewUsernameMap() *usernameMap {
	return &usernameMap{
		m: tsmap.NewTSMap[Uid, Username](20),
	}
}

func (m *usernameMap) Get(u Uid) Username {
	if c, _ := user.Current(); c.Uid == u {
		return c.Username
	}

	v, _ := m.m.GetOrInit(u, func() Groupname {
		targetUser, err := user.LookupId(u)
		if err != nil {
			targetUser.Username = "uid:" + u
		}
		return targetUser.Username
	})
	return v
}

type (
	Gid       = string
	Groupname = string
)

// groupnameMap is a map from Gid to Groupname
type groupnameMap struct {
	m *tsmap.Map[Gid, Groupname]
}

func NewGroupnameMap() *groupnameMap {
	return &groupnameMap{
		m: tsmap.NewTSMap[Gid, Groupname](20),
	}
}

func (m *groupnameMap) Get(g Gid) Groupname {
	v, _ := m.m.GetOrInit(g, func() Groupname {
		targetUser, err := user.LookupGroupId(g)
		if err != nil {
			targetUser.Name = "gid:" + g
		}
		return targetUser.Name
	})
	return v
}

var (
	mainGroupnameMap = NewGroupnameMap()
	mainUsernameMap  = NewUsernameMap()
)

func GetGroupname(g Gid) Groupname {
	return mainGroupnameMap.Get(g)
}

func GetUsername(u Uid) Groupname {
	return mainUsernameMap.Get(u)
}
