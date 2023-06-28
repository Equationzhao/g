package content

import (
	"sync"

	"github.com/Equationzhao/g/filter"
	"github.com/Equationzhao/g/item"
	"github.com/Equationzhao/g/osbased"
	"github.com/Equationzhao/g/render"
)

type InodeEnabler struct {
	*sync.WaitGroup
}

func NewInodeEnabler() *InodeEnabler {
	return &InodeEnabler{
		WaitGroup: new(sync.WaitGroup),
	}
}

const Inode = "Inode"

func (i *InodeEnabler) Enable(renderer *render.Renderer) filter.ContentOption {
	wait := func(res string) string {
		i.Wait()
		return renderer.Inode(res)
	}

	return func(info *item.FileInfo) (string, string) {
		str := osbased.Inode(info)
		return wait(str), Inode
	}
}
