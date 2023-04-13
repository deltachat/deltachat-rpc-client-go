# Delta Chat API for Go

![Latest release](https://img.shields.io/github/v/tag/deltachat/deltachat-rpc-client-go?label=release)
[![Go Reference](https://pkg.go.dev/badge/github.com/deltachat/deltachat-rpc-client-go.svg)](https://pkg.go.dev/github.com/deltachat/deltachat-rpc-client-go)
[![CI](https://github.com/deltachat/deltachat-rpc-client-go/actions/workflows/ci.yml/badge.svg)](https://github.com/deltachat/deltachat-rpc-client-go/actions/workflows/ci.yml)
![Coverage](https://img.shields.io/badge/Coverage-99.2%25-brightgreen)
[![Go Report Card](https://goreportcard.com/badge/github.com/deltachat/deltachat-rpc-client-go)](https://goreportcard.com/report/github.com/deltachat/deltachat-rpc-client-go)

Delta Chat client & bot API for Golang.

## Install

```sh
go get -u github.com/deltachat/deltachat-rpc-client-go
```

### Installing deltachat-rpc-server

This package depends on a standalone Delta Chat RPC server `deltachat-rpc-server` program that must be
available in your `PATH`. For installation instructions check:
https://github.com/deltachat/deltachat-core-rust/tree/master/deltachat-rpc-server

## Usage

Example echo-bot that will echo back any text message you send to it:

<!-- MARKDOWN-AUTO-DOCS:START (CODE:src=./examples/echobot_full/echobot.go) -->
<!-- MARKDOWN-AUTO-DOCS:END -->

Save the previous code snippet as `echobot.go` then run:

```sh
go run ./echobot.go bot@example.com PASSWORD
```

Check the [examples folder](https://github.com/deltachat/deltachat-rpc-client-go/tree/master/examples)
for more examples.

## Developing bots faster ⚡

If what you want is to develop bots, you probably should use this library together with
[deltabot-cli-go](https://github.com/deltachat-bot/deltabot-cli-go/), a library that takes away the
repetitive process of creating the bot CLI and let you focus on writing your message processing logic.

## Testing your code

The `acfactory` package is provided to help users of this library to unit-test their code.

### Local mail server

You need to have a local fake email server running. The easiest way to do that is with Docker:

```
$ docker pull ghcr.io/deltachat/mail-server-tester:release
$ docker run -it --rm -p 3025:25 -p 3110:110 -p 3143:143 -p 3465:465 -p 3993:993 ghcr.io/deltachat/mail-server-tester
```

### Using acfactory

After setting up the fake email server, create a file called `main_test.go` inside your tests folder,
and save it with the following content:

<!-- MARKDOWN-AUTO-DOCS:START (CODE:src=./examples/echobot_full/main_test.go) -->
<!-- MARKDOWN-AUTO-DOCS:END -->

Now in your other test files you can do:

<!-- MARKDOWN-AUTO-DOCS:START (CODE:src=./examples/echobot_full/echobot_test.go) -->
<!-- MARKDOWN-AUTO-DOCS:END -->

### GitHub action

To run the tests in a GitHub action with the fake mail server service:

<!-- MARKDOWN-AUTO-DOCS:START (CODE:src=./examples/echobot_full/.github/workflows/ci.yml) -->
<!-- MARKDOWN-AUTO-DOCS:END -->

Check the complete example at [examples/echobot_full](https://github.com/deltachat/deltachat-rpc-client-go/tree/master/examples/echobot_full)

## Contributing

Pull requests are welcome! check [CONTRIBUTING.md](https://github.com/deltachat/deltachat-rpc-client-go/blob/master/CONTRIBUTING.md)
