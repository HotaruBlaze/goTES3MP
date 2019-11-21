#!/bin/bash
go build \
    -o build/goTes3mp \
    src/main.go \
    src/linter.go \
    src/config.go \
    src/webserver.go \
    src/status.go \
    src/discord.go \
    src/cfgReader.go 

# mkdir tmp/
# Current=$(md5sum build/goTes3mp | awk '{print $1}')
# curl -o tmp/goTes3mp https://img.fluttershub.com/4XMhb0PLmP8mcbaQ.bin
# New=$(md5sum ./tmp/goTes3mp | awk '{print $1}')

# if [ "$Current" = "$New" ]; then
#     echo "Binarys are equal."
#     rm -Rf tmp/
# else
#     echo "Binarys are not equal."
# fi