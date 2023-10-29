package content

import (
	"os"
	"strconv"

	"github.com/Equationzhao/g/align"
	"github.com/Equationzhao/g/filter"
	"github.com/Equationzhao/g/item"
	"github.com/Equationzhao/g/render"
	"github.com/pkg/xattr"
)

const Permissions = "Permissions"

// EnableFileMode return file mode like -rwxrwxrwx/drwxrwxrwx
func EnableFileMode(renderer *render.Renderer) filter.ContentOption {
	align.Register(Permissions)
	return func(info *item.FileInfo) (string, string) {
		perm := renderer.FileMode(info.Mode().String())
		f, err := os.Open(info.FullPath)
		if err == nil {
			list, _ := xattr.FList(f)
			if len(list) != 0 {
				perm += "@"
			}
		}
		return perm, Permissions
	}
}

const OctalPermissions = "Octal"

func EnableFileOctalPermissions(renderer *render.Renderer) filter.ContentOption {
	return func(info *item.FileInfo) (string, string) {
		return renderer.OctalPerm(
			"0" + strconv.FormatUint(uint64(info.Mode().Perm()), 8),
		), OctalPermissions
	}
}
