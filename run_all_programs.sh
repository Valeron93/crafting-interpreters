#/usr/bin/env bash

set -e

echo +==========================================+

for file in ./test_programs/*.vl; do
    echo RUNNING FILE: $file

    ./build/vl "$file"
    echo +==========================================+
done