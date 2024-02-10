output="$(g --no-update -term-width 200 --no-config --icons --permission --size --reverse tests/test_data tests/test_data/)"
echo "$output" | diff - tests/multipath.stdout
