package theme

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Style struct {
	// Color of the text.
	Color string `json:"color,omitempty"`
	// unicode Icon
	Icon      string `json:"icon,omitempty"`
	Underline bool   `json:"underline,omitempty"`
	Bold      bool   `json:"bold,omitempty"`
}

func (s *Style) ToReadable() Style {
	r := *s
	r.Color = color2str(r.Color)
	return r
}

func (s *Style) FromReadable() error {
	c, err := str2color(s.Color)
	if err != nil {
		return err
	}
	s.Color = c
	return nil
}

type Theme map[string]Style

// var info = []string{"d", "l", "b", "c", "p", "s", "r", "w", "x", "-", "time", "size", "owner", "group", "git_modified_dot", "git_renamed_dot", "git_copied_dot", "git_deleted_dot", "git_added_dot", "git_untracked_dot", "git_ignored_dot", "git_modified_sym", "git_renamed_sym", "git_copied_sym", "git_deleted_sym", "git_added_sym", "git_untracked_sym", "git_ignored_sym"}
func color2str(color string) string {
	switch color {
	case Red:
		return "red"
	case Green:
		return "green"
	case Yellow:
		return "yellow"
	case Blue:
		return "blue"
	case Purple:
		return "purple"
	case Cyan:
		return "cyan"
	case White:
		return "white"
	case Black:
		return "black"
	case BrightRed:
		return "BrightPed"
	case BrightGreen:
		return "BrightPreen"
	case BrightYellow:
		return "BrightYellow"
	case BrightBlue:
		return "BrightBlue"
	case BrightPurple:
		return "BrightPurple"
	case BrightCyan:
		return "BrightCyan"
	case BrightWhite:
		return "BrightWhite"
	case BrightBlack:
		return "BrightBlack"
	case Reset:
		return "reset"
	case Underline:
		return "underline"
	default:
		// detect format:
		strReader := strings.NewReader(color)

		// 1.8bit/256 color
		var c uint8
		_, err := fmt.Fscanf(strReader, Color256Format, &c)
		if err == nil {
			return fmt.Sprintf("[%d]@256", c)
		}
		// 2.rgb
		var (
			r uint8 = 0
			g uint8 = 0
			b uint8 = 0
		)
		strReader = strings.NewReader(color)
		_, err = fmt.Fscanf(strReader, RGBFormat, &r, &g, &b)
		if err == nil {
			return fmt.Sprintf("[%d,%d,%d]@rgb", r, g, b)
		}

		color = strings.ReplaceAll(color, " ", "")
		if strings.HasPrefix(color, Underline) {
			return color2str(Underline) + " + " + color2str(color[len(Underline):])
		}
		return ""
	}
}

