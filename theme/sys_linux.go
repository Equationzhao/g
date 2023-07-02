//go:build linux

package theme

func init() {
	DefaultTheme["dir"] = Style{
		Icon:  "\ue712",
		Color: BrightBlue,
	}
}
