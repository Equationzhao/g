output="$(g --no-update -term-width 200 --no-config --icons --permission --size --group --owner --ignore-backups tests/test_data )"
echo "$output" | diff - tests/ignore-backups.stdout
