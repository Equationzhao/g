output="$(g --no-update -term-width 200 --no-config --icons --permission --size --file-type tests/test_data )"
echo "$output" | diff - tests/file-type.stdout
