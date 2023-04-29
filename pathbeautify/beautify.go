package pathbeautify

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Equationzhao/g/cached"
	"github.com/valyala/bytebufferpool"
)

// Transform ~ to $HOME
// ... -> ../..
// .... -> ../../..

func Transform(path *string) {
	switch *path {
	case ".", "..":
	case "...":
		*path = filepath.Join("..", "..")
	case "....":
		*path = filepath.Join("..", "..", "..")
	case "":
	case "~":
		if //goland:noinspection GoBoolExpressions
		runtime.GOOS == "windows" {
			*path = cached.GetUserHomeDir()
		}
	default:
		// ~/a/b/c
		if strings.HasPrefix(*path, "~") {
			home := cached.GetUserHomeDir()
			*path = home + (*path)[1:]
		}

		if strings.HasPrefix(*path, string(filepath.Separator)) {
			return
		}

		// .....?
		// start from 3, aka ...
		matchDots := true
		times := -1
		for _, dot := range *path {
			if dot != '.' {
				if dot != filepath.Separator {
					matchDots = false
				}
				break
			}
			times++
		}

		// case 1
		// .../a/b/c -> times = 2
		// ../../ + a/b/c -> ../../a/b/c
		// case 2
		// ... -> times = 2
		// ../../ + empty -> ../../
		// case 3
		// .../ -> times = 2
		// ../../ + empty -> ../../
		if matchDots {
			const parent = ".."
			buffer := bytebufferpool.Get()
			for i := 0; i < times; i++ {
				_, _ = buffer.WriteString(parent)
				_ = buffer.WriteByte(filepath.Separator)
			}
			if times+2 < len(*path) {
				_, _ = buffer.WriteString((*path)[times+2:])
			}

			*path = buffer.String()
			bytebufferpool.Put(buffer)
		}

	}
}
