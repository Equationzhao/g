package osbased

import (
	"github.com/Equationzhao/g/internal/item"
	"golang.org/x/sys/unix"
	"os"
	"syscall"
)

type flagDescriptions struct {
	flag int
	name string
}

var flags = []flagDescriptions{
	{unix.UF_NODUMP, "UF_NODUMP"},
	{unix.UF_IMMUTABLE, "UF_IMMUTABLE"},
	{unix.UF_APPEND, "UF_APPEND"},
	{unix.UF_OPAQUE, "UF_OPAQUE"},
	{unix.UF_HIDDEN, "UF_HIDDEN"},
	{unix.SF_ARCHIVED, "SF_ARCHIVED"},
	{unix.SF_IMMUTABLE, "SF_IMMUTABLE"},
	{unix.SF_APPEND, "SF_APPEND"},
}

func getFlags(filename string) uint32 {
	file, err := os.Open(filename)
	if err != nil {
		return 0
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return 0
	}

	stat := fileInfo.Sys().(*syscall.Stat_t)
	return stat.Flags
}

func CheckFlags(i *item.FileInfo) []string {
	res := make([]string, 0, 8)
	f := getFlags(i.FullPath)
	for _, flag := range flags {
		if f&uint32(flag.flag) != 0 {
			res = append(res, flag.name)
		}
	}
	return res
}
