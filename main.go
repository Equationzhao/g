package main

import (
	"compress/gzip"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"

	. "github.com/Equationzhao/g/app"
	"github.com/Equationzhao/g/config"
	"github.com/Equationzhao/g/slices"
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
		// load config
		if !slices.Contains(os.Args, "-no-config") {
			defaultArgs, err := config.Load()
			if err == nil {
				os.Args = slices.Insert(os.Args, 1, defaultArgs.Args...)
			} else if _, ok := err.(config.ErrReadConfig); ok {
				_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
			}
		}
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
