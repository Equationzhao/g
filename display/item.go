package display

import (
	"sort"

	"github.com/Equationzhao/g/util"
	"github.com/valyala/bytebufferpool"
)

type Content interface {
	String() string
}

type StringContent string

func (s StringContent) String() string {
	return string(s)
}

type ItemContent struct {
	No      int
	Content Content
}

type Item struct {
	Delimiter string
	internal  map[string]ItemContent
}

func (i *Item) GetAll() map[string]ItemContent {
	return i.internal
}

// Keys return all keys in random order
func (i *Item) Keys() []string {
	res := make([]string, 0, len(i.internal))
	for k := range i.internal {
		res = append(res, k)
	}
	return res
}

// KeysByOrder return Keys(ordered by itemContent.No, ascending)
func (i *Item) KeysByOrder() []string {
	res := make([]string, 0, len(i.internal))
	kNo := make([]struct {
		k  string
		no int
	}, 0, len(i.internal))
	for k, v := range i.internal {
		kNo = append(kNo, struct {
			k  string
			no int
		}{
			k:  k,
			no: v.No,
		})
	}

	sort.Slice(kNo, func(i, j int) bool {
		return kNo[i].no < kNo[j].no
	})

	for _, v := range kNo {
		res = append(res, v.k)
	}
	return res
}

func (i *Item) Del(key string) {
	delete(i.internal, key)
}

// Get content by key
func (i *Item) Get(key string) (ItemContent, bool) {
	c, ok := i.internal[key]
	return c, ok
}

// Set content by key
func (i *Item) Set(key string, ic ItemContent) {
	i.internal[key] = ic
}

// ExcludeOrderedContent get content in order, exclude those match given parameter(ordered by itemContent.No, ascending)
func (i *Item) ExcludeOrderedContent(key ...string) string {
	res := bytebufferpool.Get()
	ics := make([]ItemContent, 0, len(i.internal))
	for name, v := range i.internal {
		if util.SliceContains(key, name) {
			continue
		}
		ics = append(ics, v)
	}
	sort.Slice(ics, func(i, j int) bool {
		return ics[i].No < ics[j].No
	})

	for _, v := range ics {
		_, _ = res.WriteString(v.Content.String())
		_, _ = res.WriteString(i.Delimiter)
	}

	defer bytebufferpool.Put(res)
	return res.String()
}

// IncludeOrderedContent return those content inorder(ordered by itemContent.No, ascending)
func (i *Item) IncludeOrderedContent(names ...string) string {
	res := bytebufferpool.Get()
	ics := make([]ItemContent, 0, len(i.internal))
	for name, v := range i.internal {
		if !util.SliceContains(names, name) {
			continue
		}
		ics = append(ics, v)
	}
	sort.Slice(ics, func(i, j int) bool {
		return ics[i].No < ics[j].No
	})

	for _, v := range ics {
		_, _ = res.WriteString(v.Content.String())
		_, _ = res.WriteString(i.Delimiter)
	}

	defer bytebufferpool.Put(res)
	return res.String()
}

// OrderedContent return all content in order(ordered by itemContent.No, ascending)
func (i *Item) OrderedContent() string {
	res := bytebufferpool.Get()
	ics := make([]ItemContent, 0, len(i.internal))
	for _, v := range i.internal {
		ics = append(ics, v)
	}
	sort.Slice(ics, func(i, j int) bool {
		return ics[i].No < ics[j].No
	})

	for j, v := range ics {
		_, _ = res.WriteString(v.Content.String())
		if j != len(ics)-1 {
			_, _ = res.WriteString(i.Delimiter)
		}
	}

	defer bytebufferpool.Put(res)
	return res.String()
}

type ItemOptions func(*Item)

func NewItem(Ops ...ItemOptions) *Item {
	res := &Item{internal: make(map[string]ItemContent)}
	for _, op := range Ops {
		op(res)
	}
	return res
}

func WithDelimiter(d string) ItemOptions {
	return func(item *Item) {
		item.Delimiter = d
	}
}
