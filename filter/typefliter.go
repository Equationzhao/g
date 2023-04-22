package filter

import (
	"os"
	"regexp"
	"strings"
)

type TypeFilter struct {
	tfs []*TypeFunc
}

func NewTypeFilter(tfs ...*TypeFunc) *TypeFilter {
	return &TypeFilter{tfs: tfs}
}

func (tf *TypeFilter) Filter(e []os.DirEntry) (res []os.DirEntry) {
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

func alwaysTrue(e os.DirEntry) bool {
	return true
}

// TypeFunc return true -> Keep
// return false -> remove
type TypeFunc = func(e os.DirEntry) bool

var RemoveDir = func(e os.DirEntry) bool {
	return !e.IsDir()
}

var DirOnly = func(e os.DirEntry) bool {
	return e.IsDir()
}

// RemoveByExt
//
//	eg:
//		a.go b.c c.rs d.cxx dir
//		RemoveByExt("go")
//	result:
//		b.c c.rs d.cxx dir
var RemoveByExt = func(ext string) TypeFunc {
	return func(e os.DirEntry) bool {
		return !strings.HasSuffix(e.Name(), "."+ext)
	}
}

var RemoveRegexp = func(regexpression string) TypeFunc {
	compiled, err := regexp.Compile(regexpression)
	if err != nil {
		return alwaysTrue
	}

	return func(e os.DirEntry) bool {
		return !compiled.Match([]byte(regexpression))
	}
}

var RegexpOnly = func(regexpression string) TypeFunc {
	compiled, err := regexp.Compile(regexpression)
	if err != nil {
		return alwaysTrue
	}

	return func(e os.DirEntry) bool {
		return compiled.Match([]byte(regexpression))
	}
}

var RemoveHidden = func(e os.DirEntry) bool {
	return !strings.HasPrefix(e.Name(), ".")
}

func HiddenOnly(e os.DirEntry) bool {
	return strings.HasPrefix(e.Name(), ".")
}
