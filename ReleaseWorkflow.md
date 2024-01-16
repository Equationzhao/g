# how to release

## requirement

os: `macOS`

software/toolchain:

| name       | how to install             | remark                                                                            |
|------------|----------------------------|-----------------------------------------------------------------------------------|
| go >1.21.0 | `brew install go`          | or use [go.dev](https://go.dev/dl/) / [goup](https://github.com/owenthereal/goup) |
| git        | `brew install git`         | or use xcode version                                                              |
| upx        | `brew install upx`         |                                                                                   |
| dpkg-deb   | `brew install dpkg`        |                                                                                   |
| shasum     | `brew install shasum`      |                                                                                   |
| gh         | `brew install gh`          |                                                                                   |
| wget       | `brew install wget`        |                                                                                   |
| gofumpt    | `brew install gofumpt`     |                                                                                   |
| just       | `brew install just`        |                                                                                   |
| prettier   | `brew install prettier`    |                                                                                   |
| choose     | `brew install choose-rust` |                                                                                   |

```zsh
brew install go git upx dpkg shasum gh wget gofumpt just prettier choose-rust
```

## pre-check

- [ ] check code format: `just format`
- [ ] gen theme/doc file: `just theme doc`
- [ ] run test: `just test`
- [ ] check version: make sure the git tag and internal/cli/Version is the same. And git status is clean, git tag is at the current HEAD: `just check`

## build

- [ ] cleanup: `just clean`
- [ ] build: `just build`
- [ ] compress: `just compress`
- [ ] gen deb pkg: `just deb`
- [ ] gen checksum: `just checksum`

## release

- [ ] release: `just release`

## package manager

### AUR

ssh://aur@aur.archlinux.org/g-ls.git
make sure the aur repo is at '../g-ls' and 'Already up-to-date'

```zsh
just aur
```

### brew-tap

git@github.com:Equationzhao/homebrew-g.git
make sure the brew-tap repo is at '../homebrew-g' and 'Already up-to-date'

```zsh
just brew-tap
```

### scoop

the scoop manifest is at [scoop](scoop/g.json)

```zsh
git add -u && git commit -m 'ci: :construction_worker: update scoop'
git push
```

if you have no access to push to the master branch, please push to another branch and make a pull request
