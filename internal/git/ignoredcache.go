package git

import (
	"sync"

	"github.com/Equationzhao/g/internal/cached"
	"github.com/Equationzhao/pathbeautify"
	"github.com/zeebo/xxh3"
)

var (
	ignored          *cached.Map[RepoPath, *FileGits]
	IgnoredInitOnce  sync.Once
	TopLevelCache    *cached.Map[RepoPath, RepoPath]
	TopLevelInitOnce sync.Once
)

const size = 20

// your custom hash function
func hasher(s string) uintptr {
	return uintptr(xxh3.HashString(s))
}

type Cache = *cached.Map[RepoPath, *FileGits]

func GetCache() Cache {
	IgnoredInitOnce.Do(
		func() {
			ignored = cached.NewCacheMap[RepoPath, *FileGits](size)
			ignored.SetHasher(hasher)
		},
	)
	return ignored
}

func DefaultInit(repoPath RepoPath) func() *FileGits {
	return func() *FileGits {
		res := make(FileGits, 0)
		out, err := GetShortGitStatus(repoPath)
		if err == nil && out != "" {
			res = ParseShort(out)
		}
		return &res
	}
}

// GetTopLevel returns the top level of the repoPath
// the returned path is cleaned by pathbeautify.CleanSeparator
func GetTopLevel(path string) (RepoPath, error) {
	TopLevelInitOnce.Do(
		func() {
			if TopLevelCache == nil {
				TopLevelCache = cached.NewCacheMap[RepoPath, RepoPath](size)
				TopLevelCache.SetHasher(hasher)
			}
		},
	)
	var err error
	actual, _ := TopLevelCache.GetOrCompute(
		path, func() RepoPath {
			out, err_ := getTopLevel(path)
			if err_ != nil {
				err = err_
				return ""
			}
			return out
		},
	)
	actual = pathbeautify.CleanSeparator(actual)
	return actual, err
}
