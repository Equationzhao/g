output="$(g -term-width 200 --no-config --no-update --icons tests/test_data)"
echo "$output" | diff - tests/basic.stdout
