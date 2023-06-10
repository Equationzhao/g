# g

<div style="text-align: center;"><img src="logo.jpg" width="400"  alt="logo"/></div>

> a powerful ls

[![CodeFactor](https://www.codefactor.io/repository/github/equationzhao/g/badge/master)](https://www.codefactor.io/repository/github/equationzhao/g/overview/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/Equationzhao/g)](https://goreportcard.com/report/github.com/Equationzhao/g)
[![wakatime](https://wakatime.com/badge/github/Equationzhao/g.svg)](https://wakatime.com/badge/github/Equationzhao/g)
[![Go](https://github.com/Equationzhao/g/actions/workflows/go.yml/badge.svg)](https://github.com/Equationzhao/g/actions/workflows/go.yml)
![AUR license](https://img.shields.io/aur/license/g-ls)

![linux](https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black)
![macos](https://img.shields.io/badge/mac%20os-000000?style=for-the-badge&logo=apple&logoColor=white)
![windows](https://img.shields.io/badge/Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white)
![AUR version](https://img.shields.io/aur/version/g-ls?color=1793d1&label=g-ls&logo=arch-linux&style=for-the-badge)

g is a ls alternative with features:

1. display items with type-specific icons and colors that are easy to be customized
2. display in verious layouts ( grid/across/byline/tree/zero/comma )
3. user-friendly options with many alias
4. distinguish file git-status with with icons or symbol chars
5. highly customizable sort option
6. cross-platform ( Linux/Windows/MacOS )
7. option to fuzzy match path like [`zoxide`](https://github.com/ajeetdsouza/zoxide) with [`fzf`](https://github.com/junegunn/fzf) algorithm

## Screenshots

![image](./how-g-works.gif)

## install

### From source

go version required >= 1.20

```bash
go install -ldflags="-s -w -v"  github.com/Equationzhao/g@latest
```

or Clone this repo

```bash
git clone github.com/Equationzhao/g
cd g
go build -ldflags="-s -w" # use -s -w shrink size
# then add the executable file to your `PATH`
```

### Via package manager

![archlinux](https://img.shields.io/badge/Arch_Linux-1793D1?logo=arch-linux&logoColor=white)
user can install `g` from AUR

```bash
yay -S g-ls
```

windows Scoop

```powershell
scoop install https://raw.githubusercontent.com/Equationzhao/g/master/scoop/g.json
```

```powershell
# upgrade
scoop uninstall g # uninstall first
scoop install https://raw.githubusercontent.com/Equationzhao/g/master/scoop/g.json
# error msg like this can be ignored
# Move-Item: 
# Line |
#    2 |  mv g-amd64.exe g.exe
#      |  ~~~~~~~~~~~~~~~~~~~~
# Move-Item: 
# Line |
#    3 |  mv g-amd64.shim g.shim
#      |  ~~~~~~~~~~~~~~~~~~~~~~
```

homebrew: ***todo***


### Pre-built executable

just download from [release page](https://github.com/Equationzhao/g/releases) and add the executable file to your `PATH`

> binaries with -upx are compressed by [`upx`](https://github.com/upx/upx) with `--best` option

## usage

```bash
g path(s)
```

with icon

```bash
g -si    path
g -icons path
```

with mod(default) time

```bash
g -st        path
g -show-time path
g -time      path
```

with access/create/mod time

```bash
g -st -time-type="access" path
```

with fileperm

```bash
g -sp         path
g -permission path
```

with owner/group

```bash
g -so    path
g -owner path
g -sg    path
g -group path
```

with size

```bash
g -ss   path
g -size path
```

show all files, including hidden files

```bash
g -sh path
g -a  path
```

show dir only

```bash
g -sd  path
g -dir path
```

list by line

```bash
g -1             path
g -bl            path
g -oneline       path
g -single-column path
```

show file only with target ext

```bash
g -ext=<target ext(s)> path
```

list in tree

```bash
g -t    path
g -tree path
```

recurse into directories

```bash
g -R       path
g -recurse path
```

limit depth in tree/recurse (default: no limit)

```bash
g -t -depth=<level> path
g -R -depth=<level> path
```

fuzzy search

```bash
g -f     path
g -fuzzy path
# eg: g -f in
# /mnt/e/Project/gverything/index
# pathindex.go
```

disable index update

```bash
g -di            path  
g -no-update     path
g -disable-index path
```

disable color

```bash
g -nc        path
g -no-color  path
g -colorless path
```

show checksum (md5,sha1,sha224,sha256,sha384,sha512,crc32)

```bash
g -cs -ca=sha256 path
```

show git status with icon

```bash
g -gs         path
g -git        path
g -git-status path
```

show git status with char symbol

```bash
g -git -git-style=sym path
```

...

## More options

[g.md](g.md)

## custom theme

[theme](THEME.md)

## logo

created by bing

## [![Repography logo](https://images.repography.com/logo.svg)](https://repography.com) / Recent activity [![Time period](https://images.repography.com/35290882/Equationzhao/g/recent-activity/d06TKxKV8-Bc1zgTdodyAUFkmX-KdMR5ydV1GeE2jJY/r-OWQ7WewQlCCz2r7byT3_mCR0x8LTCx95ZyLfOY7CI_badge.svg)](https://repography.com)

[![Timeline graph](https://images.repography.com/35290882/Equationzhao/g/recent-activity/d06TKxKV8-Bc1zgTdodyAUFkmX-KdMR5ydV1GeE2jJY/r-OWQ7WewQlCCz2r7byT3_mCR0x8LTCx95ZyLfOY7CI_timeline.svg)](https://github.com/Equationzhao/g/commits)
[![Pull request status graph](https://images.repography.com/35290882/Equationzhao/g/recent-activity/d06TKxKV8-Bc1zgTdodyAUFkmX-KdMR5ydV1GeE2jJY/r-OWQ7WewQlCCz2r7byT3_mCR0x8LTCx95ZyLfOY7CI_prs.svg)](https://github.com/Equationzhao/g/pulls)
[![Trending topics](https://images.repography.com/35290882/Equationzhao/g/recent-activity/d06TKxKV8-Bc1zgTdodyAUFkmX-KdMR5ydV1GeE2jJY/r-OWQ7WewQlCCz2r7byT3_mCR0x8LTCx95ZyLfOY7CI_words.svg)](https://github.com/Equationzhao/g/commits)
[![Top contributors](https://images.repography.com/35290882/Equationzhao/g/recent-activity/d06TKxKV8-Bc1zgTdodyAUFkmX-KdMR5ydV1GeE2jJY/r-OWQ7WewQlCCz2r7byT3_mCR0x8LTCx95ZyLfOY7CI_users.svg)](https://github.com/Equationzhao/g/graphs/contributors)
