package theme

import (
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

var info = []string{"d", "l", "b", "c", "p", "s", "r", "w", "x", "-", "time", "size", "owner", "group"}

func str2color(str string) string {
	switch str {
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
	case "reset", "Reset":
		return Reset
	default:
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
			for _, v := range info {
				DefaultInfoTheme[v] = Style{
					Color: str2color(section.Key(v).String()),
				}
			}
		}
		names := strings.Split(section.Name(), ",")
		color := str2color(section.Key("color").String())
		icon := section.Key("icon").String()
		for _, name := range names {
			DefaultInfoTheme[name] = Style{
				Color: color,
				Icon:  icon,
			}
		}
	}
	return nil
}
