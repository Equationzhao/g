package filter

import (
	"os"
	"strings"

	"github.com/Equationzhao/g/git"
	"github.com/gobwas/glob"
)

const (
	keep   = true
	remove = false
)

type TypeFilter struct {
	tfs []*TypeFunc
}

func (tf *TypeFilter) AppendTo(typeFunc ...*TypeFunc) {
	tf.tfs = append(tf.tfs, typeFunc...)
}

func NewTypeFilter(tfs ...*TypeFunc) *TypeFilter {
	return &TypeFilter{tfs: tfs}
}

func (tf *TypeFilter) Filter(e ...os.FileInfo) (res []os.FileInfo) {
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
		for _, extI := range ext {
			if strings.HasSuffix(e.Name(), "."+extI) {
				return remove
			}
		}
		return keep
	}
}

var ExtOnly = func(ext ...string) TypeFunc {
	return func(e os.FileInfo) bool {
		for _, extI := range ext {
			if strings.HasSuffix(e.Name(), "."+extI) {
				return keep
			}
		}

		return remove
	}
}

// RemoveGlob if all pattern complied successfully, return a func and nil error,
// if match any one, the fn will return false, else return false
// if error occurred, return nil func and error
var RemoveGlob = func(globPattern ...string) (TypeFunc, error) {
	compiled := make([]glob.Glob, 0, len(globPattern))
	for _, v := range globPattern {
		compile, err := glob.Compile(v)
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, compile)
	}

	return func(e os.FileInfo) bool {
		for _, r := range compiled {
			if r.Match(e.Name()) {
				return remove
			}
		}
		return keep
	}, nil
}

// GlobOnly if all pattern complied successfully, return a func and nil error,
// if match any one, the fn will return true, else return false
// if error occurred, return nil func and error
var GlobOnly = func(globPattern ...string) (TypeFunc, error) {
	compiled := make([]glob.Glob, 0, len(globPattern))
	for _, v := range globPattern {
		compile, err := glob.Compile(v)
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, compile)
	}

	return func(e os.FileInfo) bool {
		for _, r := range compiled {
			if r.Match(e.Name()) {
				return keep
			}
		}
		return remove
	}, nil
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

var RemoveGitIgnore = func(repoPath git.GitRepoPath) TypeFunc {
	isOrIsParentOf := func(parent, child string) bool {
		if parent == child {
			return true
		}
		if strings.HasPrefix(child, parent+"/") { // should not use filepath.Separator
			return true
		}
		return false
	}
	return func(e os.FileInfo) (ok bool) {
		ignoredCache := git.GetCache()
		actual, _ := ignoredCache.GetOrInit(repoPath, git.DefaultInit(repoPath))
		ok = true
		for _, fileGit := range *actual {
			if fileGit.Status == git.Ignored {
				if isOrIsParentOf(e.Name(), fileGit.Name) {
					ok = false
				}
			}
		}
		return
	}
}
