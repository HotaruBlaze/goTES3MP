#!/bin/bash

# Commit hash
commitHash=$(git rev-parse HEAD)

go build -ldflags="-X 'main.Build=v0.0.0-Dev' -X 'main.GitCommit=$commitHash'" -o build/goTES3MP-Linux src/*.go 