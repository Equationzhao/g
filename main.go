package main

import (
	"compress/gzip"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"slices"

	"github.com/Equationzhao/g/internal/cli"
	"github.com/Equationzhao/g/internal/config"
	"github.com/Equationzhao/g/internal/const/doc"
	"github.com/Equationzhao/g/internal/util"
)

func main() {
	// catch panic and print stack trace and version info
	defer func() {
		if err := recover(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Version: v%s\n", cli.Version)
			_, _ = fmt.Fprintf(os.Stderr, "Please file an issue at %s with the following panic info\n\n", util.MakeLink("https://github.com/Equationzhao/g/issues/new/choose", "Github Repo"))
			_, _ = fmt.Fprintln(os.Stderr, cli.MakeErrorStr(fmt.Sprintf("error message:\n%v\n", err)))
			_, _ = fmt.Fprintln(os.Stderr, cli.MakeErrorStr(fmt.Sprintf("stack trace:\n%s", debug.Stack())))
			cli.ReturnCode = 2
		}
		os.Exit(cli.ReturnCode)
	}()
	// when build with tag `doc`, generate md and man file
	if doc.Enable {
		// md
		md, err := os.Create("g.md")
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
	} else {
		// normal logic
		// load config if the args do not contains -no-config
		if !slices.ContainsFunc(os.Args, match) {
			defaultArgs, err := config.Load()
			// if successfully load config and **the config.Args do not contain -no-config**
			if err == nil && !slices.ContainsFunc(defaultArgs.Args, match) {
				os.Args = slices.Insert(os.Args, 1, defaultArgs.Args...)
			} else if err != nil { // if failed to load config
				// if it's read error
				var errReadConfig config.ErrReadConfig
				if errors.As(err, &errReadConfig) {
					_, _ = fmt.Fprintln(os.Stderr, cli.MakeErrorStr(err.Error()))
				}
			}
		} else {
			// contains -no-config
			// remove it before the cli.G starts
			os.Args = slices.DeleteFunc(os.Args, match)
		}
		err := cli.G.Run(os.Args)
		if err != nil {
			if !errors.Is(err, cli.Err4Exit{}) {
				if cli.ReturnCode == 0 {
					cli.ReturnCode = 1
				}
				_, _ = fmt.Fprintln(os.Stderr, cli.MakeErrorStr(err.Error()))
			}
		}
	}
}

func match(s string) bool {
	if s == "-no-config" || s == "--no-config" {
		return true
	}
	return false
}
