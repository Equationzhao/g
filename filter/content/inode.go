package content

import (
	"os"
	"sync"

	"github.com/Equationzhao/g/filter"
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
	m := sync.RWMutex{}
	longestInode := 0

	wait := func(res string) string {
		i.Wait()
		return renderer.Inode(filter.FillBlank(res, longestInode))
	}

	done := func(name string) {
		defer i.Done()
		m.RLock()
		if len(name) > longestInode {
			m.RUnlock()
			m.Lock()
			if len(name) > longestInode {
				longestInode = len(name)
			}
			m.Unlock()
		} else {
			m.RUnlock()
		}
	}

	return func(info os.FileInfo) (string, string) {
		str := osbased.Inode(info)
		done(str)
		return wait(str), Inode
	}
}
