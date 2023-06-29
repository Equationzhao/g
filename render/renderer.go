package render

import (
	"errors"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/Equationzhao/g/theme"
	"github.com/valyala/bytebufferpool"
)

type Renderer struct {
	infoTheme, theme theme.Theme
}

func (r *Renderer) SetInfoTheme(theme theme.Theme) *Renderer {
	r.infoTheme = theme
	return r
}

func (r *Renderer) SetTheme(theme theme.Theme) *Renderer {
	r.theme = theme
	return r
}

func NewRenderer(theme, infoTheme theme.Theme) *Renderer {
	return &Renderer{infoTheme: infoTheme, theme: theme}
}

func (r *Renderer) FileMode(toRender string) string {
	// return file mode like -rwxrwxrwx/drwxrwxrwx but in color
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	for _, c := range toRender {
		switch c {
		case '-':
			_, _ = bb.WriteString(r.infoTheme["-"].Color)
		case 'L':
			_, _ = bb.WriteString(r.infoTheme["l"].Color)
		case 'd':
			_, _ = bb.WriteString(r.infoTheme["d"].Color)
		case 'r':
			_, _ = bb.WriteString(r.infoTheme["r"].Color)
		case 'w':
			_, _ = bb.WriteString(r.infoTheme["w"].Color)
		case 'x':
			_, _ = bb.WriteString(r.infoTheme["x"].Color)
		case 'S':
			_, _ = bb.WriteString(r.infoTheme["s"].Color)
		}
		_, _ = bb.WriteString(string(c))
	}
	_, _ = bb.WriteString(r.infoTheme["reset"].Color)
	return bb.String()
}

func (r *Renderer) Size(toRender, unit string) string {
	if strings.HasSuffix(toRender, "-") {
		return r.infoByName(toRender, "-")
	}
	return r.infoByName(toRender, unit)
}

func (r *Renderer) BlockSize(toRender string) string {
	if strings.HasSuffix(toRender, "-") {
		return r.infoByName(toRender, "-")
	}
	return r.infoByName(toRender, "bit")
}

func (r *Renderer) Link(toRender string) string {
	return r.infoByName(toRender, "link")
}

func (r *Renderer) Owner(toRender string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	toRenderNoSpace := strings.Replace(toRender, " ", "", -1)
	if toRenderNoSpace == "root" {
		_, _ = bb.WriteString(r.infoTheme["root"].Color)
		_, _ = bb.WriteString("\ue315")
	} else {
		_, _ = bb.WriteString(r.infoTheme["owner"].Color)
	}

	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(r.infoTheme["reset"].Color)
	return bb.String()
}

func (r *Renderer) Group(toRender string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	toRenderNoSpace := strings.Replace(toRender, " ", "", -1)
	if toRenderNoSpace == "root" {
		_, _ = bb.WriteString(r.infoTheme["root"].Color)
		_, _ = bb.WriteString("\ue315")
	} else {
		_, _ = bb.WriteString(r.infoTheme["group"].Color)
	}

	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(r.infoTheme["reset"].Color)
	return bb.String()
}

func (r *Renderer) Time(toRender string) string {
	return r.infoByName(toRender, "time")
}

func (r *Renderer) RTime(toRender string) string {
	return r.infoByName(toRender, "time")
}

func (r *Renderer) Name(toRender string) string {
	return r.infoByName(toRender, "name")
}

func (r *Renderer) infoByName(toRender string, name string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString(r.infoTheme[name].Color)
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(r.infoTheme["reset"].Color)
	return bb.String()
}

func (r *Renderer) ByExt(toRender string) string {
	// get ext
	ext := filepath.Ext(toRender)
	if len(ext) > 0 {
		ext = ext[1:]
	} else {
		ext = toRender // if no ext, try to color by name
	}
	return r.byName(toRender, ext)
}

