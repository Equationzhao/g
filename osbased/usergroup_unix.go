//go:build unix

package osbased

import (
	"os"
	"strconv"
	"syscall"

	"github.com/Equationzhao/g/cached"
)

func GroupID(a os.FileInfo) string {
	return strconv.FormatInt(int64(a.Sys().(*syscall.Stat_t).Gid), 10)
}

func Group(a os.FileInfo) string {
	return cached.GetGroupname(strconv.FormatInt(int64(a.Sys().(*syscall.Stat_t).Gid), 10))
}

func OwnerID(a os.FileInfo) string {
	return strconv.FormatInt(int64(a.Sys().(*syscall.Stat_t).Uid), 10)
}

func Owner(a os.FileInfo) string {
	return cached.GetUsername(strconv.FormatInt(int64(a.Sys().(*syscall.Stat_t).Uid), 10))
}
