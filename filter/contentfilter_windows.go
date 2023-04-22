//go:build windows

package filter

import (
	"g/render"
	"os"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

var (
	libadvapi32                    = syscall.NewLazyDLL("advapi32.dll")
	procGetFileSecurity            = libadvapi32.NewProc("GetFileSecurityW")
	procGetSecurityDescriptorOwner = libadvapi32.NewProc("GetSecurityDescriptorOwner")
)

func EnableOwner(renderer *render.Renderer) ContentOption {
	m := sync.RWMutex{}
	longestOwner := 0
	return func(info os.FileInfo) string {
		path := info.Name()

		var needed uint32
		fromString, err := syscall.UTF16PtrFromString(path)
		if err != nil {
			return ""
		}
		_, _, _ = procGetFileSecurity.Call(
			uintptr(unsafe.Pointer(fromString)),
			0x00000001, /* OWNER_SECURITY_INFORMATION */
			0,
			0,
			uintptr(unsafe.Pointer(&needed)))
		buf := make([]byte, needed)
		r1, _, err := procGetFileSecurity.Call(
			uintptr(unsafe.Pointer(fromString)),
			0x00000001, /* OWNER_SECURITY_INFORMATION */
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(needed),
			uintptr(unsafe.Pointer(&needed)))
		if r1 == 0 && err != nil {
			return renderer.Owner(fillBlank("", longestOwner))
		}
		var ownerDefaulted uint32
		var sid *syscall.SID
		r1, _, err = procGetSecurityDescriptorOwner.Call(
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(unsafe.Pointer(&sid)),
			uintptr(unsafe.Pointer(&ownerDefaulted)))
		if r1 == 0 && err != nil {
			return renderer.Owner(fillBlank("", longestOwner))
		}
		name, _, _, err := sid.LookupAccount("")
		if r1 == 0 && err != nil {
			return renderer.Owner(fillBlank("", longestOwner))
		}

		m.RLock()
		if len(name) > longestOwner {
			m.RUnlock()
			m.Lock()
			if len(name) > longestOwner {
				longestOwner = len(name)
			}
			m.Unlock()
		} else {
			m.RUnlock()
		}

		time.Sleep(time.Microsecond * 5)

		return renderer.Owner(fillBlank(name, longestOwner))

	}
}

func EnableGroup(renderer *render.Renderer) ContentOption {
	m := sync.RWMutex{}
	longestGroup := 0
	return func(info os.FileInfo) string {
		path := info.Name()

		var needed uint32
		fromString, err := syscall.UTF16PtrFromString(path)
		if err != nil {
			return ""
		}
		_, _, _ = procGetFileSecurity.Call(
			uintptr(unsafe.Pointer(fromString)),
			0x00000001, /* OWNER_SECURITY_INFORMATION */
			0,
			0,
			uintptr(unsafe.Pointer(&needed)))
		buf := make([]byte, needed)
		r1, _, err := procGetFileSecurity.Call(
			uintptr(unsafe.Pointer(fromString)),
			0x00000001, /* OWNER_SECURITY_INFORMATION */
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(needed),
			uintptr(unsafe.Pointer(&needed)))
		if r1 == 0 && err != nil {
			return renderer.Owner(fillBlank("", longestGroup))
		}
		var ownerDefaulted uint32
		var sid *syscall.SID
		r1, _, err = procGetSecurityDescriptorOwner.Call(
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(unsafe.Pointer(&sid)),
			uintptr(unsafe.Pointer(&ownerDefaulted)))
		if r1 == 0 && err != nil {
			return renderer.Owner(fillBlank("", longestGroup))
		}
		_, name, _, err := sid.LookupAccount("")
		if r1 == 0 && err != nil {
			return renderer.Owner(fillBlank("", longestGroup))
		}

		m.RLock()
		if len(name) > longestGroup {
			m.RUnlock()
			m.Lock()
			if len(name) > longestGroup {
				longestGroup = len(name)
			}
			m.Unlock()
		} else {
			m.RUnlock()
		}

		time.Sleep(time.Microsecond * 5)

		return renderer.Owner(fillBlank(name, longestGroup))

	}
}
