package tree

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/Equationzhao/g/filter"
	"github.com/Equationzhao/g/filter/content"
	"github.com/Equationzhao/g/item"
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

func NewTreeString(entry string, depthLimit int, typeFilter *filter.ItemFilter, contentFilter *filter.ContentFilter) (t *Tree, serious error, minor error) {
	stat, err := os.Stat(entry)
	if err != nil {
		return nil, err, nil
	}
	cm := sync.Mutex{}
	cm.Lock()
	abs, err := filepath.Abs(stat.Name())
	if err != nil {
		return nil, err, nil
	}
	ExtraName := contentFilter.GetDisplayItems(&item.FileInfo{
		FileInfo: stat,
		FullPath: abs,
		Meta:     make(map[string]item.Item),
	})[0]
	extra := ExtraName.ExcludeOrderedContent(content.NameName)
	name, _ := ExtraName.Get(content.NameName)
	cm.Unlock()
	n := &Tree{
		tree: NewWithExtraInfoRoot(extra, "", name.Content.String()),
	}

	var wg sync.WaitGroup

	if stat.IsDir() {
		n.stat.directory.Add(1)
	} else {
		n.stat.file.Add(1)
	}
	errChan := make(chan error, 10)
	var errSum error
	expand(n.tree, depthLimit, &wg, entry, &n.stat, typeFilter, contentFilter, &cm, errChan)
	errWg := sync.WaitGroup{}
	errWg.Add(1)
	go func() {
		for err := range errChan {
			if err != nil {
				errSum = errors.Join(errSum, err)
			}
		}
		errWg.Done()
	}()
	wg.Wait()
	close(errChan)
	errWg.Wait()
	return n, nil, errSum
}

func expand(node tree, depthLimit int, wg *sync.WaitGroup, parent string, s *statistic, typeFilter *filter.ItemFilter, contentFilter *filter.ContentFilter, cm *sync.Mutex, errSender chan<- error) {
	if depthLimit == 0 {
		return
	}

	d, err := os.ReadDir(parent)
	if err != nil {
		errSender <- err
	}

	infos := make([]*item.FileInfo, 0, len(d))
	for _, entry := range d {
		info, err := entry.Info()
		if err != nil {
			errSender <- err
			continue
		}
		abs, err := filepath.Abs(info.Name())
		if err != nil {
			errSender <- err
			continue
		}
		infos = append(infos, &item.FileInfo{
			FileInfo: info,
			FullPath: abs,
			Meta:     make(map[string]item.Item),
		})
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
				en := contentFilter.GetDisplayItems(v)[0]
				cm.Unlock()
				wrappedName, _ := en.Get(content.NameName)
				name = wrappedName.Content.String()
				extra = en.ExcludeOrderedContent(content.NameName)
			} else {
				name = v.Name()
			}
			newBranch := node.AddInfoBranch(extra, "", name)
			go func() {
				expand(newBranch, depthLimit-1, wg, filepath.Join(parent, v.Name()), s, typeFilter, contentFilter, cm, errSender)
				wg.Done()
			}()
		} else {
			s.file.Add(1)
			if contentFilter != nil {
				cm.Lock()
				en := contentFilter.GetDisplayItems(v)[0]
				cm.Unlock()
				wrappedName, _ := en.Get(content.NameName)
				extra := en.ExcludeOrderedContent(content.NameName)
				node.AddInfoNode(extra, "", wrappedName.Content.String())
			} else {
				node.AddNode(v.Name())
			}
		}
	}
}
