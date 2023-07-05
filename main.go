package main

import (
	"compress/gzip"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"

	. "github.com/Equationzhao/g/app"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(Version)
			fmt.Println(MakeErrorStr(fmt.Sprint(err)))
			fmt.Println(MakeErrorStr(string(debug.Stack())))
		}
		os.Exit(ReturnCode)
	}()

	if doc {
		md, _ := os.Create("g.md")
		s, _ := G.ToMarkdown()
		_, _ = fmt.Fprintln(md, s)
		man, _ := os.Create(filepath.Join("man", "g.1.gz"))
		s, _ = G.ToMan()
		// compress to gzip
		manGz := gzip.NewWriter(man)
		defer manGz.Close()
		_, _ = manGz.Write([]byte(s))
		_ = manGz.Flush()
	} else {
		err := G.Run(os.Args)
		if err != nil {
			if !errors.Is(err, Err4Exit{}) {
				if ReturnCode == 0 {
					ReturnCode = 1
				}
				_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
			}
		}
	}
}
