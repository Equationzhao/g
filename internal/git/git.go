package git

import (
	"os/exec"
	"strings"
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
	return "", err
}
