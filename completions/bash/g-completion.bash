_g() {
    local cur prev opts

    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    opts="
    --bug
    --duplicate
    --no-config
    --no-path-transform
    --check-new-version
    --help -h -?
    --version -v -#
    --csv
    --tsv
    --byline -1
    --classic
    --color
    --colorless
    --depth
    --format
    --file-type
    --md
    --markdown
    --table
    --table-style
    --term-width
    --theme
    --tree-style
    --zero -0 -C -F -R -T -d -j -m -x
    --init
    --sort
    --dir-first
    --group-directories-first
    --reverse -r
    --versionsort -S
    --si
    --sizesort -U
    --no-sort -X
    --sort-by-ext
    --accessed
    --all
    --birth
    --blocks
    --charset
    --checksum
    --checksum-algorithm
    --created
    --dereference
    --detect-size
    --footer
    --full-path
    --full-time
    --flags
    --gid
    --git
    --git-status
    --git-repo-branch
    --branch
    --git-repo-status
    --repo-status
    --group
    --header
    --title
    --hyperlink
    --icon
    --inode -i
    --mime
    --mime-parent
    --modified
    --mounts
    --no-dereference
    --no-icon
    --no-total-size
    --numeric
    --numeric-uid-gid
    --octal-perm
    --owner
    --perm
    --recursive-size
    --relative-to
    --relative-time
    --size
    --size-unit
    --block-size
    --smart-group
    --statistic
    --stdin
    --time
    --time-style
    --time-type
    --total-size
    --uid -G
    --no-group -H
    --link -N
    --literal -O
    --no-owner -Q
    --quote-name -g -l
    --long -o
    --extended -@"

    COMPREPLY=($(compgen -W "${opts}" -- ${cur}))
    return 0
}

complete -F _g g
