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
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Equationzhao/g/cached"
	"github.com/Equationzhao/g/git"
	"github.com/Equationzhao/g/osbased"
	"github.com/Equationzhao/g/render"
	"github.com/Equationzhao/tsmap"
	"github.com/hako/durafmt"
	"github.com/valyala/bytebufferpool"
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

type ContentOption func(info os.FileInfo) string

// EnableFileMode return file mode like -rwxrwxrwx/drwxrwxrwx
func EnableFileMode(renderer *render.Renderer) ContentOption {
	return func(info os.FileInfo) string {
		return renderer.FileMode(fillBlank(info.Mode().String(), 12))
	}
}

type InodeEnabler struct {
	*sync.WaitGroup
}

func NewInodeEnabler() *InodeEnabler {
	return &InodeEnabler{
		WaitGroup: new(sync.WaitGroup),
	}
}

func (i *InodeEnabler) Enable(renderer *render.Renderer) ContentOption {
	m := sync.RWMutex{}
	longestInode := 0

	wait := func(res string) string {
		i.Wait()
		return renderer.Inode(fillBlank(res, longestInode))
	}

	done := func(name string) {
		defer i.Done()
		m.RLock()
		if len(name) > longestInode {
			m.RUnlock()
			m.Lock()
			if len(name) > longestInode {
				longestInode = len(name)
			}
			m.Unlock()
		} else {
			m.RUnlock()
		}
	}

	return func(info os.FileInfo) string {
		str := osbased.Inode(info)
		done(str)
		return wait(str)
	}
}

type SizeUnit float64

