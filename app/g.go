package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Equationzhao/g/filter"
	"github.com/Equationzhao/g/pathbeautify"
	"github.com/Equationzhao/g/printer"
	"github.com/Equationzhao/g/render"
	"github.com/Equationzhao/g/sorter"
	"github.com/Equationzhao/g/theme"
	"github.com/Equationzhao/g/timeparse"
	"github.com/Equationzhao/g/tree"
	"github.com/urfave/cli/v2"
	versionInfo "go.szostok.io/version"
	vp "go.szostok.io/version/printer"
	"go.szostok.io/version/style"
	"go.szostok.io/version/upgrade"
)

var (
	typeFunc      = make([]*filter.TypeFunc, 0)
	contentFunc   = make([]filter.ContentOption, 0)
	r             = render.NewRenderer(theme.DefaultTheme, theme.DefaultInfoTheme)
	p             = printer.NewFitTerminal()
	timeFormat    = "02.Jan'06 15:04"
	ReturnCode    = 0
	contentFilter = filter.NewContentFilter()
	CompiledAt    = ""
	sort          = sorter.NewSorter()
	timeType      = "mod"
	sizeUint      = filter.Auto
	sizeEnabler   = filter.NewSizeEnabler()
	wgs           = make([]filter.LengthFixed, 0, 1)
)

var Version = "0.5.2"

var G *cli.App

func init() {
	typeFunc = append(typeFunc, &filter.RemoveHidden)
	if CompiledAt == "" {
		CompiledAt = time.Now().Format(timeFormat)
	} else {
		CompiledAtTime, err := time.Parse("2006-01-02-15:04:05", CompiledAt)
		if err != nil {
			CompiledAt = time.Now().Format(timeFormat)
		} else {
			CompiledAt = CompiledAtTime.UTC().Format(timeFormat)
		}
	}
	sizeEnabler.SetRenderer(r)
	wgs = append(wgs, sizeEnabler)

	G = &cli.App{
		Name:      "g",
		Usage:     "a powerful ls",
		UsageText: "g [options] [path]",
		Copyright: `Copyright (C) 2023 Equationzhao. MIT License
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.`,
		Version: Version,
		Authors: []*cli.Author{
			{
				Name:  "Equationzhao",
				Email: "equationzhao@foxmail.com",
			},
		},
		SliceFlagSeparator: ",",
		HideHelpCommand:    true,
		OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
			_, _ = fmt.Println(MakeErrorStr(err.Error()))
			return nil
		},
		Flags: make([]cli.Flag, 0, len(viewFlag)+len(filteringFlag)+len(sortingFlags)+len(displayFlag)+5),
		Action: func(context *cli.Context) error {
			var (
				minorErr   = false
				seriousErr = false
			)

			path := context.Args().Slice()

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
			if context.Bool("git-status") {
				nameToDisplay.SetGit()
			}

			{
				s := context.String("git-status-style")
				switch s {
				case "symbol", "sym":
					nameToDisplay.GitStyle = filter.GitStyleSym
				case "dot", ".":
					nameToDisplay.GitStyle = filter.GitStyleDot
				default:
					nameToDisplay.GitStyle = filter.GitStyleDefault
				}
			}

			contentFunc = append(contentFunc, nameToDisplay.Enable())
			typeFilter := filter.NewTypeFilter(typeFunc...)

			gitignore := context.Bool("hide-git-ignore")
			removeGitIgnore := new(filter.TypeFunc)
			if gitignore {
				typeFilter.AppendTo(removeGitIgnore)
			}

			// set sort func
			if sort.Len() == 0 {
				sort.AddOption(sorter.Default)
			}
			contentFilter.SetSortFunc(sort.Build())
			contentFilter.SetOptions(contentFunc...)

			// if no path, use current path
			if len(path) == 0 {
				path = append(path, ".")
			}
			contentFilter.SetOptions(contentFunc...)
			contentFilter.AppendToLengthFixed(wgs...)

			if context.Bool("tree") {
				depth := context.Int("depth")
				for i := 0; i < len(path); i++ {
					if len(path) > 1 {
						fmt.Printf("%s:\n", path[i])
					}

					pathbeautify.Transform(&path[i])

					s, err := tree.NewTreeString(path[i], depth, typeFilter, contentFilter)
					if errors.Is(err, os.ErrNotExist) {
						fmt.Print(MakeErrorStr(fmt.Sprintf("No such file or directory: %s", err.(*os.PathError).Path)))
						seriousErr = true
						continue
					} else if err != nil {
						fmt.Println(MakeErrorStr(err.Error()))
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
							_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
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
									_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
									seriousErr = true
									continue
								}
							}
						} else {
							parent := filepath.Dir(path[i])
							_ = os.Chdir(parent)
							infos = append(infos, stat)
							isFile = true
						}
					}

					var d []os.DirEntry
					var err error
					if isFile {
						goto final
					}

					d, err = os.ReadDir(".")
					if err != nil {
						goto final
					}

					// if -A(almost-all) is not set, add the "."/".." info
					if !flagA {
						statCurrent, err := os.Stat(".")
						if err != nil {
							seriousErr = true
							_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
						}
						statParent, err := os.Stat("..")
						if err != nil {
							minorErr = true
							_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
						}
						infos = append(infos, statCurrent, statParent)
					}

					for _, v := range d {
						info, err := v.Info()
						if err != nil {
							minorErr = true
							_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
						} else {
							infos = append(infos, info)
						}
					}

					if gitignore {
						*removeGitIgnore = filter.RemoveGitIgnore(path[i])
					}
					nameToDisplay.SetParent(path[i])
					// remove non-display items
					infos = typeFilter.Filter(infos...)

				final:
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
						err = os.Chdir(startDir)
						if err != nil {
							seriousErr = true
						}
						sizeEnabler.Reset()
					}
				}
			}

			if seriousErr {
				ReturnCode = 2
			} else if minorErr {
				ReturnCode = 1
			}

			return nil
		},
	}
	G.Flags = append(G.Flags, &cli.BoolFlag{
		Name:     "check-new-version",
		Usage:    "check if there's new release",
		Category: "software info",
		Action: func(context *cli.Context, b bool) error {
			if b {
				fmt.Println(context.App.Name + " - " + context.App.Usage)
				upgrade.WithUpdateCheckTimeout(1 * time.Second)
				notice := upgrade.NewGitHubDetector("Equationzhao", "g")
				_ = notice.PrintIfFoundGreater(os.Stderr, Version)
				return Err4Exit{}
			}
			return nil
		},
		DisableDefaultText: true,
	})

	G.Flags = append(G.Flags, viewFlag...)
	G.Flags = append(G.Flags, displayFlag...)
	G.Flags = append(G.Flags, filteringFlag...)
	G.Flags = append(G.Flags, sortingFlags...)

	initHelpTemp()

	initVersionHelpFlags()
}

