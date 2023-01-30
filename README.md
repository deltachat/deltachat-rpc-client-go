# Delta Chat RPC Go client

RPC client that connects to standalone Delta Chat RPC server `deltachat-rpc-server`


## Getting started

To use Delta Chat RPC client, first build: [deltachat-rpc-server](https://github.com/deltachat/deltachat-core-rust/tree/master/deltachat-rpc-server) with `cargo build -p deltachat-rpc-server`.
Install it anywhere in your `PATH`.

## Examples

To run the Echo-bot example:

```sh
# compile the bot:
go build ./examples/echobot.go
# configure and run the bot:
./echobot bot@example.com PASSWORD
```
