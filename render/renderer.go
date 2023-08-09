package render

import (
	"errors"
	"fmt"
	"io/fs"
	"math"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/Equationzhao/g/slices"
	"github.com/Equationzhao/g/theme"
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
	checkUnderlineAndBold(&s, bb)
	_, _ = bb.WriteString(octal)
	_, _ = bb.WriteString(rd.theme.InfoTheme["reset"].Color)
	return bb.String()
}

func (rd *Renderer) FileMode(toRender string) string {
	// return file mode like -rwxrwxrwx/drwxrwxrwx but in color
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	for _, c := range toRender {
		switch c {
		case '-':
			_, _ = bb.WriteString(rd.theme.Permission["-"].Color)
		case 'L':
			_, _ = bb.WriteString(rd.theme.Permission["l"].Color)
		case 'd':
			_, _ = bb.WriteString(rd.theme.Permission["d"].Color)
		case 'r':
			_, _ = bb.WriteString(rd.theme.Permission["r"].Color)
		case 'w':
			_, _ = bb.WriteString(rd.theme.Permission["w"].Color)
		case 'x', 's', 't':
			_, _ = bb.WriteString(rd.theme.Permission["x"].Color)
		case 'c':
			_, _ = bb.WriteString(rd.theme.Permission["c"].Color)
		case 'S', 'T':
			_, _ = bb.WriteString(rd.theme.Permission["s"].Color)
		case 'D':
			_, _ = bb.WriteString(rd.theme.Permission["D"].Color)
		}
		_, _ = bb.WriteString(string(c))
	}
	_, _ = bb.WriteString(rd.theme.InfoTheme["reset"].Color)
	return bb.String()
}

func (rd *Renderer) Size(toRender, unit string) string {
	s := rd.theme.Size[unit]
	if strings.HasSuffix(toRender, "-") {
		s = rd.theme.Size["-"]
	}
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString(s.Color)
	checkUnderlineAndBold(&s, bb)
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.theme.InfoTheme["reset"].Color)
	return bb.String()
}

func (rd *Renderer) BlockSize(toRender string) string {
	return rd.Size(toRender, "block")
}

func (rd *Renderer) Link(toRender string) string {
	s := rd.theme.Symlink["link"]
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString(s.Color)
	checkUnderlineAndBold(&s, bb)
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.theme.InfoTheme["reset"].Color)
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
	checkUnderlineAndBold(&style, bb)
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.theme.InfoTheme["reset"].Color)
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
	checkUnderlineAndBold(&style, bb)
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.theme.InfoTheme["reset"].Color)
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
			"%s%s ago%s", rd.calculateRTimeColor(t), dura.LimitFirstN(1).String(), rd.theme.InfoTheme["reset"].Color,
		)
	} else if t == 0 {
		return "now"
	} else {
		dura = durafmt.Parse(-t)
		return fmt.Sprintf(
			"%sin %s%s", rd.calculateRTimeColor(t), dura.LimitFirstN(1).String(), rd.theme.InfoTheme["reset"].Color,
		)
	}
}

func (rd *Renderer) ByName(toRender string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	name := strings.ToLower(filepath.Base(toRender))
	style, ok := rd.theme.Name[name]
	if !ok {
		return ""
	}
	_, _ = bb.WriteString(style.Color)
	checkUnderlineAndBold(&style, bb)
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.theme.InfoTheme["reset"].Color)
	return bb.String()
}

func (rd *Renderer) ByNameIcon(toRender string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	name := strings.ToLower(filepath.Base(toRender))
	style, ok := rd.theme.Name[name]
	if !ok {
		return ""
	}
	_, _ = bb.WriteString(style.Color)
	_, _ = bb.WriteString(style.Icon)
	_, _ = bb.WriteString(" ")
	checkUnderlineAndBold(&style, bb)
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.theme.InfoTheme["reset"].Color)
	return bb.String()
}

