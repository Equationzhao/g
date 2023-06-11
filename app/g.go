package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Equationzhao/g/display"
	"github.com/Equationzhao/g/filter"
	"github.com/Equationzhao/g/index"
	"github.com/Equationzhao/g/render"
	"github.com/Equationzhao/g/sorter"
	"github.com/Equationzhao/g/theme"
	"github.com/Equationzhao/g/timeparse"
	"github.com/Equationzhao/g/tree"
	"github.com/Equationzhao/g/util"
	"github.com/Equationzhao/pathbeautify"
	"github.com/gabriel-vasile/mimetype"
	"github.com/hako/durafmt"
	"github.com/urfave/cli/v2"
	"github.com/valyala/bytebufferpool"
	versionInfo "go.szostok.io/version"
	vp "go.szostok.io/version/printer"
	"go.szostok.io/version/style"
	"go.szostok.io/version/upgrade"
)

var (
	typeFunc      = make([]*filter.TypeFunc, 0)
	contentFunc   = make([]filter.ContentOption, 0)
	r             = render.NewRenderer(theme.DefaultTheme, theme.DefaultInfoTheme)
	p             = display.NewFitTerminal()
	timeFormat    = "02.Jan'06 15:04"
	ReturnCode    = 0
	contentFilter = filter.NewContentFilter()
	CompiledAt    = ""
	sort          = sorter.NewSorter()
	timeType      = []string{"mod"}
	sizeUint      = filter.Auto
	sizeEnabler   = filter.NewSizeEnabler()
	wgs           = make([]filter.LengthFixed, 0, 1)
	depthLimitMap = make(map[string]int)
	limitOnce     = util.Once{}
	hookOnce      = util.Once{}
)

var Version = "0.7.0"

var G *cli.App

