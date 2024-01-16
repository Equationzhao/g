package content

import (
	"github.com/Equationzhao/g/internal/item"
)

type IndexEnabler struct{}

func NewIndexEnabler() *IndexEnabler {
	return &IndexEnabler{}
}

func (i *IndexEnabler) Enable() ContentOption {
	return func(info *item.FileInfo) (string, string) {
		return "", "#"
	}
}
