package filter

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/Equationzhao/g/display"
	"github.com/Equationzhao/g/osbased"
	"github.com/Equationzhao/g/render"
)

type LengthFixed interface {
	Done()
	Wait()
	Add(delta int)
}

// fileMode size owner group time name

type ContentFilter struct {
	noOutputOptions []NoOutputOption
	options         []ContentOption
	wgs             []LengthFixed
	sortFunc        func(a, b os.FileInfo) bool
}

func (cf *ContentFilter) AppendToLengthFixed(fixed ...LengthFixed) {
	cf.wgs = append(cf.wgs, fixed...)
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
	ContentOption  func(info os.FileInfo) (stringContent string, funcName string)
	NoOutputOption func(info os.FileInfo)
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

var (
	Uid = false
	Gid = false
)

const (
	OwnerName    = "owner"
	OwnerUidName = "owner-uid"
	OwnerSID     = "owner-sid"
)

func (cf *ContentFilter) EnableOwner(renderer *render.Renderer) ContentOption {
	m := sync.RWMutex{}
	longestOwner := 0

	wg := new(sync.WaitGroup)
	cf.wgs = append(cf.wgs, wg)
	wait := func(res string) string {
		wg.Wait()
		return renderer.Owner(FillBlank(res, longestOwner))
	}

	done := func(name string) {
		defer wg.Done()
		m.RLock()
		if len(name) > longestOwner {
			m.RUnlock()
			m.Lock()
			if len(name) > longestOwner {
				longestOwner = len(name)
			}
			m.Unlock()
		} else {
			m.RUnlock()
		}
	}
	return func(info os.FileInfo) (string, string) {
		name := ""
		returnFuncName := ""
		if Uid {
			name = osbased.OwnerID(info)
			if runtime.GOOS == "windows" {
				returnFuncName = OwnerSID
			} else {
				returnFuncName = OwnerUidName
			}
		} else {
			name = osbased.Owner(info)
			returnFuncName = OwnerName
		}
		done(name)
		return wait(name), returnFuncName
	}
}

const (
	GroupName    = "group"
	GroupUidName = "group-uid"
	GroupSID     = "group-sid"
)

func (cf *ContentFilter) EnableGroup(renderer *render.Renderer) ContentOption {
	m := sync.RWMutex{}
	longestGroup := 0

	wg := new(sync.WaitGroup)
	cf.wgs = append(cf.wgs, wg)
	wait := func(name string) string {
		wg.Wait()
		return renderer.Group(FillBlank(name, longestGroup))
	}

	done := func(name string) {
		defer wg.Done()
		m.RLock()
		if len(name) > longestGroup {
			m.RUnlock()
			m.Lock()
			if len(name) > longestGroup {
				longestGroup = len(name)
			}
			m.Unlock()
		} else {
			m.RUnlock()
		}
	}

	return func(info os.FileInfo) (string, string) {
		name := ""
		returnFuncName := ""
		if Gid {
			name = osbased.GroupID(info)
			if runtime.GOOS == "windows" {
				returnFuncName = GroupSID
			} else {
				returnFuncName = GroupUidName
			}
		} else {
			name = osbased.Group(info)
			returnFuncName = GroupName
		}
		done(name)
		return wait(name), returnFuncName
	}
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

func (cf *ContentFilter) GetDisplayItems(e ...os.FileInfo) []*display.Item {
	sort.Slice(e, func(i, j int) bool {
		if cf.sortFunc != nil {
			return cf.sortFunc(e[i], e[j])
		} else {
			return true
		}
	})

	wg := sync.WaitGroup{}
	wg.Add(len(e))
	for i := range cf.wgs {
		cf.wgs[i].Add(len(e))
	}

	res := make([]*display.Item, 0, len(e))
	for i := 0; i < len(e); i++ {
		res = append(res, display.NewItem(display.WithDelimiter(" ")))
	}

	for i, entry := range e {
		go func(entry os.FileInfo, i int) {
			for j, option := range cf.options {
				stringContent, funcName := option(entry)
				content := display.ItemContent{Content: display.StringContent(stringContent), No: j}
				res[i].Set(funcName, content)
			}

			for _, option := range cf.noOutputOptions {
				option(entry)
			}
			wg.Done()
		}(entry, i)
	}
	wg.Wait()

	return res
}

type SumType int

const (
	SumTypeMd5 SumType = iota + 1
	SumTypeSha1
	SumTypeSha224
	SumTypeSha256
	SumTypeSha384
	SumTypeSha512
	SumTypeCRC32
)

const SumName = "sum"

func (cf *ContentFilter) EnableSum(sumTypes ...SumType) ContentOption {
	length := 0
	types := make([]string, 0, len(sumTypes))
	for _, t := range sumTypes {
		switch t {
		case SumTypeMd5:
			length += 32
			types = append(types, "md5")
		case SumTypeSha1:
			length += 40
			types = append(types, "sha1")
		case SumTypeSha224:
			length += 56
			types = append(types, "sha224")
		case SumTypeSha256:
			length += 64
			types = append(types, "sha256")
		case SumTypeSha384:
			length += 96
			types = append(types, "sha384")
		case SumTypeSha512:
			length += 128
			types = append(types, "sha512")
		case SumTypeCRC32:
			length += 8
			types = append(types, "crc32")
		}
	}
	length += len(sumTypes) - 1
	sumName := fmt.Sprintf("%s(%s)", SumName, strings.Join(types, ","))
	return func(info os.FileInfo) (string, string) {
		if info.IsDir() {
			return FillBlank("", length), sumName
		}

		file, err := os.Open(info.Name())
		if err != nil {
			return FillBlank("", length), sumName
		}
		defer file.Close()
		hashes := make([]hash.Hash, 0, len(sumTypes))
		writers := make([]io.Writer, 0, len(sumTypes))
		for _, t := range sumTypes {
			var hashed hash.Hash
			switch t {
			case SumTypeMd5:
				hashed = md5.New()
			case SumTypeSha1:
				hashed = sha1.New()
			case SumTypeSha224:
				hashed = sha256.New224()
			case SumTypeSha256:
				hashed = sha256.New()
			case SumTypeSha384:
				hashed = sha512.New384()
			case SumTypeSha512:
				hashed = sha512.New()
			case SumTypeCRC32:
				hashed = crc32.NewIEEE()
			}
			writers = append(writers, hashed)
			hashes = append(hashes, hashed)
		}
		multiWriter := io.MultiWriter(writers...)
		if _, err := io.Copy(multiWriter, file); err != nil {
			return FillBlank("", length), sumName
		}
		sums := make([]string, 0, len(hashes))
		for _, h := range hashes {
			sums = append(sums, fmt.Sprintf("%x", h.Sum(nil)))
		}
		sumsStr := strings.Join(sums, " ")
		return FillBlank(sumsStr, length), sumName
	}
}
