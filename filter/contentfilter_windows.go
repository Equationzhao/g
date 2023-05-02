//go:build windows

/*
	MIT License

	Copyright (c) 2019 Andrew Carlson

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all
	copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
	SOFTWARE.
*/

package filter

import (
	"os"
	"sync"
	"syscall"
	"unsafe"

	"github.com/Equationzhao/g/render"
	"github.com/Equationzhao/g/usergroup_windows"
)

var (
	libadvapi32                    = syscall.NewLazyDLL("advapi32.dll")
	procGetFileSecurity            = libadvapi32.NewProc("GetFileSecurityW")
	procGetSecurityDescriptorOwner = libadvapi32.NewProc("GetSecurityDescriptorOwner")
)

func (cf *ContentFilter) EnableOwner(renderer *render.Renderer) ContentOption {
	m := sync.RWMutex{}
	longestOwner := 0
	return func(info os.FileInfo) string {
		path := info.Name()
		wait := func(res string) string {
			cf.wgOwner.Wait()
			return renderer.Owner(fillBlank(res, longestOwner))
		}

		done := func(name string) {
			defer cf.wgOwner.Done()
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
		}

		var needed uint32
		fromString, err := syscall.UTF16PtrFromString(path)
		if err != nil {
			done("unknown")
			return wait("unknown")
		}
		_, _, _ = procGetFileSecurity.Call(
			uintptr(unsafe.Pointer(fromString)),
			0x00000001, /* OWNER_SECURITY_INFORMATION */
			0,
			0,
			uintptr(unsafe.Pointer(&needed)))
		buf := make([]byte, needed)

		if needed == 0 {
			done("unknown")
			return wait("unknown")
		}

		r1, _, err := procGetFileSecurity.Call(
			uintptr(unsafe.Pointer(fromString)),
			0x00000001, /* OWNER_SECURITY_INFORMATION */
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(needed),
			uintptr(unsafe.Pointer(&needed)))
		if r1 == 0 && err != nil {
			done("unknown")
			return wait("unknown")
		}
		var ownerDefaulted uint32
		var sid *syscall.SID
		r1, _, err = procGetSecurityDescriptorOwner.Call(
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(unsafe.Pointer(&sid)),
			uintptr(unsafe.Pointer(&ownerDefaulted)))
		if r1 == 0 && err != nil {
			done("unknown")
			return wait("unknown")
		}
		name, _, _, err := sid.LookupAccount("")
		if r1 == 0 && err != nil {
			done("unknown")
			return wait("unknown")
		}
		done(name)
		return wait(name)
	}
}

func (cf *ContentFilter) EnableGroup(renderer *render.Renderer) ContentOption {
	m := sync.RWMutex{}
	longestGroup := 0
	return func(info os.FileInfo) string {
		wait := func(name string) string {
			cf.wgGroup.Wait()
			return renderer.Group(fillBlank(name, longestGroup))
		}
		done := func(name string) {
			defer cf.wgGroup.Done()
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
		}

		// var needed uint32
		// fromString, err := syscall.UTF16PtrFromString(path)
		// if err != nil {
		// 	done("unknown")
		// 	return wait("unknown")
		// }
		// _, _, _ = procGetFileSecurity.Call(
		// 	uintptr(unsafe.Pointer(fromString)),
		// 	0x00000001, /* OWNER_SECURITY_INFORMATION */
		// 	0,
		// 	0,
		// 	uintptr(unsafe.Pointer(&needed)))
		//
		// if needed == 0 {
		// 	done("unknown")
		// 	return wait("unknown")
		// }
		//
		// buf := make([]byte, needed)
		// r1, _, err := procGetFileSecurity.Call(
		// 	uintptr(unsafe.Pointer(fromString)),
		// 	0x00000001, /* OWNER_SECURITY_INFORMATION */
		// 	uintptr(unsafe.Pointer(&buf[0])),
		// 	uintptr(needed),
		// 	uintptr(unsafe.Pointer(&needed)))
		// if r1 == 0 && err != nil {
		// 	done("unknown")
		// 	return wait("unknown")
		// }
		// var ownerDefaulted uint32
		// var sid *syscall.SID
		// r1, _, err = procGetSecurityDescriptorOwner.Call(
		// 	uintptr(unsafe.Pointer(&buf[0])),
		// 	uintptr(unsafe.Pointer(&sid)),
		// 	uintptr(unsafe.Pointer(&ownerDefaulted)))
		// if r1 == 0 && err != nil {
		// 	done("unknown")
		// 	return wait("unknown")
		// }
		// _, name, _, err := sid.LookupAccount("")
		// if r1 == 0 && err != nil {
		// 	done("unknown")
		// 	return wait("unknown")
		// }
		name := usergroup_windows.Group(info)
		done(name)
		return wait(name)
	}
}