func init() {
	typeFunc = append(typeFunc, &filter.RemoveHidden)
	if CompiledAt == "" {
		info, err := os.Stat(os.Args[0])
		if err != nil {
			CompiledAt = time.Now().Format(timeFormat)
		} else {
			CompiledAt = info.ModTime().Format(timeFormat)
		}
	} else {
		CompiledAtTime, err := time.Parse("2006-01-02-15:04:05", CompiledAt)
		if err == nil {
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
		Flags: make([]cli.Flag, 0, len(viewFlag)+len(filteringFlag)+len(sortingFlags)+len(displayFlag)+len(indexFlags)),
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

			fuzzy := context.Bool("fuzzy")
			if fuzzy {
				defer func() {
					for i := 0; i < 10; i++ {
						err := index.Close()
						if err != nil {
							continue
						}
						return
					}
				}()
			}

			disableIndex := context.Bool("di")
			wgUpdateIndex := sync.WaitGroup{}

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

			if context.Bool("Q") {
				nameToDisplay.SetQuote(`"`)
			}
			if context.Bool("N") {
				nameToDisplay.UnsetQuote()
			}

			transformEnabled := !context.Bool("np")

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
			depth := context.Int("depth")

			printPath := false
			if len(path) > 1 {
				printPath = true
			}

			if context.Bool("tree") {
				for i := 0; i < len(path); i++ {
					start := time.Now()

					if printPath {
						fmt.Printf("%s:\n", path[i])
					}

					if transformEnabled {
						path[i] = pathbeautify.Transform(path[i])
					}
					// fuzzy search
					if fuzzy {
						_, err := os.Stat(path[i])
						if err != nil {
							if newPath, b := fuzzyPath(path[i]); b != nil {
								_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
								minorErr = true
								continue
							} else {
								path[i] = newPath
								_, err = os.Stat(path[i])
								if err != nil {
									if pathErr := new(os.PathError); errors.As(err, &pathErr) {
										_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(fmt.Sprintf("No such file or directory/Can't access: %s", pathErr.Path)))
										seriousErr = true
										continue
									} else {
										_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
										seriousErr = true
										continue
									}
								}
								fmt.Println(path[i])
							}
						}
					}

					s, err, minorErrInTree := tree.NewTreeString(path[i], depth, typeFilter, contentFilter)
					if pathErr := new(os.PathError); errors.As(err, &pathErr) {
						_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(fmt.Sprintf("No such file or directory/Can't access: %s", pathErr.Path)))
						seriousErr = true
						continue
					} else if err != nil {
						_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
						seriousErr = true
						continue
					}

					if pathErr := new(os.PathError); errors.As(minorErrInTree, &pathErr) {
						_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(fmt.Sprintf("No such file or directory/Can't access: %s", pathErr.Path)))
						minorErr = true
					} else if minorErrInTree != nil {
						_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
						minorErr = true
					}

					absPath, err := filepath.Abs(path[i])
					if err != nil {
						minorErr = true
					} else {
						if !disableIndex {
							wgUpdateIndex.Add(1)
							go func() {
								if err = fuzzyUpdate(absPath); err != nil {
									minorErr = true
								}
								wgUpdateIndex.Done()
							}()
						}
					}

					fmt.Println(s.MakeTreeStr())
					fmt.Printf("\n%d directories, %d files\nunderwent %s", s.Directory(), s.File(), time.Since(start).String())

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
				flagR := context.Bool("R")

				header := context.Bool("header")
				if context.Bool("statistic") {
					nameToDisplay.SetStatistics(&filter.Statistics{})
				}

				for i := 0; i < len(path); i++ {
					start := time.Now()

					if printPath {
						fmt.Printf("%s:\n", path[i])
					}

					if transformEnabled {
						path[i] = pathbeautify.Transform(path[i])
					}

					infos := make([]os.FileInfo, 0, 20)

					isFile := false

					// get the abs path
					absPath, err := filepath.Abs(path[i])
					if err != nil {
						_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(fmt.Sprintf("Not a valid path: %s", absPath)))
					} else {
						path[i] = absPath
					}

					if path[i] != "." {
						stat, err := os.Stat(path[i])
						if err != nil {
							// no match
							if fuzzy {
								// start fuzzy search
								if newPath, err := fuzzyPath(filepath.Base(path[i])); err != nil {
									_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
									minorErr = true
								} else {
									path[i] = newPath
									stat, err = os.Stat(path[i])
									if err != nil {
										if pathErr := new(os.PathError); errors.As(err, &pathErr) {
											_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(fmt.Sprintf("No such file or directory/Can't access: %s", pathErr.Path)))
											seriousErr = true
											continue
										} else {
											_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
											seriousErr = true
											continue
										}
									}
									fmt.Println(path[i])
								}
							} else {
								// output error
								if pathErr := new(os.PathError); errors.As(err, &pathErr) {
									_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(fmt.Sprintf("No such file or directory/Can't access: %s", pathErr.Path)))
									seriousErr = true
									continue
								} else {
									_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
									seriousErr = true
									continue
								}
							}
						}
						if stat.IsDir() {
							if flagd {
								// when -d is set, treat dir as file
								infos = append(infos, stat)
								isFile = true
							}
						} else {
							infos = append(infos, stat)
							isFile = true
						}
					}

					if !disableIndex {
						wgUpdateIndex.Add(1)
						go func(i int) {
							if err = fuzzyUpdate(path[i]); err != nil {
								minorErr = true
							}
							wgUpdateIndex.Done()
						}(i)
					}

					var d []os.DirEntry
					if isFile {
						goto final
					}

					d, err = os.ReadDir(path[i])
					if err != nil {
						goto final
					}

					// if -A(almost-all) is not set, add the "."/".." info
					if !flagA {
						err := os.Chdir(path[i])
						if err != nil {
							_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
						} else {
							statCurrent, err := os.Stat(".")
							if err != nil {
								if pathErr := new(os.PathError); errors.As(err, &pathErr) {
									_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(fmt.Sprintf("No such file or directory/Can't access: %s", pathErr.Path)))
									seriousErr = true
								} else {
									_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
									seriousErr = true
								}
							} else {
								infos = append(infos, statCurrent)
							}

							statParent, err := os.Stat("..")
							if err != nil {
								if pathErr := new(os.PathError); errors.As(err, &pathErr) {
									_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(fmt.Sprintf("No such file or directory/Can't access: %s", pathErr.Path)))
									minorErr = true
								} else {
									_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
									minorErr = true
								}
							} else {
								infos = append(infos, statParent)
							}
						}
					}

					for _, v := range d {
						info, err := v.Info()
						if err != nil {
							if pathErr := new(os.PathError); errors.As(err, &pathErr) {
								_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(fmt.Sprintf("No such file or directory/Can't access: %s", pathErr.Path)))
								minorErr = true
							} else {
								minorErr = true
								_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
							}
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

					// if -R is set, add sub dir, insert into path[i+1]
					if flagR {

						// set depth
						dep, ok := depthLimitMap[path[i]]
						if !ok {
							depthLimitMap[path[i]] = depth
							dep = depth
						}
						if dep >= 2 || dep <= -1 {
							newPathLeft := make([]string, 0, len(path)-i)
							for _, info := range infos {
								if info.IsDir() {
									if info.Name() == "." || info.Name() == ".." {
										continue
									}
									newPath := filepath.Join(path[i], info.Name())
									newPathLeft = append(newPathLeft, newPath)
									depthLimitMap[newPath] = dep - 1
								}
							}
							path = append(path[:i+1], append(newPathLeft, path[i+1:]...)...)
						}
					}

				final:
					items := contentFilter.GetDisplayItems(infos...)

					{
						var i *display.Item
						// if -l/show-total-size is set, add total size
						if total, ok := sizeEnabler.Total(); ok {
							i = display.NewItem()
							i.Set("total", display.ItemContent{No: 0, Content: display.StringContent(fmt.Sprintf("  total %s", sizeEnabler.Size2String(total, 0)))})

						}
						if s := nameToDisplay.Statistics(); s != nil {
							tFormat := "\n  underwent %s"
							if i == nil {
								i = display.NewItem()
								tFormat = "  underwent %s"
							}
							i.Set("underwent", display.ItemContent{No: 1, Content: display.StringContent(fmt.Sprintf(tFormat, r.Time(durafmt.Parse(time.Since(start)).LimitToUnit("ms").String())))})
							i.Set("statistic", display.ItemContent{No: 2, Content: display.StringContent(fmt.Sprintf("\n  statistic: %s", s))})
							s.Reset()
						}
						if i != nil {
							p.Print(*i)
						}
					}

					if header {
						_ = hookOnce.Do(func() error {
							p.AddBeforePrint(func(item ...display.Item) {
								// add header
								allPart := item[0].KeysByOrder()
								longestEachPart := make(map[string]int)
								for _, it := range item {
									for _, part := range allPart {
										content, _ := it.Get(part)
										l := display.WidthLen(content.Content.String())
										if l > longestEachPart[part] {
											longestEachPart[part] = l
										}
									}
								}

								// add longest - len(header) * space
								// print header
								contentStrBuf := bytebufferpool.Get()
								for i, s := range allPart {
									if len(s) > longestEachPart[s] {
										// expand the every item's content of this part
										for _, it := range items {
											content, _ := it.Get(s)
											content.Content = display.StringContent(fmt.Sprintf("%s%s", strings.Repeat(" ", len(s)-longestEachPart[s]), content.Content.String()))
											it.Set(s, content)
										}
										_, _ = contentStrBuf.WriteString(theme.Underline)
										_, _ = contentStrBuf.WriteString(s)
										_, _ = contentStrBuf.WriteString(theme.Reset)
										if i != len(allPart)-1 {
											_, _ = contentStrBuf.WriteString(" ")
										}
									} else {
										_, _ = contentStrBuf.WriteString(theme.Underline)
										_, _ = contentStrBuf.WriteString(s)
										_, _ = contentStrBuf.WriteString(theme.Reset)
										if i != len(allPart)-1 {
											_, _ = contentStrBuf.WriteString(strings.Repeat(" ", longestEachPart[s]-len(s)+1))
										}
									}
								}
								_, _ = contentStrBuf.WriteString(theme.Reset)
								_, _ = fmt.Fprintln(display.Output, contentStrBuf.String())
								bytebufferpool.Put(contentStrBuf)
							})
							return nil
						})
					}

					itemsCopy := make([]display.Item, 0, len(items))
					for _, item := range items {
						itemsCopy = append(itemsCopy, *item)
					}

					p.Print(itemsCopy...)

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
			wgUpdateIndex.Wait()

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
	}, &cli.BoolFlag{
		Name:    "no-path-transform",
		Aliases: []string{"np"},
		Usage:   "By default, .../a/b/c will be transformed to ../../a/b/c, and ~ will be replaced by homedir, using this flag to disable this feature",
	})

	G.Flags = append(G.Flags, viewFlag...)
	G.Flags = append(G.Flags, displayFlag...)
	G.Flags = append(G.Flags, filteringFlag...)
	G.Flags = append(G.Flags, sortingFlags...)
	G.Flags = append(G.Flags, indexFlags...)

	initHelpTemp()

	initVersionHelpFlags()
}

