package display

import (
	"fmt"
	"strings"

	"github.com/Equationzhao/g/internal/align"
	constval "github.com/Equationzhao/g/internal/global"
	"github.com/Equationzhao/g/internal/item"
	"github.com/valyala/bytebufferpool"
)

type HeaderMaker struct {
	Header, Footer  bool
	IsBefore        bool
	LongestEachPart map[string]int
	AllPart         []string
}

func (h HeaderMaker) Make(p Printer, Items ...*item.FileInfo) {
	// add header
	if len(Items) == 0 {
		return
	}

	// add longest - len(header) * space
	// print header
	headerFooterStrBuf := bytebufferpool.Get()
	defer bytebufferpool.Put(headerFooterStrBuf)
	prettyPrinter, isPrettyPrinter := p.(PrettyPrinter)

	expand := func(s string, no, space int) {
		// left align
		if no != len(h.AllPart)-1 && align.IsLeftHeaderFooter(s) {
			_, _ = headerFooterStrBuf.WriteString(strings.Repeat(" ", space-1)) // remove the additional following space for right align
		}
		_, _ = headerFooterStrBuf.WriteString(constval.Underline)
		_, _ = headerFooterStrBuf.WriteString(s)
		_, _ = headerFooterStrBuf.WriteString(constval.Reset)
		if no != len(h.AllPart)-1 {
			if !align.IsLeftHeaderFooter(s) {
				_, _ = headerFooterStrBuf.WriteString(strings.Repeat(" ", space))
			} else {
				_, _ = headerFooterStrBuf.WriteString(strings.Repeat(" ", 1)) // still need the following space for left align
			}
		}
	}

	for i, s := range h.AllPart {
		if len(s) > h.LongestEachPart[s] {
			// expand the every item's content of this part
			for _, it := range Items {
				content, _ := it.Get(s)
				if s != constval.NameOfName {
					toAddNum := len(s) - WidthNoHyperLinkLen(content.String())
					if align.IsLeft(s) {
						content.AddSuffix(strings.Repeat(" ", toAddNum))
					} else {
						content.AddPrefix(strings.Repeat(" ", toAddNum))
					}
				}
				it.Set(s, content)
				h.LongestEachPart[s] = len(s)
			}
			expand(s, i, 1)
		} else {
			expand(s, i, h.LongestEachPart[s]-len(s)+1)
		}
		if isPrettyPrinter && h.IsBefore {
			if h.Header {
				prettyPrinter.AddHeader(s)
			}
			if h.Footer {
				prettyPrinter.AddFooter(s)
			}
		}
	}
	res := headerFooterStrBuf.String()
	if !isPrettyPrinter {
		if h.Header && h.IsBefore {
			_, _ = fmt.Fprintln(p, res)
		}
		if h.Footer && !h.IsBefore {
			_, _ = fmt.Fprintln(p, res)
		}
	}
}
