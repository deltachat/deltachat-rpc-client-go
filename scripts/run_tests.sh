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
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2
fi

if ! golangci-lint run
then
    exit 1
fi

if ! command -v deltachat-rpc-server &> /dev/null
then
    echo "deltachat-rpc-server could not be found, installing..."
    pip3 install -U deltachat-rpc-server
fi

if ! command -v courtney &> /dev/null
then
    echo "courtney could not be found, installing..."
    go install github.com/dave/courtney@latest
fi

for i in ./examples/*.go
do
    echo "Testing examples: $i"
    if ! go build -v "$i"
    then
        exit 1
    fi
done
echobot_full="examples/echobot_full/"
echo "Testing examples: $echobot_full"
cd $echobot_full
if ! go test -v
then
    exit 1
fi
cd ../..
echo "Done testing examples"

courtney -v -t="./..." ${TEST_EXTRA_TAGS:--t="-parallel=1"}
go tool cover -func=coverage.out -o=coverage-percent.out