type Err4Exit struct{}

func (c Err4Exit) Error() string {
	panic("it's an error defined to exit app, should not call this")
}

func initHelpTemp() {
	cli.AppHelpTemplate = fmt.Sprintf(`%s
REPO:
	https://github.com/Equationzhao/g

%s compiled at %s
`, cli.AppHelpTemplate, Version, CompiledAt)
}

func initVersionHelpFlags() {
	repos := "https://github.com/Equationzhao/g"
	info := versionInfo.Get()
	info.Version = Version
	info.BuildDate = CompiledAt
	info.ExtraFields = repos
	format := style.Formatting{
		Header: style.Header{
			Prefix: "💡 ",
			FormatPrimitive: style.FormatPrimitive{
				Color:   "Green",
				Options: []string{"Bold"},
			},
		},
	}

	c := vp.New(vp.WithPrettyFormatting(&format))
	cli.VersionPrinter = func(cCtx *cli.Context) {
		info.Meta = versionInfo.Meta{
			CLIName: cCtx.App.Name + " - " + cCtx.App.Usage,
		}
		_ = c.PrintInfo(os.Stdout, info)
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:               "version",
		Aliases:            []string{"v"},
		Usage:              "print the version",
		DisableDefaultText: true,
		Category:           "software info",
	}

	cli.HelpFlag = &cli.BoolFlag{
		Name:               "help",
		Aliases:            []string{"h"},
		Usage:              "show help",
		DisableDefaultText: true,
		Category:           "software info",
	}
}

