# NAME

g - a powerful ls

# SYNOPSIS

g

```
[-#]
[--CSV|--csv]
[--access|--ac|--accessed]
[--all|--la]
[--block|--blocks]
[--byline|-1|--oneline|--single-column]
[--charset]
[--check-new-version]
[--checksum-algorithm|--ca]=[value]
[--checksum|--cs]
[--classic]
[--color]=[value]
[--colorless|--no-color|--nocolor]
[--create|--cr|--created]
[--depth|--level]=[value]
[--dereference]
[--detect-size]=[value]
[--df|--dir-first|--group-directories-first]
[--disable-index|--di|--no-update]
[--duplicate|--dup]
[--ext]=[value]
[--footer]
[--format]=[value]
[--fp|--full-path|--fullpath]
[--ft|--file-type]
[--full-time]
[--fuzzy|--fz|-f]
[--gid]
[--git-ignore|--hide-git-ignore]
[--git|--git-status]
[--group]
[--header|--title]
[--html|--HTML]
[--hyperlink]=[value]
[--icon|--icons]
[--init]=[value]
[--inode|-i]
[--lh|--human-readable]
[--list-index|--li]
[--md|--markdown|--Markdown]
[--mime-parent|--mime-parent-type|--mimetype-parent]
[--mime|--mime-type|--mimetype]
[--modify|--mod|--modified]
[--no-config]
[--no-dereference]
[--no-dir|--nodir|--file]
[--no-ext|--noext]=[value]
[--no-icon|--noicon|--ni]
[--no-path-transform|--np]
[--no-total-size]
[--numeric|--numeric-uid-gid]
[--octal-perm|--octal-permission]
[--only-mime]=[value]
[--owner|--author]
[--perm|--permission]
[--rebuild-index|--ri|--remove-all]
[--recursive-size]
[--relative-to]=[value]
[--remove-current-path|--rcp|--rc|--rmc]
[--remove-index|--rm]=[value]
[--remove-invalid-path|--rip]
[--rt|--relative-time]
[--show-only-hidden|--hidden]
[--size-unit|--su|--block-size]=[value]
[--size]
[--sort-by-mime-descend]
[--sort-by-mime-parent-descend]
[--sort-by-mime-parent]
[--sort-by-mime]
[--sort-reverse|--reverse|-r]
[--sort|--SORT_FIELD]=[value]
[--statistic]
[--tb-style|--table-style]=[value]
[--tb|--table]
[--theme|--th]=[value]
[--time-style]=[value]
[--time-type]=[value]
[--time]
[--total-size]
[--uid]
[--versionsort|--sort-by-version]
[--width]
[--zero|-0]
[-A|--almost-all]
[-B|--ignore-backups]
[-C|--vertical]
[-D|--dir|--only-dir]
[-F|--classify]
[-G|--no-group]
[-H|--link]
[-I|--ignore]=[value]
[-M|--match]=[value]
[-N|--literal]
[-O|--no-owner]
[-Q|--quote-name]
[-R|--recurse]
[-S|--sort-by-size|--sizesort]
[-T|--tree]
[-U|--nosort|--no-sort]
[-X|--sort-by-ext]
[-a|--sh|--show-hidden]
[-d|--directory|--list-dirs]
[-g]
[-j|--json]
[-l|--long]
[-m|--comma]
[-n|--limitN|--limit|--topN|--top]=[value]
[-o]
[-x|--col|--across|--horizontal]
```

**Usage**:

```
g [options] [path]
```

# GLOBAL OPTIONS

**-#**: print entry Number for each entry

**--CSV, --csv**: output in csv format

**--access, --ac, --accessed**: accessed time

**--all, --la**: show all info/use a long listing format

**--block, --blocks**: show block size

**--byline, -1, --oneline, --single-column**: print by line

**--charset**: show charset of text file in mime type field

**--check-new-version**: check if there's new release

**--checksum, --cs**: show checksum of file with algorithm, see --checksum-algorithm

**--checksum-algorithm, --ca**="": show checksum of file with algorithm: 
	md5, sha1, sha224, sha256, sha384, sha512, crc32 (default: sha1)

**--classic**: Enable classic mode (no colors or icons)

**--color**="": when to use terminal colors [always|auto|never][basic|256|24bit] (default: auto)

**--colorless, --no-color, --nocolor**: without color

**--create, --cr, --created**: created time

**--depth, --level**="": limit recursive/tree depth, negative -> infinity (default: infinity)

**--dereference**: dereference symbolic links

**--detect-size**="": set exact size for mimetype detection 
			eg:1M/nolimit/infinity (default: 1M)

**--df, --dir-first, --group-directories-first**: List directories before other files

**--disable-index, --di, --no-update**: disable updating index

**--duplicate, --dup**: show duplicate files

**--ext**="": show file which has target ext, eg: --show-only-ext=go,java

**--footer**: add a footer row

**--format**="": across  -x,  commas  -m, horizontal -x, long -l, single-column -1,
	verbose -l, vertical -C, table -tb, HTML -html, Markdown -md, CSV -csv, json -j, tree -T (default: C)

**--fp, --full-path, --fullpath**: show full path

**--ft, --file-type**: likewise, except do not append '*'

**--full-time**: like -all/l --time-style=full-iso

**--fuzzy, --fz, -f**: fuzzy search

