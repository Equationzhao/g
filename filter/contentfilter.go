package filter

import (
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/Equationzhao/g/display"
	"github.com/Equationzhao/g/item"
)

type ContentFilter struct {
	noOutputOptions []NoOutputOption
	options         []ContentOption
	sortFunc        func(a, b os.FileInfo) bool
	LimitN          uint // <=0 means no limit
}

func (cf *ContentFilter) SortFunc() func(a, b os.FileInfo) bool {
	return cf.sortFunc
}

func (cf *ContentFilter) SetSortFunc(sortFunc func(a, b os.FileInfo) bool) {
	cf.sortFunc = sortFunc
}

func (cf *ContentFilter) AppendToOptions(options ...ContentOption) {
	cf.options = append(cf.options, options...)
}

func (cf *ContentFilter) SetOptions(options ...ContentOption) {
	cf.options = options
}

func (cf *ContentFilter) AppendToNoOutputOptions(options ...NoOutputOption) {
	cf.noOutputOptions = append(cf.noOutputOptions, options...)
}

func (cf *ContentFilter) SetNoOutputOptions(outputFunc ...NoOutputOption) {
	cf.noOutputOptions = outputFunc
}

type (
	ContentOption  func(info *item.FileInfo) (stringContent string, funcName string)
	NoOutputOption func(info *item.FileInfo)
)

// FillBlank
// if s is shorter than length, fill blank from left
// if s is longer than length, panic
func FillBlank(s string, length int) string {
	if len(s) > length {
		return s
	}
	return strings.Repeat(" ", length-len(s)) + s
}

type ContentFilterOption func(cf *ContentFilter)

func WithOptions(options ...ContentOption) ContentFilterOption {
	return func(cf *ContentFilter) {
		cf.options = options
	}
}

func WithNoOutputOptions(options ...NoOutputOption) ContentFilterOption {
	return func(cf *ContentFilter) {
		cf.noOutputOptions = options
	}
}

func NewContentFilter(options ...ContentFilterOption) *ContentFilter {
	c := &ContentFilter{
		sortFunc: nil,
	}

	for _, option := range options {
		option(c)
	}

	if c.options == nil {
		c.options = make([]ContentOption, 0)
	}
	if c.noOutputOptions == nil {
		c.noOutputOptions = make([]NoOutputOption, 0)
	}

	return c
}

func (cf *ContentFilter) GetDisplayItems(e ...*item.FileInfo) {
	sort.Slice(e, func(i, j int) bool {
		if cf.sortFunc != nil {
			return cf.sortFunc(e[i], e[j])
		}
		return true
	})

	// limit number of entries
	// 0 means no limit
	if cf.LimitN != 0 && len(e) > int(cf.LimitN) {
		e = e[:cf.LimitN]
	}

	wg := sync.WaitGroup{}
	wg.Add(len(e))

	for i, entry := range e {
		go func(entry *item.FileInfo, i int) {
			for j, option := range cf.options {
				stringContent, funcName := option(entry)
				content := display.ItemContent{Content: display.StringContent(stringContent), No: j}
				entry.Set(funcName, content)
			}

			for _, option := range cf.noOutputOptions {
				option(entry)
			}
			wg.Done()
		}(entry, i)
	}
	wg.Wait()
}
