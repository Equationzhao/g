package content

import (
	"strconv"

	"github.com/Equationzhao/g/internal/align"
	"github.com/Equationzhao/g/internal/global"
	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/render"
	"github.com/pkg/xattr"
)

const Permissions = global.NameOfPermission

// EnableFileMode return file mode like -rwxrwxrwx/drwxrwxrwx
func EnableFileMode(renderer *render.Renderer) ContentOption {
	align.Register(Permissions)
	return func(info *item.FileInfo) (string, string) {
		perm := renderer.FileMode(info.Mode().String())
		if info.Cache[Extended] != nil {
			perm += "@"
		}
		list, _ := xattr.LList(info.FullPath)
		info.Cache[Extended] = list
		if len(list) != 0 {
			perm += "@"
		}
		return perm, Permissions
	}
}

const OctalPermissions = "Octal"

func EnableFileOctalPermissions(renderer *render.Renderer) ContentOption {
	return func(info *item.FileInfo) (string, string) {
		return renderer.OctalPerm(
			"0" + strconv.FormatUint(uint64(info.Mode().Perm()), 8),
		), OctalPermissions
	}
}
