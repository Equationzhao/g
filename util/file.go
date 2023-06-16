package util

import "os"

func IsExecutable(file os.FileInfo) bool {
	return file.Mode()&0o111 != 0
}

func IsExecutableMode(mode os.FileMode) bool {
	return mode&0o111 != 0
}
