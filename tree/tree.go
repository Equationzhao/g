package tree

import (
	"os"
	"path/filepath"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/Equationzhao/g/filter"
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

func NewTreeString(entry string, depthLimit int, typeFilter *filter.TypeFilter, contentFilter *filter.ContentFilter) (*Tree, error) {
	stat, err := os.Stat(entry)
	if err != nil {
		return nil, err
	}
	cm := sync.Mutex{}
	cm.Lock()
	ExtraName := contentFilter.GetExtraAndNameStringSlice(stat)
	cm.Unlock()
	n := &Tree{tree: NewWithExtraInfoRoot(ExtraName[0].Key(), "", ExtraName[0].Value())}

	var wg sync.WaitGroup

	if stat.IsDir() {
		n.stat.directory.Add(1)
	} else {
		n.stat.file.Add(1)
	}

	expand(n.tree, depthLimit, &wg, entry, &n.stat, typeFilter, contentFilter, &cm)
	wg.Wait()
	return n, nil
}

func expand(node tree, depthLimit int, wg *sync.WaitGroup, parent string, s *statistic, typeFilter *filter.TypeFilter, contentFilter *filter.ContentFilter, cm *sync.Mutex) {
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
		infos = typeFilter.Filter(infos...)
	}

	sort.Slice(infos, func(i, j int) bool {
		if contentFilter.SortFunc() != nil {
			return contentFilter.SortFunc()(infos[i], infos[j])
		} else {
			return true
		}
	})

	for _, v := range infos {
		if v.IsDir() {
			v := v
			wg.Add(1)

			s.directory.Add(1)
			var name, extra string
			if contentFilter != nil {
				cm.Lock()
				en := contentFilter.GetExtraAndNameStringSlice(v)[0]
				cm.Unlock()
				name = en.Value()
				extra = en.Key()
			} else {
				name = v.Name()
			}
			newBranch := node.AddInfoBranch(extra, "", name)
			go func() {
				expand(newBranch, depthLimit-1, wg, filepath.Join(parent, v.Name()), s, typeFilter, contentFilter, cm)
				wg.Done()
			}()
		} else {
			s.file.Add(1)
			if contentFilter != nil {
				cm.Lock()
				en := contentFilter.GetExtraAndNameStringSlice(v)[0]
				cm.Unlock()
				node.AddInfoNode(en.Key(), "", en.Value())
			} else {
				node.AddNode(v.Name())
			}
		}
	}
}