// str2color convert string to color
// support: red, green, yellow, blue, purple, cyan, white, black, and their bright version
// Underline
// any color with underline, should be in the format of "Underline + [color]"
// [value]@256
// [values]@rgb
// [values]@hex (will be turned to rgb)
func str2color(str string) (string, error) {
	switch str {
	case "":
		return "", nil
	case "black", "Black":
		return Black, nil
	case "red", "Red":
		return Red, nil
	case "green", "Green":
		return Green, nil
	case "yellow", "Yellow":
		return Yellow, nil
	case "blue", "Blue":
		return Blue, nil
	case "purple", "Purple":
		return Purple, nil
	case "cyan", "Cyan":
		return Cyan, nil
	case "white", "White":
		return White, nil
	case "bright-red", "BrightRed":
		return BrightRed, nil
	case "bright-green", "BrightGreen":
		return BrightGreen, nil
	case "bright-yellow", "BrightYellow":
		return BrightYellow, nil
	case "bright-blue", "BrightBlue":
		return BrightBlue, nil
	case "bright-purple", "BrightPurple":
		return BrightPurple, nil
	case "bright-cyan", "BrightCyan":
		return BrightCyan, nil
	case "bright-white", "BrightWhite":
		return BrightWhite, nil
	case "bright-black", "BrightBlack":
		return BrightBlack, nil
	case "reset", "Reset":
		return Reset, nil
	case "underline", "Underline":
		return Underline, nil
	default:
		// remove spaces
		str = strings.ReplaceAll(str, " ", "")

		// 256 color
		if strings.HasSuffix(str, "@256") {
			code, err := strconv.Atoi(strings.Trim(str[:len(str)-4], "[]"))
			if err != nil {
				return "", err
			}
			colorStr, err := Color256(code)
			if err != nil {
				return "", err
			}
			return colorStr, nil
		}

		// rgb color
		if strings.HasSuffix(str, "@rgb") {
			code := strings.Trim(str[:len(str)-4], "[]")
			rgb := strings.Split(code, ",")
			if len(rgb) != 3 {
				return "", errors.New("too many or too few rgb values")
			}
			r, err1 := strconv.ParseUint(rgb[0], 10, 8)
			g, err2 := strconv.ParseUint(rgb[1], 10, 8)
			b, err3 := strconv.ParseUint(rgb[2], 10, 8)
			if err1 != nil || err2 != nil || err3 != nil {
				return "", errors.New("rgb values must be numbers")
			}
			colorStr, err := RGB(uint8(r), uint8(g), uint8(b))
			if err != nil {
				return "", err
			}
			return colorStr, nil
		}

		// hex
		if strings.HasSuffix(str, "@hex") {
			code := strings.Trim(str[:len(str)-4], "[]")
			rgb := HexToRgb(code)
			colorStr, err := RGB(rgb[0], rgb[1], rgb[2])
			if err != nil {
				return "", errors.New("rgb values must be numbers")
			}
			return colorStr, nil
		}

		// underline
		if strings.HasPrefix(str, "underline+") {
			color, err := str2color(str[len("underline+"):])
			if err != nil {
				return "", err
			}
			return Underline + color, nil
		}

		return Reset, nil
	}
}

/*
......
// if using 256 color, you can use color code like this:
[0-255]@256
// if using rgb color, you can use color code like this:
[0-255,0-255,0-255]@rgb
// if using hex color, you can use color code like this:
[hex]@hex
*/

type ErrBadColor struct {
	name string
	error
}

func (e ErrBadColor) Error() string {
	return fmt.Sprintf("bad color for %s:%s", e.name, e.error)
}

func GetTheme(path string) error {
	themeJson, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	a, err, fatal := getTheme(themeJson)
	if fatal != nil {
		return fatal
	}
	DefaultAll = a
	return err
}

func getTheme(themeJson []byte) (theme All, errSum error, fatal error) {
	err := json.Unmarshal(themeJson, &theme)
	if err != nil {
		return All{}, nil, err
	}
	// check
	errSum = nil
	check := func(m Theme) {
		for key := range m {
			style := m[key]
			err := style.FromReadable()
			if err != nil {
				errSum = errors.Join(errSum, ErrBadColor{name: key, error: err})
				continue
			}
			m[key] = style
		}
	}
	check(theme.InfoTheme)
	check(theme.Permission)
	check(theme.Size)
	check(theme.User)
	check(theme.Group)
	check(theme.Symlink)
	check(theme.Git)
	check(theme.Name)
	check(theme.Special)
	check(theme.Ext)
	return theme, errSum, nil
}

func ConvertThemeColor() {
	convert := func(m map[string]Style) {
		for key := range m {
			color, err := ConvertColorIfGreaterThanExpect(ColorLevel, m[key].Color)
			if err != nil {
				continue
			}
			m[key] = Style{
				Icon:      m[key].Icon,
				Color:     color,
				Underline: m[key].Underline,
				Bold:      m[key].Bold,
			}
		}
	}
	convert(DefaultAll.InfoTheme)
	convert(DefaultAll.Permission)
	convert(DefaultAll.Size)
	convert(DefaultAll.User)
	convert(DefaultAll.Group)
	convert(DefaultAll.Symlink)
	convert(DefaultAll.Git)
	convert(DefaultAll.Name)
	convert(DefaultAll.Special)
	convert(DefaultAll.Ext)
}
