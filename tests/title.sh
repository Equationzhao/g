output="$(g --no-update -term-width 200 --no-config --icons --permission --size --group --owner --title tests/test_data )"
echo "$output" | diff - tests/title.stdout
