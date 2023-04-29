package tree

import (
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/Equationzhao/g/filter"
	"github.com/Equationzhao/g/render"
)

type statistic struct {
	directory atomic.Uint32
	file      atomic.Uint32
}

type Tree struct {
	stat statistic
	tree tree
}

func (n *Tree) File() uint32 {
	return n.stat.file.Load()
}

func (n *Tree) Directory() uint32 {
	return n.stat.directory.Load()
}

func (n *Tree) MakeTreeStr() string {
	return n.tree.String()
}

func NewTreeString(entry string, depthLimit int, typeFilter *filter.TypeFilter, renderer *render.Renderer) (*Tree, error) {
	n := &Tree{tree: NewWithRoot(entry)}

	var wg sync.WaitGroup
	stat, err := os.Stat(entry)
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		n.stat.directory.Add(1)
	} else {
		n.stat.file.Add(1)
	}

	expand(n.tree, depthLimit, &wg, entry, &n.stat, typeFilter, renderer)
	wg.Wait()
	return n, nil
}

func expand(node tree, depthLimit int, wg *sync.WaitGroup, parent string, s *statistic, typeFilter *filter.TypeFilter, renderer *render.Renderer) {
	if depthLimit == 0 {
		return
	}

	d, err := os.ReadDir(parent)
	if err != nil {
		node.AddNode(err.Error())
	}

	infos := make([]os.FileInfo, 0, len(d))
	for _, entry := range d {
		info, err := entry.Info()
		if err != nil {
			node.AddNode(err.Error())
			return
		}
		infos = append(infos, info)
	}

	if typeFilter != nil {
		infos = typeFilter.Filter(infos)
	}

	for _, v := range d {
		if v.IsDir() {
			info, err := v.Info()
			if err != nil {
				node.AddNode(err.Error())
				return
			}
			v := v
			wg.Add(1)
			go func() {
				s.directory.Add(1)
				var name string
				if renderer != nil {
					name = renderer.DirIcon(v.Name())
				} else {
					name = v.Name()
				}
				expand(node.AddBranch(name), depthLimit-1, wg, filepath.Join(parent, info.Name()), s, typeFilter, renderer)
				wg.Done()
			}()
		} else if v.Type()&os.ModeSymlink != 0 {
			info, err := v.Info()
			if err != nil {
				node.AddNode(err.Error())
				return
			}
			s.file.Add(1)
			if renderer != nil {
				wg.Add(1)
				go func() {
					node.AddNode(renderer.Symlink(info.Name(), parent))
					wg.Done()
				}()
			} else {
				node.AddNode(v.Name())
			}
		} else {
			s.file.Add(1)
			if renderer != nil {
				node.AddNode(renderer.ByExtIcon(v.Name()))
			} else {
				node.AddNode(v.Name())
			}
		}
	}
}
