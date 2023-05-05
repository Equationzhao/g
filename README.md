# g 
<div style="text-align: center;"><img src="logo.jpg" width="400"  alt="logo"/></div>

>   a powerful ls

[![Go Report Card](https://goreportcard.com/badge/github.com/Equationzhao/g)](https://goreportcard.com/report/github.com/Equationzhao/g)
[![wakatime](https://wakatime.com/badge/github/Equationzhao/g.svg)](https://wakatime.com/badge/github/Equationzhao/g)
## Screenshots

![image](./how-g-works.gif)

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

## More options
[g.md](g.md)


## custom theme

[theme](THEME.md)

## install
go version required >= 1.20
```bash
go install github.com/Equationzhao/g@latest
```

Archlinux user can install from AUR
```bash
yay -S g-ls
```

windows Scoop
```bash
scoop install https://raw.githubusercontent.com/Equationzhao/g/master/scoop/g.json
```

or just download from release page

## logo
created by bing
