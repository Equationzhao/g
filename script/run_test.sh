#!/bin/bash

echo "This script is deprecated"
exit 0

# should run in the root directory of the project
# run all tests under tests/*.sh

if [ ! -f "script/run_test.sh" ]; then
    echo "Please run the script in the root directory of the project"
    exit 1
fi

GREEN='\033[0;32m'
RED='\033[0;31m'
NO_COLOR='\033[0m'

# Directory containing tests
TEST_DIR="tests"

pass_count=0
fail_count=0

# run g --help and get Configuration file path
# backup Configuration file
# Configuration: path (may contains space)
config_path=$(g --help | grep "Configuration:" | cut -d ":" -f 2- | sed -e 's/^[[:space:]]*//')
should_restore=0
if [ -f "$config_path" ]; then
  echo "disable config: $config_path"
  echo "backup config: $config_path.bak"
  # if success, restore it at the end
  if mv "$config_path" "$config_path.bak"; then
    should_restore=1
  fi
fi


# Run tests
for test_script in "$TEST_DIR"/*.sh; do
    # Run the script and capture the output
    output=$(bash "$test_script" 2>&1)

    # Check if output is empty
    if [ -z "$output" ]; then
        # Test passed
        echo "${GREEN}Passed:${NO_COLOR} $test_script"
        pass_count=$((pass_count+1))
    else
        # Test failed
        echo
        echo "${RED}Failed:${NO_COLOR} $test_script"
        echo "${RED}$output${NO_COLOR}"
        fail_count=$((fail_count+1))
    fi
done

echo
echo "Passed: $pass_count"
echo "Failed: $fail_count"

# Restore Configuration file
if [ "$should_restore" -eq 1 ]; then
  echo "restore config: $config_path"
  mv "$config_path.bak" "$config_path"
fi

if [ "$fail_count" -gt 0 ]; then
    exit 1
fi