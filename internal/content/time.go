package content

import (
	"runtime"
	"time"

	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/osbased"
	"github.com/Equationzhao/g/internal/render"
)

type RelativeTimeEnabler struct {
	Mode string
}

func NewRelativeTimeEnabler() *RelativeTimeEnabler {
	return &RelativeTimeEnabler{}
}

const RelativeTime = "Relative-Time"

func (r *RelativeTimeEnabler) Enable(renderer *render.Renderer) ContentOption {
	return func(info *item.FileInfo) (string, string) {
		var t time.Time
		timeType := ""
		switch r.Mode {
		case "mod":
			t = osbased.ModTime(info)
			timeType = timeModified
		case "create":
			t = osbased.CreateTime(info)
			timeType = timeCreated
		case "access":
			t = osbased.AccessTime(info)
			timeType = timeAccessed
		case "birth":
			timeType = timeBirth
			// if darwin, check birth time
			if runtime.GOOS == "darwin" {
				t = osbased.BirthTime(info)
			} else {
				t = osbased.CreateTime(info)
			}
		default:
			t = osbased.ModTime(info)
			timeType = timeModified
		}
		return renderer.RTime(time.Now(), t), RelativeTime + " " + timeType
	}
}

const (
	timeName     = "Time"
	timeModified = "Modified"
	timeCreated  = "Created"
	timeAccessed = "Accessed"
	timeBirth    = "Birth"
)

// EnableTime enables time
// accepts ['mod', 'modified', 'create', 'access', 'birth']
func EnableTime(format string, mode string, renderer *render.Renderer) ContentOption {
	return func(info *item.FileInfo) (string, string) {
		// get mod time/ create time/ access time
		var t time.Time
		timeType := ""
		switch mode {
		case "modified", "mod":
			t = osbased.ModTime(info)
			timeType = timeModified
		case "created", "cr":
			t = osbased.CreateTime(info)
			timeType = timeCreated
		case "accessed", "ac":
			t = osbased.AccessTime(info)
			timeType = timeAccessed
		case "birth":
			timeType = timeBirth
			// if darwin, check birth time
			if runtime.GOOS == "darwin" {
				t = osbased.BirthTime(info)
			} else {
				t = osbased.CreateTime(info)
			}
		default:
			t = osbased.ModTime(info)
			timeType = timeModified
		}
		return renderer.Time(t.Format(format)), timeName + " " + timeType
	}
}
