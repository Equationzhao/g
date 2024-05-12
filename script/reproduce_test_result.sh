#!/bin/bash

# Reproduce the test result from the test script
# This script will overwrite the test result in tests/*.stdout

error() {
    printf '\033[1;31m%s\033[0m\n' "$1"
}

success() {
    printf '\033[1;32m%s\033[0m\n' "$1"
}

warn() {
    printf '\033[1;33m%s\033[0m\n' "$1"
}

bye(){
  echo "Bye👋"
  exit 0
}

check_input(){
  echo
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
      bye
      exit 1
  fi
}

# PRINT WARNING
warn 'This script will overwrite the test result in tests/*.stdout'

read -p "Are you sure? (y/N) " -n 1 -r
check_input
read -p "Are you sure?? (y/N) " -n 1 -r
check_input
read -p "Are you sure??? (y/N) " -n 1 -r
check_input

echo "Well, you asked for it..."

printf 'which one do you want to reproduce?(name/all/none)\n'

read -p "Enter the name(s) of the test script(s): " -r
if [ "$REPLY" == "all" ]; then
    echo "Reproducing all test result..."
    for sh_file in tests/*.sh; do
        name="${sh_file%.*}"
        first_line=$(head -n 1 "$sh_file")
        eval "$first_line"
        # output is assigned in the test script
        echo "$output" > "$name.stdout"
    done
    success "Test result reproduced successfully.🎉"
    exit 0
fi

if [ "$REPLY" == "none" ]; then
    error "No test result will be reproduced."
    bye
    exit 0
fi

# split the input by comma
IFS=',' read -r -a test_names <<< "$REPLY"

# check if the test script exists
for test_name in "${test_names[@]}"; do
    sh_file="tests/$test_name.sh"
    if [ ! -f "$sh_file" ]; then
        error "$sh_file does not exist.😭"
        exit 1
    fi
done

# reproduce the test results
for test_name in "${test_names[@]}"; do
    echo "Reproducing $sh_file..."
    name="${sh_file%.*}"
    first_line=$(head -n 1 "$sh_file")
    eval "$first_line"
    # output is assigned in the test script
    echo "$output" > "$name.stdout"
done

success "Test result reproduced successfully.🎉"