package git

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	strftime "github.com/itchyny/timefmt-go"
)

type CommitInfo struct {
	Hash string `json:"h"`

	Comitter     string `json:"c"`
	ComitterDate string `json:"cd"`

	Author     string `json:"a"`
	AuthorDate string `json:"ad"`
}

func (c CommitInfo) GetCommiterDateInFormat(format string) string {
	t, err := time.Parse(time.RFC3339, c.ComitterDate)
	if err != nil {
		return ""
	}
	return t.Format(format)
}

func (c CommitInfo) GetAuthorDateInFormat(format string) string {
	t, err := time.Parse(goParseFormat, c.AuthorDate)
	if err != nil {
		return ""
	}
	if strings.HasPrefix(format, "+") {
		return strftime.Format(t, strings.TrimPrefix(format, "+"))
	}
	return t.Format(format)
}

var NoneCommitInfo = CommitInfo{"-", "-", "-", "-", "-"}

// https://github.com/chaqchase/lla/blob/main/plugins/last_git_commit/src/lib.rs
func GetLastCommitInfo(path string) (*CommitInfo, error) {
	return getLastCommitInfo(path)
}

const (
	gitDateFormat = `format:"%Y-%m-%d %H:%M:%S.%9N %z"`
	goParseFormat = time.RFC3339
)

func getLastCommitInfo(path string) (*CommitInfo, error) {
	cmd := exec.Command("git", "log", "-1", `--pretty=format:{"h":"%h","a":"%an","c":"%cn","ad":"%aI","cd":"%cI"}`, fmt.Sprintf(`--date=%s`, gitDateFormat), path)

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	if len(output) == 0 {
		return &NoneCommitInfo, nil
	}
	var info CommitInfo
	if err := json.Unmarshal(output, &info); err != nil {
		return nil, err
	}

	return &info, nil
}

// // todo goGetLastCommitInfo using go-git to get the last commit info
// func goGetLastCommitInfo(filePath string) (*CommitInfo, error) {
// 	repo, err := goGitOpenWithCache(filePath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	commitIter, err := repo.Log(&git.LogOptions{FileName: &filePath})
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer commitIter.Close()

// 	commit, err := commitIter.Next()
// 	if err != nil {
// 		return nil, err
// 	}

// 	relativeTime := formatRelativeTime(commit.Committer.When)
// 	return &CommitInfo{
// 		Hash:   commit.Hash.String(),
// 		Author: commit.Author.Name,
// 		Date:   relativeTime,
// 	}, nil
// }

// func formatRelativeTime(commitTime time.Time) string {
// 	return durafmt.Parse(time.Since(commitTime)).LimitFirstN(1).String()
// }
