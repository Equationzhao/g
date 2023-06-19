package app

import (
	"fmt"

	"github.com/Equationzhao/g/display"
	"github.com/Equationzhao/g/filter/content"
	"github.com/Equationzhao/g/theme"
	"github.com/urfave/cli/v2"
)

var displayFlag = []cli.Flag{
	// DISPLAY
	&cli.BoolFlag{
		Name:               "tree",
		Aliases:            []string{"t", "T"},
		Usage:              "recursively list in tree",
		DisableDefaultText: true,
		Category:           "DISPLAY",
	},
	&cli.IntFlag{
		Name:        "depth",
		Aliases:     []string{"level", "L"},
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

			_ = context.Set("header", "0")
			_ = context.Set("classic", "1")
			sizeEnabler.DisableTotal()

			return nil
		},
		Category: "DISPLAY",
	},
	&cli.StringFlag{
		Name:    "table-style",
		Aliases: []string{"tablestyle"},
		Usage:   "set table style (ascii(default)/unicode)",
		Action: func(context *cli.Context, s string) error {
			switch s {
			case "ascii", "ASCII", "Ascii":
				// no action needed
			case "unicode", "Unicode", "UNICODE":
				display.DefaultTBStyle = display.UNICODEStyle
			default:
				return fmt.Errorf("invalid table style: %s", s)
			}
			return nil
		},
		Category: "DISPLAY",
	},
	&cli.BoolFlag{
		Name:               "table",
		Aliases:            []string{"tb"},
		Usage:              "output in table format",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				if _, ok := p.(*display.TablePrinter); !ok {
					p = display.NewTablePrinter(display.DefaultTB)
				}
			}
			return nil
		},
		Category: "DISPLAY",
	},
	&cli.BoolFlag{
		Name:               "HTML",
		Aliases:            []string{"html"},
		Usage:              "output in HTML-table format",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				if _, ok := p.(*display.HTMLPrinter); !ok {
					p = display.NewHTMLPrinter()
					r.SetTheme(theme.Colorless)
					r.SetInfoTheme(theme.Colorless)
					theme.Reset = ""
					_ = context.Set("no-icon", "1")
				}
			}
			return nil
		},
		Category: "DISPLAY",
	},
	&cli.BoolFlag{
		Name:               "Markdown",
		Aliases:            []string{"md", "MD", "markdown"},
		Usage:              "output in markdown-table format",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				if _, ok := p.(*display.MDPrinter); !ok {
					p = display.NewMDPrinter()
					r.SetTheme(theme.Colorless)
					r.SetInfoTheme(theme.Colorless)
					theme.Reset = ""
					// _ = context.Set("no-icon", "1")
				}
			}
			return nil
		},
		Category: "DISPLAY",
	},
	&cli.BoolFlag{
		Name:               "CSV",
		Aliases:            []string{"csv"},
		Usage:              "output in csv format",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				if _, ok := p.(*display.CSVPrinter); !ok {
					p = display.NewCSVPrinter()
					r.SetTheme(theme.Colorless)
					r.SetInfoTheme(theme.Colorless)
					theme.Reset = ""
					_ = context.Set("no-icon", "1")
				}
			}
			return nil
		},
		Category: "DISPLAY",
	},
	&cli.StringFlag{
		Name:        "format",
		DefaultText: "C",
		Usage:       "across  -x,  commas  -m, horizontal -x, long -l, single-column -1, verbose -l, vertical -C, table -tb, HTML -html, Markdown -md, CSV -csv, json -j",
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
				contentFunc = append(contentFunc, content.EnableFileMode(r), sizeEnabler.EnableSize(sizeUint), contentFilter.EnableOwner(r), contentFilter.EnableGroup(r))
				for _, s := range timeType {
					contentFunc = append(contentFunc, content.EnableTime(timeFormat, s, r))
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
			case "table", "tb":
				if _, ok := p.(*display.TablePrinter); !ok {
					p = display.NewTablePrinter(display.DefaultTB)
				}
			case "HTML", "html":
				if _, ok := p.(*display.HTMLPrinter); !ok {
					p = display.NewHTMLPrinter()
					r.SetTheme(theme.Colorless)
					r.SetInfoTheme(theme.Colorless)
					theme.Reset = ""
					_ = context.Set("no-icon", "1")
				}
			case "Markdown", "md", "MD", "markdown":
				if _, ok := p.(*display.MDPrinter); !ok {
					p = display.NewMDPrinter()
					r.SetTheme(theme.Colorless)
					r.SetInfoTheme(theme.Colorless)
					theme.Reset = ""
					// _ = context.Set("no-icon", "1")
				}
			case "CSV", "csv":
				if _, ok := p.(*display.CSVPrinter); !ok {
					p = display.NewCSVPrinter()
					r.SetTheme(theme.Colorless)
					r.SetInfoTheme(theme.Colorless)
					theme.Reset = ""
					_ = context.Set("no-icon", "1")
				}
			case "json", "j":
				if _, ok := p.(*display.JsonPrinter); !ok {
					p = display.NewJsonPrinter()
					r.SetTheme(theme.Colorless)
					r.SetInfoTheme(theme.Colorless)
					theme.Reset = ""
					_ = context.Set("no-icon", "1")
				}
			default:
				return fmt.Errorf("unkown format option:%s", s)
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
				r.SetTheme(theme.Colorless)
				r.SetInfoTheme(theme.Colorless)
				theme.Reset = ""
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
			theme.SyncColorlessWithTheme()
			return nil
		},
		Category: "DISPLAY",
	},
	&cli.BoolFlag{
		Name:  "classic",
		Usage: "Enable classic mode (no colours or icons)",
		Action: func(context *cli.Context, b bool) error {
			if b {
				r.SetTheme(theme.Colorless)
				r.SetInfoTheme(theme.Colorless)
				theme.Reset = ""
				err := context.Set("no-icon", "1")
				if err != nil {
					return err
				}
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
