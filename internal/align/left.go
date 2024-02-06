package align

import (
	"github.com/Equationzhao/g/internal/util"
)

// left
// the default is right align
// field to align left should register here
var left = util.NewSet[string]()

func Register(names ...string) {
	for _, name := range names {
		left.Add(name)
	}
}

func IsLeft(name string) bool {
	return left.Contains(name)
}

var leftHeaderFooter = util.NewSet[string]()

func RegisterHeaderFooter(names ...string) {
	for _, name := range names {
		leftHeaderFooter.Add(name)
	}
}

func IsLeftHeaderFooter(name string) bool {
	return leftHeaderFooter.Contains(name)
}
