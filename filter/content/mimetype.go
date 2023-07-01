package content

import (
	"os"
	"strings"

	"github.com/Equationzhao/g/filter"
	"github.com/Equationzhao/g/item"
	"github.com/gabriel-vasile/mimetype"
)

type MimeFileTypeEnabler struct {
	ParentOnly, EnableCharset bool
}

func NewMimeFileTypeEnabler() *MimeFileTypeEnabler {
	return &MimeFileTypeEnabler{
		ParentOnly:    false,
		EnableCharset: false,
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
		} else if info.Mode()&os.ModeSymlink != 0 {
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

			if !e.EnableCharset {
				if strings.Contains(tn, ";") {
					// remove charset
					tn = strings.SplitN(tn, ";", 2)[0]
				}
			}

		}
		return tn, returnName
	}
}
