package content

import (
	"runtime"

	"github.com/Equationzhao/g/internal/align"

	constval "github.com/Equationzhao/g/internal/global"
	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/osbased"
	"github.com/Equationzhao/g/internal/render"
)

type GroupEnabler struct {
	Numeric bool
	Smart   bool
}

const (
	GroupName    = constval.NameOfGroupName
	GroupUidName = constval.NameOfGroupUidName
	GroupSID     = constval.NameOfGroupSID
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

func (g *GroupEnabler) EnableSmartMode() {
	g.Smart = true
}

func (g *GroupEnabler) DisableSmartMode() {
	g.Smart = false
}

func (g *GroupEnabler) EnableGroup(renderer *render.Renderer) ContentOption {
	align.RegisterHeaderFooter(GroupName)
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
			if g.Smart && name == osbased.Owner(info) {
				name = ""
			}
		}
		return renderer.Group(name), returnFuncName
	}
}
