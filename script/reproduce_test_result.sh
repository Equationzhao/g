#!/bin/bash

# Reproduce the test result from the test script
# This script will overwrite the test result in tests/*.stdout

check_input(){
  echo
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
      exit 1
  fi
}

# PRINT WARNING
printf '\033[1;33mThis script will overwrite the test result in tests/*.stdout\033[0m\n'

read -p "Are you sure? (y/N) " -n 1 -r
check_input
read -p "Are you sure?? (y/N) " -n 1 -r
check_input
read -p "Are you sure??? (y/N) " -n 1 -r
check_input

echo "Well, you asked for it..."
echo "Reproducing test result..."


for sh_file in tests/*.sh; do
    name="${sh_file%.*}"
    first_line=$(head -n 1 "$sh_file")
    eval "$first_line"
    # output is assigned in the test script
    echo "$output" > "$name.stdout"
done
