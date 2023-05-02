# NAME

g - a powerful ls

# SYNOPSIS

g

```
[--all|--la|-l]
[--byline|--bl|-1|--oneline|--single-column]
[--check-new-version]
[--depth]=[value]
[--file-type|--ft]
[--format]=[value]
[--full-time]
[--git-status-style|--gss]=[value]
[--git-status|--gs]
[--hide-git-ignore|--gi|--hgi]
[--lh|--human-readable]
[--show-group|--sg]
[--show-hidden|--sh|-a]
[--show-icon|--si|--icons]
[--show-no-dir|--nd|--nodir]
[--show-no-ext|--sne|--noext]=[value]
[--show-only-dir|--sd|--dir]
[--show-only-ext|--se|--ext]=[value]
[--show-owner|--so|--author]
[--show-perm|--sp]
[--show-size|--ss]
[--show-time|--st]
[--show-total-size|--ts]
[--theme|--th]=[value]
[--time-style]=[value]
[--tree|-t]
[--zero|-0]
[-A|--almost-all]
[-B|--ignore-backups]
[-C|--vertical]
[-F|--classify]
[-G|--no-group]
[-d|--directory|--list-dirs]
[-g]
[-m]
[-o]
[-x]
```

**Usage**:

```
g [options] [path]
```

# GLOBAL OPTIONS

**--all, --la, -l**: show all info/use a long listing format

**--byline, --bl, -1, --oneline, --single-column**: print by line

**--check-new-version**: check if there's new release

**--depth**="": tree limit depth, negative -> infinity (default: infinity)

**--file-type, --ft**: likewise, except do not append '*'

**--format**="": across  -x,  commas  -m, horizontal -x, long -l, single-column -1, verbose -l, vertical -C (default: C)

**--full-time**: like -all/l --time-style=full-iso

**--git-status, --gs**: show git status: ? untracked, + added, ! deleted, ~ modified, | renamed, = copied, $ ignored

**--git-status-style, --gss**="": git status style: colored-symbol: {? untracked, + added, - deleted, ~ modified, | renamed, = copied, ! ignored} colored-dot

**--hide-git-ignore, --gi, --hgi**: hide git ignored file/dir

**--lh, --human-readable**: show human readable size

**--show-group, --sg**: show group

**--show-hidden, --sh, -a**: show hidden files

**--show-icon, --si, --icons**: show icon

**--show-no-dir, --nd, --nodir**: do not show directory

**--show-no-ext, --sne, --noext**="": show file which doesn't have target ext

**--show-only-dir, --sd, --dir**: show directory only

**--show-only-ext, --se, --ext**="": show file which has target ext, eg: --show-only-ext=go,java

**--show-owner, --so, --author**: show owner

**--show-perm, --sp**: show permission

**--show-size, --ss**: show file/dir size

**--show-time, --st**: show time

**--show-total-size, --ts**: show total size

**--theme, --th**="": apply theme `path/to/theme`

**--time-style**="": time/date format with -l, Valid timestamp styles are `default', `iso`, `long iso`, `full-iso`, `locale`, custom `+FORMAT` like date(1). (default: +%d.%b'%y %H:%M (like 02.Jan'06 15:04))

**--tree, -t**: list in tree

**--zero, -0**: end each output line with NUL, not newline

**-A, --almost-all**: do not list implied . and ..

**-B, --ignore-backups**: do not list implied entries ending with ~

**-C, --vertical**: list entries by columns (default)

**-F, --classify**: append indicator (one of */=>@|) to entries

**-G, --no-group**: in a long listing, don't print group names

**-d, --directory, --list-dirs**: list directories themselves, not their contents

**-g**: like -all/l, but do not list owner

**-m**: fill width with a comma separated list of entries

**-o**: like -all/l, but do not list group information

**-x**: list entries by lines instead of by columns


