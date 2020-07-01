#!/bin/bash
go build \
    -ldflags="-X 'main.Build=v0.1.0' -X 'main.GitCommit=DebugBuild'"\
    -o build/goTES3MP-Linux \
    src/*.go 