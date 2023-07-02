package theme

import (
	"testing"

	colortool "github.com/gookit/color"
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
