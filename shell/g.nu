alias ls = ^g
alias ll = ^g --perm --icons --time --group --owner --size --title
alias l = ^g --perm --icons --time --group --owner --size --title --show-hidden
alias la = ^g --show-hidden

# add the following to your $nu.env-path
# ^g --init nushell | save -f ~/.g.nu
# then add the following to your $nu.config-path
# source ~/.g.nu
# if you want to replace g with nushell's g command
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
