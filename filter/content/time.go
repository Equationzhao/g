package content

import (
	"fmt"
	"sync"
	"time"

	"github.com/Equationzhao/g/filter"
	"github.com/Equationzhao/g/item"
	"github.com/Equationzhao/g/osbased"
	"github.com/Equationzhao/g/render"
	"github.com/hako/durafmt"
)

type RelativeTimeEnabler struct {
	*sync.WaitGroup
	Mode string
}

func NewRelativeTimeEnabler() *RelativeTimeEnabler {
	return &RelativeTimeEnabler{
		WaitGroup: new(sync.WaitGroup),
	}
}

const RelativeTime = "Relative-Time"

func (r *RelativeTimeEnabler) Enable(renderer *render.Renderer) filter.ContentOption {
	longestRt := 0
	m := sync.RWMutex{}
	done := func(rt string) {
		defer r.Done()
		m.RLock()
		if longestRt >= len(rt) {
			m.RUnlock()
			return
		}
		m.RUnlock()
		m.Lock()
		if longestRt < len(rt) {
			longestRt = len(rt)
		}
		m.Unlock()
	}

	wait := func(size string) string {
		r.Wait()
		return filter.FillBlank(size, longestRt)
	}

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
		default:
			t = osbased.ModTime(info)
			timeType = timeModified
		}
		rt := renderer.Time(relativeTime(time.Now(), t))
		done(rt)
		return wait(rt), RelativeTime + " " + timeType
	}
}

func relativeTime(now, modTime time.Time) string {
	if t := now.Sub(modTime); t > 0 {
		return fmt.Sprintf("%s ago", durafmt.Parse(t).LimitFirstN(1).String())
	} else if t == 0 {
		return "now"
	} else {
		return fmt.Sprintf("in %s", durafmt.Parse(-t).LimitFirstN(1).String())
	}
}

const (
	timeName     = "Time"
	timeModified = "Modified"
	timeCreated  = "Created"
	timeAccessed = "Accessed"
)

func EnableTime(format string, mode string, renderer *render.Renderer) filter.ContentOption {
	return func(info *item.FileInfo) (string, string) {
		// get mod time/ create time/ access time
		var t time.Time
		timeType := ""
		switch mode {
		case "mod":
			t = osbased.ModTime(info)
			timeType = timeModified
		case "create":
			t = osbased.CreateTime(info)
			timeType = timeCreated
		case "access":
			t = osbased.AccessTime(info)
			timeType = timeAccessed
		default:
			t = osbased.ModTime(info)
			timeType = timeModified
		}
		return renderer.Time(t.Format(format)), timeName + " " + timeType
	}
}
