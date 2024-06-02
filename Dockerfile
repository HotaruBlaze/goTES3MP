# syntax=docker/dockerfile:1
FROM golang:1.22-alpine as BUILDER

RUN mkdir /app
WORKDIR /app

ARG BUILD_VERSION
ARG GITHUB_SHA

COPY ["go.mod", "go.sum", "./"]
COPY ["src/", "./src/"]
COPY ["tes3mp/scripts/custom/IrcBridge/IrcBridge.lua", "/app/tes3mp/scripts/custom/IrcBridge/IrcBridge.lua"]

RUN apk add --no-cache protoc

RUN cd src && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go && \
    go get google.golang.org/grpc/cmd/protoc-gen-go-grpc && \
    go generate && \
    go mod download && \
    go build -ldflags="-X 'main.Build=$BUILD_VERSION' -X 'main.GitCommit=$GITHUB_SHA'" -o /app/build/goTES3MP-Linux .


FROM golang:1.20-alpine as RUNNER

RUN mkdir /app
WORKDIR /app

COPY --from=BUILDER /app/build/goTES3MP-Linux /app/goTES3MP-Linux

RUN chmod +x /app/goTES3MP-Linux
ENTRYPOINT ["/app/goTES3MP-Linux"]