func fuzzyUpdate(path string) error {
	err := index.Update(path)
	if err != nil {
		return err
	}
	return nil
}

// fuzzyPath returns the fuzzy path
// if error, return empty string and error
func fuzzyPath(path string) (newPath string, minorErr error) {
	fuzzed, err := index.FuzzySearch(path)
	if err == nil {
		return fuzzed, nil
	} else {
		return "", err
	}
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
		Name:  "header",
		Usage: "add a header row",
		Action: func(context *cli.Context, b bool) error {
			if b {
				if _, ok := p.(*display.Byline); !ok {
					p = display.NewByline()
				}
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:     "statistic",
		Usage:    "show statistic info",
		Category: "VIEW",
	},
	&cli.StringSliceFlag{
		Name:        "time-type",
		Aliases:     []string{"tt"},
		Usage:       "time type, mod(default), create, access, all",
		EnvVars:     []string{"TIME_TYPE"},
		DefaultText: "mod",
		Action: func(context *cli.Context, ss []string) error {
			timeType = make([]string, 0, len(ss))
			for _, s := range ss {
				if s == "mod" || s == "create" || s == "access" {
					timeType = append(timeType, s)
				} else if s == "all" {
					timeType = []string{"mod", "create", "access"}
				} else {
					ReturnCode = 1
					return errors.New("invalid time type")
				}
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.StringFlag{
		Name:        "size-unit",
		Aliases:     []string{"su", "block-size"},
		Usage:       "size unit, b, k, m, g, t, p, e, z, y, bb, nb, auto",
		EnvVars:     []string{"SIZE_UNIT"},
		DefaultText: "auto",
		Action: func(context *cli.Context, s string) error {
			if strings.EqualFold(s, "auto") {
				return nil
			}
			sizeUint = filter.ConvertFromSizeString(s)
			if sizeUint == filter.Unknown {
				ReturnCode = 1
				return fmt.Errorf("invalid size unit: %s", s)
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
				contentFunc = append(contentFunc, filter.EnableFileMode(r), sizeEnabler.EnableSize(sizeUint), contentFilter.EnableGroup(r))
				for _, s := range timeType {
					contentFunc = append(contentFunc, filter.EnableTime(timeFormat, s, r))
				}
				if _, ok := p.(*display.Byline); !ok {
					p = display.NewByline()
				}
			}
			return nil
		},
		Category: "VIEW",
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
				contentFunc = append(contentFunc, filter.EnableFileMode(r), sizeEnabler.EnableSize(sizeUint), contentFilter.EnableOwner(r))
				for _, s := range timeType {
					contentFunc = append(contentFunc, filter.EnableTime(timeFormat, s, r))
				}
				if _, ok := p.(*display.Byline); !ok {
					p = display.NewByline()
				}
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "G",
		DisableDefaultText: true,
		Aliases:            []string{"no-group"},
		Usage:              "in a long listing, don't print group names",
		Category:           "VIEW",
	},
	&cli.BoolFlag{
		Name:               "all",
		Aliases:            []string{"la", "l", "long"},
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
				for _, s := range timeType {
					contentFunc = append(contentFunc, filter.EnableTime(timeFormat, s, r))
				}
				if _, ok := p.(*display.Byline); !ok {
					p = display.NewByline()
				}
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "hide-git-ignore",
		Aliases:            []string{"gi", "hgi"},
		Usage:              "hide git ignored file/dir [if git is installed]",
		DisableDefaultText: true,
		Category:           "VIEW",
	},
	&cli.BoolFlag{
		Name:               "inode",
		Aliases:            []string{"i"},
		Usage:              "show inode[linux/darwin only]",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			i := filter.NewInodeEnabler()
			wgs = append(wgs, i)
			contentFunc = append(contentFunc, i.Enable(r))
			return nil
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
		Aliases:            []string{"numeric-uid-gid"},
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
	&cli.BoolFlag{
		Name:               "relative-time",
		Aliases:            []string{"rt"},
		Usage:              "show relative time",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				rt := filter.NewRelativeTimeEnabler()
				rt.Mode = timeType[0]
				contentFunc = append(contentFunc, rt.Enable(r))
				wgs = append(wgs, rt)
			}
			return nil
		},
		Category: "VIEW",
	},

	&cli.BoolFlag{
		Name:               "show-perm",
		Aliases:            []string{"sp", "permission", "perm"},
		Usage:              "show permission",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				contentFunc = append(contentFunc, filter.EnableFileMode(r))
				if _, ok := p.(*display.Byline); !ok {
					p = display.NewByline()
				}
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "show-size",
		Aliases:            []string{"ss", "size"},
		Usage:              "show file/dir size",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				contentFunc = append(contentFunc, sizeEnabler.EnableSize(sizeUint))
				if _, ok := p.(*display.Byline); !ok {
					p = display.NewByline()
				}
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "show-recursive-size",
		Aliases:            []string{"srs", "recursive-size"},
		Usage:              "show recursive size of dir, only work with --show-size",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				n := context.Int("depth")
				sizeEnabler.SetRecursive(filter.NewSizeRecursive(n))
			}
			return nil
		},
	},
	&cli.BoolFlag{
		Name:               "lh",
		Aliases:            []string{"human-readable", "hr"},
		DisableDefaultText: true,
		Usage:              "show human readable size",
		Action: func(context *cli.Context, b bool) error {
			if b {
				contentFunc = append(contentFunc, filter.EnableFileMode(r), sizeEnabler.EnableSize(sizeUint), contentFilter.EnableOwner(r), contentFilter.EnableGroup(r))
				for _, s := range timeType {
					contentFunc = append(contentFunc, filter.EnableTime(timeFormat, s, r))
				}
				if _, ok := p.(*display.Byline); !ok {
					p = display.NewByline()
				}
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "show-owner",
		Aliases:            []string{"so", "author", "owner"},
		Usage:              "show owner",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				contentFunc = append(contentFunc, contentFilter.EnableOwner(r))
				if _, ok := p.(*display.Byline); !ok {
					p = display.NewByline()
				}
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "show-group",
		Aliases:            []string{"sg", "group"},
		Usage:              "show group",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				contentFunc = append(contentFunc, contentFilter.EnableGroup(r))
				if _, ok := p.(*display.Byline); !ok {
					p = display.NewByline()
				}
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "show-time",
		Aliases:            []string{"st", "time"},
		Usage:              "show time",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				for _, s := range timeType {
					contentFunc = append(contentFunc, filter.EnableTime(timeFormat, s, r))
				}
				if _, ok := p.(*display.Byline); !ok {
					p = display.NewByline()
				}
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "show-icon",
		Usage:              "show icon",
		Aliases:            []string{"si", "icons", "icon"},
		DisableDefaultText: true,
		Category:           "VIEW",
	},
	&cli.BoolFlag{
		Name:               "show-total-size",
		Usage:              "show total size",
		Aliases:            []string{"ts", "total-size"},
		DisableDefaultText: true,
		Category:           "VIEW",
		Action: func(context *cli.Context, b bool) error {
			if b {
				sizeEnabler.SetEnableTotal()
			}
			return nil
		},
	},
	&cli.StringFlag{
		Name:        "exact-detect-size",
		Usage:       "set exact size for mimetype detection eg:1M/nolimit/infinity",
		Aliases:     []string{"eds", "detect-size", "ds"},
		Value:       "1M",
		DefaultText: "1M",
		Category:    "VIEW",
	},
	&cli.BoolFlag{
		Name:               "mime-type",
		Usage:              "show mime file type",
		Aliases:            []string{"mime", "mimetype"},
		DisableDefaultText: true,
		Category:           "VIEW",
		Action: func(context *cli.Context, b bool) error {
			if b {
				exact := filter.NewMimeFileTypeEnabler()

				err := limitOnce.Do(func() error {
					size := context.String("exact-detect-size")
					var bytes uint64 = 1024 * 1024
					if size == "0" || strings.EqualFold(size, "infinity") || strings.EqualFold(size, "nolimit") {
						bytes = 0
					} else if size != "" {
						sizeUint, err := filter.ParseSize(size)
						if err != nil {
							return err
						}
						bytes = sizeUint.Bytes
					}
					mimetype.SetLimit(uint32(bytes))
					return nil
				})
				if err != nil {
					return err
				}
				contentFunc = append(contentFunc, exact.Enable())
				wgs = append(wgs, exact)
			}
			return nil
		},
	},
	&cli.StringSliceFlag{
		Name:     "checksum-algorithm",
		Usage:    "show checksum of file with algorithm: md5, sha1, sha224, sha256, sha384, sha512, crc32",
		Aliases:  []string{"ca"},
		Value:    cli.NewStringSlice("sha1"),
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "checksum",
		Usage:              "show checksum of file with algorithm: md5, sha1(default), sha224, sha256, sha384, sha512, crc32",
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
				case "sha224":
					sums = append(sums, filter.SumTypeSha224)
				case "sha256":
					sums = append(sums, filter.SumTypeSha256)
				case "sha384":
					sums = append(sums, filter.SumTypeSha384)
				case "sha512":
					sums = append(sums, filter.SumTypeSha512)
				case "crc32":
					sums = append(sums, filter.SumTypeCRC32)
				}
			}

			if b {
				contentFunc = append(contentFunc, contentFilter.EnableSum(sums...))
				if _, ok := p.(*display.Byline); !ok {
					p = display.NewByline()
				}
			}
			return nil
		},
	},
	&cli.BoolFlag{
		Name:               "git-status",
		Usage:              "show git status: ? untracked, + added, ! deleted, ~ modified, | renamed, = copied, $ ignored [if git is installed]",
		Aliases:            []string{"gs", "git"},
		DisableDefaultText: true,
		Category:           "VIEW",
	},
	&cli.StringFlag{
		Name:     "git-status-style",
		Usage:    "git status style: colored-symbol: {? untracked, + added, - deleted, ~ modified, | renamed, = copied, ! ignored} colored-dot",
		Aliases:  []string{"gss", "git-style"},
		Category: "VIEW",
	},

	&cli.BoolFlag{
		Name:    "quote-name",
		Aliases: []string{"Q"},
		Usage:   "enclose entry names in double quotes(overridden by --literal)",
	},
	// &cli.StringFlag{
	// 	Name:    "quoting-style",
	// 	Aliases: []string{"Qs"},
	// 	Usage:   "use quoting style: literal, shell, shell-always, c, escape, locale, clocale",
	// },
	&cli.BoolFlag{
		Name:    "literal",
		Aliases: []string{"N"},
		Usage:   "print entry names without quoting",
	},
	&cli.BoolFlag{
		Name:    "link",
		Aliases: []string{"H"},
		Usage:   "list each file's number of hard links",
		Action: func(context *cli.Context, b bool) error {
			if b {
				link := filter.NewLinkEnabler()
				contentFunc = append(contentFunc, link.Enable())
				wgs = append(wgs, link)
			}
			return nil
		},
	},
}

