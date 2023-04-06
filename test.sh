#!/bin/bash

# This is a very simple and naive shell-script that traverses the various
# examples directories and ensures that they build. There is a ton of room for
# improvement here.
for dir in $(find . -type d -name "[0-9][0-9]-*"); do cd "$dir" && go build . && cd ../..; done
