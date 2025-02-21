#include <CoreFoundation/CoreFoundation.h>
#include <CoreServices/CoreServices.h>

#include <stdlib.h>
#include <string.h>
#include <limits.h>
#include <sys/xattr.h>

bool isAlias(const char *path) {
    if (path == NULL) {
        return false;
    }

    CFStringRef cfPath = CFStringCreateWithCString(NULL, path, kCFStringEncodingUTF8);
    if (cfPath == NULL) {
        return false;
    }

    CFURLRef url = CFURLCreateWithFileSystemPath(NULL, cfPath, kCFURLPOSIXPathStyle, false);
    CFRelease(cfPath);

    if (url == NULL) {
        return false;
    }

    Boolean isAlias = false;
    CFBooleanRef isAliasRef = NULL;
    if (CFURLCopyResourcePropertyForKey(url, kCFURLIsAliasFileKey, &isAliasRef, NULL)) {
        isAlias = CFBooleanGetValue(isAliasRef);
        CFRelease(isAliasRef);
    }
    CFRelease(url);

    return isAlias;
}

char *resolveAlias(const char *path) {
    if (path == NULL) {
        return NULL;
    }

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