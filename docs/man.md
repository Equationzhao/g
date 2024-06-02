## USAGE:
g [options] [files...]

## VERSION:
0.28.2

## GLOBAL OPTIONS
--bug                      report bug

--duplicate, --dup         show duplicate files

--no-config                do not load config file

--no-path-transform, --np  By default, .../a/b/c will be transformed to ../../a/b/c, and ~ will be replaced by homedir,
                           using this flag to disable this feature

--si                       use powers of 1000 for size format(default: false)

## META
--check-new-version  check if there's new release

--help, -h, -?       show help

--version, -v        print the version

## DISPLAY
-#                                 print entry Number for each entry

--CSV, --csv                       output in csv format

--TSV, --tsv                       output in tsv format

--byline, -1, --oneline            print by line

--classic                          enable classic mode(no colors or icons)

--color WHEN/LEVEL                 set terminal colors [always|auto|never][basic|256|24bit](default: auto)

--colorless, --no-color            without color

--depth NUM                        limit recursive/tree depth, negative -> infinity(default: infinity)

--format FORMAT                    across  -x,  commas  -m, horizontal -x, long -l, single-column -1,
								   verbose -l, vertical -C, table -tb, markdown -md, csv -csv, tsv -tsv, json -j, tree -T(default: C)


--file-type                        like --classify, except do not append '*'

--md, --markdown                   output in markdown-table format

--tb, --table                      output in table format

--table-style STYLE                set table style [ascii(default)/unicode]

--term-width COLS                  set screen width (default: auto)

--theme path/to/theme              apply theme path/to/theme

--tree-style STYLE                 set tree style [ascii/unicode(default)/rectangle]

--zero, -0                         end each output line with NUL, not newline

-C, --vertical                     list entries by columns(default)

-F, --classify                     append indicator (one of */=@|) to entries

-R, --recurse                      recurse into directories

-T, --tree                         recursively list in tree

-d, --directory,                   list directories themselves, not their contents

-j, --json                         output in json format

-m, --comma                        fill width with a comma separated list of entries

-x, --col, --across, --horizontal  list entries by lines instead of by columns

## FILTERING
--after TIME                  show items which was modified/access/created after given time, see --before

--before TIME                 show items which was modified/access/created before given time, the time field is determined by --time-type,
								the time will be parsed using format:
								MM-dd, MM-dd HH:mm, HH:mm, YYYY-MM-dd, YYYY-MM-dd HH:mm, and the format set by --time-style


--ext value                   show file which has target ext, eg: --ext=go,java

--git-ignore                  hide git ignored file/dir [if git is installed]

--no-dir, --file              do not show directory

--no-ext value                show file which doesn't have target ext

--only-mime value             only show file with given mime type

--show-only-hidden, --hidden  show only hidden files(overridden by --show-hidden/-a/-A)

-A, --almost-all              do not list implied . and ..

-B, --ignore-backups          do not list implied entries ending with ~

-D, --dir, --only-dir         show directory only

-I GLOBS, --ignore GLOBS      ignore Glob patterns

-M GLOBS, --match GLOBS       match Glob patterns

-a, --show-hidden             show hidden files

-n NUM, --limit NUM           limit display to a max of n items(n <=0 means unlimited)(default: unlimited)

## INDEX
--disable-index, --no-update      disable updating index

--fuzzy, -f                       fuzzy search

--list-index,                     list index

--rebuild-index                   rebuild index

--remove-current-path             remove current path from index

--remove-index value, --rm value  remove paths from index

--remove-invalid-path, --rip      remove invalid paths from index


## SHELL
--init value  show the init script for shell, support zsh, bash, fish, powershell, nushell

## SORTING
--sort SORT_FIELD                       sort by field, default: ascending and case-insensitive,

available fields:                       nature(default),none(nosort),
										name,.name(sorts by name without a leading dot),
										size,time,owner,group,extension,inode,width,mime.
										append '-descend' to sort descending
										field beginning with an Uppercase letter is case-sensitive


