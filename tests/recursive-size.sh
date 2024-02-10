output="$(g -term-width 200 --no-config --no-update --icons --permission --size --recursive-size tests/test_data )"
echo "$output" | diff - tests/recursive-size.stdout
