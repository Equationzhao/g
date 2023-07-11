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
	DefaultAll.Apply(convert)
	marshal, err := json.MarshalIndent(DefaultAll, "", "    ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(filepath.Join("theme", "default.json"), marshal, 0o644)
	if err != nil {
		panic(err)
	}
}
