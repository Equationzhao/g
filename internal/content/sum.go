package content

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"os"
	"strings"

	constval "github.com/Equationzhao/g/internal/const"
	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/render"
)

type SumType int

const (
	SumTypeMd5 SumType = iota + 1
	SumTypeSha1
	SumTypeSha224
	SumTypeSha256
	SumTypeSha384
	SumTypeSha512
	SumTypeCRC32
)

const SumName = constval.NameOfSum

type SumEnabler struct{}

// todo simd
func (s SumEnabler) EnableSum(renderer *render.Renderer, sumTypes ...SumType) ContentOption {
	length := 0
	types := make([]string, 0, len(sumTypes))
	for _, t := range sumTypes {
		switch t {
		case SumTypeMd5:
			length += 32
			types = append(types, "md5")
		case SumTypeSha1:
			length += 40
			types = append(types, "sha1")
		case SumTypeSha224:
			length += 56
			types = append(types, "sha224")
		case SumTypeSha256:
			length += 64
			types = append(types, "sha256")
		case SumTypeSha384:
			length += 96
			types = append(types, "sha384")
		case SumTypeSha512:
			length += 128
			types = append(types, "sha512")
		case SumTypeCRC32:
			length += 8
			types = append(types, "crc32")
		}
	}
	length += len(sumTypes) - 1
	sumName := fmt.Sprintf("%s(%s)", SumName, strings.Join(types, ","))
	return func(info *item.FileInfo) (string, string) {
		if info.IsDir() {
			return FillBlank("", length), sumName
		}

		var content []byte
		if content_, ok := info.Cache["content"]; ok {
			content = content_
		} else {
			file, err := os.Open(info.FullPath)
			if err != nil {
				return FillBlank("", length), sumName
			}
			content, err = io.ReadAll(file)
			if err != nil {
				return FillBlank("", length), sumName
			}
			info.Cache["content"] = content
			defer file.Close()
		}

		hashes := make([]hash.Hash, 0, len(sumTypes))
		writers := make([]io.Writer, 0, len(sumTypes))
		for _, t := range sumTypes {
			var hashed hash.Hash
			switch t {
			case SumTypeMd5:
				hashed = md5.New()
			case SumTypeSha1:
				hashed = sha1.New()
			case SumTypeSha224:
				hashed = sha256.New224()
			case SumTypeSha256:
				hashed = sha256.New()
			case SumTypeSha384:
				hashed = sha512.New384()
			case SumTypeSha512:
				hashed = sha512.New()
			case SumTypeCRC32:
				hashed = crc32.NewIEEE()
			}
			writers = append(writers, hashed)
			hashes = append(hashes, hashed)
		}
		multiWriter := io.MultiWriter(writers...)
		if _, err := io.Copy(multiWriter, bytes.NewReader(content)); err != nil {
			return FillBlank("", length), sumName
		}
		sums := make([]string, 0, len(hashes))
		for _, h := range hashes {
			sums = append(sums, fmt.Sprintf("%x", h.Sum(nil)))
		}
		sumsStr := strings.Join(sums, " ")
		return renderer.Checksum(FillBlank(sumsStr, length)), sumName
	}
}
