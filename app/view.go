package app

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Equationzhao/g/display"
	"github.com/Equationzhao/g/filter"
	"github.com/Equationzhao/g/filter/content"
	"github.com/Equationzhao/g/timeparse"
	"github.com/gabriel-vasile/mimetype"
	"github.com/urfave/cli/v2"
)

var viewFlag = []cli.Flag{
	// VIEW
	&cli.BoolFlag{
		Name:    "header",
		Aliases: []string{"title"},
		Usage:   "add a header row",
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
		Name:  "footer",
		Usage: "add a footer row",
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
	&cli.BoolFlag{
		Name:    "duplicate",
		Aliases: []string{"dup"},
		Usage:   "show duplicate files",
		Action: func(context *cli.Context, b bool) error {
			if b {
				noOutputFunc = append(noOutputFunc, duplicateDetect.Enable())
				hookAfter = append(hookAfter, func(p display.Printer, item ...display.Item) {
					duplicateDetect.Fprint(p)
					duplicateDetect.Reset()
				})
			}
			return nil
		},
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
	&cli.BoolFlag{
		Name:               "access",
		Aliases:            []string{"ac", "accessed"},
		Usage:              "accessed time",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				timeType = append(timeType, "access")
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "modify",
		Aliases:            []string{"mod", "modified"},
		Usage:              "modified time",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				timeType = append(timeType, "mod")
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "create",
		Aliases:            []string{"cr", "created"},
		Usage:              "created time",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				timeType = append(timeType, "create")
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
			sizeUint = content.ConvertFromSizeString(s)
			if sizeUint == content.Unknown {
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
				The TIME_STYLE argument can be full-iso, long-iso, iso, locale, or  +FORMAT.
				FORMAT is interpreted like in date(1).
				If FORMAT is FORMAT1<newline>FORMAT2, then FORMAT1 applies to non-recent files and FORMAT2 to recent files.
				TIME_STYLE prefixed with 'posix-' takes effect only outside the POSIX locale.
				Also, the TIME_STYLE environment variable sets the default style to use.
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
		Name:               "#",
		DisableDefaultText: true,
		Usage:              "print entry No. for each entry",
		Category:           "DISPLAY",
		Action: func(context *cli.Context, b bool) error {
			if b {
				contentFunc = append(contentFunc, content.NewIndexEnabler().Enable())
			}
			return nil
		},
	},
	&cli.BoolFlag{
		Name:               "o",
		DisableDefaultText: true,
		Usage:              "like -all/l, but do not list group information",
		Action: func(context *cli.Context, b bool) error {
			if b {
				// remove filter.RemoveHidden
				newFF := make([]*filter.ItemFilterFunc, 0, len(itemFiltetrFunc))
				for _, typeFunc := range itemFiltetrFunc {
					if typeFunc != &filter.RemoveHidden {
						newFF = append(newFF, typeFunc)
					}
				}
				itemFiltetrFunc = newFF
				contentFunc = append(contentFunc, content.EnableFileMode(r), sizeEnabler.EnableSize(sizeUint), contentFilter.EnableGroup(r))
				for _, s := range timeType {
					contentFunc = append(contentFunc, content.EnableTime(timeFormat, s, r))
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
				newFF := make([]*filter.ItemFilterFunc, 0, len(itemFiltetrFunc))
				for _, typeFunc := range itemFiltetrFunc {
					if typeFunc != &filter.RemoveHidden {
						newFF = append(newFF, typeFunc)
					}
				}
				itemFiltetrFunc = newFF
				contentFunc = append(contentFunc, content.EnableFileMode(r), sizeEnabler.EnableSize(sizeUint), contentFilter.EnableOwner(r))
				for _, s := range timeType {
					contentFunc = append(contentFunc, content.EnableTime(timeFormat, s, r))
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
				newFF := make([]*filter.ItemFilterFunc, 0, len(itemFiltetrFunc))
				for _, typeFunc := range itemFiltetrFunc {
					if typeFunc != &filter.RemoveHidden {
						newFF = append(newFF, typeFunc)
					}
				}
				itemFiltetrFunc = newFF
				sizeEnabler.SetEnableTotal()
				contentFunc = append(contentFunc, content.EnableFileMode(r), sizeEnabler.EnableSize(sizeUint), contentFilter.EnableOwner(r))
				if !context.Bool("G") {
					contentFunc = append(contentFunc, contentFilter.EnableGroup(r))
				}
				for _, s := range timeType {
					contentFunc = append(contentFunc, content.EnableTime(timeFormat, s, r))
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
		Name:               "inode",
		Aliases:            []string{"i"},
		Usage:              "show inode[linux/darwin only]",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			i := content.NewInodeEnabler()
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
				rt := content.NewRelativeTimeEnabler()
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
				contentFunc = append(contentFunc, content.EnableFileMode(r))
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
				sizeEnabler.SetRecursive(content.NewSizeRecursive(n))
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "block",
		Aliases:            []string{"blocks"},
		Usage:              "show block size",
		DisableDefaultText: true,
		Action: func(context *cli.Context, b bool) error {
			if b {
				contentFunc = append(contentFunc, blockEnabler.Enable())
				if _, ok := p.(*display.Byline); !ok {
					p = display.NewByline()
				}
			}
			return nil
		},
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:               "lh",
		Aliases:            []string{"human-readable", "hr"},
		DisableDefaultText: true,
		Usage:              "show human readable size",
		Action: func(context *cli.Context, b bool) error {
			if b {
				contentFunc = append(contentFunc, content.EnableFileMode(r), sizeEnabler.EnableSize(sizeUint), contentFilter.EnableOwner(r), contentFilter.EnableGroup(r))
				for _, s := range timeType {
					contentFunc = append(contentFunc, content.EnableTime(timeFormat, s, r))
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
					contentFunc = append(contentFunc, content.EnableTime(timeFormat, s, r))
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
		Name:               "no-icon",
		Usage:              "disable icon(always override show-icon)",
		Aliases:            []string{"noicon", "ni"},
		DisableDefaultText: true,
		Category:           "VIEW",
	},
	&cli.BoolFlag{
		Name:               "show-icon",
		Usage:              "show icon",
		Aliases:            []string{"si", "icons", "icon"},
		DisableDefaultText: true,
		Category:           "VIEW",
	},
	&cli.BoolFlag{
		Name:               "full-path",
		Usage:              "show full path",
		Aliases:            []string{"fp", "fullpath"},
		DisableDefaultText: true,
		Category:           "VIEW",
	},
	&cli.StringFlag{
		Name:        "relative-to",
		Usage:       "show relative path to the given path",
		DefaultText: "current directory",
		Category:    "VIEW",
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
	&cli.BoolFlag{
		Name:               "no-total-size",
		Usage:              "disable total size(always override show-total-size)",
		Aliases:            []string{"nts", "nototal-size"},
		DisableDefaultText: true,
		Category:           "VIEW",
		Action: func(context *cli.Context, b bool) error {
			if b {
				sizeEnabler.DisableTotal()
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
				exact := content.NewMimeFileTypeEnabler()

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
				contentFunc = append(contentFunc, exact.Enable())
				wgs = append(wgs, exact)
			}
			return nil
		},
	},
	&cli.BoolFlag{
		Name:     "mime-parent",
		Usage:    "show mime parent type",
		Aliases:  []string{"mime-p", "mime-parent-type", "mime-type-parent"},
		Category: "VIEW",
		Action: func(context *cli.Context, b bool) error {
			if b {
				exact := content.NewMimeFileTypeEnabler()
				exact.ParentOnly = true

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
		Name:     "quote-name",
		Aliases:  []string{"Q"},
		Usage:    "enclose entry names in double quotes(overridden by --literal)",
		Category: "VIEW",
	},
	// &cli.StringFlag{
	// 	Name:    "quoting-style",
	// 	Aliases: []string{"Qs"},
	// 	Usage:   "use quoting style: literal, shell, shell-always, c, escape, locale, clocale",
	// },
	&cli.BoolFlag{
		Name:     "literal",
		Aliases:  []string{"N"},
		Usage:    "print entry names without quoting",
		Category: "VIEW",
	},
	&cli.BoolFlag{
		Name:    "link",
		Aliases: []string{"H"},
		Usage:   "list each file's number of hard links",
		Action: func(context *cli.Context, b bool) error {
			if b {
				link := content.NewLinkEnabler()
				contentFunc = append(contentFunc, link.Enable())
				wgs = append(wgs, link)
			}
			return nil
		},
		Category: "VIEW",
	},
}
