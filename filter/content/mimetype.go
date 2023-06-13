package content

import (
	"os"
	"strings"
	"sync"

	"github.com/Equationzhao/g/filter"
	"github.com/gabriel-vasile/mimetype"
)

type MimeFileTypeEnabler struct {
	*sync.WaitGroup
	ParentOnly bool
}

func NewMimeFileTypeEnabler() *MimeFileTypeEnabler {
	return &MimeFileTypeEnabler{
		WaitGroup:  &sync.WaitGroup{},
		ParentOnly: false,
	}
}

const MimeTypeName = "Mime-type"

func (e *MimeFileTypeEnabler) Enable() filter.ContentOption {
	longestTypeName := 0
	m := sync.RWMutex{}
	done := func(tn string) {
		defer e.Done()
		m.RLock()
		if longestTypeName >= len(tn) {
			m.RUnlock()
			return
		}
		m.RUnlock()
		m.Lock()
		if longestTypeName < len(tn) {
			longestTypeName = len(tn)
		}
		m.Unlock()
	}

	wait := func(tn string) string {
		e.Wait()
		return filter.FillBlank(tn, longestTypeName)
	}
	return func(info os.FileInfo) (string, string) {
		tn := ""
		returnName := MimeTypeName
		if e.ParentOnly {
			returnName = "parent-" + returnName
		}
		if info.IsDir() {
			tn = "directory"
		} else if info.Mode()&os.ModeSymlink != 0 {
			tn = "symlink"
		} else if info.Mode()&os.ModeNamedPipe != 0 {
			tn = "named_pipe"
		} else if info.Mode()&os.ModeSocket != 0 {
			tn = "socket"
		} else {
			file, err := os.Open(info.Name())
			defer file.Close()
			if err != nil {
				// tn = err.Error()
				tn = "failed_to_read"
				done(tn)
				return wait(tn), returnName
			}
			mtype, err := mimetype.DetectReader(file)
			if err != nil {
				tn = err.Error()
				done(tn)
				return wait(tn), returnName
			}
			tn = mtype.String()

			if e.ParentOnly {
				tn = strings.SplitN(tn, "/", 2)[0]
			}

			if strings.Contains(tn, ";") {
				// remove charset
				tn = strings.SplitN(tn, ";", 2)[0]
			}

		}
		done(tn)
		return wait(tn), returnName
	}
}
