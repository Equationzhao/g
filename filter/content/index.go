package content

import (
	"os"

	"github.com/Equationzhao/g/filter"
)

type IndexEnabler struct{}

func NewIndexEnabler() *IndexEnabler {
	return &IndexEnabler{}
}

func (i *IndexEnabler) Enable() filter.ContentOption {
	return func(info os.FileInfo) (string, string) {
		return "", "#"
	}
}
