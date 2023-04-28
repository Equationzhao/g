package filter

import (
	"os"
	"regexp"
	"strings"
)

const (
	keep   = true
	remove = false
)

type TypeFilter struct {
	tfs []*TypeFunc
}

func NewTypeFilter(tfs ...*TypeFunc) *TypeFilter {
	return &TypeFilter{tfs: tfs}
}

func (tf *TypeFilter) Filter(e []os.FileInfo) (res []os.FileInfo) {
	for _, entry := range e {
		ok := true
		for _, funcPtr := range tf.tfs {
			if !(*funcPtr)(entry) {
				ok = false
				break
			}
		}
		if ok {
			res = append(res, entry)
		}
	}
	return res
}

func alwaysTrue(e os.FileInfo) bool {
	return keep
}

// TypeFunc return true -> Keep
// return false -> remove
type TypeFunc = func(e os.FileInfo) bool

var RemoveDir = func(e os.FileInfo) bool {
	return !e.IsDir()
}

var DirOnly = func(e os.FileInfo) bool {
	return e.IsDir()
}

// RemoveByExt
//
//	eg:
//		a.go b.c c.rs d.cxx dir
//		RemoveByExt([]string{"go", "cxx"})
//	result:
//		b.c c.rs dir
var RemoveByExt = func(ext ...string) TypeFunc {
	return func(e os.FileInfo) bool {
		for _, exti := range ext {
			if strings.HasSuffix(e.Name(), "."+exti) {
				return remove
			}
		}
		return keep
	}
}

var ExtOnly = func(ext ...string) TypeFunc {
	return func(e os.FileInfo) bool {
		for _, exti := range ext {
			if strings.HasSuffix(e.Name(), "."+exti) {
				return keep
			}
		}

		return remove
	}
}

var RemoveRegexp = func(regexpression string) TypeFunc {
	compiled, err := regexp.Compile(regexpression)
	if err != nil {
		return alwaysTrue
	}

	return func(e os.FileInfo) bool {
		return !compiled.Match([]byte(regexpression))
	}
}

var RegexpOnly = func(regexpression string) TypeFunc {
	compiled, err := regexp.Compile(regexpression)
	if err != nil {
		return alwaysTrue
	}

	return func(e os.FileInfo) bool {
		return compiled.Match([]byte(regexpression))
	}
}

var RemoveHidden = func(e os.FileInfo) bool {
	return !strings.HasPrefix(e.Name(), ".")
}

var HiddenOnly = func(e os.FileInfo) bool {
	return strings.HasPrefix(e.Name(), ".")
}

var RemoveBackups = func(e os.FileInfo) bool {
	return !strings.HasSuffix(e.Name(), "~")
}
