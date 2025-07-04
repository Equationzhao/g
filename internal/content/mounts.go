//go:build mounts

package content

import (
	"github.com/Equationzhao/g/internal/util"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/valyala/bytebufferpool"
)

func MountsOn(path string) string {
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
		if stat.Mountpoint == path {
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
	mountsOnce = util.Once{}
)