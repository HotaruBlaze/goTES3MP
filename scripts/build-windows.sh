#!/bin/bash
GOOS=windows GOARCH=amd64 go build \
    -ldflags="-X 'main.Build=v0.0.0-Dev' -X 'main.GitCommit=DebugBuild'"\
    -o build/goTES3MP-Windows.exe \
    src/*.go 