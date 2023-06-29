package content

import (
	"github.com/Equationzhao/g/filter"
	"github.com/Equationzhao/g/item"
	"github.com/Equationzhao/g/render"
)

const Permissions = "Permissions"

// EnableFileMode return file mode like -rwxrwxrwx/drwxrwxrwx
func EnableFileMode(renderer *render.Renderer) filter.ContentOption {
	return func(info *item.FileInfo) (string, string) {
		return renderer.FileMode(info.Mode().String()), Permissions
	}
}
