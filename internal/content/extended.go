package content

import (
	"encoding/xml"
	"fmt"
	"github.com/valyala/bytebufferpool"

	"github.com/Equationzhao/g/internal/global"
	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/render"
	"github.com/pkg/xattr"
	"howett.net/plist"
)

type ExtendedEnabler struct{}

const Extended = global.NameOfExtended

func formatBytes(bytes []byte) string {
	res := bytebufferpool.Get()
	defer bytebufferpool.Put(res)
	_ = res.WriteByte('[')
	for i, b := range bytes {
		if i > 0 {
			_, _ = res.WriteString(", ")
		}
		_, _ = res.WriteString(fmt.Sprintf("%02x", b))
	}
	_ = res.WriteByte(']')
	return res.String()
}

// formatXattrValue attempts to parse the xattr value and returns a human-readable string.
func formatXattrValue(value []byte) string {
	// Check if the value is a binary plist
	var plistData any
	if _, err := plist.Unmarshal(value, &plistData); err == nil {
		xmlOut, err := xml.MarshalIndent(plistData, "", "  ")
		if err == nil {
			return string(xmlOut)
		}
	}
	// Default to hex encoding
	return formatBytes(value)
}

func (e ExtendedEnabler) Enable(renderer *render.Renderer) NoOutputOption {
	return func(info *item.FileInfo) {
		var list any
		ok := false
		if list, ok = info.Cache[Extended]; !ok {
			list, _ = xattr.LList(info.FullPath)
			info.Cache[Extended] = list
		}
		if lists, ok := list.([]string); ok && len(lists) > 0 {
			for j, key := range lists {
				val, _ := xattr.LGet(info.FullPath, key)
				if j == len(lists)-1 {
					info.AfterLines = append(info.AfterLines, renderer.Extended(fmt.Sprintf("└── %s: %s", key, formatXattrValue(val))))
				} else {
					info.AfterLines = append(info.AfterLines, renderer.Extended(fmt.Sprintf("├── %s: %s", key, formatXattrValue(val))))
				}
			}
		}
	}
}
