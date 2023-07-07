package theme

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/ini.v1"
)

type Style struct {
	// Color of the text.
	Color string
	// unicode Icon
	Icon string
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
theme:

[info]
d 		= blue
l 		= purple
b 		= yellow
c 		= yellow
p 		= yellow
s 		= yellow
r 		= yellow
w 		= red
x 		= green
- 		= White
time 	= blue
size 	= green
owner 	= yellow
group 	= yellow
reset 	= reset
root 	= red

[dir]
color = blue
icon = 📁

[exec,exe]
color = green
icon = 🚀

[file]
color = White
icon = 📄

......
// if using 256 color, you can use color code like this:
[info]
d 		= [0-255]@256
// if using rgb color, you can use color code like this:
[info]
d 		= [0-255,0-255,0-255]@rgb
// if using hex color, you can use color code like this:
[info]
d 		= [hex]@hex
*/

type ErrBadColor struct {
	name string
	error
}

func (e ErrBadColor) Error() string {
	return fmt.Sprintf("bad color for %s:%s", e.name, e.error)
}

func GetTheme(path string) error {
	cfg, err := ini.Load(path)
	if err != nil {
		return err
	}
	infoTheme, theme, err := getTheme(cfg)

	for name, style := range infoTheme {
		DefaultInfoTheme[name] = style
	}
	for name, style := range theme {
		DefaultTheme[name] = style
	}

	SyncColorlessWithTheme()
	return err
}

func getTheme(cfg *ini.File) (infoTheme, theme Theme, errSum error) {
	cfg.BlockMode = false
	sections := cfg.Sections()
	infoThem, theme := make(Theme), make(Theme)
	for _, section := range sections {
		if section.Name() == "DEFAULT" || section.Name() == "info" {
			keys := section.Keys()
			for _, v := range keys {
				Color, err := str2color(v.String())
				if err != nil {
					errSum = errors.Join(errSum, ErrBadColor{v.Name(), err})
					continue
				}
				infoThem[v.Name()] = Style{
					Color: Color,
				}
			}
			continue
		}

		names := strings.Split(section.Name(), ",")
		color, err := str2color(section.Key("color").String())
		if err != nil {
			errSum = errors.Join(errSum, ErrBadColor{section.Name(), err})
		}

		icon := section.Key("icon").String()
		for _, name := range names {
			theme[name] = Style{
				Color: color,
				Icon:  icon,
			}
		}
	}
	return infoThem, theme, errSum
}

func ConvertThemeColor() {
	for key := range DefaultTheme {
		color, err := ConvertColorIfGreaterThanExpect(ColorLevel, DefaultTheme[key].Color)
		if err != nil {
			continue
		}
		DefaultTheme[key] = Style{
			Icon:  DefaultTheme[key].Icon,
			Color: color,
		}
	}

	for key := range DefaultInfoTheme {
		color, err := ConvertColorIfGreaterThanExpect(ColorLevel, DefaultInfoTheme[key].Color)
		if err != nil {
			continue
		}
		DefaultInfoTheme[key] = Style{
			Color: color,
		}
	}
}
