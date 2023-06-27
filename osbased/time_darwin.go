//go:build darwin

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
	return time.Unix(atim.Sec, atim.Nsec)
}

func CreateTime(a os.FileInfo) time.Time {
	ctim := a.Sys().(*syscall.Stat_t).Ctimespec
	return time.Unix(ctim.Sec, ctim.Nsec)
}
