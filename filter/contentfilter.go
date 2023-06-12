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
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Equationzhao/g/cached"
	"github.com/Equationzhao/g/display"
	"github.com/Equationzhao/g/git"
	"github.com/Equationzhao/g/osbased"
	"github.com/Equationzhao/g/render"
	"github.com/hako/durafmt"
	"github.com/valyala/bytebufferpool"

	mt "github.com/gabriel-vasile/mimetype"
)

type LengthFixed interface {
	Done()
	Wait()
	Add(delta int)
}

// fileMode size owner group time name

type ContentFilter struct {
	options  []ContentOption
	wgs      []LengthFixed
	sortFunc func(a, b os.FileInfo) bool
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

type ContentOption func(info os.FileInfo) (stringContent string, funcName string)

// FillBlank
// if s is shorter than length, fill blank from left
// if s is longer than length, panic
func FillBlank(s string, length int) string {
	if len(s) > length {
		return s
	}
	return strings.Repeat(" ", length-len(s)) + s
}

type RelativeTimeEnabler struct {
	*sync.WaitGroup
	Mode string
}

func NewRelativeTimeEnabler() *RelativeTimeEnabler {
	return &RelativeTimeEnabler{
		WaitGroup: new(sync.WaitGroup),
	}
}

const RelativeTime = "Relative-time"

func (r *RelativeTimeEnabler) Enable(renderer *render.Renderer) ContentOption {
	longestRt := 0
	m := sync.RWMutex{}
	done := func(rt string) {
		defer r.Done()
		m.RLock()
		if longestRt >= len(rt) {
			m.RUnlock()
			return
		}
		m.RUnlock()
		m.Lock()
		if longestRt < len(rt) {
			longestRt = len(rt)
		}
		m.Unlock()
	}

	wait := func(size string) string {
		r.Wait()
		return FillBlank(size, longestRt)
	}

	return func(info os.FileInfo) (string, string) {
		var t time.Time
		switch r.Mode {
		case "mod":
			t = osbased.ModTime(info)
		case "create":
			t = osbased.CreateTime(info)
		case "access":
			t = osbased.AccessTime(info)
		default:
			t = osbased.ModTime(info)
		}
		rt := renderer.Time(relativeTime(time.Now(), t))
		done(rt)
		return wait(rt), RelativeTime + " " + r.Mode
	}
}

func relativeTime(now, modTime time.Time) string {
	if t := now.Sub(modTime); t > 0 {
		return fmt.Sprintf("%s ago", durafmt.Parse(t).LimitFirstN(1).String())
	} else if t == 0 {
		return "now"
	} else {
		return fmt.Sprintf("in %s", durafmt.Parse(-t).LimitFirstN(1).String())
	}
}

const (
	timeName     = "Time"
	timeModified = "Modified"
	timeCreated  = "Created"
	timeAccessed = "Accessed"
)

func EnableTime(format string, mode string, renderer *render.Renderer) ContentOption {
	return func(info os.FileInfo) (string, string) {
		// get mod time/ create time/ access time
		var t time.Time
		timeType := ""
		switch mode {
		case "mod":
			t = osbased.ModTime(info)
			timeType = timeModified
		case "create":
			t = osbased.CreateTime(info)
			timeType = timeCreated
		case "access":
			t = osbased.AccessTime(info)
			timeType = timeAccessed
		}
		return renderer.Time(t.Format(format)), timeName + " " + timeType
	}
}

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

type (
	Name struct {
		Icon, Classify, FileType, git bool
		Renderer                      *render.Renderer
		GitCache                      *cached.Map[git.GitRepoPath, *git.FileGits]
		statistics                    *Statistics
		parent                        string
		GitStyle                      gitStyle
		Quote                         string
	}
)

type gitStyle int

const (
	GitStyleDot gitStyle = iota
	GitStyleSym
	GitStyleDefault = GitStyleDot
)

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

func NewContentFilter(options ...ContentOption) *ContentFilter {
	c := &ContentFilter{
		options:  options,
		sortFunc: nil,
		wgs:      make([]LengthFixed, 0),
	}

	return c
}

type ContentFunc func(entry os.FileInfo) bool

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

type MimeFileTypeEnabler struct {
	*sync.WaitGroup
	ParentOnly bool
}

func NewMimeFileTypeEnabler() *MimeFileTypeEnabler {
	return &MimeFileTypeEnabler{
		WaitGroup:  &sync.WaitGroup{},
		ParentOnly: false,
	}
}

const MimeTypeName = "mine type"

func (e *MimeFileTypeEnabler) Enable() ContentOption {
	longestTypeName := 0
	m := sync.RWMutex{}
	done := func(tn string) {
		defer e.Done()
		m.RLock()
		if longestTypeName >= len(tn) {
			m.RUnlock()
			return
		}
		m.RUnlock()
		m.Lock()
		if longestTypeName < len(tn) {
			longestTypeName = len(tn)
		}
		m.Unlock()
	}

	wait := func(tn string) string {
		e.Wait()
		return FillBlank(tn, longestTypeName)
	}
	return func(info os.FileInfo) (string, string) {
		tn := ""
		returnName := MimeTypeName
		if info.IsDir() {
			tn = "directory"
		} else if info.Mode()&os.ModeSymlink != 0 {
			tn = "symlink"
		} else if info.Mode()&os.ModeNamedPipe != 0 {
			tn = "named_pipe"
		} else if info.Mode()&os.ModeSocket != 0 {
			tn = "socket"
		} else {
			file, err := os.Open(info.Name())
			if err != nil {
				// tn = err.Error()
				tn = "failed_to_read"
				done(tn)
				return wait(tn), MimeTypeName
			}
			mtype, err := mt.DetectReader(file)
			if err != nil {
				tn = err.Error()
				done(tn)
				return wait(tn), MimeTypeName
			}
			tn = mtype.String()

			if e.ParentOnly {
				tn = strings.SplitN(tn, "/", 2)[0]
				returnName = "parent_" + returnName
			}

			if strings.Contains(tn, ";") {
				// remove charset
				tn = strings.SplitN(tn, ";", 2)[0]
			}

		}
		done(tn)
		return wait(tn), returnName
	}
}

type LinkEnabler struct {
	// List each file's number of hard links.
	*sync.WaitGroup
}

func NewLinkEnabler() *LinkEnabler {
	return &LinkEnabler{
		WaitGroup: &sync.WaitGroup{},
	}
}

func (l *LinkEnabler) Enable() ContentOption {
	var longestLinkNum string
	m := sync.RWMutex{}
	done := func(linkNumStr string) {
		defer l.Done()
		m.RLock()
		if len(longestLinkNum) >= len(linkNumStr) {
			m.RUnlock()
			return
		}
		m.RUnlock()
		m.Lock()
		if len(longestLinkNum) < len(linkNumStr) {
			longestLinkNum = linkNumStr
		}
		m.Unlock()
	}

	wait := func(linkNumStr string) string {
		l.Wait()
		return FillBlank(linkNumStr, len(longestLinkNum))
	}

	return func(info os.FileInfo) (string, string) {
		n := strconv.FormatUint(osbased.LinkCount(info), 10)
		done(n)
		return wait(n), "links"
	}
}

type IndexEnabler struct{}

func NewIndexEnabler() *IndexEnabler {
	return &IndexEnabler{}
}

func (i *IndexEnabler) Enable() ContentOption {
	return func(info os.FileInfo) (string, string) {
		return "", "#"
	}
}
