//go:build windows

package theme

func init() {
	DefaultTheme["System"] = Style{
		Icon:  "\uE70F",
		Color: dir,
	}
}
