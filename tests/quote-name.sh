output="$(g --no-update -term-width 200 --no-config --icons --permission --size --quote-name tests/test_data )"
echo "$output" | diff - tests/quote-name.stdout
