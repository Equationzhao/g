package item

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/Equationzhao/g/slices"
	"github.com/Equationzhao/tsmap"
	"github.com/valyala/bytebufferpool"
)

type FileInfo struct {
	os.FileInfo
	FullPath string
	Meta     *tsmap.Map[string, Item]
}

type Option = func(info *FileInfo) error

func WithSize(size int) Option {
	return func(info *FileInfo) error {
		info.Meta = tsmap.NewTSMap[string, Item](size)
		return nil
	}
}

func WithFileInfo(info os.FileInfo) Option {
	return func(f *FileInfo) error {
		f.FileInfo = info
		return nil
	}
}

// WithPath will get abs path of given string
// and set the full path of FileInfo
func WithPath(path string) Option {
	return func(f *FileInfo) error {
		abs, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		f.FullPath = abs
		return nil
	}
}

func NewFileInfoWithOption(opts ...Option) (*FileInfo, error) {
	f := &FileInfo{}
	var errSum error
	for _, opt := range opts {
		err := opt(f)
		if err != nil {
			errSum = errors.Join(errSum, err)
		}
	}
	if f.Meta == nil {
		f.Meta = tsmap.NewTSMap[string, Item](200)
	}
	return f, errSum
}

func NewFileInfo(name string) (*FileInfo, error) {
	info, err := os.Stat(name)
	if err != nil {
		return nil, err
	}

	abs, err := filepath.Abs(name)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		FileInfo: info,
		FullPath: abs,
		Meta:     tsmap.NewTSMap[string, Item](20),
	}, nil
}

// Keys return all keys in random order
func (i *FileInfo) Keys() []string {
	items := i.Meta.Values()
	res := make([]string, 0, len(items))
	for _, item := range items {
		res = append(res, item.String())
	}
	return res
}

// KeysByOrder return Keys(ordered by No, ascending)
func (i *FileInfo) KeysByOrder() []string {
	kNo := i.Meta.Pairs()

	slices.SortFunc(
		kNo, func(i, j tsmap.Pair[string, Item]) int {
			return i.Value().NO() - j.Value().NO()
		},
	)

	res := make([]string, 0, len(kNo))
	for _, v := range kNo {
		res = append(res, v.Key())
	}
	return res
}

// Del delete content by key
func (i *FileInfo) Del(key string) {
	i.Meta.Remove(key)
}

// Get content by key
func (i *FileInfo) Get(key string) (Item, bool) {
	return i.Meta.Get(key)
}

// Set content by key
func (i *FileInfo) Set(key string, ic Item) {
	i.Meta.Set(key, ic)
}

func (i *FileInfo) Values() []Item {
	return i.Meta.Values()
}

// ValuesByOrdered return all content (ordered by No, ascending)
func (i *FileInfo) ValuesByOrdered() []Item {
	ics := i.Meta.Values()
	slices.SortFunc(
		ics, func(i, j Item) int {
			return i.NO() - j.NO()
		},
	)

	return ics
}

func (i *FileInfo) OrderedContent(delimiter string) string {
	res := bytebufferpool.Get()
	defer bytebufferpool.Put(res)
	items := i.ValuesByOrdered()
	for j, item := range items {
		_, _ = res.WriteString(item.String())
		if j != len(items)-1 {
			_, _ = res.WriteString(delimiter)
		}
	}
	return res.String()
}
