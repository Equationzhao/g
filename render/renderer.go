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

	"github.com/Equationzhao/g/theme"
	"github.com/Equationzhao/g/util"
	"github.com/hako/durafmt"
	"github.com/valyala/bytebufferpool"
)

type Renderer struct {
	infoTheme, theme theme.Theme
}

func (rd *Renderer) SetInfoTheme(theme theme.Theme) *Renderer {
	rd.infoTheme = theme
	return rd
}

func (rd *Renderer) SetTheme(theme theme.Theme) *Renderer {
	rd.theme = theme
	return rd
}

func NewRenderer(theme, infoTheme theme.Theme) *Renderer {
	return &Renderer{infoTheme: infoTheme, theme: theme}
}

func (rd *Renderer) FileMode(toRender string) string {
	// return file mode like -rwxrwxrwx/drwxrwxrwx but in color
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	for _, c := range toRender {
		switch c {
		case '-':
			_, _ = bb.WriteString(rd.infoTheme["-"].Color)
		case 'L':
			_, _ = bb.WriteString(rd.infoTheme["l"].Color)
		case 'd':
			_, _ = bb.WriteString(rd.infoTheme["d"].Color)
		case 'r':
			_, _ = bb.WriteString(rd.infoTheme["r"].Color)
		case 'w':
			_, _ = bb.WriteString(rd.infoTheme["w"].Color)
		case 'x', 's', 't':
			_, _ = bb.WriteString(rd.infoTheme["x"].Color)
		case 'S', 'T':
			_, _ = bb.WriteString(rd.infoTheme["s"].Color)
		}
		_, _ = bb.WriteString(string(c))
	}
	_, _ = bb.WriteString(rd.infoTheme["reset"].Color)
	return bb.String()
}

func (rd *Renderer) Size(toRender, unit string) string {
	if strings.HasSuffix(toRender, "-") {
		return rd.infoByName(toRender, "-")
	}
	return rd.infoByName(toRender, unit)
}

func (rd *Renderer) BlockSize(toRender string) string {
	if strings.HasSuffix(toRender, "-") {
		return rd.infoByName(toRender, "-")
	}
	return rd.infoByName(toRender, "bit")
}

func (rd *Renderer) Link(toRender string) string {
	return rd.infoByName(toRender, "link")
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

	if util.SliceContains(root, toRender) {
		_, _ = bb.WriteString(rd.infoTheme["root"].Color)
	} else {
		if util.SliceContains(byName, toRender) {
			_, _ = bb.WriteString(rd.infoTheme[toRender].Color)
		} else if runtime.GOOS == "windows" && rootSid.MatchString(toRender) {
			_, _ = bb.WriteString(rd.infoTheme["root"].Color)
		} else {
			_, _ = bb.WriteString(rd.infoTheme["owner"].Color)
		}
	}

	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.infoTheme["reset"].Color)
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

	if util.SliceContains(root, toRender) {
		_, _ = bb.WriteString(rd.infoTheme["root"].Color)
	} else {
		if util.SliceContains(byName, toRender) {
			_, _ = bb.WriteString(rd.infoTheme["byName"].Color)
		} else if runtime.GOOS == "windows" && rootSid.MatchString(toRender) {
			_, _ = bb.WriteString(rd.infoTheme["root"].Color)
		} else {
			_, _ = bb.WriteString(rd.infoTheme["group"].Color)
		}
	}

	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.infoTheme["reset"].Color)
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

	const gUint = 10
	var r, b float64 = 215, 255

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
		code := 213
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
		return rd.infoTheme["time"].Color
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
			"%s%s ago%s", rd.calculateRTimeColor(t), dura.LimitFirstN(1).String(), rd.infoTheme["reset"].Color,
		)
	} else if t == 0 {
		return "now"
	} else {
		dura = durafmt.Parse(-t)
		return fmt.Sprintf(
			"%sin %s%s", rd.calculateRTimeColor(t), dura.LimitFirstN(1).String(), rd.infoTheme["reset"].Color,
		)
	}
}

