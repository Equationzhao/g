//go:build linux

package theme

func init() {
	DefaultTheme["sys"] = Style{
		Icon:  "\ue712",
		Color: BrightBlue,
	}
}
