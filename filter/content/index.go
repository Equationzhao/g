package content

import (
	"github.com/Equationzhao/g/filter"
	"github.com/Equationzhao/g/item"
)

type IndexEnabler struct{}

func NewIndexEnabler() *IndexEnabler {
	return &IndexEnabler{}
}

func (i *IndexEnabler) Enable() filter.ContentOption {
	return func(info *item.FileInfo) (string, string) {
		return "", "#"
	}
}
