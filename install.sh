# Get OS and architecture
os_type=$(uname -s | tr '[:upper:]' '[:lower:]')
os_arch=$(uname -m)
version="0.25.2"

interrupt_handler() {
    exit 1
}

trap interrupt_handler SIGINT

# flag
# -v: print version
# -h: print help
# -d: download binary (default)

help(){
    echo "g: Yet another ls"
    echo "Usage: install.sh [-v|-h|-d]"
    echo "  -v: print version"
    echo "  -h: print help"
    echo "  -d: download binary (default)"
}

error() {
    printf '\033[1;31m%s\033[0m\n' "$1"
}

success() {
    printf '\033[1;32m%s\033[0m\n' "$1"
}

warn() {
    printf '\033[1;33m%s\033[0m\n' "$1"
}

# Parse flags
while getopts "vhd" opt; do
    case $opt in
        v)
            echo "$version"
            exit 0
            ;;
        h)
            help
            exit 0
            ;;
        d)
            ;;
        \?)
            error "Invalid option: -$OPTARG"
            help
            exit 1
            ;;
    esac
done


compare_versions() {
    local A="$1"
    local B="$2"

    if [[ $A == $B ]]; then
        return 0
    fi

    IFS='.' read -ra A_parts <<< "$A"
    IFS='.' read -ra B_parts <<< "$B"

    for i in "${!A_parts[@]}"; do
        if (( ${A_parts[i]} > ${B_parts[i]} )); then
            return 0
        elif (( ${A_parts[i]} < ${B_parts[i]} )); then
            return 1
        fi
    done

}

# if already has g, and g --version >= version, exit
if command -v g &> /dev/null; then
    g_version=$(g --version | awk 'NR==2 {print $NF}')
    compare_versions "$g_version" "$version"
    result=$?
    if [ $result -eq 0 ]; then
        echo "g version $g_version already installed"
        exit 0
    fi
fi


# Determine file architecture based on OS architecture
case $os_arch in
    x86_64)
        file_arch="amd64"
        ;;
    arm64)
        file_arch="arm64"
        ;;
    i386)
        file_arch="386"
        ;;
    arm*)
        file_arch="arm"
        ;;
    *)
        error "Unsupported architecture: $os_arch"
        exit 1
        ;;
esac

# Determine file OS based on OS type
case $os_type in
    darwin)
        file_os="darwin"
        ;;
    linux)
        file_os="linux"
        ;;
#    msys)
#        file_os="windows"
#        ;;
    *)
        error "Unsupported OS type: $os_type"
        exit 1
        ;;
esac

file_name=g-${file_os}-${file_arch}

# Build download URL
url="https://github.com/Equationzhao/g/releases/download/v${version}/${file_name}"

## Add .exe extension for Windows
# if [ "$file_os" = "windows" ]; then
#     url="${url}.exe"
# fi

# Download the file using curl or wget
if command -v curl &> /dev/null; then
    echo "curl -fLO $url"
    curl -fLO $url
elif command -v wget &> /dev/null; then
    echo "wget $url"
    wget $url
else
    error "You need to install curl or wget to download the file."
    exit 1
fi

if [ $? -ne 0 ]; then
        error "Failed to download the file."
        exit 1
fi

# Make the file executable for Linux or Darwin
if [ "$file_os" = "linux" ] || [ "$file_os" = "darwin" ]; then
    chmod +x g-${file_os}-${file_arch}
fi

case $os_type in
    darwin)
        echo "mv ${file_name} /usr/local/bin/g"
        sudo mv ${file_name} /usr/local/bin/g
        ;;
    linux)
        echo "mv ${file_name} /usr/bin/g"
        sudo mv ${file_name} /usr/bin/g
        ;;
esac