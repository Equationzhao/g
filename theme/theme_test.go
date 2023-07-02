package theme

import (
	colortool "github.com/gookit/color"
	"testing"
)

func TestAll(t *testing.T) {
	ColorLevel = colortool.Level16
	ConvertThemeColor()
	t.Logf("info")
	for key, style := range DefaultInfoTheme {
		t.Logf("%s%s%s", style.Color, key, Reset)
	}
	t.Logf("theme")
	for key, style := range DefaultTheme {
		t.Logf("%s%s %s%s", style.Color, style.Icon, key, Reset)
	}
}
