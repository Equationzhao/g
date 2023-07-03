package content

import (
	"github.com/Equationzhao/g/filter"
	"github.com/Equationzhao/g/item"
	"github.com/Equationzhao/g/render"
	"strconv"
)

const Permissions = "Permissions"

// EnableFileMode return file mode like -rwxrwxrwx/drwxrwxrwx
func EnableFileMode(renderer *render.Renderer) filter.ContentOption {
	return func(info *item.FileInfo) (string, string) {
		return renderer.FileMode(filter.FillBlank(info.Mode().String(), 11)), Permissions
	}
}

const OctalPermissions = "Octal"

func EnableFileOctalPermissions(renderer *render.Renderer) filter.ContentOption {
	return func(info *item.FileInfo) (string, string) {
		return renderer.OctalPerm(
			" 0" + strconv.FormatUint(uint64(info.Mode().Perm()), 8),
		), OctalPermissions
	}
}
