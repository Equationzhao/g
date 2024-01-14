package render

import (
	"fmt"
	"math"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/Equationzhao/g/internal/git"
	"github.com/Equationzhao/g/internal/theme"
	"github.com/hako/durafmt"
	"github.com/valyala/bytebufferpool"
)

type Renderer struct {
	theme *theme.All
}

func (rd *Renderer) SetTheme(theme *theme.All) *Renderer {
	rd.theme = theme
	return rd
}

func NewRenderer(a *theme.All) *Renderer {
	return &Renderer{theme: a}
}

func (rd *Renderer) OctalPerm(octal string) string {
	s := rd.theme.Permission["octal"]
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString(s.Color)
	checkStyle(&s, bb)
	_, _ = bb.WriteString(octal)
	_, _ = bb.WriteString(rd.Colorend())
	return bb.String()
}

// FileMode
// -     Regular file.
// b     Block special file.
// c     Character special file.
// d     Directory.
// l     Symbolic link.
// p     FIFO.
// s     Socket.
// w     Whiteout.
/*
go std fs.FileMode
	ModeDir        FileMode = 1 << (32 - 1 - iota) // d: is a directory
	ModeAppend                                     // a: append-only
	ModeExclusive                                  // l: exclusive use
	ModeTemporary                                  // T: temporary file; Plan 9 only
	ModeSymlink                                    // L: symbolic link
	ModeDevice                                     // D: device file
	ModeNamedPipe                                  // p: named pipe (FIFO)
	ModeSocket                                     // S: Unix domain socket
	ModeSetuid                                     // u: setuid
	ModeSetgid                                     // g: setgid
	ModeCharDevice                                 // c: Unix character device, when ModeDevice is set
	ModeSticky                                     // t: sticky
	ModeIrregular                                  // ?: non-regular file; nothing else is known about this file
*/
func (rd *Renderer) FileMode(toRender string) string {
	// return file mode like -rwxrwxrwx/drwxrwxrwx but in color
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	toRenderBytes := []byte(toRender)
	prefixBytes := toRenderBytes[:len(toRenderBytes)-9]
	suffixBytes := toRenderBytes[len(toRenderBytes)-9:]
	var firstByte byte = '-'
	firstStyle := rd.theme.Permission["-"]
	for _, c := range prefixBytes {
		switch c {
		case '-', 't', 'T', 'a', 'l', '?':
			continue
		case 'u': // setuid
			if suffixBytes[2] == 'x' {
				suffixBytes[2] = 'u'
				c = '-'
				firstStyle = rd.theme.Permission["-"]
			}
			firstByte = '-'
		case 'g': // setgid
			if suffixBytes[5] == 'x' {
				suffixBytes[5] = 'g'
				firstStyle = rd.theme.Permission["-"]
			}
			firstByte = '-'
		case 'L': // symlink
			firstStyle = rd.theme.Permission["link"]
			firstByte = 'L'
		case 'd': // directory
			firstStyle = rd.theme.Permission["directory"]
			firstByte = 'd'
		case 'c': // char device
			firstStyle = rd.theme.Permission["char"]
			firstByte = 'c'
		case 'S': // socket
			firstStyle = rd.theme.Permission["socket"]
			firstByte = 'S'
		case 'D': // block device
			firstStyle = rd.theme.Permission["block"]
			firstByte = 'D'
		case 'p': // FIFO
			firstStyle = rd.theme.Permission["pipe"]
			firstByte = 'p'
		}
	}
	_, _ = bb.WriteString(firstStyle.Color)
	checkStyle(&firstStyle, bb)
	_, _ = bb.WriteString(string(trans(firstByte)))

	for _, c := range suffixBytes {
		var s theme.Style
		switch c {
		case '-':
			s = rd.theme.Permission["-"]
		case 'r': // readable
			s = rd.theme.Permission["read"]
		case 'w': // writable
			s = rd.theme.Permission["write"]
		case 'x': // executable
			s = rd.theme.Permission["exe"]
		case 'g': // setgid
			s = rd.theme.Permission["setgid"]
		case 'u': // setuid
			s = rd.theme.Permission["setuid"]
		}
		_, _ = bb.WriteString(s.Color)
		checkStyle(&s, bb)
		_, _ = bb.WriteString(string(trans(c)))
	}

	_, _ = bb.WriteString(rd.Colorend())
	return bb.String()
}