var viewFlag = []cli.Flag{
	// VIEW
	&cli.BoolFlag{
		Name:               "inode",
		Aliases:            []string{"i"},
		Usage:              "show inode[linux only]",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			i := filter.NewInodeEnabler()
			wgs = append(wgs, i)
			contentFunc = append(contentFunc, i.Enable(r))
			return nil
		},
	},
	&cli.StringFlag{
		Name:        "size-unit",
		Aliases:     []string{"su"},
		Usage:       "size unit, b, k, m, g, t, p, e, z, y, auto",
		EnvVars:     []string{"SIZE_UNIT"},
		DefaultText: "auto",
		Action: func(context *cli.Context, s string) error {
			defer func() {
				if err := recover(); err != nil {
					ReturnCode = 1
				}
			}()

			if strings.EqualFold(s, "auto") {
				return nil
			}
			sizeUint = filter.ConvertFromSizeString(s)
			if sizeUint == filter.Unknown {
				ReturnCode = 1
				return errors.New("invalid size unit")
			}
			return nil
		},
	},
	&cli.StringFlag{
		Name:        "time-type",
		Aliases:     []string{"tt"},
		Usage:       "time type, mod, create, access",
		EnvVars:     []string{"TIME_TYPE"},
		DefaultText: "mod",
		Action: func(context *cli.Context, s string) error {
			if s == "mod" || s == "create" || s == "access" {
				timeType = s
				return nil
			} else {
				ReturnCode = 1
				return errors.New("invalid time type")
			}
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "uid",
		Usage:              "show uid instead of username [sid in windows]",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				filter.Uid = true
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "gid",
		Usage:              "show gid instead of groupname [sid in windows]",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				filter.Gid = true
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "numeric",
		Usage:              " List numeric user and group IDs instead of name [sid in windows]",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				filter.Gid = true
				filter.Uid = true
			}
			return nil
		},
		Category: "VIEW",
	},
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
				s := s[1:] // remove +
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
				ReturnCode = 1
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

				contentFunc = append(contentFunc, sizeEnabler.EnableSize(sizeUint))
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
				contentFunc = append(contentFunc, filter.EnableTime(timeFormat, r, timeType))
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
	&cli.StringSliceFlag{
		Name:     "checksum-algorithm",
		Usage:    "show checksum of file with algorithm: md5, sha1, sha256, sha512",
		Aliases:  []string{"ca"},
		Value:    cli.NewStringSlice("sha1"),
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "checksum",
		Usage:              "show checksum of file with algorithm: md5, sha1, sha256, sha512",
		Aliases:            []string{"cs"},
		DisableDefaultText: true,
		Category:           "VIEW",
		Action: func(context *cli.Context, b bool) error {
			ss := context.StringSlice("checksum-algorithm")
			if ss == nil {
				ss = []string{"sha1"}
			}
			sums := make([]filter.SumType, 0, len(ss))
			for _, s := range ss {
				switch s {
				case "md5":
					sums = append(sums, filter.SumTypeMd5)
				case "sha1":
					sums = append(sums, filter.SumTypeSha1)
				case "sha256":
					sums = append(sums, filter.SumTypeSha256)
				case "sha512":
					sums = append(sums, filter.SumTypeSha512)
				}
			}

			if b {
				contentFunc = append(contentFunc, contentFilter.EnableSum(sums...))
				if _, ok := p.(*printer.Byline); !ok {
					p = printer.NewByline()
				}
			}
			return nil
		},
	},
	&cli.BoolFlag{
		Name:               "git-status",
		Usage:              "show git status: ? untracked, + added, ! deleted, ~ modified, | renamed, = copied, $ ignored [if git is installed]",
		Aliases:            []string{"gs"},
		DisableDefaultText: true,
		Category:           "VIEW",
	},
	&cli.StringFlag{
		Name:     "git-status-style",
		Usage:    "git status style: colored-symbol: {? untracked, + added, - deleted, ~ modified, | renamed, = copied, ! ignored} colored-dot",
		Aliases:  []string{"gss"},
		Category: "VIEW",
	},
}

var displayFlag = []cli.Flag{
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
		Aliases:            []string{"comma"},
		Usage:              "fill width with a comma separated list of entries",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				if _, ok := p.(*printer.CommaPrint); !ok {
					p = printer.NewCommaPrint()
				}
			}
			return nil
		},
		Category: "DISPLAY",
	},
	&cli.BoolFlag{
		Name:               "x",
		Aliases:            []string{"col", "across", "horizontal"},
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
				contentFunc = append(contentFunc, filter.EnableFileMode(r), sizeEnabler.EnableSize(sizeUint), contentFilter.EnableOwner(r), contentFilter.EnableGroup(r), filter.EnableTime(timeFormat, r, timeType))
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
}

