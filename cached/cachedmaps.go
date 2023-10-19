//go:build linux || darwin || openbsd

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
			if targetUser == nil {
				targetUser = new(user.User)
			}
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

// GroupnameMap is a map from Gid to Groupname
type GroupnameMap struct {
	m *tsmap.Map[Gid, Groupname]
}

func NewGroupnameMap() *GroupnameMap {
	return &GroupnameMap{
		m: tsmap.NewTSMap[Gid, Groupname](20),
	}
}

func (m *GroupnameMap) Get(g Gid) Groupname {
	v, _ := m.m.GetOrInit(g, func() Groupname {
		targetGroup, err := user.LookupGroupId(g)
		if err != nil {
			if targetGroup == nil {
				targetGroup = new(user.Group)
			}
			targetGroup.Name = "gid:" + g
		}
		return targetGroup.Name
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
