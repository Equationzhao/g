//go:build darwin

package osbased

/*
#cgo CFLAGS: -mmacosx-version-min=10.9
#cgo LDFLAGS: -framework CoreFoundation -framework CoreServices

#include <CoreFoundation/CoreFoundation.h>
#include <CoreServices/CoreServices.h>

#include <stdlib.h>
#include <string.h>
#include <limits.h>
#include <sys/xattr.h>

bool isAlias(const char *path) {
    CFStringRef cfPath = CFStringCreateWithCString(NULL, path, kCFStringEncodingUTF8);
    CFURLRef url = CFURLCreateWithFileSystemPath(NULL, cfPath, kCFURLPOSIXPathStyle, false);
    CFRelease(cfPath);

    Boolean isAlias = false;
    if (url) {
        CFBooleanRef isAliasRef;
        if (CFURLCopyResourcePropertyForKey(url, kCFURLIsAliasFileKey, &isAliasRef, NULL)) {
            isAlias = CFBooleanGetValue(isAliasRef);
            CFRelease(isAliasRef);
        }
        CFRelease(url);
    }

    return isAlias;
}

char *resolveAlias(const char *path) {
    CFURLRef url = CFURLCreateFromFileSystemRepresentation(NULL, (const UInt8 *)path, strlen(path), false);
    if (!url) {
        return NULL;
    }

    CFErrorRef error = NULL;
    CFDataRef bookmarkData = CFURLCreateBookmarkDataFromFile(NULL, url, &error);
    CFRelease(url);
    if (!bookmarkData) {
        if (error != NULL) {
            CFRelease(error);
        }
        return NULL;
    }

    Boolean bookmarkIsStale;
    CFURLRef resolvedURL = CFURLCreateByResolvingBookmarkData(NULL, bookmarkData, kCFBookmarkResolutionWithoutUIMask, NULL, NULL, &bookmarkIsStale, &error);
    CFRelease(bookmarkData);
    if (!resolvedURL) {
        if (error != NULL) {
            CFRelease(error);
        }
        return NULL;
    }

    UInt8 buffer[PATH_MAX];
    Boolean success = CFURLGetFileSystemRepresentation(resolvedURL, true, buffer, PATH_MAX);
    CFRelease(resolvedURL);
    if (!success) {
        return NULL;
    }

    return strdup((const char*)buffer);
}
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
