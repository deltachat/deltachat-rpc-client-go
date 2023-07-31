#!/bin/env bash

OUTPUT=`gofmt -d .`
if [ -n "$OUTPUT" ]
then
    echo "$OUTPUT"
    exit 1
fi

if ! command -v golangci-lint &> /dev/null
then
    # binary will be $(go env GOPATH)/bin/golangci-lint
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.52.2
fi

if ! golangci-lint run
then
    exit 1
fi

if ! command -v deltachat-rpc-server &> /dev/null
then
    echo "deltachat-rpc-server could not be found, installing..."
    curl -L "https://github.com/deltachat/deltachat-core-rust/releases/download/v1.118.0/deltachat-rpc-server-x86_64" -o ~/.cargo/bin/deltachat-rpc-server
    chmod +x ~/.cargo/bin/deltachat-rpc-server

fi

if ! command -v courtney &> /dev/null
then
    echo "courtney could not be found, installing..."
    go install github.com/dave/courtney@latest
fi

# test examples
for i in ./examples/*.go
do
    if ! go build -v "$i"
    then
        exit 1
    fi
done
cd examples/echobot_full/
if ! go test -v
then
    exit 1
fi
cd ../..

courtney -v -t="./..." ${TEST_EXTRA_TAGS:--t="-parallel=1"}
go tool cover -func=coverage.out -o=coverage-percent.out