**--gid**: show gid instead of groupname [sid in windows]

**--git, --git-status**: show git status [if git is installed]

**--git-ignore, --hide-git-ignore**: hide git ignored file/dir [if git is installed]

**--group**: show group

**--header, --title**: add a header row

**--html, --HTML**: output in HTML-table format

**--hyperlink**="": Attach hyperlink to filenames [auto|always|never] (default: auto)

**--icon, --icons**: show icon

**--init**="": show the init script for shell, support zsh, bash, fish, powershell, nushell

**--inode, -i**: show inode[linux/darwin only]

**--lh, --human-readable**: show human readable size

**--list-index, --li**: list index

**--md, --markdown, --Markdown**: output in markdown-table format

**--mime, --mime-type, --mimetype**: show mime file type

**--mime-parent, --mime-parent-type, --mimetype-parent**: show mime parent type

**--modify, --mod, --modified**: modified time

**--no-config**: do not load config file

**--no-dereference**: do not follow symbolic links

**--no-dir, --nodir, --file**: do not show directory

**--no-ext, --noext**="": show file which doesn't have target ext

**--no-icon, --noicon, --ni**: disable icon(always override --icon)

**--no-path-transform, --np**: By default, .../a/b/c will be transformed to ../../a/b/c, and ~ will be replaced by homedir, 
	using this flag to disable this feature

**--no-total-size**: disable total size(always override show-total-size)

**--numeric, --numeric-uid-gid**:  List numeric user and group IDs instead of name [sid in windows]

**--octal-perm, --octal-permission**: list each file's permission in octal format

**--only-mime**="": only show file with given mime type

**--owner, --author**: show owner

**--perm, --permission**: show permission

**--rebuild-index, --ri, --remove-all**: rebuild index

**--recursive-size**: show recursive size of dir, only work with --size

**--relative-to**="": show relative path to the given path (default: current directory)

**--remove-current-path, --rcp, --rc, --rmc**: remove current path from index

**--remove-index, --rm**="": remove paths from index

**--remove-invalid-path, --rip**: remove invalid paths from index

**--rt, --relative-time**: show relative time

**--show-only-hidden, --hidden**: show only hidden files(overridden by --show-hidden/-a/-A)

**--size**: show file/dir size

**--size-unit, --su, --block-size**="": size unit:
			bit, b, k, m, g, t, p,
			e, z, y, bb, nb, auto

**--sort, --SORT_FIELD**="": sort by field, default: 
	ascending and case insensitive, 
	field beginning with Uppercase is case sensitive,	
	available fields: 	
	nature(default),none(nosort),
	   name,.name(sorts by name without a leading dot),	
	   size,time,owner,group,extension,inode,width,mime. 	
	   following '-descend' to sort descending

**--sort-by-mime**: sort by mimetype

**--sort-by-mime-descend**: sort by mimetype, descending

**--sort-by-mime-parent**: sort by mimetype parent

**--sort-by-mime-parent-descend**: sort by mimetype parent

**--sort-reverse, --reverse, -r**: reverse the order of the sort

**--statistic**: show statistic info

**--tb, --table**: output in table format

**--tb-style, --table-style**="": set table style [ascii(default)/unicode]

**--theme, --th**="": apply theme `path/to/theme`

**--time**: show time

**--time-style**="": time/date format with -l, 
	Valid timestamp styles are default, iso, long-iso, full-iso, locale, 
	custom +FORMAT like date(1). 
	(default: +%d.%b'%y %H:%M ,like 02.Jan'06 15:04)

**--time-type**="": time type, mod(default), create, access, all

**--total-size**: show total size

**--uid**: show uid instead of username [sid in windows]

**--versionsort, --sort-by-version**: sort by version numbers, ascending

**--width**: sort by entry name width

**--zero, -0**: end each output line with NUL, not newline

**-A, --almost-all**: do not list implied . and ..

**-B, --ignore-backups**: do not list implied entries ending with ~

**-C, --vertical**: list entries by columns (default)

**-D, --dir, --only-dir**: show directory only

**-F, --classify**: append indicator (one of */=>@|) to entries

**-G, --no-group**: in a long listing, don't print group names

**-H, --link**: list each file's number of hard links

**-I, --ignore**="": ignore Glob patterns

**-M, --match**="": match Glob patterns

**-N, --literal**: print entry names without quoting

**-O, --no-owner**: in a long listing, don't print owner names

**-Q, --quote-name**: enclose entry names in double quotes(overridden by --literal)

**-R, --recurse**: recurse into directories

**-S, --sort-by-size, --sizesort**: sort by file size, largest first(descending)

**-T, --tree**: recursively list in tree

**-U, --nosort, --no-sort**: do not sort; list entries in directory order. 

**-X, --sort-by-ext**: sort alphabetically by entry extension

**-a, --sh, --show-hidden**: show hidden files

**-d, --directory, --list-dirs**: list directories themselves, not their contents

**-g**: like -all, but do not list owner

**-j, --json**: output in json format

**-l, --long**: use a long listing format

**-m, --comma**: fill width with a comma separated list of entries

**-n, --limitN, --limit, --topN, --top**="": Limit display to a max of n items (n <=0 means unlimited) (default: unlimited)

**-o**: like -all, but do not list group information

**-x, --col, --across, --horizontal**: list entries by lines instead of by columns


