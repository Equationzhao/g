error() {
    printf '\033[1;31m%s\033[0m\n' "$1"
}

success() {
    printf '\033[1;32m%s\033[0m\n' "$1"
}

warn() {
    printf '\033[1;33m%s\033[0m\n' "$1"
}