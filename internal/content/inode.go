package content

import (
	"github.com/Equationzhao/g/internal/align"
	constval "github.com/Equationzhao/g/internal/global"
	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/osbased"
	"github.com/Equationzhao/g/internal/render"
)

type InodeEnabler struct{}

func NewInodeEnabler() *InodeEnabler {
	return &InodeEnabler{}
}

const Inode = constval.NameOfInode

func (i *InodeEnabler) Enable(renderer *render.Renderer) ContentOption {
	align.RegisterHeaderFooter(Inode)
	return func(info *item.FileInfo) (string, string) {
		i := ""
		if m, ok := info.Cache[Inode]; ok {
			i = string(m)
		} else {
			i = osbased.Inode(info)
		}
		return renderer.Inode(i), Inode
	}
}
