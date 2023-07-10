//go:build custom

package theme

import (
	_ "embed"
)

//go:embed custom_builtin.json
var customThemeJson []byte

func init() {
	// read the first line of customThemeIni
	a, err, fatal := getTheme(customThemeJson)
	if fatal != nil {
		panic(fatal)
	}
	if err != nil {
		panic(err)
	}
	DefaultAll = a
	_init = true
}
