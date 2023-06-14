package filter

import (
	"os"
	"strings"

	"github.com/Equationzhao/g/git"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gobwas/glob"
)

const (
	keep   = true
	remove = false
)

type ItemFilter struct {
	tfs []*ItemFilterFunc
}

func (tf *ItemFilter) AppendTo(typeFunc ...*ItemFilterFunc) {
	tf.tfs = append(tf.tfs, typeFunc...)
}

func NewItemFilter(tfs ...*ItemFilterFunc) *ItemFilter {
	return &ItemFilter{tfs: tfs}
}

func (tf *ItemFilter) Filter(e ...os.FileInfo) (res []os.FileInfo) {
	for _, entry := range e {
		ok := keep
		for _, funcPtr := range tf.tfs {
			if !(*funcPtr)(entry) {
				ok = remove
				break
			}
		}
		if ok {
			res = append(res, entry)
		}
	}
	return res
}

// ItemFilterFunc return true -> Keep
// return false -> remove
type ItemFilterFunc = func(e os.FileInfo) bool

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
func RemoveByExt(ext ...string) ItemFilterFunc {
	return func(e os.FileInfo) bool {
		for _, extI := range ext {
			if strings.HasSuffix(e.Name(), "."+extI) {
				return remove
			}
		}
		return keep
	}
}

func ExtOnly(ext ...string) ItemFilterFunc {
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
func RemoveGlob(globPattern ...string) (ItemFilterFunc, error) {
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
func GlobOnly(globPattern ...string) (ItemFilterFunc, error) {
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

func RemoveGitIgnore(repoPath git.GitRepoPath) ItemFilterFunc {
	isOrIsParentOf := func(parent, child string) bool {
		if parent == child {
			return true
		}
		if strings.HasPrefix(child, parent+"/") { // should not use filepath.Separator
			return true
		}
		return false
	}
	ignoredCache := git.GetCache()
	return func(e os.FileInfo) (ok bool) {
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

func isOrIsSonOf(a, b string) bool {
	if a == b {
		return true
	}
	if strings.HasPrefix(a, b+"/") {
		return true
	}
	return false
}

func ExactFileTypeOnly(fileTypes ...string) ItemFilterFunc {
	return func(e os.FileInfo) bool {
		file, err := os.Open(e.Name())
		if err != nil {
			return keep
		}
		mtype, err := mimetype.DetectReader(file)
		if err != nil {
			return keep
		}

		for i := range fileTypes {
			if isOrIsSonOf(mtype.String(), fileTypes[i]) {
				return keep
			}
		}
		return remove
	}
}

func RemoveExactFileType(fileTypes ...string) ItemFilterFunc {
	return func(e os.FileInfo) bool {
		file, err := os.Open(e.Name())
		if err != nil {
			return keep
		}
		mtype, err := mimetype.DetectReader(file)
		if err != nil {
			return keep
		}

		for i := range fileTypes {
			if fileTypes[i] == mtype.String() {
				return remove
			}
		}
		return keep
	}
}
