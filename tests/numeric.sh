output="$(g -term-width 200 --no-config --no-update --icons --permission --size --group --owner --numeric tests/test_data )"
echo "$output" | diff - tests/numeric.stdout
