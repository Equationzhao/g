package git

import (
	"bufio"
	"os/exec"
	"strings"
)

type FileGit struct {
	Name string
	X, Y Status
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
func (f *FileGit) set(XY string) {
	set := func(s *Status, c byte) {
		switch c {
		case 'M':
			*s = Modified
		case 'A':
			*s = Added
		case 'D':
			*s = Deleted
		case 'R':
			*s = Renamed
		case 'C':
			*s = Copied
		case '?':
			*s = Untracked
		case '!':
			*s = Ignored
		}
	}
	set(&f.X, XY[0])
	set(&f.Y, XY[1])
}

type FileGits = []FileGit

type RepoPath = string

// GetShortGitStatus read the git status of the repository located at path
func GetShortGitStatus(repoPath RepoPath) (string, error) {
	out, err := exec.Command("git", "-C", repoPath, "status", "-s", "--ignored", "--porcelain").Output()
	return string(out), err
}

type Status int

const (
	Modified  Status = iota + 1 // M ~
	Added                       // A +
	Deleted                     // D -
	Renamed                     // R |
	Copied                      // C =
	Untracked                   // ? ?
	Ignored                     // ! !
)

// ParseShort parses a git status output command
// It is compatible with the short version of the git status command
// modified from https://le-gall.bzh/post/go/parsing-git-status-with-go/ author: Sébastien Le Gall
func ParseShort(r string) (res FileGits) {
	s := bufio.NewScanner(strings.NewReader(r))

	for s.Scan() {
		// Skip any empty line
		if len(s.Text()) < 1 {
			continue
		}
		break
	}

	fg := FileGit{}
	for {
		str := s.Text()
		if len(str) < 1 {
			continue
		}
		status := str[0:2]
		fg.set(status)
		fg.Name = str[2:]
		res = append(res, fg)
		if !s.Scan() {
			break
		}
	}

	return
}

func (s Status) String() string {
	switch s {
	case Modified:
		return "~"
	case Added:
		return "+"
	case Deleted:
		return "-"
	case Renamed:
		return "|"
	case Copied:
		return "="
	case Untracked:
		return "?"
	case Ignored:
		return "!"
	}
	return ""
}
