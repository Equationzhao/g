package main

import (
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
	Version = "0.5.3"

	if doc {
		md, _ := os.Create("g.md")
		s, _ := G.ToMarkdown()
		_, _ = fmt.Fprintln(md, s)
		man, _ := os.Create(filepath.Join("man", "g.1"))
		s, _ = G.ToMan()
		_, _ = fmt.Fprintln(man, s)
	} else {
		err := G.Run(os.Args)
		if err != nil {
			if !errors.Is(err, Err4Exit{}) {
				if ReturnCode == 0 {
					ReturnCode = 1
				}
				_, _ = fmt.Fprint(os.Stderr, MakeErrorStr(err.Error()))
			}
		}
	}
}