func trans(s byte) byte {
	switch s {
	case 'D':
		return 'b'
	case 'L':
		return 'l'
	case 'S':
		return 's'
	case 'T':
		return 't'
	case 'u', 'g':
		return 's'
	}
	return s
}

func (rd *Renderer) Size(toRender, unit string) string {
	s := rd.theme.Size[unit]
	if strings.HasSuffix(toRender, "-") {
		s = rd.theme.Size["-"]
	}
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString(s.Color)
	checkStyle(&s, bb)
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.Colorend())
	return bb.String()
}

func (rd *Renderer) BlockSize(toRender string) string {
	return rd.Size(toRender, "block")
}

func (rd *Renderer) Link(toRender string) string {
	s := rd.theme.Symlink["link-num"]
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString(s.Color)
	checkStyle(&s, bb)
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.Colorend())
	return bb.String()
}

const adminSidPattern = `^S-1-5-(?:\d+-)*\d+-500$`

var rootSid = regexp.MustCompile(adminSidPattern)

func (rd *Renderer) Owner(toRender string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	var root []string
	var byName []string
	switch runtime.GOOS {
	case "windows":
		root = []string{"Administrators", "SYSTEM", "TrustedInstaller", "S-1-5-32-544", "S-1-5-18"}
		byName = []string{"DevToolsUser"}
	case "darwin":
		root = []string{"root", "0"}
	default:
		root = []string{"root", "0"}
	}

	style := rd.theme.User["owner"]

	if slices.Contains(root, toRender) {
		style = rd.theme.User["root"]
	} else {
		if slices.Contains(byName, toRender) {
			style = rd.theme.User[toRender]
		} else if runtime.GOOS == "windows" && rootSid.MatchString(toRender) {
			style = rd.theme.User["root"]
		}
	}
	_, _ = bb.WriteString(style.Color)
	checkStyle(&style, bb)
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.Colorend())
	return bb.String()
}

func (rd *Renderer) Group(toRender string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	var root []string
	var byName []string
	switch runtime.GOOS {
	case "windows":
		root = []string{"Administrators", "SYSTEM", "S-1-5-32-544", "S-1-5-18"}
		byName = []string{"DevToolsUser"}
	case "darwin":
		root = []string{"wheel", "admin", "0"}
	default:
		root = []string{"root", "0"}
	}

	style := rd.theme.Group["group"]
	if slices.Contains(root, toRender) {
		style = rd.theme.Group["root"]
	} else {
		if slices.Contains(byName, toRender) {
			style = rd.theme.Group[toRender]
		} else if runtime.GOOS == "windows" && rootSid.MatchString(toRender) {
			style = rd.theme.Group["root"]
		} else {
			style = rd.theme.Group["group"]
		}
	}
	_, _ = bb.WriteString(style.Color)
	checkStyle(&style, bb)
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.Colorend())
	return bb.String()
}

func (rd *Renderer) Time(toRender string) string {
	return rd.infoByName(toRender, "time")
}

func (rd *Renderer) calculateRTimeColor(dura time.Duration) string {
	dura = dura.Abs()

	const (
		day  = time.Hour * 24
		week = day * 7
	)

	const gUint = 35
	var r, b float64 = 165, 255

	switch theme.ColorLevel {
	case theme.TrueColor:
		// calculate the radio.
		// radio must < 1
		// radio = e^(-dura)
		radio := math.Exp(-dura.Seconds() / (10 * week.Seconds()))
		r *= radio
		b *= radio
		rUint, bUint := uint8(math.Round(r)), uint8(math.Round(b))
		rgb, _ := theme.RGB(rUint, gUint, bUint)
		return rgb
	case theme.C256:
		code := 0
		if dura <= time.Hour*6 {
			code = 201
		} else if dura <= day {
			code = 165
		} else if dura <= day*3 {
			code = 129
		} else if dura <= week {
			code = 93
		} else if dura <= week*6 {
			code = 57
		} else if dura <= week*52 {
			code = 56
		} else {
			code = 55
		}
		res, _ := theme.Color256(code)
		return res
	case theme.Ascii:
		return rd.theme.InfoTheme["time"].Color
	default:
		return ""
	}
}

