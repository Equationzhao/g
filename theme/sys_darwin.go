//go:build darwin

package theme

func init() {
	DefaultTheme["System"] = Style{
		Icon:  "\uF179",
		Color: dir,
	}
}
