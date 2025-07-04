//go:build !mounts

package content

// Lite version without mounts functionality
// This reduces binary size by removing gopsutil dependency

func MountsOn(path string) string {
	// In lite build, return empty string (no mount info)
	return ""
}