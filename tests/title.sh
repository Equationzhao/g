output="$(g --no-update -term-width 200 --no-config --icons --permission --size --title tests/test_data )"
echo "$output" | diff - tests/title.stdout