func (rd *Renderer) Name(toRender string) string {
	return rd.infoByName(toRender, "name")
}

func (rd *Renderer) infoByName(toRender string, name string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString(rd.infoTheme[name].Color)
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.infoTheme["reset"].Color)
	return bb.String()
}

func (rd *Renderer) ByExt(toRender string) string {
	// get ext
	ext := filepath.Ext(toRender)
	if len(ext) > 0 {
		ext = ext[1:]
	} else {
		ext = toRender // if no ext, try to color by name
	}
	return rd.byName(toRender, strings.ToLower(ext))
}

// ByExtIcon returns the icon and the name of the file
// if the file has no icon it returns an empty string
func (rd *Renderer) ByExtIcon(toRender string) string {
	// get ext
	ext := filepath.Ext(toRender)
	if len(ext) > 0 {
		ext = ext[1:]
	} else {
		ext = toRender // if no ext, try to color by name
	}
	ext = strings.ToLower(ext)
	icon := rd.Icon(ext)
	if icon == "" {
		ext = "file"
		icon = rd.Icon("file")
	}
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)

	_, _ = bb.WriteString(rd.theme[ext].Color)
	_, _ = bb.WriteString(icon)
	_, _ = bb.WriteString(" ")
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.infoTheme["reset"].Color)
	return bb.String()
}

// SymlinkIconPlus returns the icon and the name of the file, and dereferences the symlink
func (rd *Renderer) SymlinkIconPlus(toRender string, path string, plus string, rel bool) string {
	icon := rd.Icon("symlink")
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString(rd.theme["symlink"].Color)
	_, _ = bb.WriteString(icon)
	_, _ = bb.WriteString(" ")
	symlinks, err := filepath.EvalSymlinks(path)
	if err != nil {
		var pathErr *fs.PathError
		if errors.As(err, &pathErr) {
			_, _ = bb.WriteString(toRender + plus)
			_, _ = bb.WriteString(rd.theme["symlink_arrow"].Color + rd.theme["symlink_arrow"].Icon)
			_, _ = bb.WriteString(rd.infoTheme["symlink_broken_path"].Color)
			if rel {
				symlinksRel, err := filepath.Rel(filepath.Dir(path), pathErr.Path)
				if err == nil {
					pathErr.Path = symlinksRel
				}
			}
			_, _ = bb.WriteString(pathErr.Path)
			_, _ = bb.WriteString(rd.infoTheme["reset"].Color)
			return bb.String()
		}
		symlinks = err.Error()
	}
	if rel {
		symlinksRel, err := filepath.Rel(filepath.Dir(path), symlinks)
		if err == nil {
			symlinks = symlinksRel
		}
	}
	_, _ = bb.WriteString(toRender + plus)
	_, _ = bb.WriteString(rd.theme["symlink_arrow"].Color + rd.theme["symlink_arrow"].Icon)
	_, _ = bb.WriteString(rd.infoTheme["symlink_path"].Color)
	_, _ = bb.WriteString(symlinks)
	_, _ = bb.WriteString(rd.infoTheme["reset"].Color)
	return bb.String()
}

