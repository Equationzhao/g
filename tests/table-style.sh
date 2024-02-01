output="$(g --no-update -term-width 200 --no-config --icons --permission --size --group --owner --table-style=unicode --table tests/test_data )"
echo "$output" | diff - tests/table-style.stdout