var filteringFlag = []cli.Flag{
	&cli.StringSliceFlag{
		Name:    "ignore-glob",
		Aliases: []string{"I"},
		Usage:   "ignore Glob patterns",
		Action: func(context *cli.Context, s []string) error {
			if len(s) > 0 {
				f, err := filter.RemoveGlob(s...)
				if err != nil {
					return err
				}
				typeFunc = append(typeFunc, &f)
			}
			return nil
		},
		Category: "FILTERING",
	},
	&cli.StringSliceFlag{
		Name:    "match-glob",
		Aliases: []string{"M"},
		Usage:   "match Glob patterns",
		Action: func(context *cli.Context, s []string) error {
			if len(s) > 0 {
				f, err := filter.GlobOnly(s...)
				if err != nil {
					return err
				}
				typeFunc = append(typeFunc, &f)
			}
			return nil
		},
		Category: "FILTERING",
	},
	&cli.BoolFlag{
		Name:               "lh",
		Aliases:            []string{"human-readable"},
		DisableDefaultText: true,
		Usage:              "show human readable size",
		Action: func(context *cli.Context, b bool) error {
			if b {
				contentFunc = append(contentFunc, filter.EnableFileMode(r), sizeEnabler.EnableSize(sizeUint), contentFilter.EnableOwner(r), contentFilter.EnableGroup(r), filter.EnableTime(timeFormat, r, timeType))
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
		Aliases:            []string{"nd", "nodir", "no-dir"},
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
		Aliases:            []string{"sd", "dir", "only-dir"},
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
				contentFunc = append(contentFunc, filter.EnableFileMode(r), sizeEnabler.EnableSize(sizeUint), contentFilter.EnableGroup(r), filter.EnableTime(timeFormat, r, timeType))
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
				contentFunc = append(contentFunc, filter.EnableFileMode(r), sizeEnabler.EnableSize(sizeUint), contentFilter.EnableOwner(r), filter.EnableTime(timeFormat, r, timeType))
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
				contentFunc = append(contentFunc, filter.EnableFileMode(r), sizeEnabler.EnableSize(sizeUint), contentFilter.EnableOwner(r))
				if !context.Bool("G") {
					contentFunc = append(contentFunc, contentFilter.EnableGroup(r))
				}
				contentFunc = append(contentFunc, filter.EnableTime(timeFormat, r, timeType))

				if _, ok := p.(*printer.Byline); !ok {
					p = printer.NewByline()
				}
			}
			return nil
		},
		Category: "FILTERING",
	},
	&cli.BoolFlag{
		Name:               "hide-git-ignore",
		Aliases:            []string{"gi", "hgi"},
		Usage:              "hide git ignored file/dir [if git is installed]",
		DisableDefaultText: true,
		Category:           "FILTERING",
	},
}

var sortingFlags = []cli.Flag{
	&cli.StringSliceFlag{
		Name:    "sort",
		Aliases: []string{"SORT_FIELD"},
		Usage:   "sort by field, default: ascending and case insensitive, field beginning with Uppercase is case sensitive, available fields: name,size,time,owner,group,extension. following `-descend` to sort descending",
		Action: func(context *cli.Context, slice []string) error {
			sorter.WithSize(len(slice))(sort)
			for _, s := range slice {
				switch s {
				case "name-descend":
					sort.AddOption(sorter.ByNameDescend)
				case "name":
					sort.AddOption(sorter.ByNameAscend)
				case "Name":
					sort.AddOption(sorter.ByNameCaseSensitiveAscend)
				case "Name-descend":
					sort.AddOption(sorter.ByNameCaseSensitiveDescend)
				case "size-descend":
					sort.AddOption(sorter.BySizeDescend)
				case "size":
					sort.AddOption(sorter.BySizeAscend)
				case "time-descend":
					sort.AddOption(sorter.ByTimeDescend)
				case "time":
					sort.AddOption(sorter.ByTimeAscend)
				case "extension-descend", "ext-descend":
					sort.AddOption(sorter.ByExtensionDescend)
				case "extension", "ext":
					sort.AddOption(sorter.ByExtensionAscend)
				case "Extension-descend", "Ext-descend":
					sort.AddOption(sorter.ByExtensionCaseSensitiveDescend)
				case "Extension", "Ext":
					sort.AddOption(sorter.ByExtensionCaseSensitiveAscend)
				case "group-descend":
					sort.AddOption(sorter.ByGroupDescend)
				case "group":
					sort.AddOption(sorter.ByGroupAscend)
				case "Group-descend":
					sort.AddOption(sorter.ByGroupCaseSensitiveDescend)
				case "Group":
					sort.AddOption(sorter.ByGroupCaseSensitiveAscend)
				case "owner-descend":
					sort.AddOption(sorter.ByOwnerDescend)
				case "owner":
					sort.AddOption(sorter.ByOwnerAscend)
				case "Owner-descend":
					sort.AddOption(sorter.ByOwnerCaseSensitiveDescend)
				case "Owner":
					sort.AddOption(sorter.ByOwnerCaseSensitiveAscend)
				default:
					return fmt.Errorf("unknown sort field: %s", s)
				}
			}
			return nil
		},
		Category: "SORTING",
	},
	&cli.BoolFlag{
		Name:               "sort-reverse",
		Aliases:            []string{"sr"},
		Usage:              "reverse the order of the sort",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				sort.Reverse()
			}
			return nil
		},
		Category: "SORTING",
	},
	&cli.BoolFlag{
		Name:               "dir-first",
		Aliases:            []string{"df"},
		Usage:              "List directories before other files",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				sort.DirFirst()
			} else {
				sort.UnsetDirFirst()
			}
			return nil
		},
		Category: "SORTING",
	},
}

func MakeErrorStr(msg string) string {
	return fmt.Sprintf("%s g: %s %s\n", theme.Error, msg, theme.Reset)
}
