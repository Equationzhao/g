#!/bin/bash

echo "This script is deprecated"
exit 0

# creat a new test file in tests/*.sh and its output in tests/*.stdout

# load base.sh
source "$(dirname "$0")/base.sh"

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

read -p "flag: " flag
if [ -z "$flag" ]; then
    error "flag is empty"
    exit 1
fi

command="g"
forever_base="--no-update"
base_flag="-term-width 200 --no-config --icons --permission --size"

read -p "use base_flag? [Y/n] " -n 1 -r
echo
if [[ $REPLY =~ ^[Nn]$ ]]; then
    base_flag=""
fi

running_command="$command $forever_base $base_flag $flag"

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
echo "#!/bin/bash" > "$test_script"
echo "output=\"\$($running_command tests/test_data )\"" >> "$test_script"
echo "echo \"\$output\" | diff - $test_stdout" >> "$test_script"

chmod +x "$test_script"
