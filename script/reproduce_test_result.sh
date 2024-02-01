# Reproduce the test result from the test script

for sh_file in tests/*.sh; do
    name="${sh_file%.*}"
    first_line=$(head -n 1 "$sh_file")
    eval "$first_line"
    # output is assigned in the test script
    echo "$output" > "$name.stdout"
done
