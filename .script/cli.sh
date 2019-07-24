#! /bin/bash

# Cross-compile Temporal using gox, injecting appropriate tags.
go get -u github.com/mitchellh/gox

TARGETS="linux/amd64 linux/386 linux/arm darwin/amd64 darwin/386 windows/amd64 windows/386"

mkdir -p release

gox -output="release/thc-{{.OS}}-{{.Arch}}" \
    -osarch="$TARGETS" \
    ./cmd/thc


ls ./release/thc* > files
for i in $(cat files); do
    sha256sum "$i" > "$i.sha256"
done