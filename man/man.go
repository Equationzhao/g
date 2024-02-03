package man

import (
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Equationzhao/g/internal/cli"
)

func GenMDAndMan() {
	// md
	md, err := os.Create(filepath.Join("docs", "g.md"))
	if err != nil {
		panic(err)
	}
	defer md.Close()
	s, _ := cli.G.ToMarkdown()
	_, _ = fmt.Fprintln(md, s)
	// man
	man, _ := os.Create(filepath.Join("man", "g.1.gz"))
	s, _ = cli.G.ToMan()
	// compress to gzip
	manGz := gzip.NewWriter(man)
	defer manGz.Close()
	_, _ = manGz.Write([]byte(s))
	_ = manGz.Flush()
}
