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
	"github.com/Equationzhao/g/item"
	"github.com/Equationzhao/g/render"
	"github.com/Equationzhao/g/sorter"
	"github.com/Equationzhao/g/theme"
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
	itemFilterFunc  = make([]*filter.ItemFilterFunc, 0)
	contentFunc     = make([]filter.ContentOption, 0)
	noOutputFunc    = make([]filter.NoOutputOption, 0)
	r               = render.NewRenderer(theme.DefaultTheme, theme.DefaultInfoTheme)
	p               = display.NewFitTerminal()
	timeFormat      = "02.Jan'06 15:04"
	ReturnCode      = 0
	contentFilter   = filter.NewContentFilter()
	sort            = sorter.NewSorter()
	timeType        = []string{"mod"}
	sizeUint        = filtercontent.Auto
	sizeEnabler     = filtercontent.NewSizeEnabler()
	blockEnabler    = filtercontent.NewBlockSizeEnabler()
	ownerEnabler    = filtercontent.NewOwnerEnabler()
	groupEnabler    = filtercontent.NewGroupEnabler()
	gitEnabler      = filtercontent.NewGitEnabler()
	depthLimitMap   = make(map[string]int)
	limitOnce       = util.Once{}
	hookOnce        = util.Once{}
	duplicateDetect = filtercontent.NewDuplicateDetect()
	hookPost        = make([]func(display.Printer, ...*item.FileInfo), 0)
)

var Version = "0.9.0"

var G *cli.App

