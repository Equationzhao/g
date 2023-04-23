package render

import (
	"github.com/Equationzhao/g/theme"
	"github.com/valyala/bytebufferpool"
	"path/filepath"
	"strings"
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

func (r *Renderer) Size(toRender string) string {
	return r.infoByName(toRender, "size")
}

func (r *Renderer) Owner(toRender string) string {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	toRenderNoSpace := strings.Replace(toRender, " ", "", -1)
	if toRenderNoSpace == "root" {
		_, _ = bb.WriteString(r.infoTheme["root"].Color)
	} else {
		_, _ = bb.WriteString(r.infoTheme["owner"].Color)
	}

	_, _ = bb.WriteString(toRender)
	_, _ = bb.WriteString(r.infoTheme["reset"].Color)
	return bb.String()
}

func (r *Renderer) Group(toRender string) string {
	return r.infoByName(toRender, "group")
}

func (r *Renderer) Time(toRender string) string {
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
	_, _ = bb.WriteString(r.theme["reset"].Color)
	return bb.String()
}

func (r *Renderer) SymlinkIcon(toRender string, path string) string {
	icon := r.Icon("symlink")
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	_, _ = bb.WriteString(r.theme["symlink"].Color)
	_, _ = bb.WriteString(icon)
	_, _ = bb.WriteString(" ")
	symlinks, err := filepath.EvalSymlinks(filepath.Join(path, toRender))
	if err != nil {
		symlinks = err.Error()
	}
	_, _ = bb.WriteString(toRender + " -> " + symlinks)
	_, _ = bb.WriteString(r.theme["reset"].Color)
	return bb.String()
}

func (r *Renderer) Symlink(toRender string) string {
	return r.byName(toRender, "symlink")
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
