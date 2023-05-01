package filter

import (
	"bufio"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/Equationzhao/g/cached"
	"github.com/Equationzhao/g/render"
	"github.com/valyala/bytebufferpool"
)

// fileMode size owner group time name

type ContentFilter struct {
	options          []ContentOption
	wgOwner, wgGroup *sync.WaitGroup
}

func (cf *ContentFilter) AppendTo(options ...ContentOption) {
	cf.options = append(cf.options, options...)
}

func (cf *ContentFilter) SetOptions(options ...ContentOption) {
	cf.options = options
}

type ContentOption func(info os.FileInfo) string

// EnableFileMode return file mode like -rwxrwxrwx/drwxrwxrwx
func EnableFileMode(renderer *render.Renderer) ContentOption {
	return func(info os.FileInfo) string {
		return renderer.FileMode(fillBlank(info.Mode().String(), 12))
	}
}

type SizeUnit int

const (
	Bit SizeUnit = iota
	B
	KB
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
	BB
	NB
	Auto
)

// fill blank
// if s is shorter than length, fill blank from left
// if s is longer than length, panic
func fillBlank(s string, length int) string {
	if len(s) > length {
		return s
	}
	return strings.Repeat(" ", length-len(s)) + s
}

func convert2Size(size SizeUnit) string {
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
	case YB:
		return "YB"
	case BB:
		return "BB"
	case NB:
		return "NB"
	default:
		panic("unknown size")
	}
}

type Size struct {
	total       atomic.Int64
	enableTotal bool
	sizeUint    SizeUnit
	renderer    *render.Renderer
}

func (s *Size) SizeUint() SizeUnit {
	return s.sizeUint
}

func (s *Size) SetEnableTotal() {
	s.enableTotal = true
}

func (s *Size) DisableTotal() {
	s.enableTotal = false
}

func (s *Size) Total() (size int64, ok bool) {
	if s.enableTotal {
		return s.total.Load(), s.enableTotal
	}
	return 0, false
}

func (s *Size) Reset() {
	if s.enableTotal {
		s.total.Store(0)
	}
}

func (s *Size) Size2String(n int64, blank int) string {
	var res string
	v := float64(n)
	switch s.sizeUint {
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
	case YB:
		res = strconv.FormatFloat(v/1024.0/1024.0/1024.0/1024.0/1024.0/1024.0/1024.0/1024.0, 'f', 1, 64)
	case BB:
		res = strconv.FormatFloat(v/1024.0/1024.0/1024.0/1024.0/1024.0/1024.0/1024.0/1024.0/1024.0, 'f', 1, 64)
	case NB:
		res = strconv.FormatFloat(v/1024.0/1024.0/1024.0/1024.0/1024.0/1024.0/1024.0/1024.0/1024.0/1024.0, 'f', 1, 64)
	case Auto:
		for i := B; i <= ZB; i++ {
			if v < 1000 {
				res = strconv.FormatFloat(v, 'f', 1, 64)
				if res == "0.0" {
					res = "-"
				} else {
					res += convert2Size(i)
				}
				return s.renderer.Size(fillBlank(res, blank))
			}
			v /= 1024
		}
		panic("too large")
	default:
		panic("invalid size uint" + strconv.Itoa(int(s.sizeUint)))
	}
	return s.renderer.Size(fillBlank(res, blank))
}

func (s *Size) EnableSize(size SizeUnit, renderer *render.Renderer) ContentOption {
	s.sizeUint = size
	s.renderer = renderer
	return func(info os.FileInfo) string {
		v := info.Size()
		if s.enableTotal {
			s.total.Add(v)
		}
		return s.Size2String(v, 7)
	}
}

func EnableTime(format string, renderer *render.Renderer) ContentOption {
	return func(info os.FileInfo) string {
		return renderer.Time(info.ModTime().Format(format))
	}
}

type (
	gitStyle    int
	fileGits    = []fileGit
	gitRepoPath = string
	Name        struct {
		Icon, Classify, FileType, git bool
		Renderer                      *render.Renderer
		parent                        string
		GitCache                      *cached.Map[gitRepoPath, *fileGits]
		GitStyle                      gitStyle
	}
)

