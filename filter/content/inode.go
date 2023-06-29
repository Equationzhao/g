package content

import (
	"github.com/Equationzhao/g/filter"
	"github.com/Equationzhao/g/item"
	"github.com/Equationzhao/g/osbased"
	"github.com/Equationzhao/g/render"
)

type InodeEnabler struct{}

func NewInodeEnabler() *InodeEnabler {
	return &InodeEnabler{}
}

const Inode = "Inode"

func (i *InodeEnabler) Enable(renderer *render.Renderer) filter.ContentOption {
	return func(info *item.FileInfo) (string, string) {
		return renderer.Inode(osbased.Inode(info)), Inode
	}
}
