#!/bin/bash

# This is a very simple and naïve shell-script that traverses the various
# examples directories and ensures that they build. There is a ton of room for
# improvement here, such as writing a Makefile or similar.
reporoot="$PWD"
for file in $(find . -type f -name "main.go"); do
  example_dir="$(dirname $file)"
  # We have to cd to each directory since some examples have additional
  # required files present.
  cd "$example_dir"
  go build -o example . && ./example && rm ./example && cd "$reporoot"
done