--dir-first, --group-directories-first  list directories before other files

--sort-by-mime                          sort by mimetype

--sort-by-mime-descend                  sort by mimetype, descending

--sort-by-mime-parent                   sort by mimetype parent

--sort-by-mime-parent-descend           sort by mimetype parent

--sort-reverse, --reverse, -r           reverse the order of the sort

--versionsort, --sort-by-version        sort by version numbers, ascending(default: false)

--width                                 sort by entry name width

-S, --sort-by-size, --sizesort          sort by file size, largest first(descending)

-U, --nosort, --no-sort                 do not sort; list entries in directory order.

-X, --sort-by-ext                       sort alphabetically by entry extension


## VIEW
--access, --ac, --accessed              accessed time

--all                                   show all info/use a long listing format

--birth                                 birth time[macOS only]

--block, --blocks                       show block size

--charset                               show charset of text file in mime type field

--checksum, --cs                        show checksum of file with algorithm, see --checksum-algorithm

--checksum-algorithm value, --ca value  show checksum of file with algorithm:
											md5, sha1, sha224, sha256, sha384, sha512, crc32(default: sha1)


--create, --cr, --created               created time

--dereference                           dereference symbolic links

--detect-size value                     set exact size for mimetype detection
											eg:1M/nolimit/infinity(default: 1M)


--footer                                add a footer row

--fp, --full-path, --fullpath           show full path

--full-time                             like -all/l --time-style=full-iso

--gid                                   show gid instead of groupname [sid in windows]

--git, --git-status                     show git status [if git is installed]

--git-repo-branch, --branch             list root of git-tree branch [if git is installed]

--git-repo-status, --repo-status        list root of git-tree status [if git is installed]

--group                                 show group

--header, --title                       add a header row

--hyperlink value                       attach hyperlink to filenames [auto|always|never](default: auto)

--icon, --icons                         show icon

--inode, -i                             show inode[linux/darwin only]

--lh, --human-readable                  show human readable size

--mime, --mime-type, --mimetype         show mime file type

--mime-parent, --mime-parent-type       show mime parent type

--modify, --mod, --modified             modified time

--mounts                                show mount details

--no-dereference                        do not follow symbolic links

--no-icon, --noicon, --ni               disable icon(always override --icon)

--no-total-size                         disable total size(always override --total-size)

--numeric, --numeric-uid-gid            list numeric user and group IDs instead of name [sid in windows]

--octal-perm, --octal-permission        list each file's permission in octal format

--owner, --author                       show owner

--perm, --permission                    show permission

--recursive-size                        show recursive size of dir, only work with --size

--relative-to value                     show relative path to the given path (default: current directory)

--rt, --relative-time                   show relative time

--size                                  show file/dir size

--size-unit value, --block-size value   size unit: bit, b, k, m, g, t, auto

--smart-group                           only show group if it has a different name from owner

--statistic                             show statistic info

--stdin                                 read path from stdin, split by newline

--time                                  show time

--time-style TIME_TYPE                  time/date format with -l,

valid TIME_TYPE are :
										default, iso, long-iso, full-iso, locale,
										and custom +FORMAT like date(1).
										(default: +%%d.%%b'%%y %%H:%%M ,like 02.Jan'06 15:04)


--time-type value                       time type, mod(default), create, access, all, birth[macOS only]

--total-size                            show total size

--uid                                   show uid instead of username [sid in windows]

-G, --no-group                          in a long listing, don't print group names

-H, --link                              list each file's number of hard links

-N, --literal                           print entry names without quoting

-O, --no-owner                          in a long listing, don't print owner names

-Q, --quote-name                        enclose entry names in double quotes(overridden by --literal)

-g                                      like -all, but do not list owner

-l, --long                              use a long listing format

-o                                      like -all, but do not list group information

