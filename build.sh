#!/bin/bash
for file in *.go; do
    if [ -f "$file" ]; then
        go build "$file"
    fi
done
