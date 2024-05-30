// Package global contains the global variables used in the project
// this package can't depend on other packages
package global

import (
	"golang.org/x/exp/constraints"
)

// Hashable is the type of values that may be used as map keys or set members.
// it should be exact same as haxmap.hashable (it's unexported, so we can't use it directly)
type Hashable interface {
	constraints.Integer | constraints.Float | constraints.Complex | ~string | uintptr
}

const (
	Black        = "\033[0;30m" // 0,0,0
	Red          = "\033[0;31m" // 205,0,0
	Green        = "\033[0;32m" // 0,205,0
	Yellow       = "\033[0;33m" // 205,205,0
	Blue         = "\033[0;34m" // 0,0,238
	Purple       = "\033[0;35m" // 205,0,205
	Cyan         = "\033[0;36m" // 0,205,205
	White        = "\033[0;37m" // 229,229,229
	BrightBlack  = "\033[0;90m" // 127,127,127
	BrightRed    = "\033[0;91m" // 255,0,0
	BrightGreen  = "\033[0;92m" // 0,255,0
	BrightYellow = "\033[0;93m" // 255,255,0
	BrightBlue   = "\033[0;94m" // 92,92,255
	BrightPurple = "\033[0;95m" // 255,0,255
	BrightCyan   = "\033[0;96m" // 0,255,255
	BrightWhite  = "\033[0;97m" // 255,255,255
	Success      = Green
	Error        = Red
	Warn         = Yellow
	Bold         = "\033[1m"
	Faint        = "\033[2m"
	Italics      = "\033[3m"
	Underline    = "\033[4m"
	Blink        = "\033[5m"
	Reverse      = "\033[7m"
	ReverseDone  = "\033[27m"
)

const (
	BoostThreshold = 0.7
	PrefixSize     = 4
)

const Reset = "\033[0m"

const DefaultHookLen = 5
const Space = 2

const (
	NameOfName          = "Name"
	NameOfCharset       = "Charset"
	NameOfMIME          = "Mime-type"
	NameOfLink          = "Link"
	NameOfInode         = "Inode"
	NameOfIndex         = "#"
	NameOfGroupName     = "Group"
	NameOfGroupUidName  = "Group-uid"
	NameOfGroupSID      = "Group-sid"
	NameOfOwner         = "Owner"
	NameOfOwnerUid      = "Owner-uid"
	NameOfOwnerSID      = "Owner-sid"
	NameOfSize          = "Size"
	NameOfGitStatus     = "Git"
	NameOfGitRepoBranch = "Branch"
	NameOfGitRepoStatus = "Repo-status"
	NameOfPermission    = "Permissions"
	NameOfSum           = "Sum"
	NameOfRelativeTime  = "Relative-Time"
	NameOfTime          = "Time"
	NameOfTimeModified  = "Modified"
	NameOfTimeCreated   = "Created"
	NameOfTimeAccessed  = "Accessed"
	NameOfTimeBirth     = "Birth"
)
