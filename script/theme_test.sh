#!/bin/bash

# load base.sh
source "$(dirname "$0")/base.sh"

temp_build_name="custom_theme_test_g_"$RANDOM
echo "go build -tags=custom -o $temp_build_name"
CGO_ENABLED=0 go build -tags=custom -o $temp_build_name .
if [ $? -ne 0 ]; then
    error "build failed"
    exit 1
fi
success "build success"
rm $temp_build_name