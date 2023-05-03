package sorter

import (
	"os"
	"strings"

	"github.com/Equationzhao/g/usergroup_windows"
)

func byGroupName(a, b os.FileInfo, Ascend bool) bool {
	if Ascend {
		return strings.ToLower(usergroup_windows.Group(a)) < strings.ToLower(usergroup_windows.Group(b))
	}
	return strings.ToLower(usergroup_windows.Group(a)) > strings.ToLower(usergroup_windows.Group(b))
}

func byUserName(a, b os.FileInfo, Ascend bool) bool {
	if Ascend {
		return strings.ToLower(usergroup_windows.Owner(a)) < strings.ToLower(usergroup_windows.Owner(b))
	}
	return strings.ToLower(usergroup_windows.Owner(a)) > strings.ToLower(usergroup_windows.Owner(b))
}

func byGroupCaseSensitiveName(a, b os.FileInfo, Ascend bool) bool {
	if Ascend {
		return usergroup_windows.Group(a) < usergroup_windows.Group(b)
	}
	return usergroup_windows.Group(a) > usergroup_windows.Group(b)
}

func byUserCaseSensitiveName(a, b os.FileInfo, Ascend bool) bool {
	if Ascend {
		return usergroup_windows.Owner(a) < usergroup_windows.Owner(b)
	}
	return usergroup_windows.Owner(a) > usergroup_windows.Owner(b)
}
