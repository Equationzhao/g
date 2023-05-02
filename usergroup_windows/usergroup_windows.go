package usergroup_windows

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	libadvapi32                    = syscall.NewLazyDLL("advapi32.dll")
	procGetFileSecurity            = libadvapi32.NewProc("GetFileSecurityW")
	procGetSecurityDescriptorOwner = libadvapi32.NewProc("GetSecurityDescriptorOwner")
)

func Group(info os.FileInfo) string {
	path := info.Name()

	var needed uint32
	fromString, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return "unknown"
	}
	_, _, _ = procGetFileSecurity.Call(
		uintptr(unsafe.Pointer(fromString)),
		0x00000001, /* OWNER_SECURITY_INFORMATION */
		0,
		0,
		uintptr(unsafe.Pointer(&needed)))

	if needed == 0 {
		return "unknown"
	}

	buf := make([]byte, needed)
	r1, _, err := procGetFileSecurity.Call(
		uintptr(unsafe.Pointer(fromString)),
		0x00000001, /* OWNER_SECURITY_INFORMATION */
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(needed),
		uintptr(unsafe.Pointer(&needed)))
	if r1 == 0 && err != nil {
		return "unknown"
	}
	var ownerDefaulted uint32
	var sid *syscall.SID
	r1, _, err = procGetSecurityDescriptorOwner.Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&sid)),
		uintptr(unsafe.Pointer(&ownerDefaulted)))
	if r1 == 0 && err != nil {
		return "unknown"
	}
	_, name, _, err := sid.LookupAccount("")
	if r1 == 0 && err != nil {
		return "unknown"
	}

	return name
}

func Owner(info os.FileInfo) string {
	path := info.Name()

	var needed uint32
	fromString, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return "unknown"
	}
	_, _, _ = procGetFileSecurity.Call(
		uintptr(unsafe.Pointer(fromString)),
		0x00000001, /* OWNER_SECURITY_INFORMATION */
		0,
		0,
		uintptr(unsafe.Pointer(&needed)))
	buf := make([]byte, needed)

	if needed == 0 {
		return "unknown"
	}

	r1, _, err := procGetFileSecurity.Call(
		uintptr(unsafe.Pointer(fromString)),
		0x00000001, /* OWNER_SECURITY_INFORMATION */
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(needed),
		uintptr(unsafe.Pointer(&needed)))
	if r1 == 0 && err != nil {
		return "unknown"
	}
	var ownerDefaulted uint32
	var sid *syscall.SID
	r1, _, err = procGetSecurityDescriptorOwner.Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&sid)),
		uintptr(unsafe.Pointer(&ownerDefaulted)))
	if r1 == 0 && err != nil {
		return "unknown"
	}
	name, _, _, err := sid.LookupAccount("")
	if r1 == 0 && err != nil {
		return "unknown"
	}
	return name
}
