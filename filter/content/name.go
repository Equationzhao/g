package content

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/Equationzhao/g/filter"
	"github.com/Equationzhao/g/item"
	"github.com/Equationzhao/g/render"
	"github.com/Equationzhao/g/theme"
	"github.com/Equationzhao/g/util"
	"github.com/valyala/bytebufferpool"
)

type (
	Name struct {
		icon, classify, fileType, fullPath, noDeference, hyperLink, mounts bool
		statistics                                                         *Statistics
		relativeTo                                                         string
		Quote                                                              string
	}
)

type Statistics struct {
	file, dir, link atomic.Uint64
}

func (s *Statistics) MarshalJSON() ([]byte, error) {
	Export := struct {
		File, Dir, Link uint64
	}{
		File: s.file.Load(),
		Dir:  s.dir.Load(),
		Link: s.link.Load(),
	}
	return json.Marshal(Export)
}

func (s *Statistics) Reset() {
	s.file.Store(0)
	s.dir.Store(0)
	s.link.Store(0)
}

func (s *Statistics) String() string {
	return fmt.Sprintf("%d file(s), %d dir(s), %d link(s)", s.file.Load(), s.dir.Load(), s.link.Load())
}

func (n *Name) SetNoDeference() *Name {
	n.noDeference = true
	return n
}

func (n *Name) UnsetNoDeference() *Name {
	n.noDeference = false
	return n
}

func (n *Name) SetHyperlink() *Name {
	n.hyperLink = true
	return n
}

