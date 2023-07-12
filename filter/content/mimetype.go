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
	MimeTypeName       = filter.MimeTypeName
	ParentMimeTypeName = "Parent-Mime-type"
)

func (e *MimeFileTypeEnabler) Enable() filter.ContentOption {
	return func(info *item.FileInfo) (string, string) {
		tn := ""
		returnName := MimeTypeName
		if e.ParentOnly {
			returnName = ParentMimeTypeName
		}
		if c, ok := info.Cache[MimeTypeName]; ok {
			tn = string(c)
		} else {
			if info.IsDir() {
				tn = "directory"
				return tn, returnName
			}
			if util.IsSymLink(info) {
				tn = "symlink"
				return tn, returnName
			}
			if info.Mode()&os.ModeNamedPipe != 0 {
				tn = "named_pipe"
				return tn, returnName
			}
			if info.Mode()&os.ModeSocket != 0 {
				tn = "socket"
				return tn, returnName
			}
			if m, ok := info.Cache[MimeTypeName]; ok {
				info.Cache[Charset] = m
				return string(m), returnName
			}

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
		}

		if e.ParentOnly {
			tn = strings.SplitN(tn, "/", 2)[0]
		}

		if strings.Contains(tn, ";") {
			// remove charset
			s := strings.SplitN(tn, ";", 2)
			tn = s[0]
			charset := strings.SplitN(s[1], "=", 2)[1]
			info.Cache[Charset] = []byte(charset)
		}

		return tn, returnName
	}
}