// ByExtIcon returns the icon and the name of the file
// if the file has no icon it returns an empty string
func (r *Renderer) ByExtIcon(toRender string) string {
	// get ext
	ext := filepath.Ext(toRender)
	if len(ext) > 0 {
		ext = ext[1:]
	} else {
		ext = toRender // if no ext, try to color by name
	}
	icon := r.Icon(ext)
	if icon == "" {
		ext = "file"
		icon = r.Icon("file")
	}
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)

	_, _ = bb.WriteString(r.theme[ext].Color)
	_, _ = bb.WriteString(icon)
	_, _ = bb.WriteString(" ")
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(r.infoTheme["reset"].Color)
	return bb.String()
}

// SymlinkIconPlus returns the icon and the name of the file, and dereferences the symlink
func (r *Renderer) SymlinkIconPlus(toRender string, path string, plus string) string {
	icon := r.Icon("symlink")
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString(r.theme["symlink"].Color)
	_, _ = bb.WriteString(icon)
	_, _ = bb.WriteString(" ")
	symlinks, err := filepath.EvalSymlinks(path)
	if err != nil {
		var pathErr *fs.PathError
		if errors.As(err, &pathErr) {
			_, _ = bb.WriteString(toRender + plus)
			_, _ = bb.WriteString(theme.Error)
			_, _ = bb.WriteString(" -> " + pathErr.Path)
			_, _ = bb.WriteString(r.infoTheme["reset"].Color)
			return bb.String()
		}
		symlinks = err.Error()
	}
	_, _ = bb.WriteString(toRender + plus)
	_, _ = bb.WriteString(theme.Success)
	_, _ = bb.WriteString(" -> " + symlinks)
	_, _ = bb.WriteString(r.infoTheme["reset"].Color)
	return bb.String()
}

// SymlinkIconNoDereferencePlus returns the icon and the name of the file, but does not dereference the symlink
func (r *Renderer) SymlinkIconNoDereferencePlus(toRender string, plus string) string {
	icon := r.Icon("symlink")
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString(r.theme["symlink"].Color)
	_, _ = bb.WriteString(icon)
	_, _ = bb.WriteString(" ")
	_, _ = bb.WriteString(toRender + plus)
	_, _ = bb.WriteString(r.infoTheme["reset"].Color)
	return bb.String()
}

func (r *Renderer) SymlinkIconNoDereference(toRender string) string {
	return r.SymlinkIconNoDereferencePlus(toRender, "")
}

func (r *Renderer) SymlinkIcon(toRender string, path string) string {
	return r.SymlinkIconPlus(toRender, path, "")
}

// SymlinkPlus returns the icon and the name of the file
func (r *Renderer) SymlinkPlus(toRender string, path string, plus string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString(r.theme["symlink"].Color)
	symlinks, err := filepath.EvalSymlinks(path)
	if err != nil {
		var pathErr *fs.PathError
		if errors.As(err, &pathErr) {
			_, _ = bb.WriteString(toRender + plus)
			_, _ = bb.WriteString(theme.Error)
			_, _ = bb.WriteString(" -> " + pathErr.Path)
			_, _ = bb.WriteString(r.infoTheme["reset"].Color)
			return bb.String()
		}
		symlinks = err.Error()
	}
	_, _ = bb.WriteString(toRender + plus)
	_, _ = bb.WriteString(theme.Success)
	_, _ = bb.WriteString(" -> " + symlinks)
	_, _ = bb.WriteString(r.infoTheme["reset"].Color)
	return bb.String()
}

func (r *Renderer) SymlinkNoDereferencePlus(toRender string, plus string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString(r.theme["symlink"].Color)
	_, _ = bb.WriteString(toRender + plus)
	_, _ = bb.WriteString(r.infoTheme["reset"].Color)
	return bb.String()
}

func (r *Renderer) SymlinkNoDereference(str string) string {
	return r.SymlinkNoDereferencePlus(str, "")
}

func (r *Renderer) Symlink(toRender string, path string) string {
	return r.SymlinkPlus(toRender, path, "")
}