var displayFlag = []cli.Flag{
	// DISPLAY
	&cli.BoolFlag{
		Name:               "tree",
		Aliases:            []string{"t"},
		Usage:              "recursively list in tree",
		DisableDefaultText: true,
		Category:           "DISPLAY",
	},
	&cli.IntFlag{
		Name:        "depth",
		Usage:       "limit recursive depth, negative -> infinity",
		DefaultText: "infinity",
		Value:       -1,
		Category:    "DISPLAY",
	},
	&cli.BoolFlag{
		Name:               "recurse",
		Aliases:            []string{"R"},
		Usage:              "recurse into directories",
		DisableDefaultText: true,
		Category:           "DISPLAY",
		Action: func(context *cli.Context, b bool) error {
			if b {
				if context.Args().Len() > 1 {
					return fmt.Errorf("'--recurse' should not be used with more than one directory")
				}
			}
			return nil
		},
	},
	&cli.BoolFlag{
		Name:               "byline",
		Aliases:            []string{"bl", "1", "oneline", "single-column"},
		Usage:              "print by line",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				if _, ok := p.(*display.Byline); !ok {
					p = display.NewByline()
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
				if _, ok := p.(*display.Zero); !ok {
					p = display.NewZero()
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
				if _, ok := p.(*display.CommaPrint); !ok {
					p = display.NewCommaPrint()
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
				if _, ok := p.(*display.Across); !ok {
					p = display.NewAcross()
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
				if _, ok := p.(*display.FitTerminal); !ok {
					p = display.NewFitTerminal()
				}
			}
			return nil
		},
		Category: "DISPLAY",
	},
	&cli.BoolFlag{
		Name:               "json",
		Aliases:            []string{"j"},
		Usage:              "output in json format",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				if _, ok := p.(*display.JsonPrinter); !ok {
					p = display.NewJsonPrinter()
				}
			}
			return nil
		},
	},
	&cli.StringFlag{
		Name:        "format",
		DefaultText: "C",
		Usage:       "across  -x,  commas  -m, horizontal -x, long -l, single-column -1, verbose -l, vertical -C",
		Action: func(context *cli.Context, s string) error {
			switch s {
			case "across", "x", "horizontal":
				if _, ok := p.(*display.Across); !ok {
					p = display.NewAcross()
				}
			case "commas", "m":
				if _, ok := p.(*display.CommaPrint); !ok {
					p = display.NewCommaPrint()
				}
			case "long", "l", "verbose":
				contentFunc = append(contentFunc, filter.EnableFileMode(r), sizeEnabler.EnableSize(sizeUint), contentFilter.EnableOwner(r), contentFilter.EnableGroup(r))
				for _, s := range timeType {
					contentFunc = append(contentFunc, filter.EnableTime(timeFormat, s, r))
				}
				if _, ok := p.(*display.Byline); !ok {
					p = display.NewByline()
				}
			case "single-column", "1":
				if _, ok := p.(*display.Byline); !ok {
					p = display.NewByline()
				}
			case "vertical", "C":
				if _, ok := p.(*display.FitTerminal); !ok {
					p = display.NewFitTerminal()
				}
			}
			return nil
		},
		Category: "DISPLAY",
	},

	&cli.BoolFlag{
		Name:               "colorless",
		Aliases:            []string{"nc", "no-color"},
		Usage:              "without color",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				*r = *render.NewRenderer(theme.Colorless, theme.ColorlessInfo)
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
		Name:  "classic",
		Usage: "Enable classic mode (no colours or icons)",
		Action: func(context *cli.Context, b bool) error {
			if b {
				*r = *render.NewRenderer(theme.Colorless, theme.ColorlessInfo)
			}
			err := context.Set("si", "0")
			if err != nil {
				return err
			}
			return nil
		},
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
		Name:               "show-only-hidden",
		Aliases:            []string{"soh", "hidden"},
		DisableDefaultText: true,
		Usage:              "show only hidden files(overridden by --show-hidden/-sh/-a/-A)",
		Action: func(context *cli.Context, b bool) error {
			if b {
				newFF := make([]*filter.TypeFunc, 0, len(typeFunc))
				for _, typeFunc := range typeFunc {
					if typeFunc != &filter.RemoveHidden {
						newFF = append(newFF, typeFunc)
					}
				}
				typeFunc = append(newFF, &filter.HiddenOnly)
			}
			return nil
		},
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
		Aliases:            []string{"sd", "dir", "only-dir", "D"},
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
	&cli.StringSliceFlag{
		Name:     "show-exact-file-type-only",
		Usage:    "only show file with given type",
		Aliases:  []string{"et-only", "eto"},
		Category: "FILTERING",
		Action: func(context *cli.Context, i []string) error {
			if len(i) > 0 {
				err := limitOnce.Do(func() error {
					size := context.String("exact-detect-size")
					var bytes uint64 = 1024 * 1024
					if size == "0" || size == "infinity" {
						bytes = 0
					} else if size != "" {
						sizeUint, err := filter.ParseSize(size)
						if err != nil {
							return err
						}
						bytes = sizeUint.Bytes
					}
					mimetype.SetLimit(uint32(bytes))
					return nil
				})
				if err != nil {
					return err
				}
				eft := filter.ExactFileTypeOnly(i...)
				typeFunc = append(typeFunc, &eft)
			}
			return nil
		},
	},
}

var indexFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:               "disable-index",
		Aliases:            []string{"di", "no-update"},
		Usage:              "disable updating index",
		Category:           "Index",
		DisableDefaultText: true,
	},
	&cli.BoolFlag{
		Name:               "rebuild-index",
		Aliases:            []string{"ri", "remove-all"},
		Usage:              "rebuild index",
		DisableDefaultText: true,
		Category:           "Index",
		Action: func(context *cli.Context, b bool) error {
			if b {
				err := index.RebuildIndex()
				if err != nil {
					return err
				}
			}
			return Err4Exit{}
		},
	},
	&cli.BoolFlag{
		Name:               "fuzzy",
		Aliases:            []string{"fz", "f"},
		Usage:              "fuzzy search",
		DisableDefaultText: true,
		Category:           "Index",
		EnvVars:            []string{"G_FZF"},
	},
	&cli.StringSliceFlag{
		Name:     "remove-index",
		Aliases:  []string{"rm"},
		Usage:    "remove paths from index",
		Category: "Index",
		Action: func(context *cli.Context, i []string) error {
			var errSum error = nil

			beautification := true
			if context.Bool("np") { // --no-path-transform
				beautification = false
			}

			for _, s := range i {
				if beautification {
					s = pathbeautify.Transform(s)
				}

				// get absolute path
				r, err := filepath.Abs(s)
				if err != nil {
					errSum = errors.Join(errSum, fmt.Errorf("remove-path: %w", err))
					continue
				}

				err = index.Delete(r)
				if err != nil {
					errSum = errors.Join(errSum, fmt.Errorf("remove-path: %w", err))
				}
			}
			if errSum != nil {
				return errSum
			} else {
				return Err4Exit{}
			}
		},
	},
	&cli.BoolFlag{
		Name:               "list-index",
		Aliases:            []string{"li"},
		Usage:              "list index",
		DisableDefaultText: true,
		Category:           "Index",
		Action: func(context *cli.Context, b bool) error {
			if b {
				keys, _, err := index.All()
				if err != nil {
					return err
				}
				for i := 0; i < len(keys); i++ {
					fmt.Println(keys[i])
				}
			}
			return Err4Exit{}
		},
	},
	&cli.BoolFlag{
		Name:     "remove-current-path",
		Aliases:  []string{"rcp", "rc", "rmc"},
		Usage:    "remove current path from index",
		Category: "Index",
		Action: func(context *cli.Context, b bool) error {
			if b {
				r, err := os.Getwd()
				if err != nil {
					return err
				}
				err = index.Delete(r)
				if err != nil {
					return err
				}
			}
			return Err4Exit{}
		},
	},
}

