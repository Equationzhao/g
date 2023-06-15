package sorter

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Equationzhao/g/osbased"
	mt "github.com/gabriel-vasile/mimetype"
)

//goland:noinspection GoUnusedParameter
func ByNone(a, b os.FileInfo) bool {
	return false
}

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

func ByNameWithoutALeadingDotDescend(a, b os.FileInfo) bool {
	return strings.ToLower(strings.TrimPrefix(a.Name(), ".")) > strings.ToLower(strings.TrimPrefix(b.Name(), "."))
}

func ByNameWithoutALeadingDotAscend(a, b os.FileInfo) bool {
	return strings.ToLower(strings.TrimPrefix(a.Name(), ".")) < strings.ToLower(strings.TrimPrefix(b.Name(), "."))
}

func ByNameWithoutALeadingDotCaseSensitiveDescend(a, b os.FileInfo) bool {
	return strings.TrimPrefix(a.Name(), ".") > strings.TrimPrefix(b.Name(), ".")
}

func ByNameWithoutALeadingDotCaseSensitiveAscend(a, b os.FileInfo) bool {
	return strings.TrimPrefix(a.Name(), ".") < strings.TrimPrefix(b.Name(), ".")
}

func ByInodeDescend(a, b os.FileInfo) bool {
	inodeA, _ := strconv.Atoi(osbased.Inode(a))
	inodeB, _ := strconv.Atoi(osbased.Inode(b))

	return inodeA > inodeB
}

func ByInodeAscend(a, b os.FileInfo) bool {
	inodeA, _ := strconv.Atoi(osbased.Inode(a))
	inodeB, _ := strconv.Atoi(osbased.Inode(b))

	return inodeA < inodeB
}

func BySizeDescend(a, b os.FileInfo) bool {
	return a.Size() > b.Size()
}

func BySizeAscend(a, b os.FileInfo) bool {
	return a.Size() < b.Size()
}

func ByTimeDescend(timeType string) FileSortFunc {
	switch timeType {
	case "mod":
		return func(a, b os.FileInfo) bool {
			return osbased.ModTime(a).After(osbased.ModTime(b))
		}
	case "access":
		return func(a, b os.FileInfo) bool {
			return osbased.AccessTime(a).After(osbased.AccessTime(b))
		}
	case "create":
		return func(a, b os.FileInfo) bool {
			return osbased.CreateTime(a).After(osbased.CreateTime(b))
		}
	default:
		panic("invalid time type")
	}
}

func ByTimeAscend(timeType string) FileSortFunc {
	switch timeType {
	case "mod":
		return func(a, b os.FileInfo) bool {
			return osbased.ModTime(a).Before(osbased.ModTime(b))
		}
	case "access":
		return func(a, b os.FileInfo) bool {
			return osbased.AccessTime(a).Before(osbased.AccessTime(b))
		}
	case "create":
		return func(a, b os.FileInfo) bool {
			return osbased.CreateTime(a).Before(osbased.CreateTime(b))
		}
	default:
		panic("invalid time type")
	}
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

func ByNameWidthDescend(a, b os.FileInfo) bool {
	return byNameWidth(a, b, false)
}

func ByNameWidthAscend(a, b os.FileInfo) bool {
	return byNameWidth(a, b, true)
}

func ByMimeTypeAscend(a, b os.FileInfo) bool {
	return byMimeType(a, b, true)
}

func ByMimeTypeDescend(a, b os.FileInfo) bool {
	return byMimeType(a, b, false)
}

func byMimeType(a, b os.FileInfo, ascend bool) bool {
	mimeAstr, mimeBstr := getMimeName(a, b)
	if ascend {
		return mimeAstr < mimeBstr
	}
	return mimeBstr < mimeAstr
}

func getMimeName(a os.FileInfo, b os.FileInfo) (string, string) {
	mimeAstr, mimeBstr := "", ""
	mimeA, err := mt.DetectFile(a.Name())
	if err != nil {
		if a.IsDir() {
			mimeAstr = "directory"
		} else if a.Mode()&os.ModeSymlink != 0 {
			mimeAstr = "symlink"
		} else if a.Mode()&os.ModeNamedPipe != 0 {
			mimeAstr = "named_pipe"
		} else if a.Mode()&os.ModeSocket != 0 {
			mimeAstr = "socket"
		}
	} else {
		mimeAstr = mimeA.String()
	}
	mimeB, err := mt.DetectFile(b.Name())
	if err != nil {
		if b.IsDir() {
			mimeBstr = "directory"
		} else if b.Mode()&os.ModeSymlink != 0 {
			mimeBstr = "symlink"
		} else if b.Mode()&os.ModeNamedPipe != 0 {
			mimeBstr = "named_pipe"
		} else if b.Mode()&os.ModeSocket != 0 {
			mimeBstr = "socket"
		}
	} else {
		mimeBstr = mimeB.String()
	}
	return mimeAstr, mimeBstr
}

func ByMimeTypeParentAscend(a, b os.FileInfo) bool {
	return byMimeTypeParent(a, b, true)
}

func ByMimeTypeParentDescend(a, b os.FileInfo) bool {
	return byMimeTypeParent(a, b, false)
}

func byMimeTypeParent(a, b os.FileInfo, ascend bool) bool {
	mimeAstr, mimeBstr := getMimeName(a, b)
	if ascend {
		return strings.SplitN(mimeAstr, "/", 2)[0] < strings.SplitN(mimeBstr, "/", 2)[0]
	} else {
		return strings.SplitN(mimeAstr, "/", 2)[0] > strings.SplitN(mimeBstr, "/", 2)[0]
	}
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
	dirFirst bool
	option   []FileSortFunc
}

func (s *Sorter) Reset() {
	s.reverse = false
	s.dirFirst = false
	s.option = make([]FileSortFunc, 0, 10)
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
		s.option = append(s.option, option...)
	}
}

func NewSorter(option ...Option) *Sorter {
	a := Sorter{
		reverse:  false,
		dirFirst: false,
		option:   make([]FileSortFunc, 0, 10),
	}

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
		}
		return result
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

func byGroupName(a, b os.FileInfo, Ascend bool) bool {
	if Ascend {
		return strings.ToLower(osbased.Group(a)) < strings.ToLower(osbased.Group(b))
	}
	return strings.ToLower(osbased.Group(a)) > strings.ToLower(osbased.Group(b))
}

func byUserName(a, b os.FileInfo, Ascend bool) bool {
	if Ascend {
		return strings.ToLower(osbased.Owner(a)) < strings.ToLower(osbased.Owner(b))
	}
	return strings.ToLower(osbased.Owner(a)) > strings.ToLower(osbased.Owner(b))
}

func byGroupCaseSensitiveName(a, b os.FileInfo, Ascend bool) bool {
	if Ascend {
		return osbased.Group(a) < osbased.Group(b)
	}
	return osbased.Group(a) > osbased.Group(b)
}

func byUserCaseSensitiveName(a, b os.FileInfo, Ascend bool) bool {
	if Ascend {
		return osbased.Owner(a) < osbased.Owner(b)
	}
	return osbased.Owner(a) > osbased.Owner(b)
}

func byNameWidth(a, b os.FileInfo, Ascend bool) bool {
	if Ascend {
		return len(a.Name()) < len(b.Name())
	}
	return len(a.Name()) > len(b.Name())
}
