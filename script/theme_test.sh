#!/bin/bash

temp_build_name="custom_theme_test_g_"$RANDOM
echo "go build -tags=custom -o $temp_build_name"
CGO_ENABLED=0 go build -tags=custom -o $temp_build_name .
if [ $? -ne 0 ]; then
    echo "build failed"
    exit 1
fi
echo "build success"
rm $temp_build_name