var sortingFlags = []cli.Flag{
	&cli.StringSliceFlag{
		Name:    "sort",
		Aliases: []string{"SORT_FIELD"},
		Usage:   "sort by field, default: ascending and case insensitive, field beginning with Uppercase is case sensitive, available fields: nature(default),none(nosort),name,size,time,owner,group,extension. following `-descend` to sort descending",
		Action: func(context *cli.Context, slice []string) error {
			sorter.WithSize(len(slice))(sort)
			for _, s := range slice {
				switch s {
				case "nature":
				case "none", "None", "nosort", "U":
					sort.AddOption(sorter.ByNone)
				case "name-descend":
					sort.AddOption(sorter.ByNameDescend)
				case "name":
					sort.AddOption(sorter.ByNameAscend)
				case "Name":
					sort.AddOption(sorter.ByNameCaseSensitiveAscend)
				case "Name-descend":
					sort.AddOption(sorter.ByNameCaseSensitiveDescend)
				case "size-descend", "S", "sizesort":
					sort.AddOption(sorter.BySizeDescend)
				case "size":
					sort.AddOption(sorter.BySizeAscend)
				case "time-descend":
					sort.AddOption(sorter.ByTimeDescend(timeType[0]))
				case "time":
					sort.AddOption(sorter.ByTimeAscend(timeType[0]))
				case "extension-descend", "ext-descend":
					sort.AddOption(sorter.ByExtensionDescend)
				case "extension", "ext", "x", "extentionsort":
					sort.AddOption(sorter.ByExtensionAscend)
				case "Extension-descend", "Ext-descend":
					sort.AddOption(sorter.ByExtensionCaseSensitiveDescend)
				case "Extension", "Ext", "X", "Extentionsort":
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
				case "width-descend", "Width-descend":
					sort.AddOption(sorter.ByNameWidthDescend)
				case "width", "Width":
					sort.AddOption(sorter.ByNameWidthAscend)
				case "mime", "mimetype", "Mime", "Mimetype":
					err := limitOnce.Do(func() error {
						size := context.String("exact-detect-size")
						var bytes uint64 = 1024 * 1024
						if size == "0" || strings.EqualFold(size, "infinity") || strings.EqualFold(size, "nolimit") {
							bytes = 0
						} else if size != "" {
							sizeUint, err := filter.ParseSize(size)
							if err != nil {
								return err
							}
							bytes = sizeUint.Bytes
						}
						mimetype.SetLimit(uint32(bytes))
						return nil
					})
					if err != nil {
						return err
					}
					sort.AddOption(sorter.ByMimeTypeAscend)
				case "mime-descend", "mimetype-descend", "Mime-descend", "Mimetype-descend":
					err := limitOnce.Do(func() error {
						size := context.String("exact-detect-size")
						var bytes uint64 = 1024 * 1024
						if size == "0" || strings.EqualFold(size, "infinity") || strings.EqualFold(size, "nolimit") {
							bytes = 0
						} else if size != "" {
							sizeUint, err := filter.ParseSize(size)
							if err != nil {
								return err
							}
							bytes = sizeUint.Bytes
						}
						mimetype.SetLimit(uint32(bytes))
						return nil
					})
					if err != nil {
						return err
					}
					sort.AddOption(sorter.ByMimeTypeDescend)
				//	todo
				//	case "v", "version":
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
		Aliases:            []string{"sr", "reverse"},
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
	&cli.BoolFlag{
		Name:               "S",
		Aliases:            []string{"sort-size", "sort-by-size", "sizesort"},
		Usage:              "sort by file size, largest first(descending)",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			sort.Reset()
			return nil
		},
		Category: "SORTING",
	},
	&cli.BoolFlag{
		Name:    "U",
		Aliases: []string{"nosort", "no-sort"},
		Usage:   "do not sort; list entries in directory order. ",
		Action: func(context *cli.Context, b bool) error {
			sort.AddOption(sorter.ByNone)
			return nil
		},
		Category: "SORTING",
	},
	&cli.BoolFlag{
		Name:    "X",
		Aliases: []string{"extensionsort", "Extentionsort"},
		Usage:   "sort alphabetically by entry extension",
		Action: func(context *cli.Context, b bool) error {
			sort.AddOption(sorter.ByExtensionAscend)
			return nil
		},
		Category: "SORTING",
	},
	&cli.BoolFlag{
		Name:  "width",
		Usage: "sort by entry name width",
		Action: func(context *cli.Context, b bool) error {
			sort.AddOption(sorter.ByNameWidthAscend)
			return nil
		},
		Category: "SORTING",
	},
	// todo sort by version
	// &cli.BoolFlag{
	// 	Name:  "v",
	// 	Usage: "sort by version",
	// 	Action: func(context *cli.Context, b bool) error {
	// 		sort.AddOption(sorter.ByVersionAscend)
	// 		return nil
	// 	},
	// 	Category: "SORTING",
	// },
	&cli.BoolFlag{
		Name:    "sort-by-mimetype",
		Aliases: []string{"mimetypesort", "Mimetypesort", "sort-by-mime"},
		Usage:   "sort by mimetype",
		Action: func(context *cli.Context, b bool) error {
			err := limitOnce.Do(func() error {
				size := context.String("exact-detect-size")
				var bytes uint64 = 1024 * 1024
				if size == "0" || strings.EqualFold(size, "infinity") || strings.EqualFold(size, "nolimit") {
					bytes = 0
				} else if size != "" {
					sizeUint, err := filter.ParseSize(size)
					if err != nil {
						return err
					}
					bytes = sizeUint.Bytes
				}
				mimetype.SetLimit(uint32(bytes))
				return nil
			})
			if err != nil {
				return err
			}

			sort.AddOption(sorter.ByMimeTypeAscend)
			return nil
		},
		Category: "SORTING",
	},
	&cli.BoolFlag{
		Name:    "sort-by-mimetype-descend",
		Aliases: []string{"mimetypesort-descend", "Mimetypesort-descend"},
		Usage:   "sort by mimetype, descending",
		Action: func(context *cli.Context, b bool) error {
			if b {
				err := limitOnce.Do(func() error {
					size := context.String("exact-detect-size")
					var bytes uint64 = 1024 * 1024
					if size == "0" || strings.EqualFold(size, "infinity") || strings.EqualFold(size, "nolimit") {
						bytes = 0
					} else if size != "" {
						sizeUint, err := filter.ParseSize(size)
						if err != nil {
							return err
						}
						bytes = sizeUint.Bytes
					}
					mimetype.SetLimit(uint32(bytes))
					return nil
				})
				if err != nil {
					return err
				}

				sort.AddOption(sorter.ByMimeTypeDescend)
			}
			return nil
		},
		Category: "SORTING",
	},
	&cli.BoolFlag{
		Name:    "sort-by-mimetype-parent",
		Aliases: []string{"mimetypesort-parent", "Mimetypesort-parent", "sort-by-mime-parent"},
		Usage:   "sort by mimetype parent",
		Action: func(context *cli.Context, b bool) error {
			if b {
				err := limitOnce.Do(func() error {
					size := context.String("exact-detect-size")
					var bytes uint64 = 1024 * 1024
					if size == "0" || strings.EqualFold(size, "infinity") || strings.EqualFold(size, "nolimit") {
						bytes = 0
					} else if size != "" {
						sizeUint, err := filter.ParseSize(size)
						if err != nil {
							return err
						}
						bytes = sizeUint.Bytes
					}
					mimetype.SetLimit(uint32(bytes))
					return nil
				})
				if err != nil {
					return err
				}

				sort.AddOption(sorter.ByMimeTypeParentAscend)
			}
			return nil
		},
	},
	&cli.BoolFlag{
		Name:    "sort-by-mimetype-parent-descend",
		Aliases: []string{"mimetypesort-parent-descend", "Mimetypesort-parent-descend", "sort-by-mime-parent-descend"},
		Usage:   "sort by mimetype parent",
		Action: func(context *cli.Context, b bool) error {
			if b {
				err := limitOnce.Do(func() error {
					size := context.String("exact-detect-size")
					var bytes uint64 = 1024 * 1024
					if size == "0" || strings.EqualFold(size, "infinity") || strings.EqualFold(size, "nolimit") {
						bytes = 0
					} else if size != "" {
						sizeUint, err := filter.ParseSize(size)
						if err != nil {
							return err
						}
						bytes = sizeUint.Bytes
					}
					mimetype.SetLimit(uint32(bytes))
					return nil
				})
				if err != nil {
					return err
				}

				sort.AddOption(sorter.ByMimeTypeParentDescend)
			}
			return nil
		},
	},
}

func MakeErrorStr(msg string) string {
	return fmt.Sprintf("%s g: %s %s\n", theme.Error, msg, theme.Reset)
}
