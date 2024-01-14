//go:build darwin

package theme

func init() {
	DefaultAll.Name["System"] = Style{
		Icon:  "\uF179",
		Color: dir,
	}
}