const (
	GitStyleDot = iota
	GitStyleSym
	GitStyleDefault = GitStyleDot
)

func (n *Name) UnsetGit() *Name {
	n.git = false
	n.GitCache = nil
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
	n.GitCache = cached.NewCacheMap[gitRepoPath, *fileGits](20)
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

func (n *Name) SetRenderer(Renderer *render.Renderer) *Name {
	n.Renderer = Renderer
	return n
}

func NewNameEnable() *Name {
	return &Name{}
}

// getShortGitStatus read the git status of the repository located at path
func getShortGitStatus(repoPath gitRepoPath) (string, error) {
	out, err := exec.Command("git", "-C", repoPath, "status", "-s", "--ignored", "--porcelain").Output()
	return string(out), err
}

type status int

func (s status) String() string {
	switch s {
	case GitModified:
		return "~"
	case GitAdded:
		return "+"
	case GitDeleted:
		return "-"
	case GitRenamed:
		return "|"
	case GitCopied:
		return "="
	case GitUntracked:
		return "?"
	case GitIgnored:
		return "!"
	}
	return ""
}

const (
	GitModified  status = iota + 1 // M ~
	GitAdded                       // A +
	GitDeleted                     // D -
	GitRenamed                     // R |
	GitCopied                      // C =
	GitUntracked                   // ? ?
	GitIgnored                     // ! !
)

type fileGit struct {
	name   string
	status status
}

func (f *fileGit) setYFromXY(XY string) {
	set := func(Y string) {
		switch Y {
		case "M":
			f.status = GitModified
		case "A":
			f.status = GitAdded
		case "D":
			f.status = GitDeleted
		case "R":
			f.status = GitRenamed
		case "C":
			f.status = GitCopied
		case "?":
			f.status = GitUntracked
		case "!":
			f.status = GitIgnored
		}
	}

	switch len(XY) {
	case 1:
		set(XY)
	case 2:
		Y := XY[1:]
		set(Y)
	default:
		return
	}
}

// parseShort parses a git status output command
// It is compatible with the short version of the git status command
// modified from https://le-gall.bzh/post/go/parsing-git-status-with-go/ author: Sébastien Le Gall
func parseShort(r string) (res fileGits) {
	s := bufio.NewScanner(strings.NewReader(r))

	// Extract branch name
	for s.Scan() {
		// Skip any empty line
		if len(s.Text()) < 1 {
			continue
		}
		break
	}

	fg := fileGit{}
	for true {
		if len(s.Text()) < 1 {
			continue
		}
		XyName := strings.Fields(s.Text())
		fg.setYFromXY(XyName[0])
		fg.name = XyName[1]
		res = append(res, fg)
		if !s.Scan() {
			break
		}
	}

	return
}

func (n *Name) Enable() ContentOption {
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

	getFromCache := func(repoPath gitRepoPath) *fileGits {
		value, _ := n.GitCache.GetOrInit(repoPath, func() *fileGits {
			res := make(fileGits, 0)
			out, err := getShortGitStatus(repoPath)
			if err == nil {
				res = parseShort(out)
			}
			return &res
		})
		return value
	}

	return func(info os.FileInfo) string {
		buffer := bytebufferpool.Get()
		defer bytebufferpool.Put(buffer)
		name := info.Name()
		str := name
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

		if n.git {
			FilesStatus := *getFromCache(n.parent)
			for _, status := range FilesStatus {
				if isOrIsParentOf(name, status.name) {
					str = n.GitByName(str, status.status.String(), n.GitStyle)
					break
				}
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
			} else if (!n.FileType) && (mode&0o111 != 0) {
				str += "*"
			}
		}

	end:
		return str
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

func NewContentFilter(options ...ContentOption) *ContentFilter {
	return &ContentFilter{options: options, wgGroup: new(sync.WaitGroup), wgOwner: new(sync.WaitGroup)}
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
	cf.wgOwner.Add(len(e))
	cf.wgGroup.Add(len(e))
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
