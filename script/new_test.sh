# creat a new test file in tests/*.sh and its output in tests/*.stdout

error() {
    printf '\033[1;31m%s\033[0m\n' "$1"
}

success() {
    printf '\033[1;32m%s\033[0m\n' "$1"
}

warn() {
    printf '\033[1;33m%s\033[0m\n' "$1"
}

if [ ! -f "script/run_test.sh" ]; then
    error "Please run the script in the root directory of the project"
    exit 1
fi

# input test_name
read -p "test_name: " test_name

if [ -z "$test_name" ]; then
    error "test_name is empty"
    exit 1
fi

read -p "running_command: " running_command
if [ -z "$running_command" ]; then
    error "running_command is empty"
    exit 1
fi

test_script="tests/$test_name.sh"
test_stdout="tests/$test_name.stdout"

if [ -f "$test_script" ]; then
    warn "$test_script already exists"
    read -p "Do you want to overwrite it? [y/N] " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

output="$($running_command tests/test_data)"

echo "$output" > $test_stdout
echo "output=\"\$($running_command tests/test_data )\"" > "$test_script"
echo "echo \"\$output\" | diff - $test_stdout" >> "$test_script"

chmod +x "$test_script"
