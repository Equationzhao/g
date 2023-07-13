package git

import (
	"bufio"
	"os/exec"
	"strings"

	"github.com/Equationzhao/g/util"
	"github.com/Equationzhao/pathbeautify"
)

// FileGit is an entry name with git status
// the name will not end with file separator
type FileGit struct {
	Name string
	X, Y Status
}

/*
Set sets the status of the file based on the XY string
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
func (f *FileGit) Set(XY string) {
	set := func(s *Status, c byte) {
		*s = Byte2Status(c)
	}
	set(&f.X, XY[0])
	set(&f.Y, XY[1])
}

type FileGits = []FileGit

type RepoPath = string

// GetShortGitStatus read the git status of the repository located at the path
func GetShortGitStatus(repoPath RepoPath) (string, error) {
	c := exec.Command("git", "status", "-s", "--ignored", "--porcelain", repoPath)
	c.Dir = repoPath
	out, err := c.Output()
	return string(out), err
}

func getTopLevel(path RepoPath) (string, error) {
	c := exec.Command("git", "rev-parse", "--show-toplevel", path)
	c.Dir = path
	out, err := c.Output()
	if err != nil {
		return "", err
	}
	// 	 get the first line
	lines := strings.Split(string(out), "\n")[0]
	return lines, err
}

type Status int

const (
	Unknown           Status = iota
	Unmodified               // -
	Modified                 // M
	Added                    // A
	Deleted                  // D
	Renamed                  // R
	Copied                   // C
	Untracked                // ?
	Ignored                  // !
	TypeChanged              // T
	UpdatedBuUnmerged        // U
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
		fg.Set(status)
		if fg.X == Renamed || fg.Y == Renamed || fg.X == Copied || fg.Y == Copied {
			// origin -> rename
			// the actual file name is rename
			o2r := str[3:]
			fg.Name = util.RemoveSep(pathbeautify.CleanSeparator(o2r[strings.Index(o2r, " -> ")+4:]))
		} else {
			fg.Name = util.RemoveSep(pathbeautify.CleanSeparator(str[3:]))
		}

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
		return "M"
	case Added:
		return "A"
	case Deleted:
		return "D"
	case Renamed:
		return "R"
	case Copied:
		return "C"
	case Untracked:
		return "?"
	case Ignored:
		return "!"
	case Unmodified:
		return "-"
	case TypeChanged:
		return "T"
	case UpdatedBuUnmerged:
		return "U"
	case Unknown:
		return "^"
	}
	return "^"
}

func Byte2Status(c byte) Status {
	switch c {
	case 'M':
		return Modified
	case 'A':
		return Added
	case 'D':
		return Deleted
	case 'R':
		return Renamed
	case 'C':
		return Copied
	case '?':
		return Untracked
	case '!':
		return Ignored
	case '-', ' ':
		return Unmodified
	case 'T':
		return TypeChanged
	case 'U':
		return UpdatedBuUnmerged
	case '^':
		return Unknown
	}
	return Unknown
}
