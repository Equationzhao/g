output="$(g -term-width 200 --no-config --no-update --icons --permission --size --header tests/test_data )"
echo "$output" | diff - tests/header.stdout
