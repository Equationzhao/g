//go:build custom

package theme

import (
	_ "embed"
	"gopkg.in/ini.v1"
	"strings"
)

//go:embed custom_builtin
var customThemeIni string

func init() {
	// read the first line of customThemeIni

	temp := strings.SplitN(customThemeIni, "\n", 2) // replace or merge
	Command := temp[0]
	themeContent := []byte(temp[1])

	cfg, err := ini.Load(themeContent)
	if err != nil {
		panic(err)
	}
	infoTheme, theme, _ := getTheme(cfg)

	if strings.EqualFold(Command, "replace") {
		DefaultInfoTheme = infoTheme
		DefaultTheme = theme
		SyncColorlessWithTheme()
	} else if strings.EqualFold(Command, "merge") {
		for name, style := range infoTheme {
			DefaultInfoTheme[name] = style
		}
		for name, style := range theme {
			DefaultTheme[name] = style
		}
	} else {
		panic("Command must be replace or merge")
	}
}