func (n *Name) UnsetHyperlink() *Name {
	n.hyperLink = false
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

func (n *Name) SetMounts() *Name {
	n.mounts = true
	return n
}

func (n *Name) UnsetMounts() *Name {
	n.mounts = false
	return n
}

func NewNameEnable() *Name {
	return &Name{}
}

const NameName = "Name"

func makeLink(abs string, name string) string {
	return fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", abs, name)
}

/*
Enable
color + icon + file://quote+filename/relative-name+quote + classify + color-end + dereference + mounts
color: filetype->filename->fileext->file
*/
func (n *Name) Enable(renderer *render.Renderer) filter.ContentOption {
	return func(info *item.FileInfo) (stringContent string, funcName string) {
		name := info.Name()
		color := ""
		icon := ""
		classify := ""
		dereference := bytebufferpool.Get()
		defer bytebufferpool.Put(dereference)
		mounts := ""
		mode := info.Mode()
		underline, bold := false, false

		if info.IsDir() {
			if n.statistics != nil {
				n.statistics.dir.Add(1)
			}
			style := renderer.Dir(name)
			if n.icon {
				icon = style.Icon
			}
			if n.classify {
				classify = "/"
			}
			color, underline, bold = style.Color, style.Underline, style.Bold
		} else if util.IsSymLinkMode(mode) {
			if n.statistics != nil {
				n.statistics.link.Add(1)
			}
			style := renderer.Symlink()
			if n.icon {
				icon = style.Icon
			}
			color, underline, bold = style.Color, style.Underline, style.Bold
			if n.classify {
				classify = "@"
			}
			// color + arrow + color-end + color + path + color-end
			if !n.noDeference {
				arrowStyle := renderer.SymlinkArrow()
				_, _ = dereference.WriteString(arrowStyle.Color)
				if arrowStyle.Underline {
					_, _ = dereference.WriteString(theme.Underline)
				}
				if arrowStyle.Bold {
					_, _ = dereference.WriteString(theme.Bold)
				}
				_, _ = dereference.WriteString(arrowStyle.Icon)
				_, _ = dereference.WriteString(renderer.Colorend())
				broken := false
				symlinks, err := filepath.EvalSymlinks(info.FullPath)
				if err != nil {
					broken = true
					var pathErr *fs.PathError
					if errors.As(err, &pathErr) {
						if n.fullPath {
							symlinksRel, err := filepath.Rel(info.FullPath, pathErr.Path)
							if err == nil {
								pathErr.Path = symlinksRel
							}
						}
						symlinks = pathErr.Path
					} else {
						symlinks = err.Error()
					}
				} else {
					if n.fullPath {
						symlinksRel, err := filepath.Rel(filepath.Dir(info.FullPath), symlinks)
						if err == nil {
							symlinks = symlinksRel
						}
					}
				}
				var style theme.Style
				if broken {
					style = renderer.SymlinkBroken()
				} else {
					style = renderer.SymlinkDereference()
				}
				_, _ = dereference.WriteString(style.Color)
				if style.Underline {
					_, _ = dereference.WriteString(theme.Underline)
				}
				if style.Bold {
					_, _ = dereference.WriteString(theme.Bold)
				}
				_, _ = dereference.WriteString(style.Icon)
				_, _ = dereference.WriteString(symlinks)
				_, _ = dereference.WriteString(renderer.Colorend())
			}
		} else {
			if n.statistics != nil {
				n.statistics.file.Add(1)
			}
			if mode&os.ModeNamedPipe != 0 {
				style := renderer.Pipe()
				if n.icon {
					icon = style.Icon
				}
				if n.classify {
					classify = "|"
				}
				color, underline, bold = style.Color, style.Underline, style.Bold
			} else if mode&os.ModeSocket != 0 {
				style := renderer.Socket()
				if n.icon {
					icon = style.Icon
				}
				if n.classify {
					classify = "="
				}
				color, underline, bold = style.Color, style.Underline, style.Bold
			} else {
				if s, ok := renderer.ByName(name); ok {
					if n.icon {
						icon = s.Icon
					}
					color, underline, bold = s.Color, s.Underline, s.Bold
				} else {
					s, ok = renderer.ByExt(name)
					if ok {
						if n.icon {
							icon = s.Icon
						}
						color, underline, bold = s.Color, s.Underline, s.Bold
					} else {
						s = renderer.File()
						if n.icon {
							icon = s.Icon
						}
						color, underline, bold = s.Color, s.Underline, s.Bold
					}
				}
			}
		}

		exe := util.IsExecutableMode(mode) && !util.IsSymLinkMode(mode) && !info.IsDir() && mode&os.ModeNamedPipe == 0 && mode&os.ModeSocket == 0
		if n.classify {
			if classify == "" && (!n.fileType) && exe {
				classify = "*"
			}
		}
		if exe {
			s := renderer.Executable()
			color, underline, bold = s.Color, s.Underline, s.Bold
		}

		if n.mounts {
			mounts = util.MountsOn(info)
		}

		if n.relativeTo != "" {
			relativePath, err := filepath.Rel(n.relativeTo, info.FullPath)
			if err == nil {
				name = relativePath
			}
		} else if n.fullPath {
			name = info.FullPath
		}

		b := bytebufferpool.Get()
		defer bytebufferpool.Put(b)
		if color != "" {
			_, _ = b.WriteString(color)
		}
		if icon != "" {
			_, _ = b.WriteString(icon)
			_ = b.WriteByte(' ')
		}
		if underline {
			_, _ = b.WriteString(theme.Underline)
		}
		if bold {
			_, _ = b.WriteString(theme.Bold)
		}
		if n.Quote != "" {
			_, _ = b.WriteString(n.Quote)
		}
		if n.hyperLink {
			_, _ = b.WriteString(makeLink("file://"+info.FullPath, name))
		} else {
			_, _ = b.WriteString(name)
		}
		if n.Quote != "" {
			_, _ = b.WriteString(n.Quote)
		}
		if classify == "@" {
			_, _ = b.WriteString(classify)
		}
		if color != "" {
			_, _ = b.WriteString(renderer.Colorend())
		}
		if classify != "" && classify != "@" {
			_, _ = b.WriteString(classify)
		}
		d := dereference.String()
		if d != "" {
			_, _ = b.WriteString(d)
		}
		if mounts != "" {
			_ = b.WriteByte(' ')
			_, _ = b.WriteString(renderer.Mounts(mounts))
		}
		return b.String(), NameName
	}
}
