#!/bin/bash

if [ "$(command -v g)" ]; then
    if [ "$(command -v ll)" ]; then
      unalias ll
    fi

    if [ "$(command -v l)" ]; then
      unalias l
    fi

    if [ "$(command -v la)" ]; then
      unalias la
    fi

    alias ls='g'
    alias ll='g --perm --icons --time --group --owner --size --title'
    alias l='g --perm --icons --time --group --owner --size --title --show-hidden'
    alias la='g --show-hidden'
fi

# add the following command to .bashrc
# eval "$(g --init bash)"
# then 'source ~/.bashrc'
