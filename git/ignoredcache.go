package git

import (
	"sync"

	"github.com/Equationzhao/g/cached"
	"github.com/Equationzhao/pathbeautify"
)

var (
	ignored          *cached.Map[RepoPath, *FileGits]
	IgnoredInitOnce  sync.Once
	TopLevelCache    *cached.Map[RepoPath, RepoPath]
	TopLevelInitOnce sync.Once
)

const shardSize = 20

type Cache = *cached.Map[RepoPath, *FileGits]

func GetCache() Cache {
	IgnoredInitOnce.Do(
		func() {
			ignored = cached.NewCacheMap[RepoPath, *FileGits](shardSize)
		},
	)
	return ignored
}

func FreeCache() {
	ignored.Free()
}

func DefaultInit(repoPath RepoPath) func() *FileGits {
	return func() *FileGits {
		res := make(FileGits, 0)
		out, err := GetShortGitStatus(repoPath)
		if err == nil {
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
				TopLevelCache = cached.NewCacheMap[RepoPath, RepoPath](shardSize)
			}
		},
	)
	var err error
	actual, _ := TopLevelCache.GetOrInit(
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
