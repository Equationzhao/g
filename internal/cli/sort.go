package cli

import (
	"fmt"
	"slices"

	"github.com/Equationzhao/g/internal/sorter"
	"github.com/urfave/cli/v2"
)

var sortingFlags = []cli.Flag{
	&cli.StringSliceFlag{
		Name:    "sort",
		Aliases: []string{"SORT_FIELD"},
		Usage: `sort by field, default: 
	ascending and case insensitive, 
	field beginning with Uppercase is case sensitive,	
	available fields: 	
	   nature(default),none(nosort),
	   name,.name(sorts by name without a leading dot),	
	   size,time,owner,group,extension,inode,width,mime. 	
	   following '-descend' to sort descending`,
		Action: func(context *cli.Context, slice []string) error {
			if slices.ContainsFunc(slice, func(s string) bool {
				nosort := []string{"none", "None", "nosort", "U"}
				return slices.Contains(nosort, s)
			}) {
				sort.Reset()
				return nil
			}
			sorter.WithSize(len(slice))(sort)
			for _, s := range slice {
				switch s {
				case "nature":
				case "name-descend":
					sort.AddOption(sorter.ByNameDescend)
				case "name":
					sort.AddOption(sorter.ByNameAscend)
				case "Name":
					sort.AddOption(sorter.ByNameCaseSensitiveAscend)
				case "Name-descend":
					sort.AddOption(sorter.ByNameCaseSensitiveDescend)
				case ".name-descend":
					sort.AddOption(sorter.ByNameWithoutALeadingDotDescend)
				case ".name":
					sort.AddOption(sorter.ByNameWithoutALeadingDotAscend)
				case ".Name":
					sort.AddOption(sorter.ByNameWithoutALeadingDotCaseSensitiveAscend)
				case ".Name-descend":
					sort.AddOption(sorter.ByNameWithoutALeadingDotCaseSensitiveDescend)
				case "size-descend", "S", "sizesort":
					if context.Bool("recursive-size") {
						sort.AddOption(sorter.ByRecursiveSizeDescend(context.Int("depth")))
					} else {
						sort.AddOption(sorter.BySizeDescend)
					}
				case "size":
					if context.Bool("recursive-size") {
						sort.AddOption(sorter.ByRecursiveSizeAscend(context.Int("depth")))
					} else {
						sort.AddOption(sorter.BySizeAscend)
					}
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
					err := limitOnce.Do(
						func() error {
							return setLimit(context)
						},
					)
					if err != nil {
						return err
					}
					sort.AddOption(sorter.ByMimeTypeAscend)
				case "mime-descend", "mimetype-descend", "Mime-descend", "Mimetype-descend":
					err := limitOnce.Do(
						func() error {
							return setLimit(context)
						},
					)
					if err != nil {
						return err
					}
					sort.AddOption(sorter.ByMimeTypeDescend)
				case "inode-descend":
					sort.AddOption(sorter.ByInodeDescend)
				case "inode":
					sort.AddOption(sorter.ByInodeAscend)
				case "version":
					sort.AddOption(sorter.ByVersionAscend)
				case "version-descend":
					sort.AddOption(sorter.ByVersionDescend)
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
		Aliases:            []string{"reverse", "r"},
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
		Name:               "df",
		Aliases:            []string{"dir-first", "group-directories-first"},
		Usage:              "list directories before other files",
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
		Aliases:            []string{"sort-by-size", "sizesort"},
		Usage:              "sort by file size, largest first(descending)",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if context.Bool("srs") { // recursive size
				sort.AddOption(sorter.ByRecursiveSizeDescend(context.Int("depth")))
			} else {
				sort.AddOption(sorter.BySizeDescend)
			}
			return nil
		},
		Category: "SORTING",
	},
	&cli.BoolFlag{
		Name:               "X",
		Aliases:            []string{"sort-by-ext"},
		Usage:              "sort alphabetically by entry extension",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			sort.AddOption(sorter.ByExtensionAscend)
			return nil
		},
		Category: "SORTING",
	},
	&cli.BoolFlag{
		Name:               "width",
		Usage:              "sort by entry name width",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			sort.AddOption(sorter.ByNameWidthAscend)
			return nil
		},
		Category: "SORTING",
	},
	&cli.BoolFlag{
		Name:               "sort-by-mime",
		Usage:              "sort by mimetype",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			err := limitOnce.Do(
				func() error {
					return setLimit(context)
				},
			)
			if err != nil {
				return err
			}

			sort.AddOption(sorter.ByMimeTypeAscend)
			return nil
		},
		Category: "SORTING",
	},
	&cli.BoolFlag{
		Name:               "sort-by-mime-descend",
		Usage:              "sort by mimetype, descending",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				err := limitOnce.Do(
					func() error {
						return setLimit(context)
					},
				)
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
		Name:               "sort-by-mime-parent",
		Usage:              "sort by mimetype parent",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				err := limitOnce.Do(
					func() error {
						return setLimit(context)
					},
				)
				if err != nil {
					return err
				}

				sort.AddOption(sorter.ByMimeTypeParentAscend)
			}
			return nil
		},
		Category: "SORTING",
	},
	&cli.BoolFlag{
		Name:               "sort-by-mime-parent-descend",
		Usage:              "sort by mimetype parent",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				err := limitOnce.Do(
					func() error {
						return setLimit(context)
					},
				)
				if err != nil {
					return err
				}

				sort.AddOption(sorter.ByMimeTypeParentDescend)
			}
			return nil
		},
		Category: "SORTING",
	},
	&cli.BoolFlag{
		Name:    "versionsort",
		Aliases: []string{"sort-by-version"},
		Usage:   "sort by version numbers, ascending",
		Action: func(context *cli.Context, b bool) error {
			sort.AddOption(sorter.ByVersionAscend)
			return nil
		},
		Category: "SORTING",
	},
	&cli.BoolFlag{
		Name:               "U",
		Aliases:            []string{"nosort", "no-sort"},
		Usage:              "do not sort; list entries in directory order. ",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			sort.Reset()
			return nil
		},
		Category: "SORTING",
	},
}
