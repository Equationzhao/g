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
g path
```

with icon
```bash
g -si path
```

with mod time
```bash
g -st path
```

with fileperm
```bash
g -sp path
```

with owner/group
```bash
g -so path
g -sg path
```

with size
```bash
g -ss path
```

show all files, including hidden files
```bash
g -sh path
```

show dir only
```bash
g -sd path
```
list by line
```bash
g -1 path
g -bl path
```

show file only with target ext
```bash
g -ext=<target ext> path
```

list in tree
```bash
g -t path
```

limit depth in tree
```bash
g -t -depth=<level> path
```

fuzzy search
```bash
g -f path
# eg: g -f in
# /mnt/e/Project/gverything/index
# pathindex.go
```

show checksum (md5,sha1,sha256,sha512)
```bash
g -cs -ca=sha256 path
```

## More options
[g.md](g.md)

## custom theme

[theme](THEME.md)

## logo
created by bing