func (rd *Renderer) RTime(now, modTime time.Time) string {
	t := now.Sub(modTime)
	var dura *durafmt.Durafmt
	if t > 0 {
		dura = durafmt.Parse(t)
		return fmt.Sprintf(
			"%s%s ago%s", rd.calculateRTimeColor(t), dura.LimitFirstN(1).String(), rd.Colorend(),
		)
	} else if t == 0 {
		return "now"
	} else {
		dura = durafmt.Parse(-t)
		return fmt.Sprintf(
			"%sin %s%s", rd.calculateRTimeColor(t), dura.LimitFirstN(1).String(), rd.Colorend(),
		)
	}
}

func (rd *Renderer) Charset(toRender string) string {
	return rd.infoByName(toRender, "charset")
}

func (rd *Renderer) Mime(toRender string) string {
	return rd.infoByName(toRender, "mime")
}

func (rd *Renderer) Checksum(toRender string) string {
	return rd.infoByName(toRender, "checksum")
}

func (rd *Renderer) ByName(toRender string) (s theme.Style, found bool) {
	name := strings.ToLower(filepath.Base(toRender))
	style, ok := rd.theme.Name[name]
	if !ok {
		return theme.Style{}, false
	}
	return style, true
}

func (rd *Renderer) infoByName(toRender string, name string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	style := rd.theme.InfoTheme[name]
	_, _ = bb.WriteString(style.Color)
	checkStyle(&style, bb)
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.Colorend())
	return bb.String()
}

// ByExt returns the colorized string by the file extension
// if the file has no extension it returns an empty string
func (rd *Renderer) ByExt(toRender string) (s theme.Style, found bool) {
	// get ext
	ext := filepath.Ext(toRender)
	if len(ext) > 0 {
		ext = ext[1:]
	} else {
		return theme.Style{}, false
	}
	ext = strings.ToLower(ext)
	style, ok := rd.theme.Ext[ext]
	if !ok {
		return theme.Style{}, false
	}
	return style, true
}

func (rd *Renderer) SymlinkArrow() theme.Style {
	return rd.theme.Symlink["symlink_arrow"]
}

func (rd *Renderer) Symlink() theme.Style {
	return rd.theme.Symlink["symlink"]
}

func (rd *Renderer) SymlinkDereference() theme.Style {
	return rd.theme.Symlink["symlink_path"]
}

func (rd *Renderer) SymlinkBroken() theme.Style {
	return rd.theme.Symlink["symlink_broken_path"]
}

func (rd *Renderer) Pipe() theme.Style {
	return rd.theme.Special["pipe"]
}

func (rd *Renderer) Socket() theme.Style {
	return rd.theme.Special["socket"]
}

func (rd *Renderer) Device() theme.Style {
	return rd.theme.Special["device"]
}

func (rd *Renderer) Char() theme.Style {
	return rd.theme.Special["char"]
}

func (rd *Renderer) Executable() theme.Style {
	return rd.theme.Special["exe"]
}

func (rd *Renderer) Dir(name string, empty bool) theme.Style {
	var style theme.Style

	if empty {
		style = rd.theme.Special["empty-dir"]
	} else {
		style = rd.theme.Special["dir"]
	}

	if s, ok := rd.theme.Name[strings.ToLower(name)]; ok {
		// keep color
		style.Icon = s.Icon
		style.Underline = s.Underline
		style.Bold = s.Bold
		style.Italics = s.Italics
		style.Faint = s.Faint
	}
	return style
}

func (rd *Renderer) File() theme.Style {
	return rd.theme.Special["file"]
}

func (rd *Renderer) HiddenFile() theme.Style {
	return rd.theme.Special["hidden-file"]
}

