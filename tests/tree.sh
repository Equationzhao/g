output="$(g -term-width 200 --no-config --no-update --icons --permission --size --group --owner --tree tests/test_data)"
echo "$output" | diff - tests/tree.stdout
