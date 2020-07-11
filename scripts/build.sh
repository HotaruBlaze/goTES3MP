#!/bin/bash
go build \
    -ldflags="-X 'main.Build=v0.0.0-Dev' -X 'main.GitCommit=DebugBuild'"\
    -o build/goTES3MP-Linux \
    src/*.go 