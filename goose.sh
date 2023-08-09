#!/bin/sh

echo "change to work dir"
cd terminal/goose || return
echo "execute commands and flags"
go run main.go "$1" "$2" "$3" # allowing only 5 args for now