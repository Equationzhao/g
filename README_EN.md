# g

<div style="text-align: center;"><img src="internal/logo.jpg" width="400"  alt="logo"/></div>

> a powerful ls

![linux](https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black)
![macos](https://img.shields.io/badge/mac%20os-000000?style=for-the-badge&logo=apple&logoColor=white)
![windows](https://img.shields.io/badge/Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white)
![AUR version](https://img.shields.io/aur/version/g-ls?color=1793d1&label=g-ls&logo=arch-linux&style=for-the-badge)

g is a ls alternative with features:

1. display items with type-specific icons and colors that are easy to be customized
2. display in various layouts ( grid/across/byline/zero/comma/table/html/json/markdown/tree )
3. user-friendly options with many aliases
4. check file git-status when listing entries
5. highly customizable sort option
6. cross-platform ( Linux/Windows/MacOS )
7. option to fuzzy match the path like [`zoxide`](https://github.com/ajeetdsouza/zoxide) with [`fzf`](https://github.com/junegunn/fzf) algorithm

## Screenshots

![image](asset/screenshot_3.png)

## Install

### From source

go version required >= 1.21

```bash
go install -ldflags="-s -w"  github.com/Equationzhao/g@latest
```

or Clone this repo (nightly build)

```bash
git clone github.com/Equationzhao/g
cd g
go build -ldflags="-s -w" # use -s -w to shrink size
# then add the executable file to your `PATH`
```

### Via package manager

![archlinux](https://img.shields.io/badge/Arch_Linux-1793D1?logo=arch-linux&logoColor=white)
user can install `g` from AUR

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

windows scoop:

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

### Pre-built executable

#### deb

download from [release](https://github.com/Equationzhao/g/releases) page 

```bash
sudo dpkg -i g_$version_$arch.deb
```

#### tar.gz/zip

just download from [release page](https://github.com/Equationzhao/g/releases), extract the gzip and add the executable file to your `PATH`

## Recommended terminal

macOS:
- [Iterm2](https://iterm2.com/)
- [Warp](https://www.warp.dev)

Windows:
- [Windows Terminal](https://github.com/microsoft/terminal)

cross-platform:
- [Hyper](https://hyper.is/)
- [WezTerm](https://wezfurlong.org/wezterm/index.html)

## Usage

```bash
g path(s)
```

## Shell scripts

generate shell scripts

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

use command `echo $profile` to find your profile path

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

## More options

[g.md](g.md) or check [g.equationzhao.space](g.equationzhao.space)

## Custom theme

[theme](THEME.md)

## TODO
- [x] Version sort
- [ ] Git sort
- [ ] Print security context
- [x] $OLDPWD
- [ ] Color Support for html/markdown
- [x] Support Scoop

The following are new features of `eza`, we may support them in the future
- [x] --git-repos: list each directory’s Git status, if tracked
- [x] --git-repos-no-status: list whether a directory is a Git repository, but not its status (faster)

## CONTRIBUTING

check [CONTRIBUTING](./CONTRIBUTING.md)

## Logo

created by bing

## Alternatives

this project is highly inspired by following projects that you wanna try!

- [exa](https://github.com/ogham/exa) or [eza](https://github.com/eza-community/eza)
- [lsd](https://github.com/lsd-rs/lsd)
- [ls-go](https://github.com/acarl005/ls-go)

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=Equationzhao/g&type=Date)](https://star-history.com/#Equationzhao/g&Date)
