//go:build freebsd

package osbased

import (
	"os"
	"syscall"
	"time"
)

func ModTime(a os.FileInfo) time.Time {
	return a.ModTime()
}

func AccessTime(a os.FileInfo) time.Time {
	atim := a.Sys().(*syscall.Stat_t).Atimespec
	return time.Unix(int64(atim.Sec), int64(atim.Nsec))
}

func CreateTime(a os.FileInfo) time.Time {
	ctim := a.Sys().(*syscall.Stat_t).Ctimespec
	return time.Unix(int64(ctim.Sec), int64(ctim.Nsec))
}

func BirthTime(a os.FileInfo) time.Time {
	btim := a.Sys().(*syscall.Stat_t).Birthtimespec
	return time.Unix(int64(btim.Sec), int64(btim.Nsec))
}
