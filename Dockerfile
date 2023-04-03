# syntax=docker/dockerfile:1
FROM golang:1.20-alpine

RUN mkdir /app
WORKDIR /app

ARG BUILD_VERSION
ARG GITHUB_SHA

COPY ["go.mod", "go.sum", "./"]
COPY ["src/", "./src/"]

RUN go mod download && \
    go build -ldflags="-X 'main.Build=$BUILD_VERSION' -X 'main.GitCommit=$GITHUB_SHA'" -o build/goTES3MP-Linux src/*.go