package content

import (
	"path/filepath"
	"strings"

	"github.com/Equationzhao/g/align"
	"github.com/Equationzhao/g/filter"
	"github.com/Equationzhao/g/git"
	"github.com/Equationzhao/g/item"
	"github.com/Equationzhao/g/render"
)

type GitEnabler struct {
	cache git.Cache
	Path  git.RepoPath
}

func (g *GitEnabler) InitCache(repo git.RepoPath) {
	g.cache.Set(repo, git.DefaultInit(repo)())
}

func NewGitEnabler() *GitEnabler {
	return &GitEnabler{
		cache: git.GetCache(),
	}
}

const GitStatus = "Git"

func (g *GitEnabler) Enable(renderer *render.Renderer) filter.ContentOption {
	isOrIsParentOf := func(parent, child string) bool {
		if parent == child {
			return true
		}
		if strings.HasPrefix(child, parent+string(filepath.Separator)) {
			return true
		}
		return false
	}

	return func(info *item.FileInfo) (string, string) {
		gits, ok := g.cache.Get(g.Path)
		if ok {
			topLevel, err := git.GetTopLevel(g.Path)
			if err != nil {
				return gitByName(git.Unmodified, renderer) + gitByName(git.Unmodified, renderer), GitStatus
			}
			rel, err := filepath.Rel(topLevel, info.FullPath)
			if err != nil {
				return gitByName(git.Unmodified, renderer) + gitByName(git.Unmodified, renderer), GitStatus
			}
			for _, status := range *gits {
				if status.X == git.Ignored || status.Y == git.Ignored {
					// if status is ignored,
					// and the file is or is a child of the ignored file
					if isOrIsParentOf(status.Name, rel) {
						return gitByName(status.X, renderer) + gitByName(status.Y, renderer), GitStatus
					}
				} else {
					if isOrIsParentOf(rel, status.Name) {
						return gitByName(status.X, renderer) + gitByName(status.Y, renderer), GitStatus
					}
				}
			}
		}
		return gitByName(git.Unmodified, renderer) + gitByName(git.Unmodified, renderer), GitStatus
	}
}

func gitByName(status git.Status, renderer *render.Renderer) string {
	switch status {
	case git.Unmodified:
		return renderer.GitUnmodified("-")
	case git.Modified:
		return renderer.GitModified("M")
	case git.Added:
		return renderer.GitAdded("A")
	case git.Deleted:
		return renderer.GitDeleted("D")
	case git.Renamed:
		return renderer.GitRenamed("R")
	case git.Copied:
		return renderer.GitCopied("C")
	case git.Untracked:
		return renderer.GitUntracked("?")
	case git.Ignored:
		return renderer.GitIgnored("!")
	case git.TypeChanged:
		return renderer.GitTypeChanged("T")
	case git.UpdatedBuUnmerged:
		return renderer.GitUpdatedButUnmerged("U")
	}
	return " "
}

const GitRepoBranch = "Branch"
const GitRepoStatus = "Repo-status"

type GitRepoEnabler struct{}

func (g *GitRepoEnabler) Enable(renderer *render.Renderer) filter.ContentOption {
	align.Register(GitRepoBranch)
	return func(info *item.FileInfo) (string, string) {
		// get branch name
		return renderer.GitRepoBranch(git.GetBranch(info.FullPath)), GitRepoBranch
	}
}

func (g *GitRepoEnabler) EnableStatus(renderer *render.Renderer) filter.ContentOption {
	align.Register(GitRepoStatus)
	return func(info *item.FileInfo) (string, string) {
		// get repo status
		return renderer.GitRepoStatus(git.GetRepoStatus(info.FullPath)), GitRepoStatus
	}
}

func NewGitRepoEnabler() *GitRepoEnabler {
	return &GitRepoEnabler{}
}
