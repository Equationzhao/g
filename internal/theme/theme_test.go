package theme

import (
	"testing"

	"github.com/Equationzhao/g/internal/const"

	colortool "github.com/gookit/color"
)

func TestAll(t *testing.T) {
	ColorLevel = colortool.Level16
	ConvertThemeColor()
	pl := func(m map[string]Style) {
		for key := range m {
			t.Logf("%s %s %s %s", m[key].Color, m[key].Icon, key, constval.Reset)
		}
	}
	pl(DefaultAll.InfoTheme)
	pl(DefaultAll.Permission)
	pl(DefaultAll.Size)
	pl(DefaultAll.User)
	pl(DefaultAll.Group)
	pl(DefaultAll.Symlink)
	pl(DefaultAll.Git)
	pl(DefaultAll.Name)
	pl(DefaultAll.Special)
	pl(DefaultAll.Ext)
}

func TestColor(t *testing.T) {
	println(constval.Green + "\uF48A " + constval.Underline + constval.Bold + "hello" + constval.Red + " hello" + constval.Reset)
	println(constval.Green + "\uF48A " + constval.Underline + "hello" + constval.Red + " hello" + constval.Reset)
}
