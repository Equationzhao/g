package sorter

import (
	"os"
	"path/filepath"
	"strings"
)

func ByNameDescend(a, b os.FileInfo) bool {
	return strings.ToLower(a.Name()) > strings.ToLower(b.Name())
}

func ByNameAscend(a, b os.FileInfo) bool {
	return strings.ToLower(a.Name()) < strings.ToLower(b.Name())
}

func ByNameCaseSensitiveDescend(a, b os.FileInfo) bool {
	return a.Name() > b.Name()
}

func ByNameCaseSensitiveAscend(a, b os.FileInfo) bool {
	return a.Name() < b.Name()
}

func BySizeDescend(a, b os.FileInfo) bool {
	return a.Size() > b.Size()
}

func BySizeAscend(a, b os.FileInfo) bool {
	return a.Size() < b.Size()
}

func ByTimeDescend(a, b os.FileInfo) bool {
	return a.ModTime().After(b.ModTime())
}

func ByTimeAscend(a, b os.FileInfo) bool {
	return a.ModTime().Before(b.ModTime())
}

func ByExtensionDescend(a, b os.FileInfo) bool {
	return filepath.Ext(a.Name()) > filepath.Ext(b.Name())
}

func ByExtensionAscend(a, b os.FileInfo) bool {
	return filepath.Ext(a.Name()) < filepath.Ext(b.Name())
}

func ByExtensionCaseSensitiveDescend(a, b os.FileInfo) bool {
	return strings.ToLower(filepath.Ext(a.Name())) > strings.ToLower(filepath.Ext(b.Name()))
}

func ByExtensionCaseSensitiveAscend(a, b os.FileInfo) bool {
	return strings.ToLower(filepath.Ext(a.Name())) < strings.ToLower(filepath.Ext(b.Name()))
}

func ByGroupDescend(a, b os.FileInfo) bool {
	return byGroupName(a, b, false)
}

func ByGroupAscend(a, b os.FileInfo) bool {
	return byGroupName(a, b, true)
}

func ByGroupCaseSensitiveDescend(a, b os.FileInfo) bool {
	return byGroupCaseSensitiveName(a, b, false)
}

func ByGroupCaseSensitiveAscend(a, b os.FileInfo) bool {
	return byGroupCaseSensitiveName(a, b, true)
}

func ByOwnerDescend(a, b os.FileInfo) bool {
	return byUserName(a, b, false)
}

func ByOwnerAscend(a, b os.FileInfo) bool {
	return byUserName(a, b, true)
}

func ByOwnerCaseSensitiveDescend(a, b os.FileInfo) bool {
	return byUserCaseSensitiveName(a, b, false)
}

func ByOwnerCaseSensitiveAscend(a, b os.FileInfo) bool {
	return byUserCaseSensitiveName(a, b, true)
}

func dirFirst(a, b os.FileInfo) bool {
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
	return false
}

func Default(a, b os.FileInfo) bool {
	return compareFileInfo(a, b)
}

type FileSortFunc = func(a, b os.FileInfo) bool

type Sorter struct {
	reverse  bool
	option   []FileSortFunc
	dirFirst bool
}

func (s *Sorter) DirFirst() {
	s.dirFirst = true
}

func (s *Sorter) UnsetDirFirst() {
	s.dirFirst = false
}

func (s *Sorter) Len() int {
	return len(s.option)
}

type Option = func(s *Sorter)

func WithSize(size int) Option {
	return func(s *Sorter) {
		s.option = make([]FileSortFunc, 0, size)
	}
}

func WithSortOption(option ...FileSortFunc) Option {
	return func(s *Sorter) {
		s.option = option
	}
}

func NewSorter(option ...Option) *Sorter {
	a := Sorter{}
	for _, opt := range option {
		opt(&a)
	}
	return &a
}

func (s *Sorter) Reverse() {
	s.reverse = !s.reverse
}

func (s *Sorter) AddOption(option ...FileSortFunc) {
	s.option = append(s.option, option...)
}

func (s *Sorter) Build() FileSortFunc {
	return func(a, b os.FileInfo) bool {
		result := false
		for _, sortFunc := range s.option {
			if s.dirFirst {
				if dirFirst(a, b) {
					result = true
					break
				}
				if dirFirst(b, a) {
					result = false
					break
				}
			}
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

func isHidden(info os.FileInfo) bool {
	return strings.HasPrefix(info.Name(), ".")
}

func isLink(info os.FileInfo) bool {
	return info.Mode()&os.ModeSymlink != 0
}

func isHiddenDir(info os.FileInfo) bool {
	return isHidden(info) && info.IsDir()
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
