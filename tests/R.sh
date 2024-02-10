output="$(g --no-update -term-width 200 --no-config --icons --permission --size -R tests/test_data )"
echo "$output" | diff - tests/R.stdout
