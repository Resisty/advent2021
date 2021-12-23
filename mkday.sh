#!/usr/bin/env bash
if [ -z "${1}" ]; then
  echo "Usage: mkday.sh [0-9]+"
  exit 1
fi
mkdir cmd/day${1}
touch cmd/day${1}/day${1}.go
touch cmd/day${1}/test.txt
touch cmd/day${1}/input.txt