func (rd *Renderer) infoByName(toRender string, name string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	style := rd.theme.InfoTheme[name]
	_, _ = bb.WriteString(style.Color)
	checkUnderlineAndBold(&style, bb)
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.theme.InfoTheme["reset"].Color)
	return bb.String()
}

// ByExt returns the colorized string by the file extension
// if the file has no extension it returns an empty string
func (rd *Renderer) ByExt(toRender string) string {
	// get ext
	ext := filepath.Ext(toRender)
	if len(ext) > 0 {
		ext = ext[1:]
	} else {
		return ""
	}
	ext = strings.ToLower(ext)
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	style, ok := rd.theme.Ext[ext]
	if !ok {
		return ""
	}
	_, _ = bb.WriteString(style.Color)
	checkUnderlineAndBold(&style, bb)
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.theme.InfoTheme["reset"].Color)
	return bb.String()
}

// ByExtIcon returns the icon and the name of the file
// if the file has no icon it returns an empty string
func (rd *Renderer) ByExtIcon(toRender string) string {
	// get ext
	ext := filepath.Ext(toRender)
	if len(ext) > 0 {
		ext = ext[1:]
	} else {
		return ""
	}
	ext = strings.ToLower(ext)
	style, ok := rd.theme.Ext[ext]
	if !ok {
		return ""
	}
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString(style.Color)
	_, _ = bb.WriteString(style.Icon)
	_, _ = bb.WriteString(" ")
	checkUnderlineAndBold(&style, bb)
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.theme.InfoTheme["reset"].Color)
	return bb.String()
}

// SymlinkIconPlus returns the icon and the name of the file, and dereferences the symlink
func (rd *Renderer) SymlinkIconPlus(toRender string, path string, plus string, rel bool) string {
	style := rd.theme.Symlink["symlink"]
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString(style.Color)
	_, _ = bb.WriteString(style.Icon)
	_, _ = bb.WriteString(" ")
	checkUnderlineAndBold(&style, bb)
	symlinks, err := filepath.EvalSymlinks(path)
	if err != nil {
		var pathErr *fs.PathError
		_, _ = bb.WriteString(toRender + plus)
		_, _ = bb.WriteString(rd.theme.Symlink["symlink_arrow"].Color + theme.Symlink["symlink_arrow"].Icon)
		brokenStyle := rd.theme.Symlink["symlink_broken_path"]
		_, _ = bb.WriteString(brokenStyle.Color)
		checkUnderlineAndBold(&brokenStyle, bb)
		if errors.As(err, &pathErr) {
			if rel {
				symlinksRel, err := filepath.Rel(filepath.Dir(path), pathErr.Path)
				if err == nil {
					pathErr.Path = symlinksRel
				}
			}
			symlinks = pathErr.Path
		} else {
			symlinks = err.Error()
		}
		_, _ = bb.WriteString(symlinks)
		_, _ = bb.WriteString(rd.theme.InfoTheme["reset"].Color)
		return bb.String()
	}

	if rel {
		symlinksRel, err := filepath.Rel(filepath.Dir(path), symlinks)
		if err == nil {
			symlinks = symlinksRel
		}
	}

	_, _ = bb.WriteString(toRender + plus)
	_, _ = bb.WriteString(rd.theme.Symlink["symlink_arrow"].Color + rd.theme.Symlink["symlink_arrow"].Icon)
	pathStyle := rd.theme.Symlink["symlink_path"]
	_, _ = bb.WriteString(pathStyle.Color)
	checkUnderlineAndBold(&pathStyle, bb)
	_, _ = bb.WriteString(symlinks)
	_, _ = bb.WriteString(rd.theme.InfoTheme["reset"].Color)
	return bb.String()
}

// SymlinkIconNoDereferencePlus returns the icon and the name of the file, but does not dereference the symlink
func (rd *Renderer) SymlinkIconNoDereferencePlus(toRender string, plus string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	style := rd.theme.Symlink["symlink"]
	_, _ = bb.WriteString(style.Color)
	_, _ = bb.WriteString(style.Icon)
	_, _ = bb.WriteString(" ")
	checkUnderlineAndBold(&style, bb)
	_, _ = bb.WriteString(toRender + plus)
	_, _ = bb.WriteString(rd.theme.InfoTheme["reset"].Color)
	return bb.String()
}

