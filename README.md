# g

<div style="text-align: center;"><img src="logo.jpg" width="400"  alt="logo"/></div>

> 一个强大的 ls 工具

[![CodeFactor](https://www.codefactor.io/repository/github/equationzhao/g/badge/master)](https://www.codefactor.io/repository/github/equationzhao/g/overview/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/Equationzhao/g)](https://goreportcard.com/report/github.com/Equationzhao/g)
[![wakatime](https://wakatime.com/badge/github/Equationzhao/g.svg)](https://wakatime.com/badge/github/Equationzhao/g)
[![Go](https://github.com/Equationzhao/g/actions/workflows/go.yml/badge.svg)](https://github.com/Equationzhao/g/actions/workflows/go.yml)
![AUR license](https://img.shields.io/aur/license/g-ls)

![linux](https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black)
![macos](https://img.shields.io/badge/mac%20os-000000?style=for-the-badge&logo=apple&logoColor=white)
![windows](https://img.shields.io/badge/Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white)
![AUR version](https://img.shields.io/aur/version/g-ls?color=1793d1&label=g-ls&logo=arch-linux&style=for-the-badge)

<p align="center">
<a href="README_EN.md">View this document in English</a>
</p>


g 是一个 ls 替代品，拥有下面一些功能：

1. 显示带有类型特定图标和颜色的条目，并且易于更改
2. 有丰富的输出格式  ( grid/across/byline/zero/comma/table/html/json/markdown/tree )
3.  用户友好的选项
4. 支持显示 git status 
5. 丰富且可自定义的排序选项
6. 跨平台 ( Linux/Windows/MacOS )
7. 支持使用[`fzf`](https://github.com/junegunn/fzf) 算法，像 [`zoxide`](https://github.com/ajeetdsouza/zoxide) 一样模糊匹配路径 

## 截图

![image](how-g-works.gif)

## 安装

### 源码安装

要求 go version >= 1.21

```bash
go install -ldflags="-s -w"  github.com/Equationzhao/g@latest
```

或者 clone 这个仓库 (nightly build)

```bash
git clone github.com/Equationzhao/g
cd g
go build -ldflags="-s -w" # use -s -w to shrink size
# then add the executable file to your `PATH`
```

### 通过包管理器

![archlinux](https://img.shields.io/badge/Arch_Linux-1793D1?logo=arch-linux&logoColor=white)
用户可以通过 AUR 安装 `g`

```bash
yay -S g-ls
```

homebrew:

```bash
brew tap equationzhao/core git@github.com:Equationzhao/homebrew-g.git
```

```bash
brew install g-ls
```

### 二进制文件

#### deb
```bash
sudo dpkg -i g_$version_$arch.deb
```

#### tar,gz/zip
从 [release page](https://github.com/Equationzhao/g/releases) 下载对应平台的文件, 解压 gzip 并将可执行文件添加到 `PATH`

## 推荐使用的终端

macOS:
- [Iterm2](https://iterm2.com/)
- [Warp](https://www.warp.dev)

Windows:
- [Windows Terminal](https://github.com/microsoft/terminal)

cross-platform:
- [Hyper](https://hyper.is/)
- [WezTerm](https://wezfurlong.org/wezterm/index.html)

## 用法

```bash
g path(s)
```

显示图标

```bash
g -icons
```

显示修改(默认)时间

```bash
g -time    
```

显示 访问/创建/修改 时间

```bash
g -time -time-type=access
g -time -ac/cr/mod
```

显示文件权限

```bash
g -permission 
g -octal-perm # show octal permission like 0777
```

显示用户/群组

```bash
g -owner 
g -group 
```

显示文件大小

```bash
g -size 
g -size -recusive-size # show size of dir recursively
```

显示所有文件，包括隐藏文件

```bash
g -sh 
g -show-hidden
g -a  
```

只显示目录

```bash
g -dir 
```

按行显示

```bash
g -1           
g -oneline     
g -single-column 
```

显示有指定拓展名的文件

```bash
g -ext=<target ext(s)> 
# eg:
# g -ext=go,md
```

递归显示目录

```bash
g -R     
g -recurse 
```

限制在 树/递归 模式下的 递归深度 (默认: 无限制)

```bash
g -R -depth=<level> 
```

模糊搜索

```bash
g -f   
g -fuzzy 
# eg: g -f in
# /mnt/e/Project/gverything/index
# pathindex.go
```

禁用索引更新

```bash
g -di            
g -no-update   
g -disable-index 
```

禁用颜色

```bash
g -no-color  
g -colorless 
```

设置颜色

```bash
g -color=always 
g -color=auto    # default
g -color=never   
g -color=16/basic       # 16-color
g -color=256/8bit       # 256-color
g -color=16m/24bit/true-color   # 24-bit
```

显示校验和 (md5,sha1,sha224,sha256,sha384,sha512,crc32)

```bash
g -cs -ca=sha256 
```

显示 git status

```bash
g -git     
g -git-status
```

表格式输出

```bash
g -tb
```

树状输出

```bash
g -tree
```

以 markdown 格式输出, 并用 [glow](github.com/charmbracelet/glow) 渲染
( 不支持颜色 )

```bash
g -md | glow 
```

![image](https://github.com/Equationzhao/g/assets/75521101/7ec1e0d7-03cd-4968-ba48-2ec5375086fa)

...

## Shell 脚本

生成 shell 脚本

```bash
g -init bash/zsh/fish/pwsh
```

### bash

```.bash
# add the following command to .bashrc
eval "$(g --init bash)"
# then `source ~/.bashrc`
```

### zsh

```zsh
# add the following command to .zshrc
eval "$(g --init zsh)"
# then `source ~/.zshrc`
```

### fish

```fish
#  add to fish config:
g --init fish | source
#  then `source ~/.config/fish/config.fish`
```

### powershell

```powershell
# add the following line to your profile
Invoke-Expression (& { (g --init powershell | Out-String) })
```

使用 `echo $profile`命令查找配置文件路径

### nushell

```nu
# add the following to your $nu.env-path
^g --init nushell | save -f ~/.g.nu
# then add the following to your $nu.config-path
source ~/.g.nu

# if you want to replace nushell's g command with g
# add the following definition and alias to your $nu.config-path
#
# def nug [arg?] {
#     if ($arg == null) {
#         g $arg
#     } else {
#         g
#     }
# }
# alias g = ^g
```

## 更多选项

[g.md](g.md)  或查看 [主页](g.equationzhao.space)

## 自定义主题

[theme](THEME.md)

## TODO
- [x] Version sort
- [ ] Git sort
- [ ] Print security context
- [x] $OLDPWD
- [ ] Color Support for html/markdown

以下是 eza 的新功能，后续计划支持
- [ ] --git-repos: list each directory’s Git status, if tracked
- [ ] --git-repos-no-status: list whether a directory is a Git repository, but not its status (faster)

## Logo

created by bing

## 其他选择

本项目受到以下项目的启发，你也许想试试

- [exa](https://github.com/ogham/exa) 或者 [eza](https://github.com/eza-community/eza)
- [lsd](https://github.com/lsd-rs/lsd)
- [ls-go](https://github.com/acarl005/ls-go)

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=Equationzhao/g&type=Date)](https://star-history.com/#Equationzhao/g&Date)


## 查看帖子

- [deepin bbs](https://bbs.deepin.org/post/261954)

## [![Repography logo](https://images.repography.com/logo.svg)](https://repography.com) / Recent activity [![Time period](https://images.repography.com/35290882/Equationzhao/g/recent-activity/d06TKxKV8-Bc1zgTdodyAUFkmX-KdMR5ydV1GeE2jJY/r-OWQ7WewQlCCz2r7byT3_mCR0x8LTCx95ZyLfOY7CI_badge.svg)](https://repography.com)

[![Timeline graph](https://images.repography.com/35290882/Equationzhao/g/recent-activity/d06TKxKV8-Bc1zgTdodyAUFkmX-KdMR5ydV1GeE2jJY/r-OWQ7WewQlCCz2r7byT3_mCR0x8LTCx95ZyLfOY7CI_timeline.svg)](https://github.com/Equationzhao/g/commits)
[![Pull request status graph](https://images.repography.com/35290882/Equationzhao/g/recent-activity/d06TKxKV8-Bc1zgTdodyAUFkmX-KdMR5ydV1GeE2jJY/r-OWQ7WewQlCCz2r7byT3_mCR0x8LTCx95ZyLfOY7CI_prs.svg)](https://github.com/Equationzhao/g/pulls)
[![Trending topics](https://images.repography.com/35290882/Equationzhao/g/recent-activity/d06TKxKV8-Bc1zgTdodyAUFkmX-KdMR5ydV1GeE2jJY/r-OWQ7WewQlCCz2r7byT3_mCR0x8LTCx95ZyLfOY7CI_words.svg)](https://github.com/Equationzhao/g/commits)
[![Top contributors](https://images.repography.com/35290882/Equationzhao/g/recent-activity/d06TKxKV8-Bc1zgTdodyAUFkmX-KdMR5ydV1GeE2jJY/r-OWQ7WewQlCCz2r7byT3_mCR0x8LTCx95ZyLfOY7CI_users.svg)](https://github.com/Equationzhao/g/graphs/contributors)
