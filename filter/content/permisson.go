package content

import (
	"github.com/Equationzhao/g/filter"
	"github.com/Equationzhao/g/render"
	"os"
)

const Permissions = "Permissions"

// EnableFileMode return file mode like -rwxrwxrwx/drwxrwxrwx
func EnableFileMode(renderer *render.Renderer) filter.ContentOption {
	return func(info os.FileInfo) (string, string) {
		return renderer.FileMode(filter.FillBlank(info.Mode().String(), 12)), Permissions
	}
}
