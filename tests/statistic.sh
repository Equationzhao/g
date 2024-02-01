output="$(g -term-width 200 --no-config --no-update --icons --statistic tests/test_data)"
echo "$output" | tail -n +2 | diff - <(tail -n +2 tests/statistic.stdout)