func init() {
	itemFilterFunc = append(itemFilterFunc, &filter.RemoveHidden)

	G = &cli.App{
		Name:      "g",
		Usage:     "a powerful ls",
		UsageText: "g [options] [path]",
		Version:   Version,
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

			nameToDisplay := filtercontent.NewNameEnable()
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
			git := context.Bool("git")
			if git {
				contentFunc = append(contentFunc, gitEnabler.Enable(r))
			}
			if context.Bool("no-dereference") {
				nameToDisplay.SetNoDeference()
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
				// git-status-style
				s := context.String("gss")
				switch s {
				case "symbol", "sym":
					nameToDisplay.GitStyle = filtercontent.GitStyleSym
				case "dot", ".":
					nameToDisplay.GitStyle = filtercontent.GitStyleDot
				default:
					nameToDisplay.GitStyle = filtercontent.GitStyleDefault
				}
			}

			// set quote
			if context.Bool("Q") {
				nameToDisplay.SetQuote(`"`)
			}

			// if no quote, set quote to empty
			// this will override the quote set by -Q
			if context.Bool("N") {
				nameToDisplay.UnsetQuote()
			}

			// no path transform
			transformEnabled := !context.Bool("np")
			if rp := context.String("relative-to"); rp != "" {
				if transformEnabled {
					rp = pathbeautify.Beautify(rp)
				}
				if temp, err := filepath.Abs(rp); err == nil {
					rp = temp
				}
				nameToDisplay.SetRelativeTo(rp)
			} else if context.Bool("fp") {
				nameToDisplay.SetFullPath()
			}
			contentFunc = append(contentFunc, nameToDisplay.Enable(r))
			itemFilter := filter.NewItemFilter(itemFilterFunc...)

			gitignore := context.Bool("git-ignore")
			removeGitIgnore := new(filter.ItemFilterFunc)
			if gitignore {
				itemFilter.AppendTo(removeGitIgnore)
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
			depth := context.Int("depth")

			// if context.Bool("tree") {
			// 	for i := 0; i < len(path); i++ {
			// 		start := time.Now()
			//
			// 		if len(path) > 1 {
			// 			fmt.Printf("%s:\n", path[i])
			// 		}
			//
			// 		if transformEnabled {
			// 			path[i] = pathbeautify.Transform(path[i])
			// 		}
			// 		// fuzzy search
			// 		if fuzzy {
			// 			_, err := os.Stat(path[i])
			// 			if err != nil {
			// 				newPath, b := fuzzyPath(path[i])
			// 				if b != nil {
			// 					_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
			// 					minorErr = true
			// 					continue
			// 				} else {
			// 					path[i] = newPath
			// 					_, err = os.Stat(path[i])
			// 					if err != nil {
			// 						checkErr(err)
			// 						seriousErr = true
			// 						continue
			// 					}
			// 					fmt.Println(path[i])
			// 				}
			// 			}
			// 		}
			// 		if gitignore {
			// 			*removeGitIgnore = filter.RemoveGitIgnore(path[i])
			// 		}
			//
			// 		s, err, minorErrInTree := tree.NewTreeString(path[i], depth, itemFilter, contentFilter)
			// 		if pathErr := new(os.PathError); errors.As(err, &pathErr) {
			// 			_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(fmt.Sprintf("%s: %s", pathErr.Err, pathErr.Path)))
			// 			seriousErr = true
			// 			continue
			// 		} else if err != nil {
			// 			_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
			// 			seriousErr = true
			// 			continue
			// 		}
			//
			// 		if pathErr := new(os.PathError); errors.As(minorErrInTree, &pathErr) {
			// 			_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(fmt.Sprintf("%s: %s", pathErr.Err, pathErr.Path)))
			// 			minorErr = true
			// 		} else if minorErrInTree != nil {
			// 			_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
			// 			minorErr = true
			// 		}
			//
			// 		absPath, err := filepath.Abs(path[i])
			// 		if err != nil {
			// 			minorErr = true
			// 		} else {
			// 			if !disableIndex {
			// 				wgUpdateIndex.Add(1)
			// 				go func() {
			// 					if err = fuzzyUpdate(absPath); err != nil {
			// 						minorErr = true
			// 					}
			// 					wgUpdateIndex.Done()
			// 				}()
			// 			}
			// 		}
			//
			// 		fmt.Println(s.MakeTreeStr())
			// 		fmt.Printf("\n%d directories, %d files\nunderwent %s\n", s.Directory(), s.File(), time.Since(start).String())
			//
			// 		if i != len(path)-1 {
			// 			//goland:noinspection GoPrintFunctions
			// 			fmt.Println("\n") //nolint:govet
			// 		}
			// 	}
			// } else {
			startDir, _ := os.Getwd()

			// flag: if d is set, display directory them self
			flagd := context.Bool("d")
			// flag: if A is set
			flagA := context.Bool("A")
			flagR := context.Bool("R")
			header := context.Bool("header")
			footer := context.Bool("footer")
			if context.Bool("statistic") {
				nameToDisplay.SetStatistics(&filtercontent.Statistics{})
			}

			if n := context.Uint("n"); n > 0 {
				contentFilter.LimitN = n
			}

			flagSharp := context.Bool("#")
			theme.ConvertThemeColor()

			longestEachPart := make(map[string]int)

			for i := 0; i < len(path); i++ {
				start := time.Now()

				if len(path) > 1 {
					fmt.Printf("%s:\n", path[i])
				}

				if transformEnabled {
					path[i] = pathbeautify.Transform(path[i])
				}

				infos := make([]*item.FileInfo, 0, 20)

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
									checkErr(err)
									seriousErr = true
									continue
								}
								fmt.Println(path[i])
							}
						} else {
							// output error
							if pathErr := new(os.PathError); errors.As(err, &pathErr) {
								_, _ = fmt.Fprintln(
									os.Stderr, MakeErrorStr(fmt.Sprintf("%s: %s", pathErr.Err, pathErr.Path)),
								)
								seriousErr = true
								continue
							}
							_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
							seriousErr = true
							continue
						}
					}
					if stat.IsDir() {
						if flagd {
							// when -d is set, treat dir as file
							info, err := item.NewFileInfoWithOption(item.WithFileInfo(stat), item.WithPath(path[i]))
							if err != nil {
								checkErr(err)
								seriousErr = true
								continue
							}
							infos = append(infos, info)
							isFile = true
						}
					} else {
						info, err := item.NewFileInfoWithOption(item.WithFileInfo(stat), item.WithPath(path[i]))
						if err != nil {
							checkErr(err)
							seriousErr = true
							continue
						}
						infos = append(infos, info)
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
					if pathErr := new(os.PathError); errors.As(err, &pathErr) {
						_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(fmt.Sprintf("%s: %s", pathErr.Err, pathErr.Path)))
						seriousErr = true
					} else {
						_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
						seriousErr = true
					}
					continue
				}

				// if -A(almost-all) is not set, add the "."/".." info
				if !flagA {
					err := os.Chdir(path[i])
					if err != nil {
						_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
					} else {
						FileInfoCurrent, err := item.NewFileInfo(".")
						if err != nil {
							if pathErr := new(os.PathError); errors.As(err, &pathErr) {
								_, _ = fmt.Fprintln(
									os.Stderr, MakeErrorStr(fmt.Sprintf("%s: %s", pathErr.Err, pathErr.Path)),
								)
								seriousErr = true
							} else {
								_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
								seriousErr = true
							}
						} else {
							infos = append(infos, FileInfoCurrent)
						}

						FileInfoParent, err := item.NewFileInfo("..")
						if err != nil {
							if pathErr := new(os.PathError); errors.As(err, &pathErr) {
								_, _ = fmt.Fprintln(
									os.Stderr, MakeErrorStr(fmt.Sprintf("%s: %s", pathErr.Err, pathErr.Path)),
								)
								minorErr = true
							} else {
								_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
								minorErr = true
							}
						} else {
							infos = append(infos, FileInfoParent)
						}
					}
				}

				for _, v := range d {
					info, err := v.Info()
					if err != nil {
						if pathErr := new(os.PathError); errors.As(err, &pathErr) {
							_, _ = fmt.Fprintln(
								os.Stderr, MakeErrorStr(fmt.Sprintf("%s: %s", pathErr.Err, pathErr.Path)),
							)
							minorErr = true
						} else {
							minorErr = true
							_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
						}
					} else {
						info, err := item.NewFileInfoWithOption(item.WithFileInfo(info), item.WithPath(v.Name()))
						if err != nil {
							checkErr(err)
							seriousErr = true
							continue
						}
						infos = append(infos, info)
					}
				}

				if gitignore {
					*removeGitIgnore = filter.RemoveGitIgnore(path[i])
				}
				if git {
					gitEnabler.Path = path[i]
					gitEnabler.InitCache(path[i])
				}

				// remove non-display items
				infos = itemFilter.Filter(infos...)

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
								newPath, _ := filepath.Rel(startDir, filepath.Join(path[i], info.Name()))
								newPathLeft = append(newPathLeft, newPath)
								depthLimitMap[newPath] = dep - 1
							}
						}
						path = append(path[:i+1], append(newPathLeft, path[i+1:]...)...)
					}
				}

			final:
				contentFilter.GetDisplayItems(&infos)
				if len(infos) == 0 {
					continue
				}

				// add total && statistics
				{
					var i *item.FileInfo
					// if -l/show-total-size is set, add total size
					_, isPrettyPrinter := p.(display.PrettyPrinter)

					if total, ok := sizeEnabler.Total(); ok {
						if !isPrettyPrinter {
							i, _ = item.NewFileInfoWithOption()
							s, _ := sizeEnabler.Size2String(total)
							i.Set(
								"total", &display.ItemContent{No: 0, Content: display.StringContent(
									fmt.Sprintf(
										"  total %s", s,
									),
								)},
							)
						} else {
							s, _ := sizeEnabler.Size2String(total)
							_, _ = display.RawPrint(fmt.Sprintf("  total %s\n", s))
						}
					}
					if s := nameToDisplay.Statistics(); s != nil {
						if !isPrettyPrinter {
							tFormat := "\n  underwent %s"
							if i == nil {
								i, _ = item.NewFileInfoWithOption()
								tFormat = "  underwent %s"
							}
							i.Set(
								"underwent", &display.ItemContent{No: 1, Content: display.StringContent(
									fmt.Sprintf(
										tFormat, r.Time(durafmt.Parse(time.Since(start)).LimitToUnit("ms").String()),
									),
								)},
							)
							i.Set(
								"statistic", &display.ItemContent{No: 2, Content: display.StringContent(
									fmt.Sprintf(
										"\n  statistic: %s", s,
									),
								)},
							)
						} else {
							_, _ = display.RawPrint(
								fmt.Sprintf(
									"  underwent %s",
									r.Time(durafmt.Parse(time.Since(start)).LimitToUnit("ms").String()),
								),
							)
							_, _ = display.RawPrint(fmt.Sprintf("\n  statistic: %s\n", s))
						}
						s.Reset()
					}
					if i != nil {
						p.DisablePreHook()
						p.Print(i)
						p.EnablePreHook()
					}
				}

				// if is table printer, set title
				prettyPrinter, isTablePrinter := p.(display.PrettyPrinter)
				if isTablePrinter {
					prettyPrinter.SetTitle(path[i])
				}
				l := len(strconv.Itoa(len(infos)))
				for i, info := range infos {
					// if there is #, add No
					if flagSharp {
						no := &display.ItemContent{
							No:      -1,
							Content: display.StringContent(strconv.Itoa(i)),
						}
						no.SetSuffix(strings.Repeat(" ", l-len(strconv.Itoa(i))))
						info.Set("#", no)
					}
				}

				// do scan
				// get max length for each Meta[key].Value
				allPart := infos[0].KeysByOrder()

				for _, it := range infos {
					for _, part := range allPart {
						content, _ := it.Get(part)
						l := display.WidthLen(content.String())
						if l > longestEachPart[part] {
							longestEachPart[part] = l
						}
					}
				}

				for _, it := range infos {
					for _, part := range allPart {
						content, _ := it.Get(part)
						l := display.WidthLen(content.String())
						if l < longestEachPart[part] {
							// expand
							if part != filtercontent.NameName {
								content.SetSuffix(strings.Repeat(" ", longestEachPart[part]-l))
								it.Set(part, content)
							}
						}
					}
				}

				_ = hookOnce.Do(
					func() error {
						if header || footer {
							headerFooter := func(isBefore bool) func(p display.Printer, items ...*item.FileInfo) {
								return func(p display.Printer, item ...*item.FileInfo) {
									// add header
									if len(item) == 0 {
										return
									}

									// add longest - len(header) * space
									// print header
									headerFooterStrBuf := bytebufferpool.Get()
									defer bytebufferpool.Put(headerFooterStrBuf)
									prettyPrinter, isPrettyPrinter := p.(display.PrettyPrinter)

									expand := func(s string, no, space int) {
										_, _ = headerFooterStrBuf.WriteString(theme.Underline)
										_, _ = headerFooterStrBuf.WriteString(s)
										_, _ = headerFooterStrBuf.WriteString(theme.Reset)
										if no != len(allPart)-1 {
											_, _ = headerFooterStrBuf.WriteString(strings.Repeat(" ", space))
										}
									}

									for i, s := range allPart {
										if len(s) > longestEachPart[s] {
											// expand the every item's content of this part
											for _, it := range infos {
												content, _ := it.Get(s)
												content.AddSuffix(
													strings.Repeat(
														" ", len(s)-display.WidthLen(content.String()),
													),
												)
												it.Set(s, content)
												longestEachPart[s] = len(s)
											}
											expand(s, i, 1)
										} else {
											expand(s, i, longestEachPart[s]-len(s)+1)
										}
										if isPrettyPrinter && isBefore {
											if header {
												prettyPrinter.AddHeader(s)
											}
											if footer {
												prettyPrinter.AddFooter(s)
											}
										}
									}
									res := headerFooterStrBuf.String()
									if !isPrettyPrinter {
										if header && isBefore {
											_, _ = fmt.Fprintln(p, res)
										}
										if footer && !isBefore {
											_, _ = fmt.Fprintln(p, res)
										}
									}
								}
							}
							if header {
								// pre scan
								p.AddBeforePrint(headerFooter(true))
							}
							if footer {
								if !header {
									// pre scan
									p.AddBeforePrint(headerFooter(true))
								}
								p.AddAfterPrint(headerFooter(false))
							}
						}
						p.AddAfterPrint(hookPost...)
						return nil
					},
				)

				p.Print(infos...)

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
			// }
			wgUpdateIndex.Wait()

			if seriousErr {
				ReturnCode = 2
			} else if minorErr {
				ReturnCode = 1
			}

			return nil
		},
	}
	G.Flags = append(
		G.Flags,
		&cli.BoolFlag{
			Name:     "check-new-version",
			Usage:    "check if there's new release",
			Category: "software info",
			Action: func(context *cli.Context, b bool) error {
				if b {
					upgrade.WithUpdateCheckTimeout(5 * time.Second)
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
		},
		&cli.BoolFlag{
			Name:               "duplicate",
			Aliases:            []string{"dup"},
			Usage:              "show duplicate files",
			DisableDefaultText: true,
			Action: func(context *cli.Context, b bool) error {
				if b {
					noOutputFunc = append(noOutputFunc, duplicateDetect.Enable())
					hookPost = append(hookPost, func(p display.Printer, item ...*item.FileInfo) {
						duplicateDetect.Fprint(p)
						duplicateDetect.Reset()
					})
				}
				return nil
			},
		},
	)

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
	}
	return "", err
}

