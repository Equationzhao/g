package cached

import (
	"os"
	"sync"

	"github.com/Equationzhao/tsmap"
)

var (
	syncHomedir sync.Once
	userHomeDir string
)

func GetUserHomeDir() string {
	syncHomedir.Do(func() {
		userHomeDir, _ = os.UserHomeDir()
	})
	return userHomeDir
}

type Map[k comparable, v any] struct {
	*tsmap.Map[k, v]
}

func NewCacheMap[k comparable, v any](len int) *Map[k, v] {
	return &Map[k, v]{
		Map: tsmap.NewTSMap[k, v](len),
	}
}
