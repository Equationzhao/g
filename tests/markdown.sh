output="$(g -term-width 200 --no-config --icons -l --markdown tests/test_data )"
echo "$output" | diff - tests/markdown.stdout
