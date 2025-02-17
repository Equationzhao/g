//go:build linux

package osbased

import (
	"os"
	"strconv"
	"syscall"
)

func Inode(info os.FileInfo) string {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if ok {
		return strconv.FormatUint(stat.Ino, 10)
	}
	return ""
}

func LinkCount(info os.FileInfo) uint64 {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if ok {
		return uint64(stat.Nlink)
	}
	return 0
}

func BlockSize(info os.FileInfo) int64 {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return 0
	}

	return stat.Blocks
}

// always false on Linux
func IsMacOSAlias(_ string) bool {
	return false
}

// ResolveAlias is a no-op on Linux.
func ResolveAlias(_ string) (string, error) {
	return "", nil
}
