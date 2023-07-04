#!/bin/bash

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

# add the following command to .bashrc
# eval "$(g --init bash)"
# then `source ~/.bashrc`
