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
	ctim := a.Sys().(*syscall.Win32FileAttributeData).CreationTime
	return time.Unix(0, ctim.Nanoseconds())
}

func CreateTime(a os.FileInfo) time.Time {
	atim := a.Sys().(*syscall.Win32FileAttributeData).LastAccessTime
	return time.Unix(0, atim.Nanoseconds())
}
