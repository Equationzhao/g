package main

import (
	"errors"
	"fmt"
	"github.com/Equationzhao/g/filter"
	"github.com/Equationzhao/g/pathbeautify"
	"github.com/Equationzhao/g/printer"
	"github.com/Equationzhao/g/render"
	"github.com/Equationzhao/g/theme"
	"github.com/Equationzhao/g/tree"
	"github.com/urfave/cli/v2"
	"os"
)

var typeFunc = make([]*filter.TypeFunc, 0)
var contentFunc = make([]filter.ContentOption, 0)
var r = render.NewRenderer(theme.DefaultTheme, theme.DefaultInfoTheme)
var p printer.Printer = printer.NewFitTerminal()
var timeFormat = "02.Jan'06 15:04"

const version = "v0.2.1"

func init() {
	typeFunc = append(typeFunc, &filter.RemoveHidden)
}

func main() {

	app := &cli.App{
		Name:      "gverything",
		Usage:     "a powerful ls",
		UsageText: "g [options] [path]",
		Version:   version,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "tree",
				Aliases: []string{"t"},
				Usage:   "list in tree",
			},
			&cli.IntFlag{
				Name:        "depth",
				Usage:       "tree limit depth, negative -> infinity",
				DefaultText: "infinity",
				Value:       -1,
			},
			&cli.BoolFlag{
				Name:    "show-hidden",
				Aliases: []string{"sh", "a"},
				Usage:   "show hidden files",
				Action: func(context *cli.Context, b bool) error {
					if b {
						// remove filter.RemoveHidden
						newFF := make([]*filter.TypeFunc, 0, len(typeFunc))
						for _, typeFunc := range typeFunc {
							if typeFunc != &filter.RemoveHidden {
								newFF = append(newFF, typeFunc)
							}
						}
						typeFunc = newFF
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:    "A",
				Aliases: []string{"almost-all"},
				Usage:   "do not list implied . and ..",
				Action: func(context *cli.Context, b bool) error {
					if b {
						// remove filter.RemoveHidden
						newFF := make([]*filter.TypeFunc, 0, len(typeFunc))
						for _, typeFunc := range typeFunc {
							if typeFunc != &filter.RemoveHidden {
								newFF = append(newFF, typeFunc)
							}
						}
						typeFunc = newFF
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:    "time-format",
				Usage:   "time/date format with -l",
				EnvVars: []string{"TIME_STYLE"},
				Action: func(context *cli.Context, s string) error {
					/*
						The TIME_STYLE argument can be full-iso, long-iso, iso, locale, or  +FORMAT.   FORMAT
						is  interpreted  like in date(1).  If FORMAT is FORMAT1<newline>FORMAT2, then FORMAT1
						applies to non-recent files and FORMAT2 to recent files.   TIME_STYLE  prefixed  with
						'posix-' takes effect only outside the POSIX locale.  Also the TIME_STYLE environment
						variable sets the default style to use.
					*/
					switch s {
					case "full-iso":
						timeFormat = "2006-01-02 15:04:05.000000000 -0700"
					case "long-iso":
						timeFormat = "2006-01-02 15:04"
					case "locale":
						timeFormat = "Jan 02 15:04"
					default:

					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:  "full-time",
				Usage: "like -all/l --time-style=full-iso",
				Action: func(context *cli.Context, b bool) error {
					if b {
						timeFormat = "2006-01-02 15:04:05.000000000 -0700"
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:    "show-perm",
				Aliases: []string{"sp"},
				Usage:   "show permission",
				Action: func(context *cli.Context, b bool) error {
					if b {
						contentFunc = append(contentFunc, filter.EnableFileMode(r))
						if _, ok := p.(*printer.Byline); !ok {
							p = printer.NewByline()
						}
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:    "show-size",
				Aliases: []string{"ss"},
				Usage:   "show file/dir size",
				Action: func(context *cli.Context, b bool) error {
					if b {
						contentFunc = append(contentFunc, filter.EnableSize(filter.Auto, r))
						if _, ok := p.(*printer.Byline); !ok {
							p = printer.NewByline()
						}
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:    "show-owner",
				Aliases: []string{"so", "author"},
				Usage:   "show owner",
				Action: func(context *cli.Context, b bool) error {
					if b {
						contentFunc = append(contentFunc, filter.EnableOwner(r))
						if _, ok := p.(*printer.Byline); !ok {
							p = printer.NewByline()
						}
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:    "B",
				Aliases: []string{"ignore-backups"},
				Usage:   "do not list implied entries ending with ~",
				Action: func(context *cli.Context, b bool) error {
					if b {
						typeFunc = append(typeFunc, &filter.RemoveBackups)
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:    "show-group",
				Aliases: []string{"sg"},
				Usage:   "show group",
				Action: func(context *cli.Context, b bool) error {
					if b {
						contentFunc = append(contentFunc, filter.EnableGroup(r))
						if _, ok := p.(*printer.Byline); !ok {
							p = printer.NewByline()
						}
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:    "show-time",
				Aliases: []string{"st"},
				Usage:   "show time",
				Action: func(context *cli.Context, b bool) error {
					if b {
						contentFunc = append(contentFunc, filter.EnableTime(timeFormat, r))
						if _, ok := p.(*printer.Byline); !ok {
							p = printer.NewByline()
						}
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:    "show-icon",
				Usage:   "show icon",
				Aliases: []string{"si"},
			},
			&cli.BoolFlag{
				Name:    "byline",
				Aliases: []string{"bl", "1"},
				Usage:   "print by line",
				Action: func(context *cli.Context, b bool) error {
					if b {
						if _, ok := p.(*printer.Byline); !ok {
							p = printer.NewByline()
						}
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:    "zero",
				Aliases: []string{"0"},
				Usage:   "end each output line with NUL, not newline",
				Action: func(context *cli.Context, b bool) error {
					if b {
						if _, ok := p.(*printer.Zero); !ok {
							p = printer.NewZero()
						}
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:  "m",
				Usage: "fill width with a comma separated list of entries",
				Action: func(context *cli.Context, b bool) error {
					if b {
						if _, ok := p.(*printer.Byline); !ok {
							p = printer.NewCommaPrint()
						}
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:  "x",
				Usage: "list entries by lines instead of by columns",
				Action: func(context *cli.Context, b bool) error {
					if b {
						if _, ok := p.(*printer.Across); !ok {
							p = printer.NewAcross()
						}
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"la", "l"},
				Usage:   "show all info/use a long listing format",
				Action: func(context *cli.Context, b bool) error {
					if b {
						// remove filter.RemoveHidden
						newFF := make([]*filter.TypeFunc, 0, len(typeFunc))
						for _, typeFunc := range typeFunc {
							if typeFunc != &filter.RemoveHidden {
								newFF = append(newFF, typeFunc)
							}
						}
						typeFunc = newFF
						contentFunc = append(contentFunc, filter.EnableFileMode(r), filter.EnableSize(filter.Auto, r), filter.EnableOwner(r), filter.EnableGroup(r), filter.EnableTime(timeFormat, r))
						if _, ok := p.(*printer.Byline); !ok {
							p = printer.NewByline()
						}
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:  "C",
				Usage: "list entries by columns",
				Action: func(context *cli.Context, b bool) error {
					if b {
						if _, ok := p.(*printer.FitTerminal); !ok {
							p = printer.NewFitTerminal()
						}
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:  "format",
				Usage: "across  -x,  commas  -m, horizontal -x, long -l, single-column -1, verbose -l, vertical -C",
				Action: func(context *cli.Context, s string) error {
					switch s {
					case "across", "x", "horizontal":
						if _, ok := p.(*printer.Across); !ok {
							p = printer.NewAcross()
						}
					case "commas", "m":
						if _, ok := p.(*printer.CommaPrint); !ok {
							p = printer.NewCommaPrint()
						}
					case "long", "l", "verbose":
						contentFunc = append(contentFunc, filter.EnableFileMode(r), filter.EnableSize(filter.Auto, r), filter.EnableOwner(r), filter.EnableGroup(r), filter.EnableTime(timeFormat, r))
						if _, ok := p.(*printer.Byline); !ok {
							p = printer.NewByline()
						}
					case "single-column", "1":
						if _, ok := p.(*printer.Byline); !ok {
							p = printer.NewByline()
						}
					case "vertical", "C":
						if _, ok := p.(*printer.FitTerminal); !ok {
							p = printer.NewFitTerminal()
						}
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:        "theme",
				Aliases:     []string{"th"},
				DefaultText: "default",
				Usage:       "apply theme `path/to/theme`",
				Action: func(context *cli.Context, s string) error {
					err := theme.GetTheme(s)
					if err != nil {
						return err
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:  "g",
				Usage: "like -all/l, but do not list owner",
				Action: func(context *cli.Context, b bool) error {
					if b {
						// remove filter.RemoveHidden
						newFF := make([]*filter.TypeFunc, 0, len(typeFunc))
						for _, typeFunc := range typeFunc {
							if typeFunc != &filter.RemoveHidden {
								newFF = append(newFF, typeFunc)
							}
						}
						typeFunc = newFF
						contentFunc = append(contentFunc, filter.EnableFileMode(r), filter.EnableSize(filter.Auto, r), filter.EnableOwner(r), filter.EnableTime(timeFormat, r))
						if _, ok := p.(*printer.Byline); !ok {
							p = printer.NewByline()
						}
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:  "o",
				Usage: "like -all/l, but do not list group information",
				Action: func(context *cli.Context, b bool) error {
					if b {
						// remove filter.RemoveHidden
						newFF := make([]*filter.TypeFunc, 0, len(typeFunc))
						for _, typeFunc := range typeFunc {
							if typeFunc != &filter.RemoveHidden {
								newFF = append(newFF, typeFunc)
							}
						}
						typeFunc = newFF
						contentFunc = append(contentFunc, filter.EnableFileMode(r), filter.EnableSize(filter.Auto, r), filter.EnableGroup(r), filter.EnableTime(timeFormat, r))
						if _, ok := p.(*printer.Byline); !ok {
							p = printer.NewByline()
						}
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:    "lh",
				Aliases: []string{"human-readable"},
				Usage:   "show human readable size",
				Action: func(context *cli.Context, b bool) error {
					if b {
						contentFunc = append(contentFunc, filter.EnableFileMode(r), filter.EnableSize(filter.Auto, r), filter.EnableOwner(r), filter.EnableGroup(r), filter.EnableTime(timeFormat, r))
						if _, ok := p.(*printer.Byline); !ok {
							p = printer.NewByline()
						}
					}
					return nil
				},
			},
			&cli.StringSliceFlag{
				Name:    "show-with-ext",
				Aliases: []string{"se", "ext"},
				Usage:   "show file which has target ext",
				Action: func(context *cli.Context, s []string) error {
					if len(s) > 0 {
						f := filter.ExtOnly(s...)
						typeFunc = append(typeFunc, &f)
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:    "show-only-dir",
				Aliases: []string{"sd", "dir"},
				Usage:   "show directory only",
				Action: func(context *cli.Context, b bool) error {
					if b {
						typeFunc = append(typeFunc, &filter.DirOnly)
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:    "d",
				Aliases: []string{"directory"},
				Usage:   "list directories themselves, not their contents",
			},
			&cli.BoolFlag{
				Name:    "F",
				Aliases: []string{"file"},
			},
		},
		HideHelpCommand: true,
		Action: func(context *cli.Context) error {
			path := context.Args().Slice()
			if len(path) == 0 {
				path = append(path, ".")
			}

			if context.Bool("tree") {
				depth := context.Int("depth")
				for i := 0; i < len(path); i++ {
					if len(path) > 1 {
						fmt.Printf("%s:\n", path[i])
					}

					pathbeautify.Transform(&path[i])

					s, err := tree.NewTreeString(path[i], depth, filter.NewTypeFilter(typeFunc...), r)
					if errors.Is(err, os.ErrNotExist) {
						fmt.Printf("%s g: No such file or directory: %s %s\n", theme.Error, err.(*os.PathError).Path, theme.Reset)
						continue
					} else if err != nil {
						fmt.Println(theme.Error+err.Error()+theme.Reset, err)
						continue
					}
					fmt.Println(s.MakeTreeStr())
					fmt.Printf("\n%d directories, %d files\n", s.Directory(), s.File())

					if i != len(path)-1 {
						//goland:noinspection GoPrintFunctions
						fmt.Println("\n") //nolint:govet
					}
				}
			} else {
				startDir, _ := os.Getwd()
				// flag: if show-icon or all is set
				flagsi := context.Bool("show-icon") || context.Bool("all")
				// flag: if d is set
				flagd := context.Bool("d")
				// flag: if A is set
				flagA := context.Bool("A")
				for i := 0; i < len(path); i++ {
					if len(path) > 1 {
						fmt.Printf("%s:\n", path[i])
					}

					pathbeautify.Transform(&path[i])

					infos := make([]os.FileInfo, 0, 20)

					isFile := false
					// switch to target dir
					// or get target file info
					if path[i] != "." {
						stat, err := os.Stat(path[i])
						if err != nil {
							fmt.Printf("%s g: %s %s\n", theme.Error, err.Error(), theme.Reset)
							continue
						}
						if stat.IsDir() {
							if flagd {
								// when -d is set, treat dir as file
								infos = append(infos, stat)
								isFile = true
							} else {
								_ = os.Chdir(path[i])
								if err != nil {
									fmt.Printf("%s g: %s %s\n", theme.Error, err.Error(), theme.Reset)
									continue
								}
								path[i] = "."
							}
						} else {
							infos = append(infos, stat)
							isFile = true
						}
					}

					// if flagsi->add flagsi, else->add default
					if flagsi {
						contentFunc = append(contentFunc, filter.EnableIconName(r, path[i]))
					} else {
						contentFunc = append(contentFunc, filter.EnableName(r))
					}

					var d []os.DirEntry
					var err error
					if isFile {
						goto final
					}

					d, err = os.ReadDir(path[i])
					if err != nil {
						goto final
					}

					// if -A(almost-all) is not set, add the "." info
					if !flagA {
						statCurrent, err := os.Stat(".")
						if err != nil {
							fmt.Println(err)
						}
						statParent, err := os.Stat("..")
						if err != nil {
							fmt.Println(err)
						}
						infos = append(infos, statCurrent, statParent)
					}

					for _, v := range d {
						info, err := v.Info()
						if err != nil {
							fmt.Println(err)
						} else {
							infos = append(infos, info)
						}
					}

					// remove non-display items
					infos = filter.NewTypeFilter(typeFunc...).Filter(infos)

				final:
					stringSlice := filter.NewContentFilter(contentFunc...).GetStringSlice(infos)
					p.Print(stringSlice...)

					// remove the last func
					if i != len(path)-1 {
						//goland:noinspection GoPrintFunctions
						fmt.Println("\n") //nolint:govet
						contentFunc = contentFunc[:len(contentFunc)-1]
						_ = os.Chdir(startDir)
					}
				}
			}
			return nil
		},
		Authors: []*cli.Author{
			{
				Name:  "Equationzhao",
				Email: "equationzhao@foxmail.com",
			},
		},
	}

	_ = app.Run(os.Args)
}
