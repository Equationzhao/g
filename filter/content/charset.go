package content

import (
	"bytes"
	"os"
	"strings"

	"github.com/Equationzhao/g/filter"
	"github.com/Equationzhao/g/item"
	"github.com/Equationzhao/g/util"
	"github.com/gabriel-vasile/mimetype"
	"github.com/saintfish/chardet"
)

type CharsetEnabler struct{}

func NewCharsetEnabler() *CharsetEnabler {
	return &CharsetEnabler{}
}

const (
	charsetIdentifier = "charset"
	Charset           = "Charset"
)

func (c *CharsetEnabler) Enable() filter.ContentOption {
	det := chardet.NewTextDetector()
	return func(info *item.FileInfo) (string, string) {
		if c, ok := info.Cache[charsetIdentifier]; ok {
			return string(c), Charset
		}
		if info.IsDir() {
			return "", Charset
		} else if util.IsSymLink(info) {
			return "", Charset
		} else if info.Mode()&os.ModeNamedPipe != 0 {
			return "", Charset
		} else if info.Mode()&os.ModeSocket != 0 {
			return "", Charset
		} else {
			file, err := os.Open(info.FullPath)
			defer file.Close()
			if err != nil {
				return "failed_to_read", Charset
			}
			content := make([]byte, 1024*1024)
			_, err = file.Read(content)
			if err != nil {
				return err.Error(), Charset
			}

			mtype, err := mimetype.DetectReader(bytes.NewReader(content))
			if err != nil {
				return err.Error(), Charset
			}

			charset := ""
			if tn := mtype.String(); strings.Contains(tn, ";") {
				s := strings.SplitN(tn, ";", 2)
				charset = strings.SplitN(s[1], "=", 2)[1]
			} else if p := mtype.Parent(); p != nil && strings.Contains(p.String(), "text") {
				best, err := det.DetectBest(content)
				if err != nil {
					return "failed_to_detect", Charset
				}
				charset = best.Charset
				info.Cache[charsetIdentifier] = []byte(charset)
			}
			return charset, Charset
		}
	}
}
