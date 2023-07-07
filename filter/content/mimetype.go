package content

import (
	"os"
	"strings"

	"github.com/Equationzhao/g/filter"
	"github.com/Equationzhao/g/item"
	"github.com/Equationzhao/g/util"
	"github.com/gabriel-vasile/mimetype"
)

type MimeFileTypeEnabler struct {
	ParentOnly bool
}

func NewMimeFileTypeEnabler() *MimeFileTypeEnabler {
	return &MimeFileTypeEnabler{
		ParentOnly: false,
	}
}

const (
	MimeTypeName       = "Mime-type"
	ParentMimeTypeName = "Parent-Mime-type"
)

func (e *MimeFileTypeEnabler) Enable() filter.ContentOption {
	return func(info *item.FileInfo) (string, string) {
		tn := ""
		returnName := MimeTypeName
		if e.ParentOnly {
			returnName = ParentMimeTypeName
		}
		if info.IsDir() {
			tn = "directory"
		} else if util.IsSymLink(info) {
			tn = "symlink"
		} else if info.Mode()&os.ModeNamedPipe != 0 {
			tn = "named_pipe"
		} else if info.Mode()&os.ModeSocket != 0 {
			tn = "socket"
		} else {
			file, err := os.Open(info.FullPath)
			defer file.Close()
			if err != nil {
				return "failed_to_read", returnName
			}
			mtype, err := mimetype.DetectReader(file)
			if err != nil {
				return err.Error(), returnName
			}
			tn = mtype.String()

			if e.ParentOnly {
				tn = strings.SplitN(tn, "/", 2)[0]
			}

			if strings.Contains(tn, ";") {
				// remove charset
				s := strings.SplitN(tn, ";", 2)
				tn = s[0]
				charset := strings.SplitN(s[1], "=", 2)[1]
				info.Cache[charsetIdentifier] = []byte(charset)
			}

		}
		return tn, returnName
	}
}
