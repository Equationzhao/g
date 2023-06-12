package content

import (
	"os"
	"strconv"
	"sync"

	"github.com/Equationzhao/g/filter"
	"github.com/Equationzhao/g/osbased"
)

type LinkEnabler struct {
	// List each file's number of hard links.
	*sync.WaitGroup
}

func NewLinkEnabler() *LinkEnabler {
	return &LinkEnabler{
		WaitGroup: &sync.WaitGroup{},
	}
}

func (l *LinkEnabler) Enable() filter.ContentOption {
	var longestLinkNum string
	m := sync.RWMutex{}
	done := func(linkNumStr string) {
		defer l.Done()
		m.RLock()
		if len(longestLinkNum) >= len(linkNumStr) {
			m.RUnlock()
			return
		}
		m.RUnlock()
		m.Lock()
		if len(longestLinkNum) < len(linkNumStr) {
			longestLinkNum = linkNumStr
		}
		m.Unlock()
	}

	wait := func(linkNumStr string) string {
		l.Wait()
		return filter.FillBlank(linkNumStr, len(longestLinkNum))
	}

	return func(info os.FileInfo) (string, string) {
		n := strconv.FormatUint(osbased.LinkCount(info), 10)
		done(n)
		return wait(n), "links"
	}
}