func (rd *Renderer) gitByStatus(name string, status string) string {
	style, ok := rd.theme.Git[status]
	if !ok {
		panic("no such git status:" + status)
	}
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString(style.Color)
	checkStyle(&style, bb)
	_, _ = bb.WriteString(name)
	_, _ = bb.WriteString(rd.Colorend())
	return bb.String()
}

func (rd *Renderer) GitUnmodified(name string) string {
	return rd.gitByStatus(name, "git_unmodified")
}

func (rd *Renderer) GitModified(name string) string {
	return rd.gitByStatus(name, "git_modified")
}

func (rd *Renderer) GitUntracked(name string) string {
	return rd.gitByStatus(name, "git_untracked")
}

func (rd *Renderer) GitAdded(name string) string {
	return rd.gitByStatus(name, "git_added")
}

func (rd *Renderer) GitRenamed(name string) string {
	return rd.gitByStatus(name, "git_renamed")
}

func (rd *Renderer) GitDeleted(name string) string {
	return rd.gitByStatus(name, "git_deleted")
}

func (rd *Renderer) GitIgnored(name string) string {
	return rd.gitByStatus(name, "git_ignored")
}

func (rd *Renderer) GitCopied(name string) string {
	return rd.gitByStatus(name, "git_copied")
}

func (rd *Renderer) GitTypeChanged(s string) string {
	return rd.gitByStatus(s, "git_type_changed")
}

func (rd *Renderer) GitUpdatedButUnmerged(s string) string {
	return rd.gitByStatus(s, "git_updated_but_unmerged")
}

func (rd *Renderer) GitRepoBranch(branch string) string {
	var style theme.Style
	switch branch {
	case "master", "main":
		style = rd.theme.Git["git-branch-master"]
	case "":
		style = rd.theme.Git["git-branch-none"]
		branch = style.Icon // if branch is empty, set it to the icon of none
	default:
		style = rd.theme.Git["git-branch"]
	}
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)

	_, _ = bb.WriteString(style.Color)
	checkStyle(&style, bb)
	_, _ = bb.WriteString(branch)
	_, _ = bb.WriteString(rd.Colorend())
	return bb.String()
}

func (rd *Renderer) GitRepoStatus(status git.RepoStatus) string {
	var style theme.Style
	switch status {
	case git.RepoStatusClean:
		style = rd.theme.Git["git-repo-clean"]
	case git.RepoStatusDirty:
		style = rd.theme.Git["git-repo-dirty"]
	case git.RepoStatusSkip:
		style = rd.theme.Git["git-repo-skip"]
	}
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	bb.WriteString(style.Color)
	checkStyle(&style, bb)
	_, _ = bb.WriteString(style.Icon)
	_, _ = bb.WriteString(rd.Colorend())
	return bb.String()
}

func (rd *Renderer) Inode(inode string) string {
	return rd.infoByName(inode, "inode")
}

func (rd *Renderer) DirPrompt(dir string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	style := rd.theme.Special["dir-prompt"]
	_, _ = bb.WriteString(style.Color)
	checkStyle(&style, bb)
	_, _ = bb.WriteString(style.Icon)
	_, _ = bb.WriteString(dir)
	_, _ = bb.WriteString(rd.Colorend())
	return bb.String()
}

func (rd *Renderer) Mounts(mounts string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	style := rd.theme.Special["mounts"]
	_, _ = bb.WriteString(style.Color)
	checkStyle(&style, bb)
	_, _ = bb.WriteString(style.Icon)
	_, _ = bb.WriteString(mounts)
	_, _ = bb.WriteString(rd.Colorend())
	return bb.String()
}

func (rd *Renderer) Colorend() string {
	return rd.theme.InfoTheme["reset"].Color
}

func checkStyle(style *theme.Style, bb *bytebufferpool.ByteBuffer) {
	if style.Underline {
		_, _ = bb.WriteString(theme.Underline)
	}
	if style.Bold {
		_, _ = bb.WriteString(theme.Bold)
	}
	if style.Italics {
		_, _ = bb.WriteString(theme.Italics)
	}
	if style.Faint {
		_, _ = bb.WriteString(theme.Faint)
	}
	if style.Blink {
		_, _ = bb.WriteString(theme.Blink)
	}
}
