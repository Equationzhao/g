//go:build linux

package sorter

import (
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/Equationzhao/g/cached"
)

func byGroupID(a, b os.FileInfo, Ascend bool) bool {
	if Ascend {
		return a.Sys().(*syscall.Stat_t).Gid < b.Sys().(*syscall.Stat_t).Gid
	}
	return a.Sys().(*syscall.Stat_t).Gid > b.Sys().(*syscall.Stat_t).Gid
}

func byOwnerID(a, b os.FileInfo, Ascend bool) bool {
	if Ascend {
		return a.Sys().(*syscall.Stat_t).Uid < b.Sys().(*syscall.Stat_t).Uid
	}
	return a.Sys().(*syscall.Stat_t).Uid > b.Sys().(*syscall.Stat_t).Uid
}

func byGroupName(a, b os.FileInfo, Ascend bool) bool {
	if Ascend {
		return strings.ToLower(cached.GetGroupname(strconv.FormatInt(int64(a.Sys().(*syscall.Stat_t).Gid), 10))) < strings.ToLower(cached.GetGroupname(strconv.FormatInt(int64(b.Sys().(*syscall.Stat_t).Gid), 10)))
	}
	return strings.ToLower(cached.GetGroupname(strconv.FormatInt(int64(a.Sys().(*syscall.Stat_t).Gid), 10))) > strings.ToLower(cached.GetGroupname(strconv.FormatInt(int64(b.Sys().(*syscall.Stat_t).Gid), 10)))
}

func byUserName(a, b os.FileInfo, Ascend bool) bool {
	if Ascend {
		return strings.ToLower(cached.GetUsername(strconv.FormatInt(int64(a.Sys().(*syscall.Stat_t).Uid), 10))) < strings.ToLower(cached.GetGroupname(strconv.FormatInt(int64(b.Sys().(*syscall.Stat_t).Uid), 10)))
	}
	return strings.ToLower(cached.GetUsername(strconv.FormatInt(int64(a.Sys().(*syscall.Stat_t).Uid), 10))) > strings.ToLower(cached.GetGroupname(strconv.FormatInt(int64(b.Sys().(*syscall.Stat_t).Uid), 10)))
}

func byGroupCaseSensitiveName(a, b os.FileInfo, Ascend bool) bool {
	if Ascend {
		return cached.GetGroupname(strconv.FormatInt(int64(a.Sys().(*syscall.Stat_t).Gid), 10)) < cached.GetGroupname(strconv.FormatInt(int64(b.Sys().(*syscall.Stat_t).Gid), 10))
	}
	return cached.GetGroupname(strconv.FormatInt(int64(a.Sys().(*syscall.Stat_t).Gid), 10)) > cached.GetGroupname(strconv.FormatInt(int64(b.Sys().(*syscall.Stat_t).Gid), 10))
}

func byUserCaseSensitiveName(a, b os.FileInfo, Ascend bool) bool {
	if Ascend {
		return cached.GetUsername(strconv.FormatInt(int64(a.Sys().(*syscall.Stat_t).Uid), 10)) < cached.GetGroupname(strconv.FormatInt(int64(b.Sys().(*syscall.Stat_t).Uid), 10))
	}
	return cached.GetUsername(strconv.FormatInt(int64(a.Sys().(*syscall.Stat_t).Uid), 10)) > cached.GetGroupname(strconv.FormatInt(int64(b.Sys().(*syscall.Stat_t).Uid), 10))
}
