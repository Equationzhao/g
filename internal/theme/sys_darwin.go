//go:build darwin

package theme

func init() {
	DefaultAll.Name["system"] = Style{
		Icon:  "\uF179",
		Color: dir,
	}
}
