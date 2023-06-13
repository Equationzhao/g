package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Equationzhao/g/display"
	"github.com/Equationzhao/g/filter"
	filtercontent "github.com/Equationzhao/g/filter/content"
	"github.com/Equationzhao/g/index"
	"github.com/Equationzhao/g/render"
	"github.com/Equationzhao/g/sorter"
	"github.com/Equationzhao/g/theme"
	"github.com/Equationzhao/g/tree"
	"github.com/Equationzhao/g/util"
	"github.com/Equationzhao/pathbeautify"
	"github.com/hako/durafmt"
	"github.com/urfave/cli/v2"
	"github.com/valyala/bytebufferpool"
	versionInfo "go.szostok.io/version"
	vp "go.szostok.io/version/printer"
	"go.szostok.io/version/style"
	"go.szostok.io/version/upgrade"
)

var (
	typeFunc        = make([]*filter.TypeFunc, 0)
	contentFunc     = make([]filter.ContentOption, 0)
	noOutputFunc    = make([]filter.NoOutputOption, 0)
	r               = render.NewRenderer(theme.DefaultTheme, theme.DefaultInfoTheme)
	p               = display.NewFitTerminal()
	timeFormat      = "02.Jan'06 15:04"
	ReturnCode      = 0
	contentFilter   = filter.NewContentFilter()
	CompiledAt      = ""
	sort            = sorter.NewSorter()
	timeType        = []string{"mod"}
	sizeUint        = filtercontent.Auto
	sizeEnabler     = filtercontent.NewSizeEnabler()
	blockEnabler    = filtercontent.NewBlockSizeEnabler()
	wgs             = make([]filter.LengthFixed, 0, 1)
	depthLimitMap   = make(map[string]int)
	limitOnce       = util.Once{}
	hookOnce        = util.Once{}
	duplicateDetect = filtercontent.NewDuplicateDetect()
	hookAfter       = make([]func(display.Printer, ...display.Item), 0)
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
	blockEnabler.SetRenderer(r)
	wgs = append(wgs, sizeEnabler)
	wgs = append(wgs, blockEnabler)

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

			nameToDisplay := filtercontent.NewNameEnable().SetRenderer(r)
			if !context.Bool("no-icon") && (context.Bool("show-icon") || context.Bool("all")) {
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
			if context.Bool("fp") {
				nameToDisplay.SetFullPath()
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
					nameToDisplay.GitStyle = filtercontent.GitStyleSym
				case "dot", ".":
					nameToDisplay.GitStyle = filtercontent.GitStyleDot
				default:
					nameToDisplay.GitStyle = filtercontent.GitStyleDefault
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
			contentFilter.SetNoOutputOptions(noOutputFunc...)

			// if no path, use the current path
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
					fmt.Printf("\n%d directories, %d files\nunderwent %s\n", s.Directory(), s.File(), time.Since(start).String())

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
					nameToDisplay.SetStatistics(&filtercontent.Statistics{})
				}

				flagSharp := context.Bool("#")

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
						_, isPrettyPrinter := p.(display.PrettyPrinter)

						if total, ok := sizeEnabler.Total(); ok {
							if !isPrettyPrinter {
								i = display.NewItem()
								i.Set("total", display.ItemContent{No: 0, Content: display.StringContent(fmt.Sprintf("  total %s", sizeEnabler.Size2String(total, 0)))})
							} else {
								_, _ = display.RawPrint(fmt.Sprintf("  total %s\n", sizeEnabler.Size2String(total, 0)))
							}
						}
						if s := nameToDisplay.Statistics(); s != nil {
							if !isPrettyPrinter {
								tFormat := "\n  underwent %s"
								if i == nil {
									i = display.NewItem()
									tFormat = "  underwent %s"
								}
								i.Set("underwent", display.ItemContent{No: 1, Content: display.StringContent(fmt.Sprintf(tFormat, r.Time(durafmt.Parse(time.Since(start)).LimitToUnit("ms").String())))})
								i.Set("statistic", display.ItemContent{No: 2, Content: display.StringContent(fmt.Sprintf("\n  statistic: %s", s))})
							} else {
								_, _ = display.RawPrint(fmt.Sprintf("  underwent %s", r.Time(durafmt.Parse(time.Since(start)).LimitToUnit("ms").String())))
								_, _ = display.RawPrint(fmt.Sprintf("\n  statistic: %s\n", s))
							}
							s.Reset()
						}
						if i != nil {
							p.DisableHookBefore()
							p.Print(*i)
							p.EnableHookBefore()
						}
					}

					_ = hookOnce.Do(func() error {
						p.AddAfterPrint(hookAfter...)
						if header {
							p.AddBeforePrint(func(p display.Printer, item ...display.Item) {
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
								prettyPrinter, isPrettyPrinter := p.(display.PrettyPrinter)
								for i, s := range allPart {
									if len(s) > longestEachPart[s] {
										// expand the every item's content of this part
										for _, it := range items {
											content, _ := it.Get(s)
											content.Content = display.StringContent(fmt.Sprintf("%s%s", strings.Repeat(" ", len(s)-longestEachPart[s]), content.Content.String()))

											it.Set(s, content)
										}
										if !isPrettyPrinter {
											_, _ = contentStrBuf.WriteString(theme.Underline)
											_, _ = contentStrBuf.WriteString(s)
											_, _ = contentStrBuf.WriteString(theme.Reset)
											if i != len(allPart)-1 {
												_, _ = contentStrBuf.WriteString(" ")
											}
										} else {
											prettyPrinter.AddHeader(s)
										}
									} else {
										if !isPrettyPrinter {
											_, _ = contentStrBuf.WriteString(theme.Underline)
											_, _ = contentStrBuf.WriteString(s)
											_, _ = contentStrBuf.WriteString(theme.Reset)
											if i != len(allPart)-1 {
												_, _ = contentStrBuf.WriteString(strings.Repeat(" ", longestEachPart[s]-len(s)+1))
											}
										} else {
											prettyPrinter.AddHeader(s)
										}
									}
								}
								_, _ = contentStrBuf.WriteString(theme.Reset)
								_, _ = fmt.Fprintln(display.Output, contentStrBuf.String())
								bytebufferpool.Put(contentStrBuf)
							})
						}
						return nil
					})

					itemsCopy := make([]display.Item, 0, len(items))
					{
						// if is table printer, set title
						prettyPrinter, isTablePrinter := p.(display.PrettyPrinter)
						if isTablePrinter {
							prettyPrinter.SetTitle(path[i])
						}
						l := len(strconv.Itoa(len(items)))
						for i, item := range items {
							// if there is #, add No
							if flagSharp {
								item.Set("#", display.ItemContent{
									No:      -1,
									Content: display.StringContent(fmt.Sprintf("%d%s", i, strings.Repeat(" ", l-len(strconv.Itoa(i))))),
								})
							}
							itemsCopy = append(itemsCopy, *item)
						}
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
	G.Flags = append(G.Flags,
		&cli.BoolFlag{
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
		},
		&cli.BoolFlag{
			Name:               "no-path-transform",
			Aliases:            []string{"np"},
			DisableDefaultText: true,
			Usage:              "By default, .../a/b/c will be transformed to ../../a/b/c, and ~ will be replaced by homedir, using this flag to disable this feature",
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
		Aliases:            []string{"h", "?"},
		Usage:              "show help",
		DisableDefaultText: true,
		Category:           "software info",
	}
}

func MakeErrorStr(msg string) string {
	return fmt.Sprintf("%s g: %s %s\n", theme.Error, msg, theme.Reset)
}
