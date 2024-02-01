output="$(g -term-width 200 --no-config --icons -l -j tests/test_data )"
echo "$output" | diff - tests/json.stdout
