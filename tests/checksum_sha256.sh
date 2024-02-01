output="$(g -term-width 200 --no-config --icons --checksum -ca sha256 tests/test_data )"
echo "$output" | diff - tests/checksum_sha256.stdout
