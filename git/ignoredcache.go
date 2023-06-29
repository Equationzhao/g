package git

import (
	"sync"

	"github.com/Equationzhao/g/cached"
)

var (
	ignored         *cached.Map[RepoPath, *FileGits]
	IgnoredInitOnce sync.Once
)

const shardSize = 20

func GetCache() *cached.Map[RepoPath, *FileGits] {
	IgnoredInitOnce.Do(func() {
		ignored = cached.NewCacheMap[RepoPath, *FileGits](shardSize)
	})
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
