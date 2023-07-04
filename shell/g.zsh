#!/bin/zsh

# alias for g
if [ "$(command -v g)" ]; then
    unalias -m 'll'
    unalias -m 'l'
    unalias -m 'la'
    unalias -m 'ls'
    alias ls='g'
    alias ll='g --perm --icons --time --group --owner --size --title'
    alias l='g --perm --icons --time --group --owner --size --title --show-hidden'
    alias la='g --show-hidden'
fi

# add the following command to .zshrc
# eval "$(g --init zsh)"
# then `source ~/.zshrc`
