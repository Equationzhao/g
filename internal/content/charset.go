package content

import (
	"os"
	"strings"

	constval "github.com/Equationzhao/g/internal/global"
	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/render"
	"github.com/Equationzhao/g/internal/util"
	"github.com/gabriel-vasile/mimetype"
	"github.com/saintfish/chardet"
)

type CharsetEnabler struct{}

func NewCharsetEnabler() *CharsetEnabler {
	return &CharsetEnabler{}
}

const (
	Charset = constval.NameOfCharset
)

func (c *CharsetEnabler) Enable(renderer *render.Renderer) ContentOption {
	det := chardet.NewTextDetector()
	return func(info *item.FileInfo) (string, string) {
		res, returnName := func() (string, string) {
			if c, ok := info.Cache[Charset]; ok {
				return string(c), Charset
			}
			// only text file has charset
			if info.IsDir() {
				return "-", Charset
			} else if util.IsSymLink(info) {
				return "-", Charset
			} else if info.Mode()&os.ModeNamedPipe != 0 {
				return "-", Charset
			} else if info.Mode()&os.ModeSocket != 0 {
				return "-", Charset
			} else {
				mtype, err := mimetype.DetectFile(info.FullPath)
				if err != nil {
					return "failed_to_read", Charset
				}
				charset := "-"
				if tn := mtype.String(); strings.Contains(tn, ";") {
					s := strings.SplitN(tn, ";", 2)
					charset = strings.SplitN(s[1], "=", 2)[1]
				} else if p := mtype.Parent(); p != nil && strings.Contains(p.String(), "text") {
					file, err := os.Open(info.FullPath)
					if err != nil {
						return "failed_to_read", Charset
					}
					defer file.Close()
					content := make([]byte, 1024*1024)
					_, err = file.Read(content)
					if err != nil {
						return err.Error(), Charset
					}
					best, err := det.DetectBest(content)
					if err != nil {
						return "failed_to_detect", Charset
					}
					charset = best.Charset
					info.Cache[Charset] = []byte(charset)
				}
				return charset, Charset
			}
		}()
		return renderer.Charset(res), returnName
	}
}
