package osbased

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

func Inode(info os.FileInfo) string {
	return "-"
}

var (
	kernel32                   = syscall.NewLazyDLL("kernel32.dll")
	getFileInformationByHandle = kernel32.NewProc("GetFileInformationByHandle")
)

type BY_HANDLE_FILE_INFORMATION struct {
	FileAttributes     uint32
	CreationTime       syscall.Filetime
	LastAccessTime     syscall.Filetime
	LastWriteTime      syscall.Filetime
	VolumeSerialNumber uint32
	FileSizeHigh       uint32
	FileSizeLow        uint32
	NumberOfLinks      uint32
	FileIndexHigh      uint32
	FileIndexLow       uint32
}

func GetNumberOfHardLinks(info os.FileInfo) (uint64, error) {
	path := info.Name()
	utf16PtrFromString, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return 0, err
	}
	handle, err := syscall.CreateFile(
		utf16PtrFromString,
		0,
		0,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_FLAG_BACKUP_SEMANTICS,
		0,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %v", err)
	}
	defer func() {
		_ = syscall.CloseHandle(handle)
	}()

	var fileInfo BY_HANDLE_FILE_INFORMATION
	ret, _, err := getFileInformationByHandle.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&fileInfo)),
	)
	if ret == 0 {
		return 0, fmt.Errorf("failed to get file information: %w", err)
	}

	return uint64(fileInfo.NumberOfLinks), nil
}

func LinkCount(info os.FileInfo) uint64 {
	n, err := GetNumberOfHardLinks(info)
	if err != nil {
		return 0
	}
	return n
}
