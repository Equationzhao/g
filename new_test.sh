# creat a new test file in tests/*.sh and its output in tests/*.stdout

test_name=$1
if [ -z "$test_name" ]; then
    echo "test_name is empty"
    exit 1
fi

running_command=$2
if [ -z "$running_command" ]; then
    echo "running_command is empty"
    exit 1
fi

test_script="tests/$test_name.sh"
test_stdout="tests/$test_name.stdout"

output="$($running_command tests/test_data)"

echo "$output" > $test_stdout
echo "output=\"\$($running_command tests/test_data )\"" > "$test_script"
echo "echo \"\$output\" | diff - $test_stdout" >> "$test_script"

chmod +x "$test_script"
