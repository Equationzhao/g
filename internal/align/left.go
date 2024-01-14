package align

import (
	"github.com/Equationzhao/g/internal/util"
)

// left
// the default is right align
// field to align left should register here
var left = util.NewSet[string]()

func Register(name string) {
	left.Add(name)
}

func IsLeft(name string) bool {
	return left.Contains(name)
}
