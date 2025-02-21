//go:build darwin

package osbased

/*
#cgo CFLAGS: -mmacosx-version-min=10.9
#cgo LDFLAGS: -framework CoreFoundation -framework CoreServices
#include "macos_alias.h"
*/
import "C"

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"unsafe"
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

func IsMacOSAlias(fullPath string) bool {
	fi, err := os.Lstat(fullPath)
	if err != nil {
		return false
	}

	if fi.Mode()&os.ModeSymlink != 0 {
		return false
	}

	cPath := C.CString(fullPath)
	defer C.free(unsafe.Pointer(cPath))
	return bool(C.isAlias(cPath))
}

func ResolveAlias(fullPath string) (string, error) {
	cPath := C.CString(fullPath)
	defer C.free(unsafe.Pointer(cPath))

	resolved := C.resolveAlias(cPath)
	if resolved == nil {
		return "", fmt.Errorf("failed to resolve macOS alias for %s", fullPath)
	}
	defer C.free(unsafe.Pointer(resolved))

	return C.GoString(resolved), nil
}
