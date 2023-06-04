# NAME

g - a powerful ls

# SYNOPSIS

g

```
[--all|--la|-l]
[--byline|--bl|-1|--oneline|--single-column]
[--check-new-version]
[--checksum-algorithm|--ca]=[value]
[--checksum|--cs]
[--depth]=[value]
[--dir-first|--df]
[--file-type|--ft]
[--format]=[value]
[--full-time]
[--fuzzy|--fz|-f]
[--gid]
[--git-status-style|--gss]=[value]
[--git-status|--gs]
[--hide-git-ignore|--gi|--hgi]
[--ignore-glob|-I]=[value]
[--inode|-i]
[--lh|--human-readable]
[--list-index|--li]
[--match-glob|-M]=[value]
[--no-path-transform|--np]
[--numeric]
[--rebuild-index|--ri]
[--recurse|-R]
[--relative-time|--rt]
[--remove-index]=[value]
[--show-group|--sg]
[--show-hidden|--sh|-a]
[--show-icon|--si|--icons]
[--show-no-dir|--nd|--nodir|--no-dir]
[--show-no-ext|--sne|--noext]=[value]
[--show-only-dir|--sd|--dir|--only-dir]
[--show-only-ext|--se|--ext]=[value]
[--show-owner|--so|--author]
[--show-perm|--sp]
[--show-size|--ss]
[--show-time|--st]
[--show-total-size|--ts]
[--size-unit|--su]=[value]
[--sort-reverse|--sr]
[--sort|--SORT_FIELD]=[value]
[--theme|--th]=[value]
[--time-style]=[value]
[--time-type|--tt]=[value]
[--tree|-t]
[--uid]
[--zero|-0]
[-A|--almost-all]
[-B|--ignore-backups]
[-C|--vertical]
[-F|--classify]
[-G|--no-group]
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

**--all, --la, -l**: show all info/use a long listing format

**--byline, --bl, -1, --oneline, --single-column**: print by line

**--check-new-version**: check if there's new release

**--checksum, --cs**: show checksum of file with algorithm: md5, sha1, sha256, sha512

**--checksum-algorithm, --ca**="": show checksum of file with algorithm: md5, sha1, sha256, sha512 (default: "sha1")

**--depth**="": limit recursive depth, negative -> infinity (default: infinity)

**--dir-first, --df**: List directories before other files

**--file-type, --ft**: likewise, except do not append '*'

**--format**="": across  -x,  commas  -m, horizontal -x, long -l, single-column -1, verbose -l, vertical -C (default: C)

**--full-time**: like -all/l --time-style=full-iso

**--fuzzy, --fz, -f**: fuzzy search

**--gid**: show gid instead of groupname [sid in windows]

**--git-status, --gs**: show git status: ? untracked, + added, ! deleted, ~ modified, | renamed, = copied, $ ignored [if git is installed]

**--git-status-style, --gss**="": git status style: colored-symbol: {? untracked, + added, - deleted, ~ modified, | renamed, = copied, ! ignored} colored-dot

**--hide-git-ignore, --gi, --hgi**: hide git ignored file/dir [if git is installed]

**--ignore-glob, -I**="": ignore Glob patterns

**--inode, -i**: show inode[linux only]

**--lh, --human-readable**: show human readable size

**--list-index, --li**: list index

**--match-glob, -M**="": match Glob patterns

**--no-path-transform, --np**: By default, .../a/b/c will be transformed to ../../a/b/c, and ~ will be replaced by homedir, using this flag to disable this feature

**--numeric**:  List numeric user and group IDs instead of name [sid in windows]

**--rebuild-index, --ri**: rebuild index

**--recurse, -R**: recurse into directories

**--relative-time, --rt**: show relative time

**--remove-index**="": remove paths from index

**--show-group, --sg**: show group

**--show-hidden, --sh, -a**: show hidden files

**--show-icon, --si, --icons**: show icon

**--show-no-dir, --nd, --nodir, --no-dir**: do not show directory

**--show-no-ext, --sne, --noext**="": show file which doesn't have target ext

**--show-only-dir, --sd, --dir, --only-dir**: show directory only

**--show-only-ext, --se, --ext**="": show file which has target ext, eg: --show-only-ext=go,java

**--show-owner, --so, --author**: show owner

**--show-perm, --sp**: show permission

**--show-size, --ss**: show file/dir size

**--show-time, --st**: show time

**--show-total-size, --ts**: show total size

**--size-unit, --su**="": size unit, b, k, m, g, t, p, e, z, y, auto (default: auto)

**--sort, --SORT_FIELD**="": sort by field, default: ascending and case insensitive, field beginning with Uppercase is case sensitive, available fields: name,size,time,owner,group,extension. following `-descend` to sort descending

**--sort-reverse, --sr**: reverse the order of the sort

**--theme, --th**="": apply theme `path/to/theme`

**--time-style**="": time/date format with -l, Valid timestamp styles are `default', `iso`, `long iso`, `full-iso`, `locale`, custom `+FORMAT` like date(1). (default: +%d.%b'%y %H:%M (like 02.Jan'06 15:04))

**--time-type, --tt**="": time type, mod, create, access (default: mod)

**--tree, -t**: recursively list in tree

**--uid**: show uid instead of username [sid in windows]

**--zero, -0**: end each output line with NUL, not newline

**-A, --almost-all**: do not list implied . and ..

**-B, --ignore-backups**: do not list implied entries ending with ~

**-C, --vertical**: list entries by columns (default)

**-F, --classify**: append indicator (one of */=>@|) to entries

**-G, --no-group**: in a long listing, don't print group names

**-d, --directory, --list-dirs**: list directories themselves, not their contents

**-g**: like -all/l, but do not list owner

**-m, --comma**: fill width with a comma separated list of entries

**-o**: like -all/l, but do not list group information

**-x, --col, --across, --horizontal**: list entries by lines instead of by columns


