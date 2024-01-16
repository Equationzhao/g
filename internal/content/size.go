package content

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/osbased"
	"github.com/Equationzhao/g/internal/render"
	"github.com/Equationzhao/g/internal/util"
)

const Unknown SizeUnit = -1

const (
	Auto SizeUnit = iota
	Bit  SizeUnit = 1.0 << (10 * iota)
	B
	KB
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
	BB
	NB
)

const SizeName = "Size"

type Size struct {
	Bytes uint64
}

var allSizeUnits = []SizeUnit{
	Bit,
	B,
	KB,
	MB,
	GB,
	TB,
	PB,
	EB,
	ZB,
	YB,
	BB,
	NB,
}

func ParseSize(size string) (Size, error) {
	for _, sizeUint := range allSizeUnits {
		sets := sizeStringSets(sizeUint)
		for _, sizeString := range sets {
			if strings.HasSuffix(size, sizeString) {
				size = strings.TrimSuffix(size, sizeString)
				sizeFloat, err := strconv.ParseFloat(size, 64)
				if err != nil {
					continue
				}
				return Size{
					Bytes: uint64(sizeFloat * float64(sizeUint)),
				}, nil
			}
		}
	}
	return Size{}, fmt.Errorf("unknown size unit")
}

type SizeUnit float64

func sizeStringSets(size SizeUnit) []string {
	switch size {
	case Unknown:
		return []string{"unknown"}
	case Auto:
		return []string{""}
	case Bit:
		return []string{"bit", "Bit"}
	case B:
		return []string{"B", "b", "byte", "Byte", "BYTE"}
	case KB:
		return []string{"KB", "kb", "K", "k"}
	case MB:
		return []string{"MB", "mb", "M", "m"}
	case GB:
		return []string{"GB", "gb", "G", "g"}
	case TB:
		return []string{"TB", "tb", "T", "t"}
	case PB:
		return []string{"PB", "pb", "P", "p"}
	case EB:
		return []string{"EB", "eb", "E", "e"}
	case ZB:
		return []string{"ZB", "zb", "Z", "z"}
	case YB:
		return []string{"YB", "yb", "Y", "y"}
	case BB:
		return []string{"BB", "bb"}
	case NB:
		return []string{"NB", "nb", "N", "n"}
	default:
		return []string{"unknown"}
	}
}

func Convert2SizeString(size SizeUnit) string {
	return sizeStringSets(size)[0]
}

func ConvertFromSizeString(size string) SizeUnit {
	switch size {
	case "bit", "Bit", "BIT":
		return Bit
	case "B", "b", "byte", "Byte", "BYTE":
		return B
	case "KB", "kb", "Kb", "k":
		return KB
	case "MB", "mb", "Mb", "M", "m":
		return MB
	case "GB", "gb", "Gb", "G", "g":
		return GB
	case "TB", "tb", "Tb", "T", "t":
		return TB
	case "PB", "pb", "Pb", "P", "p":
		return PB
	case "EB", "eb", "Eb", "E", "e":
		return EB
	case "ZB", "zb", "Zb", "Z", "z":
		return ZB
	case "YB", "yb", "Yb", "Y", "y":
		return YB
	case "BB", "bb", "Bb":
		return BB
	case "NB", "nb", "Nb", "N", "n":
		return NB
	default:
		return Unknown
	}
}

type SizeEnabler struct {
	total       atomic.Int64
	enableTotal bool
	sizeUint    SizeUnit
	recursive   *SizeRecursive
}

func (s *SizeEnabler) Recursive() *SizeRecursive {
	return s.recursive
}

func (s *SizeEnabler) SetRecursive(sr *SizeRecursive) {
	s.recursive = sr
}

type SizeRecursive struct {
	depth int
}

func NewSizeRecursive(depth int) *SizeRecursive {
	return &SizeRecursive{depth: depth}
}

func NewSizeEnabler() *SizeEnabler {
	return &SizeEnabler{
		total:       atomic.Int64{},
		enableTotal: false,
		sizeUint:    Auto,
		recursive:   nil,
	}
}

func (s *SizeEnabler) SizeUint() SizeUnit {
	return s.sizeUint
}

func (s *SizeEnabler) SetEnableTotal() {
	s.enableTotal = true
}

func (s *SizeEnabler) DisableTotal() {
	s.enableTotal = false
}

func (s *SizeEnabler) Total() (size int64, ok bool) {
	if s.enableTotal {
		return s.total.Load(), s.enableTotal
	}
	return 0, false
}

func (s *SizeEnabler) Reset() {
	if s.enableTotal {
		s.total.Store(0)
	}
}

func (s *SizeEnabler) Size2String(b int64) (string, SizeUnit) {
	var res string
	actualUnit := s.sizeUint
	v := float64(b)
	switch s.sizeUint {
	case Bit:
		res = strconv.FormatInt(int64(v*8), 10)
	case B:
		res = strconv.FormatInt(int64(v), 10)
	case KB, MB, GB, TB, PB, EB, ZB, YB, BB, NB:
		res = fmt.Sprintf("%.1f", v*float64(B)/float64(s.sizeUint))
	case Auto:
		for i := B; i <= NB; i *= 1024 {
			if v < 1024 {
				res = strconv.FormatFloat(v, 'f', 1, 64)
				if res == "0.0" {
					// make align
					return "       -", actualUnit
				}
				res += " " + Convert2SizeString(i)
				actualUnit = i
				return FillBlank(res, 8), actualUnit
			}
			v /= 1024
		}
		panic("too large")
	default:
		panic("invalid size uint" + strconv.Itoa(int(s.sizeUint)))
	}

	if res == "0" {
		res = "-"
	} else {
		res += " " + Convert2SizeString(s.sizeUint)
	}
	return res, actualUnit
}

const RecursiveSizeName = "recursive_size"

func (s *SizeEnabler) EnableSize(size SizeUnit, renderer *render.Renderer) ContentOption {
	s.sizeUint = size
	return func(info *item.FileInfo) (string, string) {
		var v int64
		if s.recursive != nil {
			if r, ok := info.Cache[RecursiveSizeName]; ok {
				// convert []byte to int64
				v, _ = strconv.ParseInt(string(r), 10, 64)
			} else {
				v = util.RecursivelySizeOf(info, s.recursive.depth)
			}
		} else {
			v = info.Size()
		}
		if s.enableTotal {
			s.total.Add(v)
		}
		res, unit := s.Size2String(v)
		return renderer.Size(res, Convert2SizeString(unit)), SizeName
	}
}

type BlockSizeEnabler struct{}

func NewBlockSizeEnabler() *BlockSizeEnabler {
	return &BlockSizeEnabler{}
}

const BlockSizeName = "Blocks"

func (b *BlockSizeEnabler) Enable(renderer *render.Renderer) ContentOption {
	return func(info *item.FileInfo) (string, string) {
		res := ""
		bs := osbased.BlockSize(info)
		if bs == 0 {
			res = "-"
		} else {
			res = strconv.FormatInt(bs, 10)
		}
		return renderer.BlockSize(res), BlockSizeName
	}
}
