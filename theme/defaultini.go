//go:build theme

package theme

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func init() {
	convert := func(theme Theme) {
		for k, style := range theme {
			(theme)[k] = style.ToReadable()
		}
	}
	convert(DefaultAll.InfoTheme)
	convert(DefaultAll.Permission)
	convert(DefaultAll.Size)
	convert(DefaultAll.User)
	convert(DefaultAll.Group)
	convert(DefaultAll.Symlink)
	convert(DefaultAll.Git)
	convert(DefaultAll.Name)
	convert(DefaultAll.Special)
	convert(DefaultAll.Ext)
	marshal, err := json.MarshalIndent(DefaultAll, "", "    ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(filepath.Join("theme", "default.json"), marshal, 0o644)
	if err != nil {
		panic(err)
	}
}
