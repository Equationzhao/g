package content

import (
	"strings"

	"github.com/Equationzhao/g/internal/align"
	constval "github.com/Equationzhao/g/internal/global"
	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/osbased"
)

type FlagsEnabler struct{}

func NewFlagsEnabler() *FlagsEnabler {
	return &FlagsEnabler{}
}

const (
	Flags = constval.NameOfFlags
)

func (f FlagsEnabler) Enable() ContentOption {
	align.Register(Flags)
	return func(info *item.FileInfo) (string, string) {
		flags := osbased.CheckFlags(info)
		if len(flags) == 0 {
			return "-", Flags
		}
		return strings.Join(flags, ","), Flags
	}
}
