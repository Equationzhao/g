package content

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Equationzhao/g/cached"
	"github.com/Equationzhao/g/filter"
	"github.com/Equationzhao/g/git"
	"github.com/Equationzhao/g/render"
	"github.com/valyala/bytebufferpool"
)

type (
	Name struct {
		Icon, Classify, FileType, git, fullPath bool
		Renderer                                *render.Renderer
		GitCache                                *cached.Map[git.GitRepoPath, *git.FileGits]
		statistics                              *Statistics
		parent                                  string
		Quote                                   string
		GitStyle                                gitStyle
	}
)

type Statistics struct {
	file, dir, link uint64
}

func (s *Statistics) Reset() {
	s.file = 0
	s.dir = 0
	s.link = 0
}

func (s *Statistics) String() string {
	return fmt.Sprintf("%d file(s), %d dir(s), %d link(s)", s.file, s.dir, s.link)
}

type gitStyle int

const (
	GitStyleDot gitStyle = iota
	GitStyleSym
	GitStyleDefault = GitStyleDot
)

func (n *Name) FullPath() bool {
	return n.fullPath
}

func (n *Name) SetFullPath() {
	n.fullPath = true
}

func (n *Name) UnsetFullPath() {
	n.fullPath = false
}

func (n *Name) Statistics() *Statistics {
	return n.statistics
}

func (n *Name) SetStatistics(Statistics *Statistics) {
	n.statistics = Statistics
}

func (n *Name) SetQuote(quote string) *Name {
	n.Quote = quote
	return n
}

func (n *Name) UnsetQuote() *Name {
	n.Quote = ""
	return n
}

func (n *Name) UnsetGit() *Name {
	n.git = false
	git.FreeCache()
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

func (n *Name) SetGit() *Name {
	n.git = true
	n.GitCache = git.GetCache()
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

func (n *Name) SetParent(parent string) *Name {
	n.parent = parent
	return n
}

// SetFileType set file type, should set Classify first
// if Classify is false, FileType will be ignored
func (n *Name) SetFileType() *Name {
	n.FileType = true
	return n
}

func (n *Name) SetRenderer(renderer *render.Renderer) *Name {
	n.Renderer = renderer
	return n
}

func NewNameEnable() *Name {
	return &Name{}
}

const NameName = "name"

func (n *Name) Enable() filter.ContentOption {
	/*
		 -F      Display a slash (`/`) immediately after each pathname that is a
				 directory, an asterisk (`*`) after each that is executable, an at
				 sign (`@`) after each symbolic link, a percent sign (`%`) after
				 each whiteout, an equal sign (`=`) after each socket, and a
				 vertical bar (`|`) after each that is a FIFO.
	*/

	isOrIsParentOf := func(parent, child string) bool {
		if parent == child {
			return true
		}
		if strings.HasPrefix(child, parent+"/") { // should not use filepath.Separator
			return true
		}
		return false
	}

	getFromCache := func(repoPath git.GitRepoPath) *git.FileGits {
		value, _ := n.GitCache.GetOrInit(repoPath, git.DefaultInit(repoPath))
		return value
	}

	return func(info os.FileInfo) (string, string) {
		buffer := bytebufferpool.Get()
		defer bytebufferpool.Put(buffer)
		name := info.Name()
		str := name
		mode := info.Mode()

		char := ""

		if n.Icon {
			if info.IsDir() {
				if n.statistics != nil {
					n.statistics.dir++
				}
				str = n.Renderer.DirIcon(str)
				char = "/"
			} else if mode&os.ModeSymlink != 0 {
				if n.statistics != nil {
					n.statistics.link++
				}
				if n.Classify {
					str = n.Renderer.SymlinkIconPlus(str, n.parent, "@")
				} else {
					str = n.Renderer.SymlinkIcon(str, n.parent)
				}
			} else {
				if n.statistics != nil {
					n.statistics.file++
				}
				if mode&os.ModeNamedPipe != 0 {
					str = n.Renderer.PipeIcon(str)
					char = "|"
				} else if mode&os.ModeSocket != 0 {
					str = n.Renderer.SocketIcon(str)
					char = "="
				} else {
					str = n.Renderer.ByExtIcon(str)
				}
			}
		} else {
			if info.IsDir() {
				if n.statistics != nil {
					n.statistics.dir++
				}
				str = n.Renderer.Dir(str)
				char = "/"
			} else if mode&os.ModeSymlink != 0 {
				if n.statistics != nil {
					n.statistics.link++
				}
				if n.Classify {
					str = n.Renderer.SymlinkPlus(str, n.parent, "@")
				} else {
					str = n.Renderer.Symlink(str, n.parent)
				}
			} else {
				if n.statistics != nil {
					n.statistics.file++
				}
				if mode&os.ModeNamedPipe != 0 {
					str = n.Renderer.Pipe(str)
					char = "|"
				} else if mode&os.ModeSocket != 0 {
					str = n.Renderer.Socket(str)
					char = "="
				} else {
					str = n.Renderer.ByExt(str)
				}
			}
		}

		if n.git {
			FilesStatus := *getFromCache(n.parent)
			for _, status := range FilesStatus {
				if isOrIsParentOf(name, status.Name) {
					str = n.GitByName(str, status.Status.String(), n.GitStyle)
					break
				}
			}
		}

		if n.Classify {
			if (!n.FileType) && (mode&0o111 != 0) {
				str += "*"
			} else {
				str += char
			}
		}

		if n.Quote != "" {
			str = strings.Replace(str, name, n.Quote+name+n.Quote, 1)
		}

		if n.fullPath {
			fullPath := filepath.Join(n.parent, name)
			str = strings.Replace(str, name, fullPath, 1)
		}

		return str, NameName
	}
}

func (n *Name) GitByName(name string, status string, style gitStyle) string {
	switch status {
	case "~":
		switch style {
		case GitStyleDot:
			return n.Renderer.GitModified(name)
		case GitStyleSym:
			return n.Renderer.GitModifiedSym(name)
		}
	case "?":
		switch style {
		case GitStyleDot:
			return n.Renderer.GitUntracked(name)
		case GitStyleSym:
			return n.Renderer.GitUntrackedSym(name)
		}
	case "+":
		switch style {
		case GitStyleDot:
			return n.Renderer.GitAdded(name)
		case GitStyleSym:
			return n.Renderer.GitAddedSym(name)
		}
	case "|":
		switch style {
		case GitStyleDot:
			return n.Renderer.GitRenamed(name)
		case GitStyleSym:
			return n.Renderer.GitRenamedSym(name)
		}
	case "-":
		switch style {
		case GitStyleDot:
			return n.Renderer.GitDeleted(name)
		case GitStyleSym:
			return n.Renderer.GitDeletedSym(name)
		}
	case "=":
		switch style {
		case GitStyleDot:
			return n.Renderer.GitCopied(name)
		case GitStyleSym:
			return n.Renderer.GitCopiedSym(name)
		}
	case "!":
		switch style {
		case GitStyleDot:
			return n.Renderer.GitIgnored(name)
		case GitStyleSym:
			return n.Renderer.GitIgnoredSym(name)
		}
	}
	return ""
}
