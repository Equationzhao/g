output="$(g -term-width 200 --no-config --icons --table tests/test_data )"
echo "$output" | diff - tests/table.stdout
