output="$(g --no-update -term-width 200 --no-config --icons --permission --size --group --owner --sizesort tests/test_data )"
echo "$output" | diff - tests/sizesort.stdout
