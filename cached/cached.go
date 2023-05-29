package cached

import (
	"github.com/Equationzhao/pathbeautify"
	"github.com/Equationzhao/tsmap"
)

func GetUserHomeDir() string {
	return pathbeautify.GetUserHomeDir()
}

type Map[k comparable, v any] struct {
	*tsmap.Map[k, v]
}

func NewCacheMap[k comparable, v any](len int) *Map[k, v] {
	return &Map[k, v]{
		Map: tsmap.NewTSMap[k, v](len),
	}
}
