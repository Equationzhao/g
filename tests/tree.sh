output="$(g -term-width 200 --no-config --icons --tree tests/test_data)"
echo "$output" | diff - tests/tree.stdout
