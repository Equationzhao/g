output="$(g -term-width 200 --no-config --icons -R tests/test_data)"
echo "$output" | diff - tests/basic_reverse.stdout