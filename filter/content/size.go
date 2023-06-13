package content

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/Equationzhao/g/filter"
	"github.com/Equationzhao/g/osbased"
	"github.com/Equationzhao/g/render"
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

var allSizeUints = []SizeUnit{
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
	for _, sizeUint := range allSizeUints {
		for _, sizeString := range sizeStringSets(sizeUint) {
			if strings.HasSuffix(size, sizeString) {
				size = strings.TrimSuffix(size, sizeString)
				sizeFloat, err := strconv.ParseFloat(size, 64)
				if err != nil {
					return Size{}, err
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
	renderer    *render.Renderer
	*sync.WaitGroup
}

func (s *SizeEnabler) Recursive() *SizeRecursive {
	return s.recursive
}

func (s *SizeEnabler) SetRecursive(sr *SizeRecursive) {
	s.recursive = sr
}

func (s *SizeEnabler) SetRenderer(renderer *render.Renderer) {
	s.renderer = renderer
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
		renderer:    nil,
		recursive:   nil,
		WaitGroup:   new(sync.WaitGroup),
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

func (s *SizeEnabler) Size2String(b int64, blank int) string {
	var res string
	v := float64(b)
	switch s.sizeUint {
	case Bit:
		res = strconv.FormatInt(int64(v*8), 10)
	case B:
		res = strconv.FormatInt(int64(v), 10)
	case KB:
		fallthrough
	case MB:
		fallthrough
	case GB:
		fallthrough
	case TB:
		fallthrough
	case PB:
		fallthrough
	case EB:
		fallthrough
	case ZB:
		fallthrough
	case YB:
		fallthrough
	case BB:
		fallthrough
	case NB:
		res = fmt.Sprintf("%g", v*float64(B)/float64(s.sizeUint))

	case Auto:
		for i := B; i <= NB; i *= 1024 {
			if v < 1024 {
				res = strconv.FormatFloat(v, 'f', 1, 64)
				if res == "0.0" {
					res = "-"
				} else {
					res += Convert2SizeString(i)
				}
				return s.renderer.Size(filter.FillBlank(res, blank))
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
		res += Convert2SizeString(s.sizeUint)
	}
	return s.renderer.Size(filter.FillBlank(res, blank))
}

func recursivelySizeOf(info os.FileInfo, depth int) int64 {
	currentDepth := 0
	if info.IsDir() {
		totalSize := info.Size()
		if depth < 0 {
			// -1 means no limit
			_ = filepath.WalkDir(info.Name(), func(path string, dir os.DirEntry, err error) error {
				if err != nil {
					return err
				}

				if !dir.IsDir() {
					info, err := dir.Info()
					if err == nil {
						totalSize += info.Size()
					}
				}

				return nil
			})
		} else {
			_ = filepath.WalkDir(info.Name(), func(path string, dir os.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if currentDepth > depth {
					if dir.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}

				if !dir.IsDir() {
					info, err := dir.Info()
					if err == nil {
						totalSize += info.Size()
					}
				}

				if dir.IsDir() {
					currentDepth++
				}

				return nil
			})
		}

		return totalSize
	}
	return info.Size()
}

func (s *SizeEnabler) EnableSize(size SizeUnit) filter.ContentOption {
	s.sizeUint = size

	if size != Auto {
		longestSize := 0
		m := sync.RWMutex{}
		done := func(size string) {
			defer s.Done()
			m.RLock()
			if longestSize >= len(size) {
				m.RUnlock()
				return
			}
			m.RUnlock()
			m.Lock()
			if longestSize < len(size) {
				longestSize = len(size)
			}
			m.Unlock()
		}

		wait := func(size string) string {
			s.Wait()
			return filter.FillBlank(size, longestSize)
		}

		return func(info os.FileInfo) (string, string) {
			var v int64
			if s.recursive != nil {
				v = recursivelySizeOf(info, s.recursive.depth)
			} else {
				v = info.Size()
			}
			if s.enableTotal {
				s.total.Add(v)
			}
			size := s.Size2String(v, 0)
			done(size)
			return wait(size), SizeName
		}
	}

	return func(info os.FileInfo) (string, string) {
		var v int64
		if s.recursive != nil {
			v = recursivelySizeOf(info, s.recursive.depth)
		} else {
			v = info.Size()
		}
		if s.enableTotal {
			s.total.Add(v)
		}
		return s.Size2String(v, 7), SizeName
	}
}

type BlockSizeEnabler struct {
	renderer *render.Renderer
	*sync.WaitGroup
}

func NewBlockSizeEnabler() *BlockSizeEnabler {
	return &BlockSizeEnabler{
		renderer:  nil,
		WaitGroup: new(sync.WaitGroup),
	}
}

func (b *BlockSizeEnabler) SetRenderer(r *render.Renderer) {
	b.renderer = r
}

const BlockSizeName = "Blocks"

func (b *BlockSizeEnabler) Enable() filter.ContentOption {

	longestSize := 0
	m := sync.RWMutex{}
	done := func(size string) {
		defer b.Done()
		m.RLock()
		if longestSize >= len(size) {
			m.RUnlock()
			return
		}
		m.RUnlock()
		m.Lock()
		if longestSize < len(size) {
			longestSize = len(size)
		}
		m.Unlock()
	}

	wait := func(size string) string {
		b.Wait()
		return filter.FillBlank(size, longestSize)
	}

	return func(info os.FileInfo) (string, string) {
		bs := osbased.BlockSize(info)
		res := b.renderer.Size(strconv.FormatInt(bs, 10))
		done(res)
		return wait(res), BlockSizeName
	}
}
