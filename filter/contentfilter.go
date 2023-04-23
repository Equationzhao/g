package filter

import (
	"github.com/Equationzhao/g/render"
	"github.com/valyala/bytebufferpool"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// fileMode size owner group time name

type ContentFilter struct {
	options []ContentOption
}

type ContentOption func(info os.FileInfo) string

// EnableFileMode return file mode like -rwxrwxrwx/drwxrwxrwx
func EnableFileMode(renderer *render.Renderer) ContentOption {
	return func(info os.FileInfo) string {
		return renderer.FileMode(fillBlank(info.Mode().String(), 12))
	}
}

type Size int

const (
	Bit Size = iota
	B
	KB
	MB
	GB
	TB
	PB
	EB
	ZB
	Auto
)

// fill blank
// if s is shorter than length, fill blank from left
// if s is longer than length, panic
func fillBlank(s string, length int) string {
	return strings.Repeat(" ", length-len(s)) + s
}

func convert2Size(size Size) string {
	switch size {
	case Bit:
		return "bit"
	case B:
		return "B"
	case KB:
		return "KB"
	case MB:
		return "MB"
	case GB:
		return "GB"
	case TB:
		return "TB"
	case PB:
		return "PB"
	case EB:
		return "EB"
	case ZB:
		return "ZB"
	default:
		panic("unknown size")
	}
}

func EnableSize(size Size, renderer *render.Renderer) ContentOption {
	return func(info os.FileInfo) string {
		v := float64(info.Size())
		var res string
		switch size {
		case Bit:
			res = strconv.FormatInt(int64(v*8), 10) + "bit"
		case B:
			res = strconv.FormatInt(int64(v), 10) + "B"
		case KB:
			res = strconv.FormatFloat(v/1024.0, 'f', 0, 64) + "KB"
		case MB:
			res = strconv.FormatFloat(v/1024.0/1024.0, 'f', 1, 64) + "MB"
		case GB:
			res = strconv.FormatFloat(v/1024.0/1024.0/1024.0, 'f', 1, 64) + "GB"
		case TB:
			res = strconv.FormatFloat(v/1024.0/1024.0/1024.0/1024.0, 'f', 1, 64) + "TB"
		case PB:
			res = strconv.FormatFloat(v/1024.0/1024.0/1024.0/1024.0/1024.0, 'f', 1, 64) + "PB"
		case EB:
			res = strconv.FormatFloat(v/1024.0/1024.0/1024.0/1024.0/1024.0/1024.0, 'f', 1, 64) + "EB"
		case ZB:
			res = strconv.FormatFloat(v/1024.0/1024.0/1024.0/1024.0/1024.0/1024.0/1024.0, 'f', 1, 64) + "ZB"
		case Auto:
			for i := B; i <= ZB; i++ {
				if v < 1000 {
					res = strconv.FormatFloat(v, 'f', 1, 64)
					if res == "0.0" {
						res = ""
					} else {
						res += convert2Size(i)
					}
					return renderer.Size(fillBlank(res, 7))
				}
				v /= 1024
			}
			panic("too large")
		default:
			panic("invalid " + strconv.Itoa(int(size)))
		}
		return renderer.Size(res)
	}
}

func EnableTime(format string, renderer *render.Renderer) ContentOption {
	return func(info os.FileInfo) string {
		return renderer.Time(info.ModTime().Format(format))
	}
}

func EnableName(renderer *render.Renderer) ContentOption {
	return func(info os.FileInfo) string {
		if info.IsDir() {
			return renderer.Dir(info.Name())
		} else if info.Mode()&os.ModeSymlink != 0 {
			return renderer.Symlink(info.Name())
		} else {
			return renderer.ByExt(info.Name())
		}
	}
}

func EnableIconName(renderer *render.Renderer, s string) ContentOption {
	return func(info os.FileInfo) string {
		if info.IsDir() {
			return renderer.DirIcon(info.Name())
		} else if info.Mode()&os.ModeSymlink != 0 {
			return renderer.SymlinkIcon(info.Name(), s)
		} else {
			return renderer.ByExtIcon(info.Name())
		}
	}
}

func NewContentFilter(options ...ContentOption) *ContentFilter {
	return &ContentFilter{options: options}
}

type ContentFunc func(entry os.FileInfo) bool

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

func (cf *ContentFilter) GetStringSlice(e []os.FileInfo) []string {
	resBuffers := make([]*bytebufferpool.ByteBuffer, len(e))

	for i := range resBuffers {
		resBuffers[i] = bytebufferpool.Get()
	}

	defer func() {
		for i := range resBuffers {
			bytebufferpool.Put(resBuffers[i])
		}
	}()

	sort.Slice(e, func(i, j int) bool {
		return compareFileInfo(e[i], e[j])
	})

	wg := sync.WaitGroup{}
	wg.Add(len(e))
	for i, entry := range e {
		entry := entry
		i := i
		go func() {
			for _, option := range cf.options {
				_, _ = resBuffers[i].WriteString(option(entry))
				_ = resBuffers[i].WriteByte(' ')
			}
			wg.Done()
		}()
	}
	res := make([]string, 0, len(e))
	wg.Wait()
	for _, builder := range resBuffers {
		res = append(res, builder.String())
	}

	return res
}
