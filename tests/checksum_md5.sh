output="$(g -term-width 200 --no-config --no-update --icons --checksum -ca md5 tests/test_data )"
echo "$output" | diff - tests/checksum_md5.stdout
