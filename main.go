package main

import (
	"compress/gzip"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"slices"

	"github.com/Equationzhao/g/app"
	"github.com/Equationzhao/g/config"
	"github.com/Equationzhao/g/util"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Version: v%s\n", app.Version)
			fmt.Printf("Please file an issue at %s with the following panic info\n\n", util.MakeLink("https://github.com/Equationzhao/g/issues/new/choose", "Github Repo"))
			fmt.Println(app.MakeErrorStr(fmt.Sprintf("error message:\n%v\n", err)))
			fmt.Println(app.MakeErrorStr(fmt.Sprintf("stack trace:\n%s", debug.Stack())))
			if app.ReturnCode == 0 {
				app.ReturnCode = 2
			}
		}
		os.Exit(app.ReturnCode)
	}()
	if doc {
		md, _ := os.Create("g.md")
		s, _ := app.G.ToMarkdown()
		_, _ = fmt.Fprintln(md, s)
		man, _ := os.Create(filepath.Join("man", "g.1.gz"))
		s, _ = app.G.ToMan()
		// compress to gzip
		manGz := gzip.NewWriter(man)
		defer manGz.Close()
		_, _ = manGz.Write([]byte(s))
		_ = manGz.Flush()
	} else {
		// load config
		if !slices.ContainsFunc(
			os.Args, match,
		) {
			defaultArgs, err := config.Load()
			if err == nil && !slices.ContainsFunc(
				defaultArgs.Args, match,
			) {
				os.Args = slices.DeleteFunc(
					os.Args, match,
				)
				os.Args = slices.Insert(os.Args, 1, defaultArgs.Args...)
			} else {
				var errReadConfig config.ErrReadConfig
				if errors.As(err, &errReadConfig) {
					_, _ = fmt.Fprintln(os.Stderr, app.MakeErrorStr(err.Error()))
				}
			}
		} else {
			// contains -no-config
			// remove it
			os.Args = slices.DeleteFunc(
				os.Args, match,
			)
		}
		err := app.G.Run(os.Args)
		if err != nil {
			if !errors.Is(err, app.Err4Exit{}) {
				if app.ReturnCode == 0 {
					app.ReturnCode = 1
				}
				_, _ = fmt.Fprintln(os.Stderr, app.MakeErrorStr(err.Error()))
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