// SymlinkIconNoDereferencePlus returns the icon and the name of the file, but does not dereference the symlink
func (rd *Renderer) SymlinkIconNoDereferencePlus(toRender string, plus string) string {
	icon := rd.Icon("symlink")
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString(rd.theme["symlink"].Color)
	_, _ = bb.WriteString(icon)
	_, _ = bb.WriteString(" ")
	_, _ = bb.WriteString(toRender + plus)
	_, _ = bb.WriteString(rd.infoTheme["reset"].Color)
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
	_, _ = bb.WriteString(rd.theme["symlink"].Color)
	symlinks, err := filepath.EvalSymlinks(path)
	if err != nil {
		var pathErr *fs.PathError
		if errors.As(err, &pathErr) {
			_, _ = bb.WriteString(toRender + plus)
			_, _ = bb.WriteString(rd.theme["symlink_arrow"].Color + rd.theme["symlink_arrow"].Icon)
			_, _ = bb.WriteString(rd.infoTheme["symlink_broken_path"].Color)
			if rel {
				symlinksRel, err := filepath.Rel(filepath.Dir(path), pathErr.Path)
				if err == nil {
					pathErr.Path = symlinksRel
				}
			}
			_, _ = bb.WriteString(pathErr.Path)
			_, _ = bb.WriteString(rd.infoTheme["reset"].Color)
			return bb.String()
		}
		symlinks = err.Error()
	}
	if rel {
		symlinksRel, err := filepath.Rel(filepath.Dir(path), symlinks)
		if err == nil {
			symlinks = symlinksRel
		}
	}
	_, _ = bb.WriteString(toRender + plus)
	_, _ = bb.WriteString(rd.theme["symlink_arrow"].Color + rd.theme["symlink_arrow"].Icon)
	_, _ = bb.WriteString(rd.infoTheme["symlink_path"].Color)
	_, _ = bb.WriteString(symlinks)
	_, _ = bb.WriteString(rd.infoTheme["reset"].Color)
	return bb.String()
}

func (rd *Renderer) SymlinkNoDereferencePlus(toRender string, plus string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString(rd.theme["symlink"].Color)
	_, _ = bb.WriteString(toRender + plus)
	_, _ = bb.WriteString(rd.infoTheme["reset"].Color)
	return bb.String()
}

func (rd *Renderer) SymlinkNoDereference(str string) string {
	return rd.SymlinkNoDereferencePlus(str, "")
}

func (rd *Renderer) Symlink(toRender string, path string, rel bool) string {
	return rd.SymlinkPlus(toRender, path, "", rel)
}

func (rd *Renderer) PipeIcon(toRender string) string {
	icon := rd.Icon("pipe")
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString(rd.theme["pipe"].Color)
	_, _ = bb.WriteString(icon)
	_, _ = bb.WriteString(" ")
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.infoTheme["reset"].Color)
	return bb.String()
}

func (rd *Renderer) Pipe(toRender string) string {
	return rd.byName(toRender, "symlink")
}

func (rd *Renderer) SocketIcon(toRender string) string {
	icon := rd.Icon("socket")
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString(rd.theme["socket"].Color)
	_, _ = bb.WriteString(icon)
	_, _ = bb.WriteString(" ")
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.infoTheme["reset"].Color)
	return bb.String()
}

func (rd *Renderer) Socket(toRender string) string {
	return rd.byName(toRender, "socket")
}

func (rd *Renderer) Executable(toRender string) string {
	return rd.byName(toRender, "exe")
}

func (rd *Renderer) ExecutableIcon(toRender string) string {
	icon := rd.Icon("exe")
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString(rd.theme["exe"].Color)
	_, _ = bb.WriteString(icon)
	_, _ = bb.WriteString(" ")
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.infoTheme["reset"].Color)
	return bb.String()
}

func (rd *Renderer) RegularFile(toRender string) string {
	return rd.byName(toRender, "file")
}

func (rd *Renderer) Dir(toRender string) string {
	return rd.byName(toRender, "dir")
}

func (rd *Renderer) DirIcon(toRender string) string {
	ext := filepath.Ext(toRender)
	if len(ext) > 0 {
		ext = ext[1:]
	} else {
		ext = toRender // if no ext, try to color by name
	}
	ext = strings.ToLower(ext)
	icon := rd.Icon(ext)
	if icon == "" {
		icon = rd.Icon("dir")
	}
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString(rd.theme["dir"].Color)
	_, _ = bb.WriteString(icon)
	_, _ = bb.WriteString(" ")
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.infoTheme["reset"].Color)

	return bb.String()
}

func (rd *Renderer) byName(toRender string, name string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString(rd.theme[name].Color)
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(rd.infoTheme["reset"].Color) // IT IS INFO THEME
	return bb.String()
}

func (rd *Renderer) Icon(name string) string {
	return rd.theme[name].Icon
}

func (rd *Renderer) gitByStatus(name string, status string) string {
	return rd.infoTheme[status].Color + name + rd.infoTheme["reset"].Color
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
