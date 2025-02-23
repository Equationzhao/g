package git

import (
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5"
)

func getTopLevel(path RepoPath) (string, error) {
	c := exec.Command("git", "rev-parse", "--show-toplevel", path)
	c.Dir = path
	out, err := c.Output()
	if err == nil {
		// get the first line
		lines := strings.Split(string(out), "\n")
		if len(lines) > 0 {
			return lines[0], nil
		}
	}

	// if failed, try go-git
	return goGitTopLevel(path)
}

func goGitTopLevel(path RepoPath) (string, error) {
	r, err := goGitOpenWithCache(path)
	if err != nil {
		return "", err
	}
	w, err := r.Worktree()
	if err != nil {
		return "", err
	}
	return w.Filesystem.Root(), nil
}

var goGitRepoCache = make(map[RepoPath]any)

func goGitOpenWithCache(path RepoPath) (*git.Repository, error) {
	if repo, ok := goGitRepoCache[path]; ok {
		if err, ok := repo.(error); ok {
			return nil, err
		}
		return repo.(*git.Repository), nil
	}

	repo, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		goGitRepoCache[path] = err
		return nil, err
	}
	goGitRepoCache[path] = repo
	return repo, nil
}
