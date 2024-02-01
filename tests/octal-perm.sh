output="$(g -term-width 200 --no-config --no-update --icons --permission --size --group --owner --octal-perm tests/test_data )"
echo "$output" | diff - tests/octal-perm.stdout
