package content

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"os"

	"github.com/Equationzhao/g/filter"
	"github.com/Equationzhao/g/util"
	"github.com/Equationzhao/tsmap"
)

type (
	filenameList = *util.Slice[string]
	hashStr      = string
)

type DuplicateDetect struct {
	IsThrough bool
	hashTb    *tsmap.Map[hashStr, filenameList]
}

type DOption func(d *DuplicateDetect)

const defaultTbSize = 200

func NewDuplicateDetect(options ...DOption) *DuplicateDetect {
	d := &DuplicateDetect{}

	for _, option := range options {
		option(d)
	}

	if d.hashTb == nil {
		d.hashTb = tsmap.NewTSMap[hashStr, filenameList](defaultTbSize)
	}

	return d
}

func DuplicateWithTbSize(size int) DOption {
	return func(d *DuplicateDetect) {
		d.hashTb = tsmap.NewTSMap[hashStr, filenameList](size)
	}
}

func DetectorFallthrough(d *DuplicateDetect) {
	d.IsThrough = true
}

func (d *DuplicateDetect) Enable() filter.NoOutputOption {
	return func(info os.FileInfo) {
		afterHash, err := fileHash(info, d.IsThrough)
		if err != nil {
			return
		}
		actual, _ := d.hashTb.GetOrInit(afterHash, func() filenameList {
			return util.NewSlice[string](10)
		})
		actual.AppendTo(info.Name())
	}
}

type Duplicate struct {
	Filenames []string
}

func (d *DuplicateDetect) Result() []Duplicate {
	list := d.hashTb.Values()
	res := make([]Duplicate, 0, len(list))
	for _, i := range list {
		if l := i.Len(); l > 1 {
			res = append(res, Duplicate{Filenames: i.GetCopy()})
		}
	}
	return res
}

func (d *DuplicateDetect) Reset() {
	filenameLists := d.hashTb.Values()
	for _, list := range filenameLists {
		list.Clear()
	}
}

func (d *DuplicateDetect) Fprint(w io.Writer) {
	r := d.Result()
	if len(r) != 0 {
		_, _ = fmt.Fprintln(w, "Duplicates:")
		for _, i := range r {
			for _, filename := range i.Filenames {
				_, _ = fmt.Fprint(w, "    ", filename)
			}
			_, _ = fmt.Fprintln(w)
		}
	}
}

var thresholdFileSize = int64(16 * KB)

// fileHash calculates the hash of the file provided.
// If isThorough is true, then it uses SHA256 of the entire file.
// Otherwise, it uses CRC32 of "crucial bytes" of the file.
func fileHash(fileInfo os.FileInfo, isThorough bool) (string, error) {
	if !fileInfo.Mode().IsRegular() {
		return "", fmt.Errorf("can't compute hash of non-regular file")
	}
	var prefix string
	var bytes []byte
	var fileReadErr error
	if isThorough {
		bytes, fileReadErr = os.ReadFile(fileInfo.Name())
	} else if fileInfo.Size() <= thresholdFileSize {
		prefix = "f"
		bytes, fileReadErr = os.ReadFile(fileInfo.Name())
	} else {
		prefix = "s"
		bytes, fileReadErr = readCrucialBytes(fileInfo.Name(), fileInfo.Size())
	}
	if fileReadErr != nil {
		return "", fmt.Errorf("couldn't calculate hash: %w", fileReadErr)
	}
	var h hash.Hash
	if isThorough {
		h = sha256.New()
	} else {
		h = crc32.NewIEEE()
	}
	_, hashErr := h.Write(bytes)
	if hashErr != nil {
		return "", fmt.Errorf("error while computing hash: %w", hashErr)
	}
	hashBytes := h.Sum(nil)
	return prefix + hex.EncodeToString(hashBytes), nil
}

// readCrucialBytes reads the first few bytes, middle bytes and last few bytes of the file
func readCrucialBytes(filePath string, fileSize int64) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	firstBytes := make([]byte, thresholdFileSize/2)
	_, fErr := file.ReadAt(firstBytes, 0)
	if fErr != nil {
		return nil, fmt.Errorf("couldn't read first few bytes (maybe file is corrupted?): %w", fErr)
	}
	middleBytes := make([]byte, thresholdFileSize/4)
	_, mErr := file.ReadAt(middleBytes, fileSize/2)
	if mErr != nil {
		return nil, fmt.Errorf("couldn't read middle bytes (maybe file is corrupted?): %w", mErr)
	}
	lastBytes := make([]byte, thresholdFileSize/4)
	_, lErr := file.ReadAt(lastBytes, fileSize-thresholdFileSize/4)
	if lErr != nil {
		return nil, fmt.Errorf("couldn't read end bytes (maybe file is corrupted?): %w", lErr)
	}
	bytes := append(append(firstBytes, middleBytes...), lastBytes...)
	return bytes, nil
}
