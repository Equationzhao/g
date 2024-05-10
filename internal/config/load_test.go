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

func TestTreeStyle_IsEnabled(t1 *testing.T) {
	type fields struct {
		Child     string
		LastChild string
		Mid       string
		Empty     string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "test",
			fields: fields{
				Child:     "├── ",
				LastChild: "╰── ",
				Mid:       "│   ",
				Empty:     "    ",
			},
			want: true,
		},
		{
			name: "test",
			fields: fields{
				Child:     "",
				LastChild: "",
				Mid:       "",
				Empty:     "",
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			ts := TreeStyle{
				Child:     tt.fields.Child,
				LastChild: tt.fields.LastChild,
				Mid:       tt.fields.Mid,
				Empty:     tt.fields.Empty,
			}
			if got := ts.IsEnabled(); got != tt.want {
				t1.Errorf("IsEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}