type Err4Exit struct{}

func (c Err4Exit) Error() string {
	panic("it's an error defined to exit app, should not call this")
}

func initHelpTemp() {
	cli.AppHelpTemplate = `NAME:
   {{template "helpNameTemplate" .}}

USAGE:
   {{if .UsageText}}{{wrap .UsageText 3}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}

VERSION:
   {{.Version}}{{end}}{{end}}{{if .Description}}

DESCRIPTION:
   {{template "descriptionTemplate" .}}{{end}}
{{- if len .Authors}}

AUTHOR{{template "authorsTemplate" .}}{{end}}{{if .VisibleCommands}}

COMMANDS:{{template "visibleCommandCategoryTemplate" .}}{{end}}{{if .VisibleFlagCategories}}

GLOBAL OPTIONS:{{template "visibleFlagCategoryTemplate" .}}{{else if .VisibleFlags}}

GLOBAL OPTIONS:{{template "visibleFlagTemplate" .}}{{end}}
`
}

func initVersionHelpFlags() {
	info := versionInfo.Get()
	info.Version = Version
	config := &style.Config{
		Formatting: style.Formatting{
			Header: style.Header{
				Prefix: "💡 ",
				FormatPrimitive: style.FormatPrimitive{
					Color:   "Green",
					Options: []string{"Bold"},
				},
			},
			Key: style.Key{
				FormatPrimitive: style.FormatPrimitive{
					Color:      "Yellow",
					Background: "",
					Options:    nil,
				},
			},
		},
		Layout: style.Layout{
			GoTemplate: `{{ Header .Meta.CLIName }}
| {{ Key "Version"     }}        {{ .Version                     | Val   }}
| {{ Key "Go Version"  }}        {{ .GoVersion  | trimPrefix "go"| Val   }}
| {{ Key "Compiler"    }}        {{ .Compiler                    | Val   }}
| {{ Key "Platform"    }}        {{ .Platform                    | Val   }}
   
Copyright (C) 2023 Equationzhao. MIT License
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.
`,
		},
	}

	c := vp.New(vp.WithPrettyStyle(config))
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
	return fmt.Sprintf("%s × %s %s\n", theme.Error, msg, theme.Reset)
}

func checkErr(err error) {
	if pathErr := new(os.PathError); errors.As(err, &pathErr) {
		_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(fmt.Sprintf("%s: %s", pathErr.Err, pathErr.Path)))
		return
	}
	_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
}
