//go:build linux

package theme

import "github.com/Equationzhao/g/internal/global"

func init() {
	DefaultAll.Name["sys"] = Style{
		Icon:  "\ue712",
		Color: global.BrightBlue,
	}
}
