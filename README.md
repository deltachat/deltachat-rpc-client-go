# Delta Chat API for Go

Delta Chat client API for Go programming language.

## Installing Dependencies

This package depends on a standalone Delta Chat RPC server `deltachat-rpc-server` program that must be
available in your `PATH`. To install it check:
https://github.com/deltachat/deltachat-core-rust/tree/master/deltachat-rpc-server

## Usage

To see how to use this module, check the examples folder. To run the Echo-bot example:

```sh
# configure and run the bot:
go run ./examples/echobot.go bot@example.com PASSWORD
```
