output="$(g --no-update -term-width 200 --no-config --icons --permission --size --group --owner --classic tests/test_data )"
echo "$output" | diff - tests/classic.stdout
