output="$(g --no-update -term-width 200 --no-config --icons --permission --size --header --footer tests/test_data )"
echo "$output" | diff - tests/header-footer.stdout
