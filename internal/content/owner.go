package content

import (
	"github.com/Equationzhao/g/internal/align"
	"runtime"

	constval "github.com/Equationzhao/g/internal/const"
	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/osbased"
	"github.com/Equationzhao/g/internal/render"
)

const (
	OwnerName    = constval.NameOfOwner
	OwnerUidName = constval.NameOfOwnerUid
	OwnerSID     = constval.NameOfOwnerSID
)

type OwnerEnabler struct {
	Numeric bool
}

func NewOwnerEnabler() *OwnerEnabler {
	return &OwnerEnabler{}
}

func (o *OwnerEnabler) EnableNumeric() {
	o.Numeric = true
}

func (o *OwnerEnabler) DisableNumeric() {
	o.Numeric = false
}

func (o *OwnerEnabler) EnableOwner(renderer *render.Renderer) ContentOption {
	align.RegisterHeaderFooter(OwnerName)
	return func(info *item.FileInfo) (string, string) {
		name, returnFuncName := "", ""
		if o.Numeric {
			name = osbased.OwnerID(info)
			if runtime.GOOS == "windows" {
				returnFuncName = OwnerSID
			} else {
				returnFuncName = OwnerUidName
			}
		} else {
			name = osbased.Owner(info)
			returnFuncName = OwnerName
		}
		return renderer.Owner(name), returnFuncName
	}
}
