//go:build windows

package theme

func init() {
	DefaultAll.Special["system"] = Style{
		Icon:  "\uE70F",
		Color: dir,
	}
	DefaultAll.Group["devtoolsuser"] = Style{
		Color: color256(202),
	}
	DefaultAll.Name["program files"] = Style{
		Icon: "\ueb44",
	}
	DefaultAll.Name["program files (x86)"] = Style{
		Icon: "\ueb44",
	}
	DefaultAll.Name["windows"] = Style{
		Icon: "\uE70F",
	}
}
