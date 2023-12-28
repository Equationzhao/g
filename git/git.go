package git

import (
	"bufio"
	"os/exec"
	"strings"

	"github.com/Equationzhao/g/util"
	"github.com/Equationzhao/pathbeautify"
	"github.com/go-git/go-git/v5"
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
	if err == nil {
		return string(out), err
	}

	// if failed, try go-git
	return goGitStatus(repoPath)
}

func goGitStatus(repoPath RepoPath) (string, error) {
	r, err := git.PlainOpenWithOptions(repoPath, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return "", err
	}
	w, err := r.Worktree()
	if err != nil {
		return "", err
	}
	status, err := w.Status()
	if err != nil {
		return "", err
	}
	return status.String(), nil
}

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
	r, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return "", err
	}
	w, err := r.Worktree()
	if err != nil {
		return "", err
	}
	return w.Filesystem.Root(), nil
}

type Status uint8

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

func (r RepoStatus) String() string {
	switch r {
	case RepoStatusClean:
		return "+"
	case RepoStatusDirty:
		return "|"
	case RepoStatusSkip:
		return ""
	}
	return ""
}

const (
	RepoStatusSkip RepoStatus = iota
	RepoStatusClean
	RepoStatusDirty
)

type RepoStatus uint8

// GetBranch returns the branch of the repository
// only return the branch when the path is the root of the repository
func GetBranch(repoPath RepoPath) string {
	if root, _ := getTopLevel(repoPath); root != repoPath {
		return ""
	}

	c := exec.Command("git", "branch", "--show-current")
	c.Dir = repoPath
	out, err := c.Output()
	if err == nil {
		return strings.TrimSpace(string(out))
	}
	return goBranch(repoPath)
}

func goBranch(repoPath RepoPath) string {
	r, err := git.PlainOpenWithOptions(repoPath, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return ""
	}
	ref, err := r.Head()
	if err != nil {
		return ""
	}
	return ref.Name().Short()
}

// GetRepoStatus returns the status of the repository
// only return the status when the path is the root of the repository
func GetRepoStatus(repoPath RepoPath) RepoStatus {
	if root, _ := getTopLevel(repoPath); root != repoPath {
		return RepoStatusSkip
	}

	c := exec.Command("git", "status", "--porcelain")
	c.Dir = repoPath
	out, err := c.Output()
	if err == nil {
		if len(out) == 0 {
			return RepoStatusClean
		}
		return RepoStatusDirty
	}
	return goRepoStatus(repoPath)
}

func goRepoStatus(repoPath RepoPath) RepoStatus {
	r, err := git.PlainOpenWithOptions(repoPath, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return RepoStatusSkip
	}
	w, err := r.Worktree()
	if err != nil {
		return RepoStatusSkip
	}
	status, err := w.Status()
	if err != nil {
		return RepoStatusSkip
	}
	if status.IsClean() {
		return RepoStatusClean
	}
	return RepoStatusDirty
}
