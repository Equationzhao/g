# [Deprecated] Test Workflow 

## requirement

the tests directory should be updated with `git submodule update --remote`.

> ⚠️ there are some test files whose name contains invalid characters in Windows, so the tests can't be updated in Windows.

the test data contains a series of testfile/dir(under test_data dir), scripts(*.sh), and expected outputs(*.stdout). The scripts are used to generate the outputs, and the expected outputs are used to compare with the outputs.

make sure you have `just` installed.

## pass script/run_test.sh(just test)
`just test`

check internal/theme/theme_test.go TestAll and TestColor 

> make sure running in a terminal supporting those features

## steps to add test case

### test flag

run `just newtest`, and follow the instructions

example:
```zsh
> just newtest
test_name: zero
flag: --zero
use base_flag? [Y/n] Y
```

the generated script will be `tests/zero.sh`:
```sh
output="$(g --no-update -term-width 200 --no-config --icons --permission --size --group --owner --zero tests/test_data )"
echo "$output" | diff - tests/zero.stdout
```

and the output will be zero.stdout

### test data

create files/directories in `tests/test_data`

run `just reproducetest` to generate the expected output


