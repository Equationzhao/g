//go:build windows

package theme

func init() {
	DefaultTheme["System"] = Style{
		Icon:  "\uE70F",
		Color: dir,
	}
	Group["DevToolsUser"] = Style{
		Color: color256(202),
	}
	Name["program files"] = Style{
		Icon: "\ueb44",
	}
	Name["program files (x86)"] = Style{
		Icon: "\ueb44",
	}
}
