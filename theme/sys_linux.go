//go:build linux

package theme

func init() {
	DefaultAll.Name["sys"] = Style{
		Icon:  "\ue712",
		Color: BrightBlue,
	}
}
