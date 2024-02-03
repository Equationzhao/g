# g

---

A feature-rich, customizable, and cross-platform `ls` alternative.

Experience enhanced visuals with type-specific icons, various layout options, and git status integration.

---

## Key Features

1. **Customizable Display**: Icons and colors specific to file types, easy to customize.
2. **Multiple Layouts**: Choose from grid, across, byline, zero, comma, table, json, markdown, and tree layouts.
3. **Git Integration**: View file git-status/repo-status/repo-branch directly in your listings.
4. **Advanced Sorting**: Highly customizable sorting options like version-sort.
5. **Cross-Platform Compatibility**: Works seamlessly on Linux, Windows, and MacOS.
6. **Fuzzy Path Matching**: [`zoxide`](https://github.com/ajeetdsouza/zoxide) and [`fzf`](https://github.com/junegunn/fzf) like fuzzy path matching.
7. **Hyperlink support**: Open files/directories with a single click.

## Screenshots

![image](asset/screenshot_3.png)

## Usage

```bash
g path(s)
```

```bash
g --icon --long path(s) # show icons and long format
```

```bash
g --tree --long path(s) # show tree layout
```

## More options

[g.md](docs/g.md)

## Installation Guide

### Via package manager

#### Arch Linux (AUR)

```bash
yay -S g-ls
```

#### Homebrew

```bash
brew tap equationzhao/core git@github.com:Equationzhao/homebrew-g.git
```

```bash
brew install g-ls
```

#### Windows Scoop

windows scoop:

```powershell
scoop install https://raw.githubusercontent.com/Equationzhao/g/master/scoop/g.json
```

```powershell
# upgrade
scoop uninstall g # uninstall first
scoop install https://raw.githubusercontent.com/Equationzhao/g/master/scoop/g.json
```

#### Winget

TODO, see [issue](https://github.com/Equationzhao/g/issues/119)

### Pre-built executable

#### curl

```sh
bash -c "$(curl -fsSLk https://raw.githubusercontent.com/Equationzhao/g/master/script/install.sh)"
```

#### deb

download from [release](https://github.com/Equationzhao/g/releases) page

```bash
sudo dpkg -i g_$version_$arch.deb
```

#### tar.gz/zip

just download from [release page](https://github.com/Equationzhao/g/releases), extract the gzip and add the executable file to your `PATH`

### From source

Requires Go version >= 1.21

```bash
go install -ldflags="-s -w"  github.com/Equationzhao/g@latest
```

Alternatively, clone the repo for a dev version:

```bash
git clone github.com/Equationzhao/g
cd g
go build -ldflags="-s -w" 
# then add the executable file to your `PATH`
```

## Recommended terminal

macOS:
- [Iterm2](https://iterm2.com/)
- [Warp](https://www.warp.dev)

Windows:
- [Windows Terminal](https://github.com/microsoft/terminal)

cross-platform:
- [Hyper](https://hyper.is/)
- [WezTerm](https://wezfurlong.org/wezterm/index.html)


## Shell Integration (alias)

Generate initialization scripts(alias) for various shells:

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

the nushell has a nice built-in ls command, but if you wanna try `g` in nushell, you can do the following:

ps: the script is not guaranteed to work, if you have any problem, please [file an issue](https://github.com/Equationzhao/g/issues/new/choose)

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

## Custom theme

[theme](THEME.md)

## TODO
- [x] Version sort
- [ ] Git sort
- [ ] Print security context
- [x] $OLDPWD
- [x] Support Scoop

The following are new features of `eza`, we may support them in the future
- [x] --git-repos: list each directory’s Git status, if tracked
- [x] --git-repos-no-status: list whether a directory is a Git repository, but not its status (faster)

## CONTRIBUTING

Interested in contributing? Check out the [contributing guidelines](./CONTRIBUTING.md).

## Alternatives

`g` is highly inspired by following projects that you wanna try!

- [exa](https://github.com/ogham/exa) or [eza](https://github.com/eza-community/eza)
- [lsd](https://github.com/lsd-rs/lsd)
- [ls-go](https://github.com/acarl005/ls-go)

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=Equationzhao/g&type=Date)](https://star-history.com/#Equationzhao/g&Date)