const Unknown SizeUnit = -1
const (
	Auto SizeUnit = iota
	Bit  SizeUnit = 1.0 << (10 * iota)
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

func Convert2SizeString(size SizeUnit) string {
	switch size {
	case Unknown:
		return "?"
	case Auto:
		return ""
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
		return "unknown"
	}
}

func ConvertFromSizeString(size string) SizeUnit {
	switch size {
	case "bit", "Bit", "BIT":
		return Bit
	case "B", "b":
		return B
	case "KB", "kb", "Kb":
		return KB
	case "MB", "mb", "Mb":
		return MB
	case "GB", "gb", "Gb":
		return GB
	case "TB", "tb", "Tb":
		return TB
	case "PB", "pb", "Pb":
		return PB
	case "EB", "eb", "Eb":
		return EB
	case "ZB", "zb", "Zb":
		return ZB
	case "YB", "yb", "Yb":
		return YB
	case "BB", "bb", "Bb":
		return BB
	case "NB", "nb", "Nb":
		return NB
	default:
		return Unknown
	}
}

type SizeEnabler struct {
	total       atomic.Int64
	enableTotal bool
	sizeUint    SizeUnit
	renderer    *render.Renderer
	*sync.WaitGroup
}

func (s *SizeEnabler) SetRenderer(renderer *render.Renderer) {
	s.renderer = renderer
}

func NewSizeEnabler() *SizeEnabler {
	return &SizeEnabler{
		total:       atomic.Int64{},
		enableTotal: false,
		sizeUint:    Auto,
		renderer:    nil,
		WaitGroup:   new(sync.WaitGroup),
	}
}

func (s *SizeEnabler) SizeUint() SizeUnit {
	return s.sizeUint
}

func (s *SizeEnabler) SetEnableTotal() {
	s.enableTotal = true
}

func (s *SizeEnabler) DisableTotal() {
	s.enableTotal = false
}

func (s *SizeEnabler) Total() (size int64, ok bool) {
	if s.enableTotal {
		return s.total.Load(), s.enableTotal
	}
	return 0, false
}

func (s *SizeEnabler) Reset() {
	if s.enableTotal {
		s.total.Store(0)
	}
}

func (s *SizeEnabler) Size2String(b int64, blank int) string {
	var res string
	v := float64(b)
	switch s.sizeUint {
	case Bit:
		res = strconv.FormatInt(int64(v*8), 10)
	case B:
		res = strconv.FormatInt(int64(v), 10)
	case KB:
		fallthrough
	case MB:
		fallthrough
	case GB:
		fallthrough
	case TB:
		fallthrough
	case PB:
		fallthrough
	case EB:
		fallthrough
	case ZB:
		fallthrough
	case YB:
		fallthrough
	case BB:
		fallthrough
	case NB:
		res = fmt.Sprintf("%g", v*float64(B)/float64(s.sizeUint))

	case Auto:
		for i := B; i <= ZB; i *= 1024 {
			if v < 1024 {
				res = strconv.FormatFloat(v, 'f', 1, 64)
				if res == "0.0" {
					res = "-"
				} else {
					res += Convert2SizeString(i)
				}
				return s.renderer.Size(fillBlank(res, blank))
			}
			v /= 1024
		}
		panic("too large")
	default:
		panic("invalid size uint" + strconv.Itoa(int(s.sizeUint)))
	}

	if res == "0" {
		res = "-"
	} else {
		res += Convert2SizeString(s.sizeUint)
	}
	return s.renderer.Size(fillBlank(res, blank))
}

func (s *SizeEnabler) EnableSize(size SizeUnit) ContentOption {
	s.sizeUint = size

	if size != Auto {
		longestSize := 0
		m := sync.RWMutex{}
		done := func(size string) {
			defer s.Done()
			m.RLock()
			if longestSize >= len(size) {
				m.RUnlock()
				return
			}
			m.RUnlock()
			m.Lock()
			if longestSize < len(size) {
				longestSize = len(size)
			}
			m.Unlock()
		}

		wait := func(size string) string {
			s.Wait()
			return fillBlank(size, longestSize)
		}

		return func(info os.FileInfo) string {
			v := info.Size()
			if s.enableTotal {
				s.total.Add(v)
			}
			size := s.Size2String(v, 0)
			done(size)
			return wait(size)
		}
	}

	return func(info os.FileInfo) string {
		v := info.Size()
		if s.enableTotal {
			s.total.Add(v)
		}
		return s.Size2String(v, 7)
	}
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
		return fillBlank(size, longestRt)
	}

	return func(info os.FileInfo) string {
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
		return wait(rt)
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

func EnableTime(format string, renderer *render.Renderer, mod string) ContentOption {
	return func(info os.FileInfo) string {
		// get mod time/ create time/ access time
		var t time.Time
		switch mod {
		case "mod":
			t = osbased.ModTime(info)
		case "create":
			t = osbased.CreateTime(info)
		case "access":
			t = osbased.AccessTime(info)
		}
		return renderer.Time(t.Format(format))
	}
}

type (
	Name struct {
		Icon, Classify, FileType, git bool
		Renderer                      *render.Renderer
		parent                        string
		GitCache                      *cached.Map[git.GitRepoPath, *git.FileGits]
		GitStyle                      gitStyle
	}
)

type gitStyle int

const (
	GitStyleDot gitStyle = iota
	GitStyleSym
	GitStyleDefault = GitStyleDot
)

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
				if isOrIsParentOf(name, status.Name) {
					str = n.GitByName(str, status.Status.String(), n.GitStyle)
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

var (
	Uid = false
	Gid = false
)

func (cf *ContentFilter) EnableOwner(renderer *render.Renderer) ContentOption {
	m := sync.RWMutex{}
	longestOwner := 0

	wg := new(sync.WaitGroup)
	cf.wgs = append(cf.wgs, wg)
	wait := func(res string) string {
		wg.Wait()
		return renderer.Owner(fillBlank(res, longestOwner))
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
	return func(info os.FileInfo) string {
		name := ""
		if Uid {
			name = osbased.OwnerID(info)
		} else {
			name = osbased.Owner(info)
		}
		done(name)
		return wait(name)
	}
}

func (cf *ContentFilter) EnableGroup(renderer *render.Renderer) ContentOption {
	m := sync.RWMutex{}
	longestGroup := 0

	wg := new(sync.WaitGroup)
	cf.wgs = append(cf.wgs, wg)
	wait := func(name string) string {
		wg.Wait()
		return renderer.Group(fillBlank(name, longestGroup))
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

	return func(info os.FileInfo) string {
		name := ""
		if Gid {
			name = osbased.GroupID(info)
		} else {
			name = osbased.Group(info)
		}
		done(name)
		return wait(name)
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

// GetExtraAndNameStringSlice returns a slice of string which contains extra information and name of the file
/*
	buffers layout:
		extra extra extra ... |half|  name name name ...
*/
func (cf *ContentFilter) GetExtraAndNameStringSlice(e ...os.FileInfo) []tsmap.Pair[string, string] {
	resBuffers := make([]*bytebufferpool.ByteBuffer, 2*len(e))

	for i := range resBuffers {
		resBuffers[i] = bytebufferpool.Get()
	}

	defer func() {
		for i := range resBuffers {
			bytebufferpool.Put(resBuffers[i])
		}
	}()

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
	for i, entry := range e {
		go func(entry os.FileInfo, i int) {
			options := cf.options[:len(cf.options)-1]
			for j := range options {
				_, _ = resBuffers[i].WriteString(options[j](entry))
				_ = resBuffers[i].WriteByte(' ')
			}
			// the last one should not follow by space
			_, _ = resBuffers[i+len(e)].WriteString(cf.options[len(cf.options)-1](entry))
			wg.Done()
		}(entry, i)
	}
	res := make([]tsmap.Pair[string, string], 0, len(e))
	wg.Wait()

	/*
		buffers layout:
			extra extra extra ... |half|  name name name ...
	*/
	bufLen := len(resBuffers)
	for i := 0; i < bufLen/2; i++ {
		res = append(res, tsmap.MakePair(resBuffers[i].String(), resBuffers[i+bufLen/2].String()))
	}

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

func (cf *ContentFilter) EnableSum(sumTypes ...SumType) ContentOption {
	length := 0
	for _, t := range sumTypes {
		switch t {
		case SumTypeMd5:
			length += 32
		case SumTypeSha1:
			length += 40
		case SumTypeSha224:
			length += 56
		case SumTypeSha256:
			length += 64
		case SumTypeSha384:
			length += 96
		case SumTypeSha512:
			length += 128
		case SumTypeCRC32:
			length += 8
		}
	}
	length += len(sumTypes) - 1

	return func(info os.FileInfo) string {
		if info.IsDir() {
			return fillBlank("", length)
		}

		file, err := os.Open(info.Name())
		if err != nil {
			return fillBlank("", length)
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
			return fillBlank("", length)
		}
		sums := make([]string, 0, len(hashes))
		for _, h := range hashes {
			sums = append(sums, fmt.Sprintf("%x", h.Sum(nil)))
		}
		sumsStr := strings.Join(sums, " ")
		return fillBlank(sumsStr, length)
	}
}
