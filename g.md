# NAME

g - a powerful ls

# SYNOPSIS

g

```
[-#]
[--CSV|--csv]
[--HTML|--html]
[--Markdown|--md|--MD|--markdown]
[--access|--ac|--accessed]
[--all|--la|-l|--long]
[--block|--blocks]
[--byline|--bl|-1|--oneline|--single-column]
[--check-new-version]
[--checksum-algorithm|--ca]=[value]
[--checksum|--cs]
[--classic]
[--color]=[value]
[--colorless|--nc|--no-color|--nocolor]
[--create|--cr|--created]
[--depth|--level|-L]=[value]
[--dereference]
[--dir-first|--df|--group-directories-first]
[--disable-index|--di|--no-update]
[--duplicate|--dup]
[--exact-detect-size|--eds|--detect-size|--ds]=[value]
[--file-type|--ft]
[--footer]
[--format]=[value]
[--full-path|--fp|--fullpath]
[--full-time]
[--fuzzy|--fz|-f]
[--gid]
[--git-status|--gs|--git]
[--header|--title]
[--hide-git-ignore|--gi|--hgi|--git-ignore]
[--hyperlink]=[value]
[--ignore-glob|-I|--ignore|--ig]=[value]
[--init]=[value]
[--inode|-i]
[--json|-j]
[--lh|--human-readable|--hr]
[--limitN|-n|--limit|--topN|--top]=[value]
[--link|-H]
[--list-index|--li]
[--literal|-N]
[--match-glob|-M|--glob|--match]=[value]
[--mime-charset|--charset]
[--mime-parent|--mime-p|--mime-parent-type|--mime-type-parent]
[--mime-type|--mime|--mimetype]
[--modify|--mod|--modified]
[--no-dereference]
[--no-icon|--noicon|--ni]
[--no-path-transform|--np]
[--no-total-size|--nts|--nototal-size]
[--numeric|--numeric-uid-gid]
[--quote-name|-Q]
[--rebuild-index|--ri|--remove-all]
[--recurse|-R]
[--relative-time|--rt]
[--relative-to]=[value]
[--remove-current-path|--rcp|--rc|--rmc]
[--remove-index|--rm]=[value]
[--remove-invalid-path|--rip]
[--show-group|--sg|--group]
[--show-hidden|--sh|-a]
[--show-icon|--si|--icons|--icon]
[--show-mime-file-type-only|--mime-only]=[value]
[--show-no-dir|--nd|--nodir|--no-dir|--file]
[--show-no-ext|--sne|--noext]=[value]
[--show-octal-perm|--octal-perm|--octal-permission|--octal-permissions]
[--show-only-dir|--sd|--dir|--only-dir|-D]
[--show-only-ext|--se|--ext]=[value]
[--show-only-hidden|--soh|--hidden]
[--show-owner|--so|--author|--owner]
[--show-perm|--sp|--permission|--perm]
[--show-recursive-size|--srs|--recursive-size]
[--show-size|--ss|--size]
[--show-time|--st|--time]
[--show-total-size|--ts|--total-size]
[--size-unit|--su|--block-size]=[value]
[--sort-by-mimetype-descend|--mimetypesort-descend|--Mimetypesort-descend]
[--sort-by-mimetype-parent-descend|--mimetypesort-parent-descend|--Mimetypesort-parent-descend|--sort-by-mime-parent-descend]
[--sort-by-mimetype-parent|--mimetypesort-parent|--Mimetypesort-parent|--sort-by-mime-parent]
[--sort-by-mimetype|--mimetypesort|--Mimetypesort|--sort-by-mime]
[--sort-reverse|--sr|--reverse|-r]
[--sort|--SORT_FIELD]=[value]
[--statistic]
[--table-style|--tablestyle|--tb-style]=[value]
[--table|--tb]
[--theme|--th]=[value]
[--time-style]=[value]
[--time-type|--tt]=[value]
[--uid]
[--width]
[--zero|-0]
[-A|--almost-all]
[-B|--ignore-backups]
[-C|--vertical]
[-F|--classify]
[-G|--no-group]
[-S|--sort-size|--sort-by-size|--sizesort]
[-U|--nosort|--no-sort]
[-X|--extensionsort|--Extentionsort]
[-d|--directory|--list-dirs]
[-g]
[-m|--comma]
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

**--HTML, --html**: output in HTML-table format

**--Markdown, --md, --MD, --markdown**: output in markdown-table format

**--access, --ac, --accessed**: accessed time

**--all, --la, -l, --long**: show all info/use a long listing format

**--block, --blocks**: show block size

**--byline, --bl, -1, --oneline, --single-column**: print by line

**--check-new-version**: check if there's new release

**--checksum, --cs**: show checksum of file with algorithm: md5, sha1(default), sha224, sha256, sha384, sha512, crc32

**--checksum-algorithm, --ca**="": show checksum of file with algorithm: md5, sha1, sha224, sha256, sha384, sha512, crc32 (default: "sha1")

**--classic**: Enable classic mode (no colours or icons)

**--color**="": when to use terminal colours[always|auto|never][basic|256|24bit] (default: auto)

**--colorless, --nc, --no-color, --nocolor**: without color

**--create, --cr, --created**: created time

**--depth, --level, -L**="": limit recursive depth, negative -> infinity (default: infinity)

**--dereference**: dereference symbolic links

**--dir-first, --df, --group-directories-first**: List directories before other files

**--disable-index, --di, --no-update**: disable updating index

**--duplicate, --dup**: show duplicate files

**--exact-detect-size, --eds, --detect-size, --ds**="": set exact size for mimetype detection eg:1M/nolimit/infinity (default: 1M)

**--file-type, --ft**: likewise, except do not append '*'

**--footer**: add a footer row

**--format**="": across  -x,  commas  -m, horizontal -x, long -l, single-column -1, verbose -l, vertical -C, table -tb, HTML -html, Markdown -md, CSV -csv, json -j (default: C)

**--full-path, --fp, --fullpath**: show full path

**--full-time**: like -all/l --time-style=full-iso

**--fuzzy, --fz, -f**: fuzzy search

**--gid**: show gid instead of groupname [sid in windows]

**--git-status, --gs, --git**: show git status [if git is installed]

**--header, --title**: add a header row

**--hide-git-ignore, --gi, --hgi, --git-ignore**: hide git ignored file/dir [if git is installed]

**--hyperlink**="": Attach hyperlink to filenames [auto|always|never] (default: auto)

**--ignore-glob, -I, --ignore, --ig**="": ignore Glob patterns

**--init**="": init the config file, default path is ~/.config/g/config.yaml

**--inode, -i**: show inode[linux/darwin only]

**--json, -j**: output in json format

**--lh, --human-readable, --hr**: show human readable size

**--limitN, -n, --limit, --topN, --top**="": limit n items(n <=0 means unlimited) (default: unlimited)

**--link, -H**: list each file's number of hard links

**--list-index, --li**: list index

**--literal, -N**: print entry names without quoting

**--match-glob, -M, --glob, --match**="": match Glob patterns

**--mime-charset, --charset**: show charset of text file

**--mime-parent, --mime-p, --mime-parent-type, --mime-type-parent**: show mime parent type

**--mime-type, --mime, --mimetype**: show mime file type

**--modify, --mod, --modified**: modified time

**--no-dereference**: do not follow symbolic links

**--no-icon, --noicon, --ni**: disable icon(always override show-icon)

**--no-path-transform, --np**: By default, .../a/b/c will be transformed to ../../a/b/c, and ~ will be replaced by homedir, using this flag to disable this feature

**--no-total-size, --nts, --nototal-size**: disable total size(always override show-total-size)

**--numeric, --numeric-uid-gid**:  List numeric user and group IDs instead of name [sid in windows]

**--quote-name, -Q**: enclose entry names in double quotes(overridden by --literal)

**--rebuild-index, --ri, --remove-all**: rebuild index

**--recurse, -R**: recurse into directories

**--relative-time, --rt**: show relative time

**--relative-to**="": show relative path to the given path (default: current directory)

**--remove-current-path, --rcp, --rc, --rmc**: remove current path from index

**--remove-index, --rm**="": remove paths from index

**--remove-invalid-path, --rip**: remove invalid paths from index

**--show-group, --sg, --group**: show group

**--show-hidden, --sh, -a**: show hidden files

**--show-icon, --si, --icons, --icon**: show icon

**--show-mime-file-type-only, --mime-only**="": only show file with given mime type

**--show-no-dir, --nd, --nodir, --no-dir, --file**: do not show directory

**--show-no-ext, --sne, --noext**="": show file which doesn't have target ext

**--show-octal-perm, --octal-perm, --octal-permission, --octal-permissions**: list each file's permission in octal format

**--show-only-dir, --sd, --dir, --only-dir, -D**: show directory only

**--show-only-ext, --se, --ext**="": show file which has target ext, eg: --show-only-ext=go,java

**--show-only-hidden, --soh, --hidden**: show only hidden files(overridden by --show-hidden/-sh/-a/-A)

**--show-owner, --so, --author, --owner**: show owner

**--show-perm, --sp, --permission, --perm**: show permission

**--show-recursive-size, --srs, --recursive-size**: show recursive size of dir, only work with --show-size

**--show-size, --ss, --size**: show file/dir size

**--show-time, --st, --time**: show time

**--show-total-size, --ts, --total-size**: show total size

**--size-unit, --su, --block-size**="": size unit, b, k, m, g, t, p, e, z, y, bb, nb, auto (default: auto)

**--sort, --SORT_FIELD**="": sort by field, default: ascending and case insensitive, field beginning with Uppercase is case sensitive, available fields: nature(default),none(nosort),name,.name(sorts by name without a leading dot),size,time,owner,group,extension,inode,width,mime. following '-descend' to sort descending

**--sort-by-mimetype, --mimetypesort, --Mimetypesort, --sort-by-mime**: sort by mimetype

**--sort-by-mimetype-descend, --mimetypesort-descend, --Mimetypesort-descend**: sort by mimetype, descending

**--sort-by-mimetype-parent, --mimetypesort-parent, --Mimetypesort-parent, --sort-by-mime-parent**: sort by mimetype parent

**--sort-by-mimetype-parent-descend, --mimetypesort-parent-descend, --Mimetypesort-parent-descend, --sort-by-mime-parent-descend**: sort by mimetype parent

**--sort-reverse, --sr, --reverse, -r**: reverse the order of the sort

**--statistic**: show statistic info

**--table, --tb**: output in table format

**--table-style, --tablestyle, --tb-style**="": set table style (ascii(default)/unicode)

**--theme, --th**="": apply theme `path/to/theme`

**--time-style**="": time/date format with -l, Valid timestamp styles are default, iso, long iso, full-iso, locale, custom +FORMAT like date(1). (default: +%d.%b'%y %H:%M (like 02.Jan'06 15:04))

**--time-type, --tt**="": time type, mod(default), create, access, all (default: mod)

**--uid**: show uid instead of username [sid in windows]

**--width**: sort by entry name width

**--zero, -0**: end each output line with NUL, not newline

**-A, --almost-all**: do not list implied . and ..

**-B, --ignore-backups**: do not list implied entries ending with ~

**-C, --vertical**: list entries by columns (default)

**-F, --classify**: append indicator (one of */=>@|) to entries

**-G, --no-group**: in a long listing, don't print group names

**-S, --sort-size, --sort-by-size, --sizesort**: sort by file size, largest first(descending)

**-U, --nosort, --no-sort**: do not sort; list entries in directory order. 

**-X, --extensionsort, --Extentionsort**: sort alphabetically by entry extension

**-d, --directory, --list-dirs**: list directories themselves, not their contents

**-g**: like -all/l, but do not list owner

**-m, --comma**: fill width with a comma separated list of entries

**-o**: like -all/l, but do not list group information

**-x, --col, --across, --horizontal**: list entries by lines instead of by columns


