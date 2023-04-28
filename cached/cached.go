package cached

import (
	"os"
	"sync"
)

var syncHomedir sync.Once
var userHomeDir string

func GetUserHomeDir() string {
	syncHomedir.Do(func() {
		userHomeDir, _ = os.UserHomeDir()
	})
	return userHomeDir
}
