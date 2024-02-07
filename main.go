package main

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"slices"

	"github.com/Equationzhao/g/internal/cli"
	"github.com/Equationzhao/g/internal/config"
	debugSetting "github.com/Equationzhao/g/internal/const/debug"
	"github.com/Equationzhao/g/internal/const/doc"
	"github.com/Equationzhao/g/internal/util"
	"github.com/Equationzhao/g/man"
)

func main() {
	// catch panic and print stack trace and version info
	defer func() {
		if !debugSetting.Enable {
			catchPanic(recover())
		}
	}()
	// when build with tag `doc`, generate md and man file
	if doc.Enable {
		man.GenMDAndMan()
	} else {
		preprocessArgs()
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

func catchPanic(err any) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Version: v%s\n", cli.Version)
		_, _ = fmt.Fprintf(os.Stderr, "Please file an issue at %s with the following panic info\n\n", util.MakeLink("https://github.com/Equationzhao/g/issues/new/choose", "Github Repo"))
		_, _ = fmt.Fprintln(os.Stderr, cli.MakeErrorStr(fmt.Sprintf("error message:\n%v\n", err)))
		_, _ = fmt.Fprintln(os.Stderr, cli.MakeErrorStr(fmt.Sprintf("stack trace:\n%s", debug.Stack())))
		cli.ReturnCode = 1
	}
	os.Exit(cli.ReturnCode)
}

func preprocessArgs() {
	// normal logic
	// load config if the args do not contains -no-config
	if !slices.ContainsFunc(os.Args, hasNoConfig) {
		defaultArgs, err := config.Load()
		// if successfully load config and **the config.Args do not contain -no-config**
		if err == nil && !slices.ContainsFunc(defaultArgs.Args, hasNoConfig) {
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
		os.Args = slices.DeleteFunc(os.Args, hasNoConfig)
	}
}

func hasNoConfig(s string) bool {
	if s == "-no-config" || s == "--no-config" {
		return true
	}
	return false
}
