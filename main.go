package main

import (
	"errors"
	"fmt"
	"github.com/Equationzhao/g/filter"
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
var p printer.Printer = printer.FitTerminal{}

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
				Aliases:     []string{"d"},
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
				Name:    "show-perm",
				Aliases: []string{"sp"},
				Usage:   "show permission",
				Action: func(context *cli.Context, b bool) error {
					if b {
						contentFunc = append(contentFunc, filter.EnableFileMode(r))
						if _, ok := p.(printer.Byline); !ok {
							p = printer.Byline{}
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
						if _, ok := p.(printer.Byline); !ok {
							p = printer.Byline{}
						}
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:    "show-owner",
				Aliases: []string{"so"},
				Usage:   "show owner",
				Action: func(context *cli.Context, b bool) error {
					if b {
						contentFunc = append(contentFunc, filter.EnableOwner(r))
						if _, ok := p.(printer.Byline); !ok {
							p = printer.Byline{}
						}
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
						if _, ok := p.(printer.Byline); !ok {
							p = printer.Byline{}
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
						contentFunc = append(contentFunc, filter.EnableTime("02.Jan'06 15:04", r))
						if _, ok := p.(printer.Byline); !ok {
							p = printer.Byline{}
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
						if _, ok := p.(printer.Byline); !ok {
							p = printer.Byline{}
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
				Name:    "all",
				Aliases: []string{"la"},
				Usage:   "show all info",
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
						contentFunc = append(contentFunc, filter.EnableFileMode(r), filter.EnableSize(filter.Auto, r), filter.EnableOwner(r), filter.EnableGroup(r), filter.EnableTime("06 Jan 02 15:04", r))
						if _, ok := p.(printer.Byline); !ok {
							p = printer.Byline{}
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
						contentFunc = append(contentFunc, filter.EnableFileMode(r), filter.EnableSize(filter.Auto, r), filter.EnableOwner(r), filter.EnableGroup(r), filter.EnableTime("06 Jan 02 15:04", r))
						if _, ok := p.(printer.Byline); !ok {
							p = printer.Byline{}
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

					if path[i] == "~" {
						home, _ := os.UserHomeDir()
						path[i] = home
					}

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
				si := context.Bool("show-icon") || context.Bool("all")
				for i := 0; i < len(path); i++ {
					if len(path) > 1 {
						fmt.Printf("%s:\n", path[i])
					}
					d, err := os.ReadDir(path[i])
					if err != nil {
						fmt.Println(err)
					}
					if si {
						contentFunc = append(contentFunc, filter.EnableIconName(r, path[i]))
					} else {
						contentFunc = append(contentFunc, filter.EnableName(r))
					}
					res := filter.NewTypeFilter(typeFunc...).Filter(d)
					infos := make([]os.FileInfo, 0, len(res))
					for _, v := range res {
						info, err := v.Info()
						if err != nil {
							fmt.Println(err)
						} else {
							infos = append(infos, info)
						}
					}
					stringSlice := filter.NewContentFilter(contentFunc...).GetStringSlice(infos)
					p.Print(stringSlice...)

					// remove the last func
					if i != len(path)-1 {
						//goland:noinspection GoPrintFunctions
						fmt.Println("\n") //nolint:govet
						contentFunc = contentFunc[:len(contentFunc)-1]
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
