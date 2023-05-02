package sorter

import (
	"os"

	"github.com/Equationzhao/g/usergroup_windows"
)

func byGroupName(a, b os.FileInfo, Ascend bool) bool {
	if Ascend {
		return usergroup_windows.Group(a) < usergroup_windows.Group(b)
	}
	return usergroup_windows.Group(a) > usergroup_windows.Group(b)
}

func byUserName(a, b os.FileInfo, Ascend bool) bool {
	if Ascend {
		return usergroup_windows.Owner(a) < usergroup_windows.Owner(b)
	}
	return usergroup_windows.Owner(a) > usergroup_windows.Owner(b)
}
