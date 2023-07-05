//go:build theme

package theme

import (
	"path/filepath"

	"github.com/Equationzhao/g/slices"
	"github.com/Equationzhao/g/util/cmp"
	"gopkg.in/ini.v1"
)

type kv struct {
	key   string
	value any
}

func init() {
	defaultThemeIni := ini.Empty()
	defaultThemeIni.BlockMode = false

	info := defaultThemeIni.Section("info")
	infoArray := make([]kv, 0, len(DefaultInfoTheme))
	for k, v := range DefaultInfoTheme {
		infoArray = append(infoArray, kv{k, color2str(v.Color)})
	}
	slices.SortFunc(
		infoArray, func(i, j kv) int {
			return cmp.Compare(i.key, j.key)
		},
	)
	for _, v := range infoArray {
		_, _ = info.NewKey(v.key, v.value.(string))
	}

	default_ := make([]kv, 0, len(DefaultTheme))
	for k, v := range DefaultTheme {
		default_ = append(default_, kv{k, v})
	}

	slices.SortFunc(
		default_, func(i, j kv) int {
			return cmp.Compare(i.key, j.key)
		},
	)

	for _, k := range default_ {
		section := defaultThemeIni.Section(k.key)
		if c := k.value.(Style).Color; c != "" {
			_, err := section.NewKey("color", color2str(k.value.(Style).Color))
			if err != nil {
				println(err.Error())
			}
		}
		_, err := section.NewKey("icon", k.value.(Style).Icon)
		if err != nil {
			println(err.Error())
		}
	}

	err := defaultThemeIni.SaveTo(filepath.Join("theme", "default.ini"))
	if err != nil {
		println(err.Error())
	}
}
