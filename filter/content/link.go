package content

import (
	"strconv"

	"github.com/Equationzhao/g/filter"
	"github.com/Equationzhao/g/item"
	"github.com/Equationzhao/g/osbased"
	"github.com/Equationzhao/g/render"
)

// LinkEnabler List each file's number of hard links.
type LinkEnabler struct{}

func NewLinkEnabler() *LinkEnabler {
	return &LinkEnabler{}
}

const Link = "Link"

func (l *LinkEnabler) Enable(renderer *render.Renderer) filter.ContentOption {
	return func(info *item.FileInfo) (string, string) {
		return renderer.Link(strconv.FormatUint(osbased.LinkCount(info), 10)), Link
	}
}