func (rd *Renderer) SymlinkIconNoDereference(toRender string) string {
	return rd.SymlinkIconNoDereferencePlus(toRender, "")
}

func (rd *Renderer) SymlinkIcon(toRender string, path string, rel bool) string {
	return rd.SymlinkIconPlus(toRender, path, "", rel)
}

// SymlinkPlus returns the icon and the name of the file
func (rd *Renderer) SymlinkPlus(toRender string, path string, plus string, rel bool) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	style := rd.theme.Symlink["symlink"]
	_, _ = bb.WriteString(style.Color)
	checkUnderlineAndBold(&style, bb)
	symlinks, err := filepath.EvalSymlinks(path)
	if err != nil {
		var pathErr *fs.PathError
		_, _ = bb.WriteString(toRender + plus)
		_, _ = bb.WriteString(rd.theme.Symlink["symlink_arrow"].Color + theme.Symlink["symlink_arrow"].Icon)
		brokenStyle := rd.theme.Symlink["symlink_broken_path"]
		_, _ = bb.WriteString(brokenStyle.Color)
		checkUnderlineAndBold(&brokenStyle, bb)
		if errors.As(err, &pathErr) {
			if rel {
				symlinksRel, err := filepath.Rel(filepath.Dir(path), pathErr.Path)
				if err == nil {
					pathErr.Path = symlinksRel
				}
			}
			symlinks = pathErr.Path
		} else {
			symlinks = err.Error()
		}
		_, _ = bb.WriteString(symlinks)
		_, _ = bb.WriteString(rd.theme.InfoTheme["reset"].Color)
		return bb.String()
	}
	if rel {
		symlinksRel, err := filepath.Rel(filepath.Dir(path), symlinks)
		if err == nil {
			symlinks = symlinksRel
		}
	}
	_, _ = bb.WriteString(toRender + plus)
	_, _ = bb.WriteString(rd.theme.Symlink["symlink_arrow"].Color + rd.theme.Symlink["symlink_arrow"].Icon)
	pathStyle := rd.theme.Symlink["symlink_path"]
	_, _ = bb.WriteString(pathStyle.Color)
	checkUnderlineAndBold(&pathStyle, bb)
	_, _ = bb.WriteString(symlinks)
	_, _ = bb.WriteString(rd.theme.InfoTheme["reset"].Color)
	return bb.String()
}

func (rd *Renderer) SymlinkNoDereferencePlus(toRender string, plus string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString(rd.theme.Symlink["symlink"].Color)
	_, _ = bb.WriteString(toRender + plus)
	_, _ = bb.WriteString(rd.theme.InfoTheme["reset"].Color)
	return bb.String()
}

func (rd *Renderer) SymlinkNoDereference(str string) string {
	return rd.SymlinkNoDereferencePlus(str, "")
}

func (rd *Renderer) Symlink(toRender string, path string, rel bool) string {
	return rd.SymlinkPlus(toRender, path, "", rel)
}

func (rd *Renderer) PipeIcon(toRender string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	style := rd.theme.Special["pipe"]
	_, _ = bb.WriteString(style.Color)
	_, _ = bb.WriteString(style.Icon)
	_, _ = bb.WriteString(" ")
	checkUnderlineAndBold(&style, bb)
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.theme.InfoTheme["reset"].Color)
	return bb.String()
}

func (rd *Renderer) Pipe(toRender string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	style := rd.theme.Special["pipe"]
	_, _ = bb.WriteString(style.Color)
	checkUnderlineAndBold(&style, bb)
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.theme.InfoTheme["reset"].Color)
	return bb.String()
}

func (rd *Renderer) SocketIcon(toRender string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	style := rd.theme.Special["socket"]
	_, _ = bb.WriteString(style.Color)
	_, _ = bb.WriteString(style.Icon)
	_, _ = bb.WriteString(" ")
	checkUnderlineAndBold(&style, bb)
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.theme.InfoTheme["reset"].Color)
	return bb.String()
}

