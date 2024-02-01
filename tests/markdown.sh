output="$(g -term-width 200 --no-config --no-update --icons --permission  --size --group --owner --markdown tests/test_data)"
echo "$output" | diff - tests/markdown.stdout
