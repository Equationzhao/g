output="$(g --no-update -term-width 200 --no-config --icons --permission --size --zero tests/test_data )"
echo "$output" | diff - tests/zero.stdout
