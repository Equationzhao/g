alias ls = ^g
alias ll = ^g --perm --icons --time --group --owner --size --title
alias l = ^g --perm --icons --time --group --owner --size --title --show-hidden
alias la = ^g --show-hidden

# ^g --init nushell | save -f ~/.g.nu
# source ~/.g.nu
# if you want to replace g with nushell's g 
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
