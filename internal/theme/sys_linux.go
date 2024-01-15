//go:build linux

package theme

import "github.com/Equationzhao/g/internal/const"

func init() {
	DefaultAll.Name["sys"] = Style{
		Icon:  "\ue712",
		Color: constval.BrightBlue,
	}
}
