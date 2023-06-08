# g 
<div style="text-align: center;"><img src="logo.jpg" width="400"  alt="logo"/></div>

>   a powerful ls

[![Go Report Card](https://goreportcard.com/badge/github.com/Equationzhao/g)](https://goreportcard.com/report/github.com/Equationzhao/g)
[![wakatime](https://wakatime.com/badge/github/Equationzhao/g.svg)](https://wakatime.com/badge/github/Equationzhao/g)
## Screenshots

![image](./how-g-works.gif)

## install
go version required >= 1.20
```bash
go install -ldflags="-s -w -v"  github.com/Equationzhao/g@latest
```

Archlinux user can install `g` from AUR
```bash
yay -S g-ls
```

windows Scoop
```bash
scoop install https://raw.githubusercontent.com/Equationzhao/g/master/scoop/g.json
```
```bash
# upgrade
scoop uninstall g # uninstall first
scoop install https://raw.githubusercontent.com/Equationzhao/g/master/scoop/g.json
```

or just download from release page

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

