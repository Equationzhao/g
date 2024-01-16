package content

import (
	"runtime"

	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/osbased"
	"github.com/Equationzhao/g/internal/render"
)

const (
	OwnerName    = "Owner"
	OwnerUidName = "Owner-uid"
	OwnerSID     = "Owner-sid"
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
