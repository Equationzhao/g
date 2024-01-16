package content

import (
	"strconv"

	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/osbased"
	"github.com/Equationzhao/g/internal/render"
)

// LinkEnabler List each file's number of hard links.
type LinkEnabler struct{}

func NewLinkEnabler() *LinkEnabler {
	return &LinkEnabler{}
}

const Link = "Link"

func (l *LinkEnabler) Enable(renderer *render.Renderer) ContentOption {
	return func(info *item.FileInfo) (string, string) {
		return renderer.Link(strconv.FormatUint(osbased.LinkCount(info), 10)), Link
	}
}
