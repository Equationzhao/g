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

	"github.com/Equationzhao/g/config"
	"github.com/Equationzhao/g/display"
	"github.com/Equationzhao/g/filter"
	filtercontent "github.com/Equationzhao/g/filter/content"
	"github.com/Equationzhao/g/index"
	"github.com/Equationzhao/g/item"
	"github.com/Equationzhao/g/render"
	"github.com/Equationzhao/g/shell"
	"github.com/Equationzhao/g/slices"
	"github.com/Equationzhao/g/sorter"
	"github.com/Equationzhao/g/theme"
	"github.com/Equationzhao/g/util"
	"github.com/Equationzhao/pathbeautify"
	"github.com/hako/durafmt"
	"github.com/savioxavier/termlink"
	"github.com/urfave/cli/v2"
	"github.com/valyala/bytebufferpool"
	"github.com/xrash/smetrics"
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
	depthLimitMap   map[string]int
	limitOnce       = util.Once{}
	hookOnce        = util.Once{}
	duplicateDetect = filtercontent.NewDuplicateDetect()
	hookPost        = make([]func(display.Printer, ...*item.FileInfo), 0)
)

var Version = "0.10.0"

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
		Suggest:            true,
		OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
			ReturnCode = 1
			str := err.Error()
			const prefix = "flag provided but not defined: "
			if strings.HasPrefix(str, prefix) {
				suggest := suggestFlag(cCtx.App.Flags, strings.TrimLeft(strings.TrimPrefix(str, prefix), "-"))
				if suggest != "" {
					str = fmt.Sprintf("%s, Did you mean %s?", str, suggest)
				}
			}
			_, _ = fmt.Println(MakeErrorStr(str))
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
			if !context.Bool("no-icon") && (context.Bool("icon") || context.Bool("all")) {
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

			// flag: if d is set, display directory them self
			flagd := context.Bool("d")
			// flag: if A is set
			flagA := context.Bool("A")
			flagR := context.Bool("R")
			if flagR {
				depthLimitMap = make(map[string]int)
			}
			header := context.Bool("header")
			footer := context.Bool("footer")
			if context.Bool("statistic") {
				nameToDisplay.SetStatistics(&filtercontent.Statistics{})
			}

			hyperlink := context.String("hyperlink")
			switch hyperlink {
			case "never":
			case "always":
				nameToDisplay.SetHyperlink()
				display.IncludeHyperlink = true
			default:
				fallthrough
			case "auto":
				if termlink.SupportsHyperlinks() {
					switch p.(type) {
					case display.PrettyPrinter:
					case *display.JsonPrinter:
					default:
						nameToDisplay.SetHyperlink()
						display.IncludeHyperlink = true
					}
				}
			}

			if n := context.Uint("n"); n > 0 {
				contentFilter.LimitN = n
			}

			flagSharp := context.Bool("#")
			longestEachPart := make(map[string]int)
			startDir, _ := os.Getwd()
			dereference := context.Bool("dereference")

			theme.ConvertThemeColor()

			for i := 0; i < len(path); i++ {
				start := time.Now()

				if len(path) > 1 {
					fmt.Printf("%s:\n", path[i])
				}

				if transformEnabled {
					_, err := os.Stat(path[i])
					if err != nil {
						path[i] = pathbeautify.Transform(path[i])
					}
				}
				originPath := path[i]

				infos := make([]*item.FileInfo, 0, 20)

				isFile := false

				// get the abs path
				absPath, err := filepath.Abs(path[i])
				if err != nil {
					_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(fmt.Sprintf("failed to get abs path: %s", absPath)))
					continue
				}
				path[i] = absPath

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
								fmt.Printf("%s:\n", path[i])
								if err != nil {
									checkErr(err, "")
									seriousErr = true
									continue
								}
							}
						} else {
							// output error
							seriousErr = true
							checkErr(err, originPath)
							continue
						}
					}
					if stat.IsDir() {
						if flagd {
							// when -d is set, treat dir as file
							info, err := item.NewFileInfoWithOption(item.WithFileInfo(stat), item.WithPath(path[i]))
							if err != nil {
								checkErr(err, originPath)
								seriousErr = true
								continue
							}
							infos = append(infos, info)
							isFile = true
						}
					} else {
						info, err := item.NewFileInfoWithOption(item.WithFileInfo(stat), item.WithPath(path[i]))
						if err != nil {
							checkErr(err, originPath)
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
					if gitignore {
						*removeGitIgnore = filter.RemoveGitIgnore(filepath.Dir(path[i]))
					}

					// remove non-display items
					infos = itemFilter.Filter(infos...)

					goto final
				}

				d, err = os.ReadDir(path[i])
				if err != nil {
					seriousErr = true
					checkErr(err, originPath)
					continue
				}

				// if -A(almost-all) is not set, add the "."/".." info
				if !flagA {
					err := os.Chdir(path[i])
					if err != nil {
						_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(err.Error()))
						seriousErr = true
					} else {
						FileInfoCurrent, err := item.NewFileInfo(".")
						if err != nil {
							seriousErr = true
							checkErr(err, ".")
						} else {
							infos = append(infos, FileInfoCurrent)
						}

						FileInfoParent, err := item.NewFileInfo("..")
						if err != nil {
							minorErr = true
							checkErr(err, "..")
						} else {
							infos = append(infos, FileInfoParent)
						}
					}
				}

				for _, v := range d {
					info, err := v.Info()
					if err != nil {
						minorErr = true
						checkErr(err, "")
					} else {
						info, err := item.NewFileInfoWithOption(
							item.WithFileInfo(info), item.WithAbsPath(filepath.Join(path[i], v.Name())),
						)
						if err != nil {
							checkErr(err, "")
							seriousErr = true
							continue
						}
						infos = append(infos, info)
					}
				}

				if gitignore {
					*removeGitIgnore = filter.RemoveGitIgnore(path[i])
				}

				// remove non-display items
				infos = itemFilter.Filter(infos...)

				// dereference
				if dereference {
					for i := range infos {
						if util.IsSymLink(infos[i]) {
							symlinks, err := filepath.EvalSymlinks(infos[i].FullPath)
							if err != nil {
								continue
							}
							info, err := os.Stat(symlinks)
							if err != nil {
								continue
							}
							infos[i].FileInfo = info
							infos[i].FullPath = symlinks
						}
					}
				}

				// if -R is set, add sub dir, insert into path[i+1]
				if flagR {

					// set depth
					dep, ok := depthLimitMap[path[i]]
					if !ok {
						depthLimitMap[path[i]] = depth
						dep = depth
					}
					if dep >= 2 || dep <= -1 {
						var j int
						for _, info := range infos {
							if info.IsDir() {
								if info.Name() == "." || info.Name() == ".." {
									continue
								}
								newPath, _ := filepath.Rel(startDir, filepath.Join(path[i], info.Name()))
								path = slices.Insert(path, i+1+j, newPath)
								j++
								depthLimitMap[newPath] = dep - 1
							}
						}
					}
				}

			final:
				if git {
					repo := path[i]
					if isFile {
						repo = filepath.Dir(path[i])
					}
					gitEnabler.Path = repo
					gitEnabler.InitCache(repo)
				}

				contentFilter.GetDisplayItems(&infos)
				if len(infos) == 0 {
					goto clean
				}
				{
					// add total && statistics
					{

						// if -l/show-total-size is set, add total size
						jp, isJsonPrinter := p.(*display.JsonPrinter)

						if total, ok := sizeEnabler.Total(); ok {
							s, unit := sizeEnabler.Size2String(total)
							s = r.Size(s, filtercontent.Convert2SizeString(unit))

							if isJsonPrinter {
								jp.Extra = append(
									jp.Extra, struct {
										Total string `json:"total"`
									}{
										Total: s,
									},
								)
							} else {
								_, _ = display.RawPrint(fmt.Sprintf("  total %s\n", s))
							}
						}
						if s := nameToDisplay.Statistics(); s != nil {
							t := r.Time(durafmt.Parse(time.Since(start)).LimitToUnit("ms").String())
							if isJsonPrinter {
								jp.Extra = append(
									jp.Extra, struct {
										Time      string                    `json:"underwent"`
										Statistic *filtercontent.Statistics `json:"statistic"`
									}{
										Time:      t,
										Statistic: s,
									},
								)
							} else {
								_, _ = display.RawPrint(
									fmt.Sprintf(
										"  underwent %s", t,
									),
								)
								_, _ = display.RawPrint(fmt.Sprintf("\n  statistic: %s\n", s))
							}
							s.Reset()
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
							l := 0
							if part != filtercontent.NameName {
								l = display.WidthNoHyperLinkLen(content.String())
							} else {
								l = display.WidthLen(content.String())
							}
							if l > longestEachPart[part] {
								longestEachPart[part] = l
							}
						}
					}

					// after the first time, expand the length of each part
					for _, it := range infos {
						for _, part := range allPart {
							if part != filtercontent.NameName {
								content, _ := it.Get(part)
								if part != filtercontent.NameName {
									l = display.WidthNoHyperLinkLen(content.String())
								} else {
									l = display.WidthLen(content.String())
								}
								if l < longestEachPart[part] {
									// expand
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
													toAddNum := 0
													if s != filtercontent.NameName {
														toAddNum = len(s) - display.WidthNoHyperLinkLen(content.String())
													} else {
														toAddNum = len(s) - display.WidthLen(content.String())
													}
													content.AddSuffix(
														strings.Repeat(
															" ", toAddNum,
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
				}

			clean:
				if i != len(path)-1 {
					//goland:noinspection GoPrintFunctions
					fmt.Println("\n") //nolint:govet
					// switch back to start dir
					_ = os.Chdir(startDir)
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
			Category: "\b\b\b   META", // add \b to ensure the category is the first one to show
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
			Usage: `By default, .../a/b/c will be transformed to ../../a/b/c, and ~ will be replaced by homedir, 
	using this flag to disable this feature`,
		},
		&cli.BoolFlag{
			Name:               "duplicate",
			Aliases:            []string{"dup"},
			Usage:              "show duplicate files",
			DisableDefaultText: true,
			Action: func(context *cli.Context, b bool) error {
				if b {
					noOutputFunc = append(noOutputFunc, duplicateDetect.Enable())
					hookPost = append(
						hookPost, func(p display.Printer, item ...*item.FileInfo) {
							duplicateDetect.Fprint(p)
							duplicateDetect.Reset()
						},
					)
				}
				return nil
			},
		},
		&cli.StringFlag{
			Name:  "init",
			Usage: `init the config file, default path is ~/.config/g/config.yaml`,
			Action: func(context *cli.Context, s string) error {
				switch s {
				case "zsh":
					_, _ = G.Writer.Write(shell.ZSHContent)
				case "bash":
					_, _ = G.Writer.Write(shell.BASHContent)
				case "fish":
					_, _ = G.Writer.Write(shell.FISHContent)
				case "powershell", "pwsh":
					_, _ = G.Writer.Write(shell.PSContent)
				case "nushell", "nu":
					_, _ = G.Writer.Write(shell.NUContent)
				default:
					return fmt.Errorf("unsupported shell: %s \n %s[zsh|bash|fish|powershell|nushell]", s, theme.Success)
				}
				return Err4Exit{}
			},
			Category: "SHELL",
		},
		&cli.BoolFlag{
			Name:  "no-config",
			Usage: "do not load config file",
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
	configDir, err := config.GetUserConfigDir()
	if err != nil {
		configDir = filepath.Join("$UserConfigDir", "g")
	}
	cli.AppHelpTemplate = fmt.Sprintf(
		`NAME:
	{{template "helpNameTemplate" .}}

USAGE:
	{{if .UsageText}}{{wrap .UsageText 3}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}

VERSION:
	{{.Version}}{{end}}{{end}}{{if .Description}}

DESCRIPTION:
	{{template "descriptionTemplate" .}}{{end}}

CONFIG:
	Configuration: %s
{{- if len .Authors}}

AUTHOR{{template "authorsTemplate" .}}{{end}}{{if .VisibleCommands}}

COMMANDS:{{template "visibleCommandCategoryTemplate" .}}{{end}}{{if .VisibleFlagCategories}}

GLOBAL OPTIONS:{{template "visibleFlagCategoryTemplate" .}}{{else if .VisibleFlags}}

GLOBAL OPTIONS:{{template "visibleFlagTemplate" .}}{{end}}
`, filepath.Join(configDir, "g.yaml"),
	)
}

func initVersionHelpFlags() {
	info := versionInfo.Get()
	info.Version = Version
	s := &style.Config{
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
   
 | Copyright (C) 2023 Equationzhao. MIT License
 | This is free software: you are free to change and redistribute it.
 | There is NO WARRANTY, to the extent permitted by law.
`,
		},
	}

	c := vp.New(vp.WithPrettyStyle(s))
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
		Category:           "\b\b\b   META",
	}

	cli.HelpFlag = &cli.BoolFlag{
		Name:               "help",
		Aliases:            []string{"h", "?"},
		Usage:              "show help",
		DisableDefaultText: true,
		Category:           "\b\b\b   META",
	}
}

func MakeErrorStr(msg string) string {
	return fmt.Sprintf("%s × %s %s", theme.Error, msg, theme.Reset)
}

func checkErr(err error, start string) {
	var toPrint string
	if pathErr := new(os.PathError); errors.As(err, &pathErr) {
		if start != "" {
			pathErr.Path = start
		}
		toPrint = fmt.Sprintf("%s: %s (os error %d)", pathErr.Err, pathErr.Path, pathErr.Err)
	} else {
		toPrint = err.Error()
	}
	_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(toPrint))
}

// suggestFlag returns the suggested flag name
// modified from cli.suggestFlag
func suggestFlag(flags []cli.Flag, provided string) string {
	const (
		boostThreshold = 0.7
		prefixSize     = 4
	)
	distance := 0.0
	suggestion := ""

	for _, flag := range flags {
		flagNames := flag.Names()
		for _, name := range flagNames {
			newDistance := smetrics.JaroWinkler(name, provided, boostThreshold, prefixSize)
			if newDistance > distance {
				distance = newDistance
				suggestion = name
			}
		}
	}

	if len(suggestion) == 1 {
		suggestion = "-" + suggestion
	} else if len(suggestion) > 1 {
		suggestion = "--" + suggestion
	}

	return suggestion
}
