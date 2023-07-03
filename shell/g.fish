#!/usr/bin/fish

if command -v g >/dev/null 2>&1
    functions -e ll
    functions -e l
    functions -e la
    functions -e ls
    alias ls 'g'
    alias ll 'g --perm --icons --time --group --owner --size --title'
    alias l 'g --perm --icons --time --group --owner --size --title --show-hidden'
    alias la 'g --show-hidden'
end

#  add to config: g --init fish | source
# `source ~/.config/fish/config.fish`
