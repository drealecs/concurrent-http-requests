#!/usr/bin/env sh

script_dir="$(realpath "$(dirname "$0")")"

cd "$script_dir"

mkdir -p "$script_dir/build"

GOOS=windows GOARCH=amd64 go build
GOOS=linux GOARCH=amd64 go build

mv concurrent-http* build

cd "$script_dir/concurrent-ssl-only"

GOOS=windows GOARCH=amd64 go build
GOOS=linux GOARCH=amd64 go build

mv concurrent-ssl* ../build
