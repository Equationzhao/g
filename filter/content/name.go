package content

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/Equationzhao/g/util"

	"github.com/Equationzhao/g/filter"
	"github.com/Equationzhao/g/item"
	"github.com/Equationzhao/g/render"
	"github.com/valyala/bytebufferpool"
)

type (
	Name struct {
		icon, classify, fileType, git, fullPath, noDeference bool
		statistics                                           *Statistics
		relativeTo                                           string
		Quote                                                string
		GitStyle                                             gitStyle
	}
)

type Statistics struct {
	file, dir, link atomic.Uint64
}

func (s *Statistics) Reset() {
	s.file.Store(0)
	s.dir.Store(0)
	s.link.Store(0)
}

func (s *Statistics) String() string {
	return fmt.Sprintf("%d file(s), %d dir(s), %d link(s)", s.file.Load(), s.dir.Load(), s.link.Load())
}

type gitStyle int

const (
	GitStyleDot gitStyle = iota
	GitStyleSym
	GitStyleDefault = GitStyleDot
)

func (n *Name) SetNoDeference() *Name {
	n.noDeference = true
	return n
}

func (n *Name) UnsetNoDeference() *Name {
	n.noDeference = false
	return n
}

func (n *Name) FullPath() bool {
	return n.fullPath
}

func (n *Name) SetFullPath() *Name {
	n.fullPath = true
	return n
}

func (n *Name) UnsetFullPath() *Name {
	n.fullPath = false
	return n
}

func (n *Name) RelativeTo() string {
	return n.relativeTo
}

func (n *Name) SetRelativeTo(relativeTo string) {
	n.relativeTo = relativeTo
}

func (n *Name) Statistics() *Statistics {
	return n.statistics
}

func (n *Name) SetStatistics(Statistics *Statistics) *Name {
	n.statistics = Statistics
	return n
}

func (n *Name) SetQuote(quote string) *Name {
	n.Quote = quote
	return n
}

func (n *Name) UnsetQuote() *Name {
	n.Quote = ""
	return n
}

func (n *Name) UnsetIcon() *Name {
	n.icon = false
	return n
}

func (n *Name) UnsetClassify() *Name {
	n.classify = false
	return n
}

func (n *Name) UnsetFileType() *Name {
	n.fileType = false
	return n
}

func (n *Name) SetIcon() *Name {
	n.icon = true
	return n
}

func (n *Name) SetClassify() *Name {
	n.classify = true
	return n
}

// SetFileType set file type
// should set classify first
// if classify is false, fileType will be ignored
func (n *Name) SetFileType() *Name {
	n.fileType = true
	return n
}

func NewNameEnable() *Name {
	return &Name{}
}

const NameName = "Name"

func (n *Name) Enable(renderer *render.Renderer) filter.ContentOption {
	/*
		 -F      Display a slash (`/`) immediately after each pathname that is a
				 directory, an asterisk (`*`) after each that is executable, an at
				 sign (`@`) after each symbolic link, a percent sign (`%`) after
				 each whiteout, an equal sign (`=`) after each socket, and a
				 vertical bar (`|`) after each that is a FIFO.
	*/

	return func(info *item.FileInfo) (string, string) {
		buffer := bytebufferpool.Get()
		defer bytebufferpool.Put(buffer)
		str := info.Name()
		if n.FullPath() {
			str = info.FullPath
		}
		name := str
		mode := info.Mode()

		char := ""

		if n.icon {
			if info.IsDir() {
				if n.statistics != nil {
					n.statistics.dir.Add(1)
				}
				str = renderer.DirIcon(str)
				char = "/"
			} else if mode&os.ModeSymlink != 0 {
				if n.statistics != nil {
					n.statistics.link.Add(1)
				}
				if n.classify {
					if n.noDeference {
						str = renderer.SymlinkIconNoDereferencePlus(str, "@")
					} else {
						str = renderer.SymlinkIconPlus(str, info.FullPath, "@", !n.fullPath)
					}
				} else {
					if n.noDeference {
						str = renderer.SymlinkIconNoDereference(str)
					} else {
						str = renderer.SymlinkIcon(str, info.FullPath, !n.fullPath)
					}
				}
			} else {
				if n.statistics != nil {
					n.statistics.file.Add(1)
				}
				if mode&os.ModeNamedPipe != 0 {
					str = renderer.PipeIcon(str)
					char = "|"
				} else if mode&os.ModeSocket != 0 {
					str = renderer.SocketIcon(str)
					char = "="
				} else {
					str = renderer.ByExtIcon(str)
				}
			}
		} else {
			if info.IsDir() {
				if n.statistics != nil {
					n.statistics.dir.Add(1)
				}
				str = renderer.Dir(str)
				char = "/"
			} else if mode&os.ModeSymlink != 0 {
				if n.statistics != nil {
					n.statistics.link.Add(1)
				}
				if n.classify {
					if n.noDeference {
						str = renderer.SymlinkNoDereferencePlus(str, "@")
					} else {
						str = renderer.SymlinkPlus(str, info.FullPath, "@", !n.fullPath)
					}
				} else {
					if n.noDeference {
						str = renderer.SymlinkNoDereference(str)
					} else {
						str = renderer.Symlink(str, info.FullPath, !n.fullPath)
					}
				}
			} else {
				if n.statistics != nil {
					n.statistics.file.Add(1)
				}
				if mode&os.ModeNamedPipe != 0 {
					str = renderer.Pipe(str)
					char = "|"
				} else if mode&os.ModeSocket != 0 {
					str = renderer.Socket(str)
					char = "="
				} else {
					str = renderer.ByExt(str)
				}
			}
		}

		if n.classify {
			if char == "" && (!n.fileType) && util.IsExecutableMode(mode) && mode&os.ModeSymlink == 0 {
				str += "*"
			} else {
				str += char
			}
		}

		if n.Quote != "" {
			str = strings.Replace(str, name, n.Quote+name+n.Quote, 1)
		}

		if n.relativeTo != "" {
			relativePath, err := filepath.Rel(n.relativeTo, info.FullPath)
			if err != nil {
				return str, NameName
			}
			str = strings.Replace(str, name, relativePath, 1)
		}

		return str, NameName
	}
}
