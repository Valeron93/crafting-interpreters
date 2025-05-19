#/usr/bin/env bash

echo +==========================================+

for file in ./test_programs/*.vl; do
    echo RUNNING FILE: $file

    go run . "$file"
    echo +==========================================+
done