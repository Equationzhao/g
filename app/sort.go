package app

import (
	"fmt"
	"strings"

	"github.com/Equationzhao/g/filter/content"
	"github.com/Equationzhao/g/sorter"
	"github.com/gabriel-vasile/mimetype"
	"github.com/urfave/cli/v2"
)

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
					sort.AddOption(sorter.ByMimeTypeAscend)
				case "mime-descend", "mimetype-descend", "Mime-descend", "Mimetype-descend":
					err := limitOnce.Do(func() error {
						size := context.String("exact-detect-size")
						var bytes uint64 = 1024 * 1024
						if size == "0" || strings.EqualFold(size, "infinity") || strings.EqualFold(size, "nolimit") {
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
		Aliases:            []string{"df", "group-directories-first"},
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

				sort.AddOption(sorter.ByMimeTypeParentAscend)
			}
			return nil
		},
		Category: "SORTING",
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

				sort.AddOption(sorter.ByMimeTypeParentDescend)
			}
			return nil
		},
		Category: "SORTING",
	},
}
