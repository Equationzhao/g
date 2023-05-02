package sorter

import (
	"os"
	"path/filepath"
	"strings"
)

type SortOption int

const (
	ByNameDescend SortOption = iota
	ByNameAscend
	BySizeDescend
	BySizeAscend
	ByTimeDescend
	ByTimeAscend
	ByExtensionDescend
	ByExtensionAscend
	ByGroupDescend
	ByGroupAscend
	ByOwnerDescend
	ByOwnerAscend
	Default
)

type FileSortFunc = func(a, b os.FileInfo) bool

type Sorter struct {
	reverse bool
	option  []SortOption
}

func (s *Sorter) Len() int {
	return len(s.option)
}

type Option = func(s *Sorter)

func WithSize(size int) Option {
	return func(s *Sorter) {
		s.option = make([]SortOption, 0, size)
	}
}

func WithSortOption(option ...SortOption) Option {
	return func(s *Sorter) {
		s.option = option
	}
}

func NewSorter(option ...Option) *Sorter {
	return &Sorter{}
}

func (s *Sorter) Reverse() {
	s.reverse = !s.reverse
}

func (s *Sorter) AddOption(option ...SortOption) {
	s.option = append(s.option, option...)
}

func (s *Sorter) Build() FileSortFunc {
	sortFuncs := make([]FileSortFunc, 0, len(s.option))
	for _, option := range s.option {
		sortFuncs = append(sortFuncs, s.buildSortFunc(option))
	}
	return func(a, b os.FileInfo) bool {
		result := false
		for _, sortFunc := range sortFuncs {
			if sortFunc(a, b) {
				result = true
				break
			}
			if sortFunc(b, a) {
				result = false
				break
			}
		}

		if s.reverse {
			return !result
		} else {
			return result
		}
	}
}

func (s *Sorter) buildSortFunc(option SortOption) FileSortFunc {
	switch option {
	case ByNameDescend:
		return func(a, b os.FileInfo) bool {
			return a.Name() > b.Name()
		}
	case ByNameAscend:
		return func(a, b os.FileInfo) bool {
			return a.Name() < b.Name()
		}
	case BySizeDescend:
		return func(a, b os.FileInfo) bool {
			return a.Size() > b.Size()
		}
	case BySizeAscend:
		return func(a, b os.FileInfo) bool {
			return a.Size() < b.Size()
		}
	case ByTimeDescend:
		return func(a, b os.FileInfo) bool {
			return a.ModTime().After(b.ModTime())
		}
	case ByTimeAscend:
		return func(a, b os.FileInfo) bool {
			return a.ModTime().Before(b.ModTime())
		}
	case ByExtensionDescend:
		return func(a, b os.FileInfo) bool {
			return filepath.Ext(a.Name()) > filepath.Ext(b.Name())
		}
	case ByExtensionAscend:
		return func(a, b os.FileInfo) bool {
			return filepath.Ext(a.Name()) < filepath.Ext(b.Name())
		}
	case ByGroupDescend:
		return func(a, b os.FileInfo) bool {
			return byGroupName(a, b, false)
		}
	case ByGroupAscend:
		return func(a, b os.FileInfo) bool {
			return byGroupName(a, b, true)
		}
	case ByOwnerDescend:
		return func(a, b os.FileInfo) bool {
			return byUserName(a, b, false)
		}
	case ByOwnerAscend:
		return func(a, b os.FileInfo) bool {
			return byUserName(a, b, true)
		}
	case Default:
		return defaultSort
	default:
		return defaultSort
	}
}

func isHidden(info os.FileInfo) bool {
	return strings.HasPrefix(info.Name(), ".")
}

func isLink(info os.FileInfo) bool {
	return info.Mode()&os.ModeSymlink != 0
}

func isHiddenDir(info os.FileInfo) bool {
	return info.IsDir() && isHidden(info)
}

func compareFileInfo(a, b os.FileInfo) bool {
	hdA := isHiddenDir(a)
	hdB := isHiddenDir(b)
	if hdA != hdB {
		// hidden dir comes first
		return hdA
	}
	// same hidden dir status
	dA := a.IsDir()
	dB := b.IsDir()
	if dA != dB {
		// dir comes first
		return dA
	}
	// same dir status
	lA := isLink(a)
	lB := isLink(b)
	switch {
	case lA && lB:
		// both are links, compare name
		return a.Name() < b.Name()
	case lA:
		// a is link, b is not link, b comes first
		return false
	case lB:
		// a is not link, b is link, a comes first
		return true
	default:
		// neither are links, compare name
		return a.Name() < b.Name()
	}
}

func defaultSort(a, b os.FileInfo) bool {
	return compareFileInfo(a, b)
}
