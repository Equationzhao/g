output="$(g -term-width 200 --no-config --no-update --icons -R tests/test_data)"
echo "$output" | diff - tests/basic_reverse.stdout