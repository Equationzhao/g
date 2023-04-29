package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Equationzhao/g/filter"
	"github.com/Equationzhao/g/pathbeautify"
	"github.com/Equationzhao/g/printer"
	"github.com/Equationzhao/g/render"
	"github.com/Equationzhao/g/theme"
	"github.com/Equationzhao/g/timeparse"
	"github.com/Equationzhao/g/tree"
	"github.com/urfave/cli/v2"
)

var (
	typeFunc      = make([]*filter.TypeFunc, 0)
	contentFunc   = make([]filter.ContentOption, 0)
	r             = render.NewRenderer(theme.DefaultTheme, theme.DefaultInfoTheme)
	p             = printer.NewFitTerminal()
	timeFormat    = "02.Jan'06 15:04"
	returnCode    = 0
	sizeEnabler   = filter.Size{}
	contentFilter = filter.NewContentFilter()
)

const version = "v0.3.1"

func init() {
	typeFunc = append(typeFunc, &filter.RemoveHidden)
}

func main() {
	defer func() {
		os.Exit(returnCode)
	}()

	cli.AppHelpTemplate = fmt.Sprintf(`%s
REPO:
	https://github.com/Equationzhao/g

%s
`, cli.AppHelpTemplate, version)

	app := &cli.App{
		Name:      "g",
		Usage:     "a powerful ls",
		UsageText: "g [options] [path]",
		Copyright: `Copyright (C) 2023 Equationzhao. MIT License
This is free software: you are free to change and redistribute it.
There is NO  WARRANTY, to the extent permitted by law.`,
		Version: version,
		Authors: []*cli.Author{
			{
				Name:  "Equationzhao",
				Email: "equationzhao@foxmail.com",
			},
		},
		SliceFlagSeparator: ",",
		HideHelpCommand:    true,
		HideVersion:        true,
		Writer:             os.Stdout,
		OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
			fmt.Fprintf(cCtx.App.Writer, "%s %s: %s %s\n", theme.Error, cCtx.App.Name, err, theme.Reset)
			return nil
		},
		Flags: []cli.Flag{
			// VIEW
			&cli.StringFlag{
				Name:        "time-style",
				Usage:       "time/date format with -l, Valid timestamp styles are `default', `iso`, `long iso`, `full-iso`, `locale`, custom `+FORMAT` like date(1).",
				EnvVars:     []string{"TIME_STYLE"},
				DefaultText: "+%d.%b'%y %H:%M (like 02.Jan'06 15:04)",
				Action: func(context *cli.Context, s string) error {
					/*
						The TIME_STYLE argument can be full-iso, long-iso, iso, locale, or  +FORMAT.   FORMAT
						is  interpreted  like in date(1).  If FORMAT is FORMAT1<newline>FORMAT2, then FORMAT1
						applies to non-recent files and FORMAT2 to recent files.   TIME_STYLE  prefixed  with
						'posix-' takes effect only outside the POSIX locale.  Also the TIME_STYLE environment
						variable sets the default style to use.
					*/
					if strings.HasPrefix(s, "+") {
						timeFormat = timeparse.Transform(s)
						return nil
					}

					switch s {
					case "full-iso":
						timeFormat = "2006-01-02 15:04:05.000000000 -0700"
					case "long-iso":
						timeFormat = "2006-01-02 15:04"
					case "locale":
						timeFormat = "Jan 02 15:04"
					case "iso":
						timeFormat = "01-02 15:04"
					case "default":
						timeFormat = "02.Jan'06 15:04"
					default:
						returnCode = 2
						return errors.New("invalid time-style")
					}
					return nil
				},
				Category: "VIEW",
			},
			&cli.BoolFlag{
				Name:               "full-time",
				Usage:              "like -all/l --time-style=full-iso",
				DisableDefaultText: true,
				Action: func(context *cli.Context, b bool) error {
					if b {
						timeFormat = "2006-01-02 15:04:05.000000000 -0700"
					}
					return nil
				},
				Category: "VIEW",
			},
			&cli.BoolFlag{
				Name:               "show-perm",
				Aliases:            []string{"sp"},
				Usage:              "show permission",
				DisableDefaultText: true,
				Action: func(context *cli.Context, b bool) error {
					if b {
						contentFunc = append(contentFunc, filter.EnableFileMode(r))
						if _, ok := p.(*printer.Byline); !ok {
							p = printer.NewByline()
						}
					}
					return nil
				},
				Category: "VIEW",
			},
			&cli.BoolFlag{
				Name:               "show-size",
				Aliases:            []string{"ss"},
				Usage:              "show file/dir size",
				DisableDefaultText: true,
				Action: func(context *cli.Context, b bool) error {
					if b {
						contentFunc = append(contentFunc, sizeEnabler.EnableSize(filter.Auto, r))
						if _, ok := p.(*printer.Byline); !ok {
							p = printer.NewByline()
						}
					}
					return nil
				},
				Category: "VIEW",
			},
			&cli.BoolFlag{
				Name:               "show-owner",
				Aliases:            []string{"so", "author"},
				Usage:              "show owner",
				DisableDefaultText: true,
				Action: func(context *cli.Context, b bool) error {
					if b {
						contentFunc = append(contentFunc, contentFilter.EnableOwner(r))
						if _, ok := p.(*printer.Byline); !ok {
							p = printer.NewByline()
						}
					}
					return nil
				},
				Category: "VIEW",
			},

			&cli.BoolFlag{
				Name:               "show-group",
				Aliases:            []string{"sg"},
				Usage:              "show group",
				DisableDefaultText: true,
				Action: func(context *cli.Context, b bool) error {
					if b {
						contentFunc = append(contentFunc, contentFilter.EnableGroup(r))
						if _, ok := p.(*printer.Byline); !ok {
							p = printer.NewByline()
						}
					}
					return nil
				},
				Category: "VIEW",
			},
			&cli.BoolFlag{
				Name:               "show-time",
				Aliases:            []string{"st"},
				Usage:              "show time",
				DisableDefaultText: true,
				Action: func(context *cli.Context, b bool) error {
					if b {
						contentFunc = append(contentFunc, filter.EnableTime(timeFormat, r))
						if _, ok := p.(*printer.Byline); !ok {
							p = printer.NewByline()
						}
					}
					return nil
				},
				Category: "VIEW",
			},
			&cli.BoolFlag{
				Name:               "show-icon",
				Usage:              "show icon",
				Aliases:            []string{"si", "icons"},
				DisableDefaultText: true,
				Category:           "VIEW",
			},
			&cli.BoolFlag{
				Name:               "show-total-size",
				Usage:              "show total size",
				Aliases:            []string{"ts"},
				DisableDefaultText: true,
				Category:           "VIEW",
				Action: func(context *cli.Context, b bool) error {
					if b {
						sizeEnabler.SetEnableTotal()
					}
					return nil
				},
			},

			// DISPLAY
			&cli.BoolFlag{
				Name:               "tree",
				Aliases:            []string{"t"},
				Usage:              "list in tree",
				DisableDefaultText: true,
				Category:           "DISPLAY",
			},
			&cli.IntFlag{
				Name:        "depth",
				Usage:       "tree limit depth, negative -> infinity",
				DefaultText: "infinity",
				Value:       -1,
				Category:    "DISPLAY",
			},
			&cli.BoolFlag{
				Name:               "byline",
				Aliases:            []string{"bl", "1", "oneline", "single-column"},
				Usage:              "print by line",
				DisableDefaultText: true,
				Action: func(context *cli.Context, b bool) error {
					if b {
						if _, ok := p.(*printer.Byline); !ok {
							p = printer.NewByline()
						}
					}
					return nil
				},
				Category: "DISPLAY",
			},
			&cli.BoolFlag{
				Name:               "zero",
				Aliases:            []string{"0"},
				Usage:              "end each output line with NUL, not newline",
				DisableDefaultText: true,
				Action: func(context *cli.Context, b bool) error {
					if b {
						if _, ok := p.(*printer.Zero); !ok {
							p = printer.NewZero()
						}
					}
					return nil
				},
				Category: "DISPLAY",
			},
			&cli.BoolFlag{
				Name:               "m",
				Usage:              "fill width with a comma separated list of entries",
				DisableDefaultText: true,
				Action: func(context *cli.Context, b bool) error {
					if b {
						if _, ok := p.(*printer.Byline); !ok {
							p = printer.NewCommaPrint()
						}
					}
					return nil
				},
				Category: "DISPLAY",
			},
			&cli.BoolFlag{
				Name:               "x",
				Usage:              "list entries by lines instead of by columns",
				DisableDefaultText: true,
				Action: func(context *cli.Context, b bool) error {
					if b {
						if _, ok := p.(*printer.Across); !ok {
							p = printer.NewAcross()
						}
					}
					return nil
				},
				Category: "DISPLAY",
			},
			&cli.BoolFlag{
				Name:               "C",
				Aliases:            []string{"vertical"},
				Usage:              "list entries by columns (default)",
				DisableDefaultText: true,
				Action: func(context *cli.Context, b bool) error {
					if b {
						if _, ok := p.(*printer.FitTerminal); !ok {
							p = printer.NewFitTerminal()
						}
					}
					return nil
				},
				Category: "DISPLAY",
			},
			&cli.StringFlag{
				Name:        "format",
				DefaultText: "C",
				Usage:       "across  -x,  commas  -m, horizontal -x, long -l, single-column -1, verbose -l, vertical -C",
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
						contentFunc = append(contentFunc, filter.EnableFileMode(r), sizeEnabler.EnableSize(filter.Auto, r), contentFilter.EnableOwner(r), contentFilter.EnableGroup(r), filter.EnableTime(timeFormat, r))
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
				Category: "DISPLAY",
			},

			&cli.StringFlag{
				Name:    "theme",
				Aliases: []string{"th"},
				Usage:   "apply theme `path/to/theme`",
				Action: func(context *cli.Context, s string) error {
					err := theme.GetTheme(s)
					if err != nil {
						return err
					}
					return nil
				},
				Category: "DISPLAY",
			},
			&cli.BoolFlag{
				Name:               "d",
				Aliases:            []string{"directory", "list-dirs"},
				DisableDefaultText: true,
				Usage:              "list directories themselves, not their contents",
				Category:           "DISPLAY",
			},
			&cli.BoolFlag{
				Name:               "F",
				Aliases:            []string{"classify"},
				DisableDefaultText: true,
				Usage:              "append indicator (one of */=>@|) to entries",
				Category:           "DISPLAY",
			},
			&cli.BoolFlag{
				Name:               "file-type",
				Aliases:            []string{"ft"},
				DisableDefaultText: true,
				Usage:              "likewise, except do not append '*'",
				Category:           "DISPLAY",
			},

			&cli.BoolFlag{
				Name:               "lh",
				Aliases:            []string{"human-readable"},
				DisableDefaultText: true,
				Usage:              "show human readable size",
				Action: func(context *cli.Context, b bool) error {
					if b {
						contentFunc = append(contentFunc, filter.EnableFileMode(r), sizeEnabler.EnableSize(filter.Auto, r), contentFilter.EnableOwner(r), contentFilter.EnableGroup(r), filter.EnableTime(timeFormat, r))
						if _, ok := p.(*printer.Byline); !ok {
							p = printer.NewByline()
						}
					}
					return nil
				},
				Category: "FILTERING",
			},
			&cli.BoolFlag{
				Name:               "show-hidden",
				Aliases:            []string{"sh", "a"},
				DisableDefaultText: true,
				Usage:              "show hidden files",
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
				Category: "FILTERING",
			},
			&cli.StringSliceFlag{
				Name:    "show-only-ext",
				Aliases: []string{"se", "ext"},
				Usage:   "show file which has target ext, eg: --show-only-ext=go,java",
				Action: func(context *cli.Context, s []string) error {
					if len(s) > 0 {
						f := filter.ExtOnly(s...)
						typeFunc = append(typeFunc, &f)
					}
					return nil
				},
				Category: "FILTERING",
			},
			&cli.StringSliceFlag{
				Name:    "show-no-ext",
				Aliases: []string{"sne", "noext"},
				Usage:   "show file which doesn't have target ext",
				Action: func(context *cli.Context, s []string) error {
					if len(s) > 0 {
						f := filter.RemoveByExt(s...)
						typeFunc = append(typeFunc, &f)
					}
					return nil
				},
				Category: "FILTERING",
			},
			&cli.BoolFlag{
				Name:               "show-no-dir",
				Aliases:            []string{"nd", "nodir"},
				DisableDefaultText: true,
				Usage:              "do not show directory",
				Action: func(context *cli.Context, b bool) error {
					if b {
						typeFunc = append(typeFunc, &filter.RemoveDir)
					}
					return nil
				},
				Category: "FILTERING",
			},
			&cli.BoolFlag{
				Name:               "show-only-dir",
				Aliases:            []string{"sd", "dir"},
				DisableDefaultText: true,
				Usage:              "show directory only",
				Action: func(context *cli.Context, b bool) error {
					if b {
						typeFunc = append(typeFunc, &filter.DirOnly)
					}
					return nil
				},
				Category: "FILTERING",
			},
			&cli.BoolFlag{
				Name:               "B",
				Aliases:            []string{"ignore-backups"},
				DisableDefaultText: true,
				Usage:              "do not list implied entries ending with ~",
				Action: func(context *cli.Context, b bool) error {
					if b {
						typeFunc = append(typeFunc, &filter.RemoveBackups)
					}
					return nil
				},
				Category: "FILTERING",
			},
			&cli.BoolFlag{
				Name:               "A",
				Aliases:            []string{"almost-all"},
				DisableDefaultText: true,
				Usage:              "do not list implied . and ..",
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
				Category: "FILTERING",
			},
			&cli.BoolFlag{
				Name:               "o",
				DisableDefaultText: true,
				Usage:              "like -all/l, but do not list group information",
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
						contentFunc = append(contentFunc, filter.EnableFileMode(r), sizeEnabler.EnableSize(filter.Auto, r), contentFilter.EnableGroup(r), filter.EnableTime(timeFormat, r))
						if _, ok := p.(*printer.Byline); !ok {
							p = printer.NewByline()
						}
					}
					return nil
				},
				Category: "FILTERING",
			},
			&cli.BoolFlag{
				Name:               "g",
				DisableDefaultText: true,
				Usage:              "like -all/l, but do not list owner",
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
						contentFunc = append(contentFunc, filter.EnableFileMode(r), sizeEnabler.EnableSize(filter.Auto, r), contentFilter.EnableOwner(r), filter.EnableTime(timeFormat, r))
						if _, ok := p.(*printer.Byline); !ok {
							p = printer.NewByline()
						}
					}
					return nil
				},
				Category: "FILTERING",
			},
			&cli.BoolFlag{
				Name:               "G",
				DisableDefaultText: true,
				Aliases:            []string{"no-group"},
				Usage:              "in a long listing, don't print group names",
				Category:           "FILTERING",
			},
			&cli.BoolFlag{
				Name:               "all",
				Aliases:            []string{"la", "l"},
				Usage:              "show all info/use a long listing format",
				DisableDefaultText: true,
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
						sizeEnabler.SetEnableTotal()
						contentFunc = append(contentFunc, filter.EnableFileMode(r), sizeEnabler.EnableSize(filter.Auto, r), contentFilter.EnableOwner(r))
						if !context.Bool("G") {
							contentFunc = append(contentFunc, contentFilter.EnableGroup(r))
						}
						contentFunc = append(contentFunc, filter.EnableTime(timeFormat, r))

						if _, ok := p.(*printer.Byline); !ok {
							p = printer.NewByline()
						}
					}
					return nil
				},
				Category: "FILTERING",
			},
		},
		Action: func(context *cli.Context) error {
			var (
				minorErr   = false
				seriousErr = false
			)

			path := context.Args().Slice()

			// if no path, use current path
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
						seriousErr = true
						continue
					} else if err != nil {
						fmt.Println(theme.Error+err.Error()+theme.Reset, err)
						seriousErr = true
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

				// flag: if d is set
				flagd := context.Bool("d")
				// flag: if A is set
				flagA := context.Bool("A")

				nameToDisplay := filter.NewNameEnable().SetRenderer(r)
				if context.Bool("show-icon") || context.Bool("all") {
					nameToDisplay.SetIcon()
				}
				if context.Bool("F") {
					nameToDisplay.SetClassify()
				}
				if context.Bool("file-type") {
					nameToDisplay.SetClassify()
					nameToDisplay.SetFileType()
				}
				contentFunc = append(contentFunc, nameToDisplay.Enable())
				typeFilter := filter.NewTypeFilter(typeFunc...)
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
							seriousErr = true
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
									seriousErr = true
									continue
								}
								path[i] = "."
							}
						} else {
							infos = append(infos, stat)
							isFile = true
						}
					}

					nameToDisplay.SetParent(path[i])

					var d []os.DirEntry
					var err error
					if isFile {
						goto final
					}

					d, err = os.ReadDir(path[i])
					if err != nil {
						goto final
					}

					// if -A(almost-all) is not set, add the "."/".." info
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
					infos = typeFilter.Filter(infos)

				final:
					contentFilter.SetOptions(contentFunc...)
					stringSlice := contentFilter.GetStringSlice(infos)

					// if -l/show-total-size is set, add total size
					if total, ok := sizeEnabler.Total(); ok {
						p.Print(fmt.Sprintf("  total %s", sizeEnabler.Size2String(total, 0)))
					}
					p.Print(stringSlice...)

					// switch back to start dir
					if i != len(path)-1 {
						//goland:noinspection GoPrintFunctions
						fmt.Println("\n") //nolint:govet
						_ = os.Chdir(startDir)
						sizeEnabler.Reset()
					}
				}
			}

			if seriousErr {
				returnCode = 2
			} else if minorErr {
				returnCode = 1
			}

			return nil
		},
	}

	if doc {
		md, _ := os.Create("g.md")
		s, _ := app.ToMarkdown()
		_, _ = fmt.Fprintln(md, s)
		man, _ := os.Create(filepath.Join("man", "g.1"))
		s, _ = app.ToMan()
		_, _ = fmt.Fprintln(man, s)
	} else {
		err := app.Run(os.Args)
		if err != nil {
			fmt.Printf("%s g: %s %s\n", theme.Error, err.Error(), theme.Reset)
		}
	}
}
