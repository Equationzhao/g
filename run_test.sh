# run all tests under tests/*.sh

GREEN='\033[0;32m'
RED='\033[0;31m'
NO_COLOR='\033[0m'

# Directory containing tests
TEST_DIR="tests"

# Run tests
for test_script in "$TEST_DIR"/*.sh; do
    # Run the script and capture the output
    output=$(bash "$test_script" 2>&1)

    # Check if output is empty
    if [ -z "$output" ]; then
        # Test passed
        echo "${GREEN}Passed:${NO_COLOR} $test_script"
    else
        # Test failed
        echo
        echo "${RED}Failed:${NO_COLOR} $test_script"
        echo "${RED}$output${NO_COLOR}"
    fi
done