func (rd *Renderer) Socket(toRender string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	style := rd.theme.Special["socket"]
	_, _ = bb.WriteString(style.Color)
	checkUnderlineAndBold(&style, bb)
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.theme.InfoTheme["reset"].Color)
	return bb.String()
}

func (rd *Renderer) Executable(toRender string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	style := rd.theme.Special["exe"]
	_, _ = bb.WriteString(style.Color)
	checkUnderlineAndBold(&style, bb)
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.theme.InfoTheme["reset"].Color)
	return bb.String()
}

func (rd *Renderer) ExecutableIcon(toRender string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	style := rd.theme.Special["exe"]
	_, _ = bb.WriteString(style.Color)
	_, _ = bb.WriteString(style.Icon)
	_, _ = bb.WriteString(" ")
	checkUnderlineAndBold(&style, bb)
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.theme.InfoTheme["reset"].Color)
	return bb.String()
}

func (rd *Renderer) Dir(toRender string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	style := rd.theme.Special["dir"]
	_, _ = bb.WriteString(style.Color)
	checkUnderlineAndBold(&style, bb)
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.theme.InfoTheme["reset"].Color)
	return bb.String()
}

func (rd *Renderer) DirIcon(toRender string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	style := rd.theme.Special["dir"]
	_, _ = bb.WriteString(style.Color)
	if s, ok := rd.theme.Name[strings.ToLower(toRender)]; ok {
		style = s
	} else {
		ext := filepath.Ext(toRender)
		if len(ext) > 0 {
			ext = ext[1:]
			ext = strings.ToLower(ext)
			if s, ok := rd.theme.Ext[ext]; ok {
				style = s
			}
		}
	}
	_, _ = bb.WriteString(style.Icon)
	_, _ = bb.WriteString(" ")
	checkUnderlineAndBold(&style, bb)
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.theme.InfoTheme["reset"].Color)

	return bb.String()
}

func (rd *Renderer) File(toRender string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	style := rd.theme.Special["file"]
	_, _ = bb.WriteString(style.Color)
	if s, ok := rd.theme.Name[toRender]; ok {
		style = s
	}
	checkUnderlineAndBold(&style, bb)
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.theme.InfoTheme["reset"].Color)

	return bb.String()
}

func (rd *Renderer) FileIcon(toRender string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	style := rd.theme.Special["file"]
	_, _ = bb.WriteString(style.Color)
	if s, ok := rd.theme.Name[toRender]; ok {
		style = s
	}
	_, _ = bb.WriteString(style.Icon)
	_, _ = bb.WriteString(" ")
	checkUnderlineAndBold(&style, bb)
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.theme.InfoTheme["reset"].Color)

	return bb.String()
}

func (rd *Renderer) gitByStatus(name string, status string) string {
	style, ok := rd.theme.Git[status]
	if !ok {
		panic("no such git status:" + status)
	}
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString(style.Color)
	checkUnderlineAndBold(&style, bb)
	_, _ = bb.WriteString(name)
	_, _ = bb.WriteString(rd.theme.InfoTheme["reset"].Color)
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

func (rd *Renderer) Inode(inode string) string {
	return rd.infoByName(inode, "inode")
}

func (rd *Renderer) DirPrompt(dir string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	style := rd.theme.Special["dir-prompt"]
	_, _ = bb.WriteString(style.Color)
	checkUnderlineAndBold(&style, bb)
	_, _ = bb.WriteString(style.Icon)
	_, _ = bb.WriteString(dir)
	_, _ = bb.WriteString(rd.theme.InfoTheme["reset"].Color)
	return bb.String()
}

func checkUnderlineAndBold(style *theme.Style, bb *bytebufferpool.ByteBuffer) {
	if style.Underline {
		_, _ = bb.WriteString(theme.Underline)
	}
	if style.Bold {
		_, _ = bb.WriteString(theme.Bold)
	}
}
