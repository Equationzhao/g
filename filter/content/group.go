package content

import (
	"runtime"

	"github.com/Equationzhao/g/filter"
	"github.com/Equationzhao/g/item"
	"github.com/Equationzhao/g/osbased"
	"github.com/Equationzhao/g/render"
)

type GroupEnabler struct {
	Numeric bool
}

const (
	GroupName    = "Group"
	GroupUidName = "Group-uid"
	GroupSID     = "Group-sid"
)

func NewGroupEnabler() *GroupEnabler {
	return &GroupEnabler{}
}

func (g *GroupEnabler) EnableNumeric() {
	g.Numeric = true
}

func (g *GroupEnabler) DisableNumeric() {
	g.Numeric = false
}

func (g *GroupEnabler) EnableGroup(renderer *render.Renderer) filter.ContentOption {
	return func(info *item.FileInfo) (string, string) {
		name, returnFuncName := "", GroupName
		if g.Numeric {
			name = osbased.GroupID(info)
			if runtime.GOOS == "windows" {
				returnFuncName = GroupSID
			} else {
				returnFuncName = GroupUidName
			}
		} else {
			name = osbased.Group(info)
		}
		return renderer.Group(name), returnFuncName
	}
}
