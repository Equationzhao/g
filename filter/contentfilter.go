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
			res = strconv.FormatInt(int64(v*8), 10)
		case B:
			res = strconv.FormatInt(int64(v), 10)
		case KB:
			res = strconv.FormatFloat(v/1024.0, 'f', 0, 64)
		case MB:
			res = strconv.FormatFloat(v/1024.0/1024.0, 'f', 1, 64)
		case GB:
			res = strconv.FormatFloat(v/1024.0/1024.0/1024.0, 'f', 1, 64)
		case TB:
			res = strconv.FormatFloat(v/1024.0/1024.0/1024.0/1024.0, 'f', 1, 64)
		case PB:
			res = strconv.FormatFloat(v/1024.0/1024.0/1024.0/1024.0/1024.0, 'f', 1, 64)
		case EB:
			res = strconv.FormatFloat(v/1024.0/1024.0/1024.0/1024.0/1024.0/1024.0, 'f', 1, 64)
		case ZB:
			res = strconv.FormatFloat(v/1024.0/1024.0/1024.0/1024.0/1024.0/1024.0/1024.0, 'f', 1, 64)
		case Auto:
			for i := B; i <= ZB; i++ {
				if v < 1000 {
					res = strconv.FormatFloat(v, 'f', 1, 64)
					if res == "0.0" {
						res = "-"
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
		return renderer.Size(fillBlank(res, 7))
	}
}

func EnableTime(format string, renderer *render.Renderer) ContentOption {
	return func(info os.FileInfo) string {
		return renderer.Time(info.ModTime().Format(format))
	}
}

type Name struct {
	Icon     bool
	Classify bool // --classify / -F
	FileType bool // --file-type
	Renderer *render.Renderer
	parent   string
}

func (n *Name) SetParent(parent string) *Name {
	n.parent = parent
	return n
}

func (n *Name) UnsetIcon() *Name {
	n.Icon = false
	return n
}

func (n *Name) UnsetClassify() *Name {
	n.Classify = false
	return n
}

func (n *Name) UnsetFileType() *Name {
	n.FileType = false
	return n
}

func (n *Name) SetIcon() *Name {
	n.Icon = true
	return n
}

func (n *Name) SetClassify() *Name {
	n.Classify = true
	return n
}

// SetFileType set file type, should set Classify first
// if Classify is false, FileType will be ignored
func (n *Name) SetFileType() *Name {
	n.FileType = true
	return n
}

func (n *Name) SetRenderer(Renderer *render.Renderer) *Name {
	n.Renderer = Renderer
	return n
}

func NewNameEnable() *Name {
	return &Name{}
}

func (n *Name) Enable() ContentOption {
	/*
		 -F      Display a slash (`/') immediately after each pathname that is a
				 directory, an asterisk (`*') after each that is executable, an at
				 sign (`@') after each symbolic link, a percent sign (`%') after
				 each whiteout, an equal sign (`=') after each socket, and a
				 vertical bar (`|') after each that is a FIFO.
	*/

	return func(info os.FileInfo) string {
		buffer := bytebufferpool.Get()
		defer bytebufferpool.Put(buffer)
		str := info.Name()
		mode := info.Mode()

		if n.Icon {
			if info.IsDir() {
				str = n.Renderer.DirIcon(str)
			} else if mode&os.ModeSymlink != 0 {
				if n.Classify {
					str = n.Renderer.SymlinkIconPlus(str, n.parent, "@")
				} else {
					str = n.Renderer.SymlinkIcon(str, n.parent)
				}
			} else if mode&os.ModeNamedPipe != 0 {
				str = n.Renderer.PipeIcon(str)
			} else if mode&os.ModeSocket != 0 {
				str = n.Renderer.SocketIcon(str)
			} else {
				str = n.Renderer.ByExtIcon(str)
			}
		} else {
			if info.IsDir() {
				str = n.Renderer.Dir(str)
			} else if mode&os.ModeSymlink != 0 {
				if n.Classify {
					str = n.Renderer.SymlinkPlus(str, n.parent, "@")
				} else {
					str = n.Renderer.Symlink(str, n.parent)
				}
			} else if mode&os.ModeNamedPipe != 0 {
				str = n.Renderer.Pipe(str)
			} else if mode&os.ModeSocket != 0 {
				str = n.Renderer.Socket(str)
			} else {
				str = n.Renderer.ByExt(str)
			}
		}

		if n.Classify {
			if info.IsDir() {
				str += "/"
			} else if mode&os.ModeSymlink != 0 {
				goto end
			} else if mode&os.ModeNamedPipe != 0 {
				str += "|"
			} else if mode&os.ModeSocket != 0 {
				str += "="
			} else if (!n.FileType) && (mode&0111 != 0) {
				str += "*"
			}
		}

	end:
		return str
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
		go func(entry os.FileInfo, i int) {
			options := cf.options[:len(cf.options)-1]
			for _, option := range options {
				_, _ = resBuffers[i].WriteString(option(entry))
				_ = resBuffers[i].WriteByte(' ')
			}
			// the last one should not follow by space
			_, _ = resBuffers[i].WriteString(cf.options[len(cf.options)-1](entry))
			wg.Done()
		}(entry, i)
	}
	res := make([]string, 0, len(e))
	wg.Wait()
	for _, buffer := range resBuffers {
		res = append(res, buffer.String())
	}

	return res
}
