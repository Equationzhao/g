package content

import (
	"github.com/Equationzhao/g/filter"
	"github.com/Equationzhao/g/item"
)

const Permissions = "Permissions"

// EnableFileMode return file mode like -rwxrwxrwx/drwxrwxrwx
func EnableFileMode() filter.ContentOption {
	return func(info *item.FileInfo) (string, string) {
		return info.Mode().String(), Permissions
	}
}
