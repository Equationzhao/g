output="$(g -term-width 200 --no-config --no-update --icons --table --permission --size tests/test_data)"
echo "$output" | diff - tests/table.stdout
