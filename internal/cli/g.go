package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Equationzhao/g/internal/align"
	"github.com/Equationzhao/g/internal/config"
	"github.com/Equationzhao/g/internal/const"
	contents "github.com/Equationzhao/g/internal/content"
	"github.com/Equationzhao/g/internal/display"
	"github.com/Equationzhao/g/internal/filter"
	"github.com/Equationzhao/g/internal/index"
	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/render"
	"github.com/Equationzhao/g/internal/shell"
	"github.com/Equationzhao/g/internal/sorter"
	"github.com/Equationzhao/g/internal/theme"
	"github.com/Equationzhao/g/internal/util"
	"github.com/Equationzhao/pathbeautify"
	"github.com/hako/durafmt"
	"github.com/panjf2000/ants/v2"
	"github.com/savioxavier/termlink"
	"github.com/urfave/cli/v2"
	"github.com/xrash/smetrics"
	"go.szostok.io/version/upgrade"
)

var (
	itemFilterFunc = make([]*filter.ItemFilterFunc, 0)
	contentFunc    = make([]contents.ContentOption, 0)
	noOutputFunc   = make([]contents.NoOutputOption, 0)
	r              = render.NewRenderer(&theme.DefaultAll)
	p              = display.NewFitTerminal()
	timeFormat     = "02.Jan'06 15:04"
	// ReturnCode - Exit status:
	//  0  if OK,
	//  1  if minor problems (e.g., cannot access subdirectory),
	//  2  if serious trouble (e.g., cannot access command-line argument).
	ReturnCode      = 0
	contentFilter   = contents.NewContentFilter()
	sort            = sorter.NewSorter()
	timeType        = []string{"mod"}
	sizeUint        = contents.Auto
	sizeEnabler     = contents.NewSizeEnabler()
	blockEnabler    = contents.NewBlockSizeEnabler()
	ownerEnabler    = contents.NewOwnerEnabler()
	groupEnabler    = contents.NewGroupEnabler()
	gitEnabler      = contents.NewGitEnabler()
	gitRepoEnabler  = contents.NewGitRepoEnabler()
	depthLimitMap   map[string]int
	limitOnce       = util.Once{}
	hookOnce        = util.Once{}
	duplicateDetect = contents.NewDuplicateDetect()
	hookPost        = make([]func(display.Printer, ...*item.FileInfo), 0)
	pool            *ants.Pool
	allPart         []string
)

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
			ReturnCode = 2
			str := err.Error()
			const prefix = "flag provided but not defined: "
			if strings.HasPrefix(str, prefix) {
				suggest := suggestFlag(cCtx.App.Flags, strings.TrimLeft(strings.TrimPrefix(str, prefix), "-"))
				if suggest != "" {
					str = fmt.Sprintf("%s, Did you mean %s?", str, suggest)
				}
			}
			_, _ = fmt.Fprintln(os.Stderr, MakeErrorStr(str))
			return nil
		},
		Flags:  make([]cli.Flag, 0, len(viewFlag)+len(filteringFlag)+len(sortingFlags)+len(displayFlag)+len(indexFlags)),
		Action: logic,
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
			Usage: `show the init script for shell, support zsh, bash, fish, powershell, nushell`,
			Action: func(context *cli.Context, s string) error {
				init, err := shell.Init(s)
				if err != nil {
					return err
				}
				_, _ = fmt.Println(init)
				return Err4Exit{}
			},
			Category: "SHELL",
		},
		&cli.BoolFlag{
			Name:               "no-config",
			Usage:              "do not load config file",
			DisableDefaultText: true,
		},
		&cli.BoolFlag{
			Name:               "bug",
			Usage:              "report bug",
			DisableDefaultText: true,
			Action: func(context *cli.Context, b bool) error {
				_, _ = fmt.Println("please report bug to equationzhao@foxmail.com\nor file an issue at https://github.com/Equationzhao/g/issues")
				return Err4Exit{}
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

// dive
// for generating file tree
func dive(
	parent string, depth, limit int, infos *util.Slice[*item.FileInfo], errSlice *util.Slice[error],
	wg *sync.WaitGroup, itemFilter *filter.ItemFilter,
) {
	defer wg.Done()
	if limit > 0 && depth > limit {
		return
	}
	dir, err := os.ReadDir(parent)
	if err != nil {
		errSlice.AppendTo(err)
		return
	}
	for _, entry := range dir {
		f, err := entry.Info()
		if err != nil {
			errSlice.AppendTo(err)
			continue
		}
		nowAbs := filepath.Join(parent, f.Name())
		info, _ := item.NewFileInfoWithOption(item.WithAbsPath(nowAbs), item.WithFileInfo(f))
		// check filter
		if !itemFilter.Match(info) {
			continue
		}
		// store its parent and level/depth
		info.Cache["parent"] = []byte(parent)
		info.Cache["level"] = []byte(strconv.Itoa(depth))
		infos.AppendTo(info)
		if entry.IsDir() {
			wg.Add(1)
			err := pool.Submit(
				func() {
					dive(info.FullPath, depth+1, limit, infos, errSlice, wg, itemFilter)
				},
			)
			if err != nil {
				panic(err)
			}
		}
	}
}

func fuzzyUpdate(path string) error {
	err := index.Update(path)
	if err != nil {
		return err
	}
	return nil
}

// fuzzyPath find the fuzzy path in index
// if error occurs, return empty string and error
func fuzzyPath(path string) (newPath string, minorErr error) {
	fuzzed, err := index.FuzzySearch(path)
	if err == nil {
		return fuzzed, nil
	}
	return "", err
}

// Err4Exit used for exiting without error print
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
	See More at: g.equationzhao.space

{{- if len .Authors}}

AUTHOR{{template "authorsTemplate" .}}{{end}}{{if .VisibleCommands}}

COMMANDS:{{template "visibleCommandCategoryTemplate" .}}{{end}}{{if .VisibleFlagCategories}}

GLOBAL OPTIONS:{{template "visibleFlagCategoryTemplate" .}}{{else if .VisibleFlags}}

GLOBAL OPTIONS:{{template "visibleFlagTemplate" .}}{{end}}
`, filepath.Join(configDir, "g.yaml"),
	)
}

func initVersionHelpFlags() {
	cli.VersionPrinter = func(cCtx *cli.Context) {
		_, _ = fmt.Fprintf(os.Stdout, `💡 g - a powerful ls
 | Version                %s
 | Go Version             %s
 | Compiler               %s
 | Platform               %s

 | Copyright (C) 2024 Equationzhao. MIT License
 | This is free software: you are free to change and redistribute it.
 | There is NO WARRANTY, to the extent permitted by law.
`, Version, runtime.Version(), runtime.Compiler, runtime.GOOS+"/"+runtime.GOARCH)
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
	return fmt.Sprintf("%s × %s %s", constval.Error, msg, constval.Reset)
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
	distance := 0.0
	suggestion := ""
	for _, flag := range flags {
		flagNames := flag.Names()
		for _, name := range flagNames {
			newDistance := smetrics.JaroWinkler(name, provided, constval.BoostThreshold, constval.PrefixSize)
			if newDistance > distance {
				distance = newDistance
				suggestion = name
			}
		}
	}

	// one dash for short flags
	// two dashes for long flags
	if len(suggestion) == 1 {
		suggestion = "-" + suggestion
	} else if len(suggestion) > 1 {
		suggestion = "--" + suggestion
	}

	return suggestion
}

var logic = func(context *cli.Context) error {
	var (
		minorErr   = false
		seriousErr = false
	)

	path := context.Args().Slice()

	nameToDisplay := contents.NewNameEnable()
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
	if _, ok := p.(*display.JsonPrinter); ok {
		nameToDisplay.SetJson()
	}
	git := context.Bool("git")
	if git {
		contentFunc = append(contentFunc, gitEnabler.Enable(r))
	}

	gitBranch := context.Bool("git-repo-branch")
	if gitBranch {
		contentFunc = append(contentFunc, gitRepoEnabler.Enable(r))
	}
	gitRepoStatus := context.Bool("git-repo-status")
	if gitRepoStatus {
		contentFunc = append(contentFunc, gitRepoEnabler.EnableStatus(r))
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

	nameToDisplay.SetQuoteString(`'`)
	// set quote to always
	if context.Bool("Q") {
		nameToDisplay.SetQuote()
	}

	// if no quote, set quote to never
	// this will override the quote set by -Q
	if context.Bool("N") {
		nameToDisplay.UnsetQuote()
	}

	if context.Bool("mounts") {
		nameToDisplay.SetMounts()
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
		nameToDisplay.SetStatistics(&contents.Statistics{})
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

	flagSharp := context.Bool("#")
	tree := context.Bool("tree")
	if tree {
		if _, ok := p.(*display.TreePrinter); !ok {
			p = display.NewTreePrinter()
			if flagSharp {
				p.(*display.TreePrinter).NO = true
			}
		}
	}

	smartGroup := context.Bool("smart-group")
	if smartGroup {
		groupEnabler.EnableSmartMode()
	}

	if n := context.Uint("n"); n > 0 && !tree {
		contentFilter.LimitN = n
	}

	longestEachPart := make(map[string]int)
	startDir, _ := os.Getwd()
	dereference := context.Bool("dereference")

	if !context.Bool("colorless") && !context.Bool("classic") && context.String("theme") == "" {
		if config.Default.ThemeLocation != "" {
			err := theme.GetTheme(config.Default.ThemeLocation)
			if err != nil {
				return err
			}
		}
	}

	theme.ConvertThemeColor()
	{
		var err error
		pool, err = ants.NewPool(ants.DefaultAntsPoolSize)
		if err != nil {
			panic(err)
		}
		contents.Pool = pool
	}

	for i := 0; i < len(path); i++ {
		start := time.Now()

		if len(path) > 1 {
			fmt.Println(r.DirPrompt(path[i]), ":")
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

		stat, err := os.Stat(path[i])
		if err != nil {
			base := filepath.Base(path[i])
			if base == "-" {
				path[i] = os.Getenv("OLDPWD")
				stat, err = os.Stat(path[i])
				if err != nil {
					checkErr(err, "")
					seriousErr = true
					continue
				}
			} else if fuzzy { // no match
				// start fuzzy search
				if newPath, err := fuzzyPath(base); err != nil {
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

		if !disableIndex {
			wgUpdateIndex.Add(1)
			i := i
			err := pool.Submit(
				func() {
					func() {
						if err = fuzzyUpdate(path[i]); err != nil {
							minorErr = true
						}
						wgUpdateIndex.Done()
					}()
				},
			)
			if err != nil {
				panic(err)
			}
		}
		if gitignore {
			*removeGitIgnore = filter.RemoveGitIgnore(path[i])
		}
		if isFile {
			// remove non-display items
			infos = itemFilter.Filter(infos...)

			if tree {
				infos[0].Cache["level"] = []byte("0")
			}
			goto final
		}

		if tree { // visit the dir recursively
			info, err := item.NewFileInfo(path[i])
			if err != nil {
				seriousErr = true
				checkErr(err, originPath)
				continue
			}
			infos = append(
				infos, info,
			)
			infos[0].Cache["level"] = []byte("0")
			if depth >= 1 || depth < 0 {
				wg := sync.WaitGroup{}
				infoSlice := util.NewSlice[*item.FileInfo](10)
				errSlice := util.NewSlice[error](10)
				wg.Add(1)
				pool.Tune(ants.DefaultAntsPoolSize)
				dive(
					path[i], 1, depth, infoSlice, errSlice, &wg,
					itemFilter,
				)
				wg.Wait()
				infos = append(infos, *infoSlice.GetRaw()...)
				for _, err := range *errSlice.GetRaw() {
					if err != nil {
						minorErr = true
						checkErr(err, "")
					}
				}
			}
		} else {
			var d []os.DirEntry
			d, err = os.ReadDir(path[i])
			if err != nil {
				seriousErr = true
				checkErr(err, originPath)
				continue
			}

			if !flagA && !tree { // if -A(almost-all) is not set, add the "."/".." info
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

			// remove non-display items
			infos = itemFilter.Filter(infos...)
		}

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
		if flagR && !tree {

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
						abs := filepath.Join(path[i], info.Name())
						newPath, err := filepath.Rel(startDir, abs)
						if err == nil {
							// if the path is relative, use it
							path = slices.Insert(path, i+1+j, newPath)
						} else {
							path = slices.Insert(path, i+1+j, abs)
						}
						j++
						depthLimitMap[abs] = dep - 1
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

		pool.Tune(1000)
		contentFilter.GetDisplayItems(&infos)

		if len(infos) == 0 {
			goto clean
		}
		{
			// add total && statistics
			addTotalAndStatistic(nameToDisplay, start)

			// if is table printer, set title
			setTitleForPrettyPrinter(path[i])
			if flagSharp {
				setNumber(infos, tree)
			}

			// do scan
			// get max length for each Meta[key].Value
			if len(allPart) == 0 {
				allPart = infos[0].KeysByOrder()
			}

			// make longestEachPart empty
			for s := range longestEachPart {
				longestEachPart[s] = 0
			}

			if _, ok := p.(*display.JsonPrinter); !ok {
				for _, it := range infos {
					for _, part := range allPart {
						content, ok := it.Get(part)
						if ok && part != contents.NameName {
							l := display.WidthNoHyperLinkLen(content.String())
							if l > longestEachPart[part] {
								longestEachPart[part] = l
							}
						}
					}
				}

				// expand the length of each part using the scan result
				for _, it := range infos {
					for _, part := range allPart {
						if part != contents.NameName {
							content, _ := it.Get(part)
							l := display.WidthNoHyperLinkLen(content.String())
							if l < longestEachPart[part] {
								expand := content.SetPrefix
								if align.IsLeft(part) {
									expand = content.SetSuffix
								}
								// expand
								expand(strings.Repeat(" ", longestEachPart[part]-l))
								it.Set(part, content)
							}
						}
					}
				}
			}

			_ = hookOnce.Do(
				func() error {
					if header || footer {
						headerFooter := display.HeaderMaker{
							AllPart:         allPart,
							LongestEachPart: longestEachPart,
						}
						if header {
							headerFooter.Header = true
							headerFooter.IsBefore = true
							// pre scan
							p.AddBeforePrint(headerFooter.Make)
						}
						if footer {
							headerFooter.Footer = true
							if !header {
								// pre scan
								headerFooter.Header = false
								headerFooter.IsBefore = true
								p.AddBeforePrint(headerFooter.Make)
							}
							headerFooter.IsBefore = false
							p.AddAfterPrint(headerFooter.Make)
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
			fmt.Print("\n\n")
			// switch back to start dir
			if err = os.Chdir(startDir); err != nil {
				seriousErr = true
			}
			sizeEnabler.Reset()
		}
	}
	pool.Release()
	wgUpdateIndex.Wait()

	if seriousErr {
		ReturnCode = 2
	} else if minorErr {
		ReturnCode = 1
	}

	return nil
}

func setNumber(infos []*item.FileInfo, isTree bool) {
	l := len(strconv.Itoa(len(infos)))
	for i, info := range infos {
		// if there is #, add No
		var no *display.ItemContent
		if !isTree {
			no = &display.ItemContent{
				No:      -1,
				Content: display.StringContent(strconv.Itoa(i)),
			}
			no.SetSuffix(strings.Repeat(" ", l-len(strconv.Itoa(i))))
		} else {
			no = &display.ItemContent{
				No:      -1,
				Content: display.StringContent(""),
			}
			no.SetSuffix(strings.Repeat(" ", l))
		}
		info.Set("#", no)
	}
}

func setTitleForPrettyPrinter(path string) {
	prettyPrinter, isPrettyPrinter := p.(display.PrettyPrinter)
	if isPrettyPrinter {
		switch p.(type) {
		case *display.CSVPrinter, *display.TSVPrinter:
			break
		default:
			prettyPrinter.SetTitle(path)
		}
	}
}

func addTotalAndStatistic(nameToDisplay *contents.Name, start time.Time) {
	{

		// if -l/show-total-size is set, add total size
		jp, isJsonPrinter := p.(*display.JsonPrinter)

		if isJsonPrinter {
			jp.Extra = make([]any, 0, 2)
		}

		if total, ok := sizeEnabler.Total(); ok {
			s, unit := sizeEnabler.Size2String(total)
			s = r.Size(s, contents.Convert2SizeString(unit))

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
						Time      string               `json:"underwent"`
						Statistic *contents.Statistics `json:"statistic"`
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
}
