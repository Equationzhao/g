package util

import (
	"os"
	"path/filepath"

	"github.com/Equationzhao/g/item"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/valyala/bytebufferpool"
)

func IsSymLink(file os.FileInfo) bool {
	return file.Mode()&os.ModeSymlink != 0
}

func IsSymLinkMode(mode os.FileMode) bool {
	return mode&os.ModeSymlink != 0
}

func IsExecutable(file os.FileInfo) bool {
	return file.Mode()&0o111 != 0
}

func IsExecutableMode(mode os.FileMode) bool {
	return mode&0o111 != 0
}

// RecursivelySizeOf returns the size of the file or directory
// depth < 0 means no limit
func RecursivelySizeOf(info os.FileInfo, depth int) int64 {
	currentDepth := 0
	if info.IsDir() {
		totalSize := info.Size()
		if depth < 0 {
			// -1 means no limit
			_ = filepath.WalkDir(
				info.Name(), func(path string, dir os.DirEntry, err error) error {
					if err != nil {
						return err
					}

					if !dir.IsDir() {
						info, err := dir.Info()
						if err == nil {
							totalSize += info.Size()
						}
					}

					return nil
				},
			)
		} else {
			_ = filepath.WalkDir(
				info.Name(), func(path string, dir os.DirEntry, err error) error {
					if err != nil {
						return err
					}
					if currentDepth > depth {
						if dir.IsDir() {
							return filepath.SkipDir
						}
						return nil
					}

					if !dir.IsDir() {
						info, err := dir.Info()
						if err == nil {
							totalSize += info.Size()
						}
					}

					if dir.IsDir() {
						currentDepth++
					}

					return nil
				},
			)
		}

		return totalSize
	}
	return info.Size()
}

func MountsOn(info *item.FileInfo) string {
	err := mountsOnce.Do(func() error {
		mount, err := disk.Partitions(true)
		if err != nil {
			return err
		}
		mounts = mount
		return nil
	})
	if err != nil {
		return ""
	}
	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)
	for _, stat := range mounts {
		if stat.Mountpoint == info.FullPath {
			_ = b.WriteByte('[')
			_, _ = b.WriteString(stat.Device)
			_, _ = b.WriteString(" (")
			_, _ = b.WriteString(stat.Fstype)
			_, _ = b.WriteString(")]")
			return b.String()
		}
	}
	return ""
}

var (
	mounts     = make([]disk.PartitionStat, 10)
	mountsOnce = Once{}
)
