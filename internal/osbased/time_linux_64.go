//go:build (amd64 || arm64 || loong64) && linux

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
	atim := a.Sys().(*syscall.Stat_t).Atim
	return time.Unix(atim.Sec, atim.Nsec)
}

func CreateTime(a os.FileInfo) time.Time {
	ctim := a.Sys().(*syscall.Stat_t).Ctim
	return time.Unix(ctim.Sec, ctim.Nsec)
}

func BirthTime(a os.FileInfo) time.Time {
	return CreateTime(a)
}
