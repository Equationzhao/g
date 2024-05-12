#!/bin/bash
# install development requirements

lists=(
  "git"
  "upx"
  "dpkg"
  "gh"
  "wget"
  "gofumpt"
  "just"
  "prettier"
  "choose-rust"
  "ripgrep"
  "golangci-lint"
)

command=(
  "git"
  "upx"
  "dpkg"
  "gh"
  "wget"
  "gofumpt"
  "just"
  "prettier"
  "choose"
  "rg"
  "golangci-lint"
)

error() {
    printf '\033[1;31m%s\033[0m\n' "$1"
}

success() {
    printf '\033[1;32m%s\033[0m\n' "$1"
}

warn() {
    printf '\033[1;33m%s\033[0m\n' "$1"
}

# check if brew is installed
if ! command -v brew &> /dev/null; then
  error "brew is not installed"
  exit 1
fi

echo "brew update"
brew update > /dev/null
if [ $? -ne 0 ]; then
  error "brew update failed"
fi

if ! command -v go &> /dev/null; then
  echo "brew install go"
  brew install go
else # check go version >= 1.21.0
  go_version=$(go version | awk '{print $3}')
  go_version=${go_version:2}
  if [ "$(printf '%s\n' "1.21.0" "$go_version" | sort -V | head -n1)" != "1.21.0" ]; then
    # check if go is installed by brew
    if brew list --versions go &> /dev/null; then
      echo "brew upgrade go"
      brew upgrade go
    else
      echo "please upgrade go to 1.21.0 or later"
    fi
  fi
fi


for i in "${!lists[@]}"; do
  if ! command -v "${command[i]}" &> /dev/null; then
    echo "brew install ${lists[i]}..."
    HOMEBREW_NO_AUTO_UPDATE=1 brew install "${lists[i]}" > /dev/null
    if [ $? -ne 0 ]; then
      error "brew install ${lists[i]} failed"
    else
      success "${lists[i]} installed"
    fi
  else
    success "${lists[i]} already installed"
  fi
done