output="$(g -term-width 200 --no-config --no-update --icons --permission --size --markdown tests/test_data)"
echo "$output" | diff - tests/markdown.stdout
