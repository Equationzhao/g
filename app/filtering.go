package app

import (
	"github.com/Equationzhao/g/filter"
	"github.com/Equationzhao/g/filter/content"
	"github.com/gabriel-vasile/mimetype"
	"github.com/urfave/cli/v2"
)

var filteringFlag = []cli.Flag{
	&cli.UintFlag{
		Name:        "limitN",
		Aliases:     []string{"n", "limit"},
		Usage:       "limit n items(n <=0 means unlimited)",
		Value:       0,
		DefaultText: "unlimited",
		Category:    "FILTERING",
	},
	&cli.StringSliceFlag{
		Name:    "ignore-glob",
		Aliases: []string{"I", "ignore", "ig"},
		Usage:   "ignore Glob patterns",
		Action: func(context *cli.Context, s []string) error {
			if len(s) > 0 {
				f, err := filter.RemoveGlob(s...)
				if err != nil {
					return err
				}
				itemFilterFunc = append(itemFilterFunc, &f)
			}
			return nil
		},
		Category: "FILTERING",
	},
	&cli.StringSliceFlag{
		Name:    "match-glob",
		Aliases: []string{"M", "glob", "match"},
		Usage:   "match Glob patterns",
		Action: func(context *cli.Context, s []string) error {
			if len(s) > 0 {
				f, err := filter.GlobOnly(s...)
				if err != nil {
					return err
				}
				itemFilterFunc = append(itemFilterFunc, &f)
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
				newFF := make([]*filter.ItemFilterFunc, 0, len(itemFilterFunc))
				for _, typeFunc := range itemFilterFunc {
					if typeFunc != &filter.RemoveHidden {
						newFF = append(newFF, typeFunc)
					}
				}
				itemFilterFunc = append(newFF, &filter.HiddenOnly)
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
				newFF := make([]*filter.ItemFilterFunc, 0, len(itemFilterFunc))
				for _, typeFunc := range itemFilterFunc {
					if typeFunc != &filter.RemoveHidden {
						newFF = append(newFF, typeFunc)
					}
				}
				itemFilterFunc = newFF
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
				itemFilterFunc = append(itemFilterFunc, &f)
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
				itemFilterFunc = append(itemFilterFunc, &f)
			}
			return nil
		},
		Category: "FILTERING",
	},
	&cli.BoolFlag{
		Name:               "show-no-dir",
		Aliases:            []string{"nd", "nodir", "no-dir", "file"},
		DisableDefaultText: true,
		Usage:              "do not show directory",
		Action: func(context *cli.Context, b bool) error {
			if b {
				itemFilterFunc = append(itemFilterFunc, &filter.RemoveDir)
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
				itemFilterFunc = append(itemFilterFunc, &filter.DirOnly)
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
				itemFilterFunc = append(itemFilterFunc, &filter.RemoveBackups)
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
				newFF := make([]*filter.ItemFilterFunc, 0, len(itemFilterFunc))
				for _, typeFunc := range itemFilterFunc {
					if typeFunc != &filter.RemoveHidden {
						newFF = append(newFF, typeFunc)
					}
				}
				itemFilterFunc = newFF
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
						sizeUint, err := content.ParseSize(size)
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
				eft := filter.MimeTypeOnly(i...)
				itemFilterFunc = append(itemFilterFunc, &eft)
			}
			return nil
		},
	},
	&cli.BoolFlag{
		Name:               "hide-git-ignore",
		Aliases:            []string{"gi", "hgi", "git-ignore"},
		Usage:              "hide git ignored file/dir [if git is installed]",
		DisableDefaultText: true,
		Category:           "FILTERING",
	},
}
