package config

import (
	"os"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/zeebo/assert"
)

func TestLoad(t *testing.T) {
	p := gomonkey.NewPatches()
	p.ApplyFunc(os.UserConfigDir, func() (string, error) {
		return "/home/user", nil
	}).ApplyFunc(os.MkdirAll, func(path string, perm os.FileMode) error {
		return nil
	}).ApplyFunc(os.ReadFile, func(name string) ([]byte, error) {
		return []byte(`Args:
  - hyperlink=never
  - icons
  - fuzzy

CustomTreeStyle:
  Child: "├── "
  LastChild: "╰── "
  Mid: "│   "
  Empty: "    "`), nil
	})
	defer p.Reset()

	load, err := Load()
	assert.NoError(t, err)
	assert.DeepEqual(t, load.Args, []string{"--hyperlink=never", "--icons", "--fuzzy"})
	assert.DeepEqual(t, load.CustomTreeStyle, TreeStyle{
		Child:     "├── ",
		LastChild: "╰── ",
		Mid:       "│   ",
		Empty:     "    ",
	})
}
