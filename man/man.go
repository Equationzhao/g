package man

import (
	"compress/gzip"
	"path/filepath"

	"github.com/Equationzhao/g/internal/cli"
	"github.com/spf13/afero"
)

func GenMDAndMan(fs afero.Fs) {
	// man
	man, _ := fs.Create(filepath.Join("man", "g.1.gz"))
	s, _ := cli.G.ToMan()
	// compress to gzip
	manGz := gzip.NewWriter(man)
	defer manGz.Close()
	_, _ = manGz.Write([]byte(s))
	_ = manGz.Flush()
}
