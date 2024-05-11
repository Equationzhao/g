package man

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Equationzhao/g/internal/cli"
)

func GenMDAndMan() {
	// md
	g, err := os.Create(filepath.Join("docs", "g.md"))
	if err != nil {
		panic(err)
	}
	defer g.Close()
	md, err := os.Create(filepath.Join("docs", "man.md"))
	if err != nil {
		panic(err)
	}
	defer md.Close()
	s, _ := cli.G.ToMarkdown()
	_, _ = fmt.Fprintln(io.MultiWriter(md, g), s)
	// man
	man, _ := os.Create(filepath.Join("man", "g.1.gz"))
	s, _ = cli.G.ToMan()
	// compress to gzip
	manGz := gzip.NewWriter(man)
	defer manGz.Close()
	_, _ = manGz.Write([]byte(s))
	_ = manGz.Flush()
}
