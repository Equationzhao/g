package theme

import (
	"strings"

	"gopkg.in/ini.v1"
)

const (
	Black     = "\033[1;30m"
	Red       = "\033[1;31m"
	Green     = "\033[1;32m"
	Yellow    = "\033[1;33m"
	Blue      = "\033[1;34m"
	Purple    = "\033[1;35m"
	Cyan      = "\033[1;36m"
	White     = "\033[1;37m"
	Reset     = "\033[0m"
	Success   = "\033[1;32m"
	Error     = "\033[1;31m"
	Warn      = "\033[1;33m"
	Underline = "\033[4m"
	Bold      = "\033[1m"
	Reverse   = "\033[7m"
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
	case Reset:
		return "reset"
	case Underline:
		return "underline"
	case Bold:
		return "bold"
	case Reverse:
		return "reverse"
	default:
		return ""
	}
}

func str2color(str string) string {
	switch str {
	case "":
		return ""
	case "red", "Red":
		return Red
	case "green", "Green":
		return Green
	case "yellow", "Yellow":
		return Yellow
	case "blue", "Blue":
		return Blue
	case "purple", "Purple":
		return Purple
	case "cyan", "Cyan":
		return Cyan
	case "white", "White":
		return White
	case "black", "Black":
		return Black
	case "reset", "Reset":
		return Reset
	case "underline", "Underline":
		return Underline
	case "bold", "Bold":
		return Bold
	case "reverse", "Reverse":
		return Reverse
	default:
		str = strings.ReplaceAll(str, " ", "")
		if strings.HasPrefix(str, "Reverse+") || strings.HasPrefix(str, "reverse+") {
			return Reverse + str2color(str[8:])
		}
		// if strings.HasPrefix(str, "Bold+") || strings.HasPrefix(str, "bold+") {
		// 	return Bold + str2color(str[8:])
		// }
		// if strings.HasPrefix(str, "Underline+") || strings.HasPrefix(str, "underline+") {
		// 	return Bold + str2color(str[8:])
		// }
		return Reset
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
*/

func GetTheme(path string) error {
	cfg, err := ini.Load(path)
	if err != nil {
		return err
	}
	cfg.BlockMode = false

	sections := cfg.Sections()
	for _, section := range sections {
		if section.Name() == "DEFAULT" || section.Name() == "info" {
			keys := section.Keys()
			for _, v := range keys {
				o := DefaultInfoTheme[v.Name()]
				o.Color = str2color(v.String())
				DefaultInfoTheme[v.Name()] = o
			}
			continue
		}

		names := strings.Split(section.Name(), ",")
		color := str2color(section.Key("color").String())
		icon := section.Key("icon").String()
		for _, name := range names {
			DefaultTheme[name] = Style{
				Color: color,
				Icon:  icon,
			}
		}
	}
	return nil
}
