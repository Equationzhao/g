output="$(g -term-width 200 --no-config --no-update --icons --permission --size --group --owner --recursive-size tests/test_data )"
echo "$output" | diff - tests/recursive-size.stdout
