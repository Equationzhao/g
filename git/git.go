package git

import (
	"bufio"
	"os/exec"
	"strings"
)

type FileGit struct {
	Name   string
	Status Status
}

/*
X          Y     Meaning
-------------------------------------------------

	[AMD]   not updated

M        [ MTD]  updated in index
T        [ MTD]  type changed in index
A        [ MTD]  added to index
D                deleted from index
R        [ MTD]  renamed in index
C        [ MTD]  copied in index
[MTARC]          index and work tree matches
[ MTARC]    M    work tree changed since index
[ MTARC]    T    type changed in work tree since index
[ MTARC]    D    deleted in work tree

//			R    renamed in work tree
//			C    copied in work tree

-------------------------------------------------
D           D    unmerged, both deleted
A           U    unmerged, added by us
U           D    unmerged, deleted by them
U           A    unmerged, added by them
D           U    unmerged, deleted by us
A           A    unmerged, both added
U           U    unmerged, both modified
-------------------------------------------------
?           ?    untracked
!           !    ignored
-------------------------------------------------
*/
func (f *FileGit) setYFromXY(XY string) {
	set := func(Y string) {
		switch Y {
		case "M":
			f.Status = GitModified
		case "A":
			f.Status = GitAdded
		case "D":
			f.Status = GitDeleted
		case "R":
			f.Status = GitRenamed
		case "C":
			f.Status = GitCopied
		case "?":
			f.Status = GitUntracked
		case "!":
			f.Status = GitIgnored
		}
	}

	switch len(XY) {
	case 1:
		set(XY)
	case 2:
		Y := XY[0:1]
		set(Y)
	default:
		return
	}
}

type FileGits = []FileGit

type GitRepoPath = string

// GetShortGitStatus read the git status of the repository located at path
func GetShortGitStatus(repoPath GitRepoPath) (string, error) {
	out, err := exec.Command("git", "-C", repoPath, "status", "-s", "--ignored", "--porcelain").Output()
	return string(out), err
}

type Status int

const (
	GitModified  Status = iota + 1 // M ~
	GitAdded                       // A +
	GitDeleted                     // D -
	GitRenamed                     // R |
	GitCopied                      // C =
	GitUntracked                   // ? ?
	GitIgnored                     // ! !
)

// ParseShort parses a git status output command
// It is compatible with the short version of the git status command
// modified from https://le-gall.bzh/post/go/parsing-git-status-with-go/ author: Sébastien Le Gall
func ParseShort(r string) (res FileGits) {
	s := bufio.NewScanner(strings.NewReader(r))

	// Extract branch name
	for s.Scan() {
		// Skip any empty line
		if len(s.Text()) < 1 {
			continue
		}
		break
	}

	fg := FileGit{}
	for true {
		if len(s.Text()) < 1 {
			continue
		}
		XyName := strings.Fields(s.Text())
		fg.setYFromXY(XyName[0])
		fg.Name = XyName[1]
		res = append(res, fg)
		if !s.Scan() {
			break
		}
	}

	return
}

func (s Status) String() string {
	switch s {
	case GitModified:
		return "~"
	case GitAdded:
		return "+"
	case GitDeleted:
		return "-"
	case GitRenamed:
		return "|"
	case GitCopied:
		return "="
	case GitUntracked:
		return "?"
	case GitIgnored:
		return "!"
	}
	return ""
}
