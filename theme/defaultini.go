//go:build theme

package theme

import (
	"path/filepath"
	"sort"

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
	sort.Slice(infoArray, func(i, j int) bool {
		return infoArray[i].key < infoArray[j].key
	})
	for _, v := range infoArray {
		_, _ = info.NewKey(v.key, v.value.(string))
	}

	default_ := make([]kv, 0, len(DefaultTheme))
	for k, v := range DefaultTheme {
		default_ = append(default_, kv{k, v})
	}

	sort.Slice(default_, func(i, j int) bool {
		return default_[i].key < default_[j].key
	})

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
