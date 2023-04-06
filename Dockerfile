# syntax=docker/dockerfile:1
FROM golang:1.20-alpine as BUILDER

RUN mkdir /app
WORKDIR /app

ARG BUILD_VERSION
ARG GITHUB_SHA

COPY ["go.mod", "go.sum", "./"]
COPY ["src/", "./src/"]

RUN go mod download && \
    go build -ldflags="-X 'main.Build=$BUILD_VERSION' -X 'main.GitCommit=$GITHUB_SHA'" -o /app/build/goTES3MP-Linux src/*.go


FROM golang:1.20-alpine as RUNNER

RUN mkdir /app
WORKDIR /app

COPY --from=BUILDER /app/build/goTES3MP-Linux /app/goTES3MP-Linux

RUN chmod +x /app/goTES3MP-Linux
ENTRYPOINT ["/app/goTES3MP-Linux"]