func (r *Renderer) PipeIcon(toRender string) string {
	icon := r.Icon("pipe")
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString(r.theme["pipe"].Color)
	_, _ = bb.WriteString(icon)
	_, _ = bb.WriteString(" ")
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(r.infoTheme["reset"].Color)
	return bb.String()
}

func (r *Renderer) Pipe(toRender string) string {
	return r.byName(toRender, "symlink")
}

func (r *Renderer) SocketIcon(toRender string) string {
	icon := r.Icon("socket")
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString(r.theme["socket"].Color)
	_, _ = bb.WriteString(icon)
	_, _ = bb.WriteString(" ")
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(r.infoTheme["reset"].Color)
	return bb.String()
}

func (r *Renderer) Socket(toRender string) string {
	return r.byName(toRender, "socket")
}

func (r *Renderer) Executable(toRender string) string {
	return r.byName(toRender, "exec")
}

func (r *Renderer) RegularFile(toRender string) string {
	return r.byName(toRender, "file")
}

func (r *Renderer) Dir(toRender string) string {
	return r.byName(toRender, "dir")
}

func (r *Renderer) DirIcon(toRender string) string {
	ext := filepath.Ext(toRender)
	if len(ext) > 0 {
		ext = ext[1:]
	} else {
		ext = toRender // if no ext, try to color by name
	}
	icon := r.Icon(ext)
	if icon == "" {
		icon = r.Icon("dir")
	}
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString(r.theme["dir"].Color)
	_, _ = bb.WriteString(icon)
	_, _ = bb.WriteString(" ")
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(r.infoTheme["reset"].Color)

	return bb.String()
}

func (r *Renderer) byName(toRender string, name string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString(r.theme[name].Color)
	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(r.infoTheme["reset"].Color) // IT IS INFO THEME
	return bb.String()
}

func (r *Renderer) Icon(name string) string {
	return r.theme[name].Icon
}

func (r *Renderer) gitByStatus(name string, status string) string {
	return r.infoTheme[status].Color + r.infoTheme[status].Icon + r.infoTheme["reset"].Color + " " + name
}

func (r *Renderer) GitModified(name string) string {
	return r.gitByStatus(name, "git_modified_dot")
}

func (r *Renderer) GitUntracked(name string) string {
	return r.gitByStatus(name, "git_untracked_dot")
}

func (r *Renderer) GitAdded(name string) string {
	return r.gitByStatus(name, "git_added_dot")
}

func (r *Renderer) GitRenamed(name string) string {
	return r.gitByStatus(name, "git_renamed_dot")
}

func (r *Renderer) GitDeleted(name string) string {
	return r.gitByStatus(name, "git_deleted_dot")
}

func (r *Renderer) GitIgnored(name string) string {
	return r.gitByStatus(name, "git_ignored_dot")
}

func (r *Renderer) GitCopied(name string) string {
	return r.gitByStatus(name, "git_copied_dot")
}

func (r *Renderer) GitModifiedSym(name string) string {
	return r.gitByStatus(name, "git_modified_sym")
}

func (r *Renderer) GitUntrackedSym(name string) string {
	return r.gitByStatus(name, "git_untracked_sym")
}

func (r *Renderer) GitAddedSym(name string) string {
	return r.gitByStatus(name, "git_added_sym")
}

func (r *Renderer) GitRenamedSym(name string) string {
	return r.gitByStatus(name, "git_renamed_sym")
}

func (r *Renderer) GitDeletedSym(name string) string {
	return r.gitByStatus(name, "git_deleted_sym")
}

func (r *Renderer) GitIgnoredSym(name string) string {
	return r.gitByStatus(name, "git_ignored_sym")
}

func (r *Renderer) GitCopiedSym(name string) string {
	return r.gitByStatus(name, "git_copied_sym")
}

func (r *Renderer) Inode(inode string) string {
	return r.infoByName(inode, "inode")
}
