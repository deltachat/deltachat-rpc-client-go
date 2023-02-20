# Delta Chat API for Go [![CI](https://github.com/deltachat/deltachat-rpc-client-go/actions/workflows/ci.yml/badge.svg)](https://github.com/deltachat/deltachat-rpc-client-go/actions/workflows/ci.yml)

Delta Chat client & bot API for Golang.

## Install

```sh
go get -u github.com/deltachat/deltachat-rpc-client-go
```

### Installing deltachat-rpc-server

This package depends on a standalone Delta Chat RPC server `deltachat-rpc-server` program that must be
available in your `PATH`. To install it run:

```sh
cargo install --git https://github.com/deltachat/deltachat-core-rust/ deltachat-rpc-server
```

For more info check:
https://github.com/deltachat/deltachat-core-rust/tree/master/deltachat-rpc-server

## Usage

Example echo-bot that will echo back any text message you send to it:

```go
package main

import (
    "github.com/deltachat/deltachat-rpc-client-go/deltachat"
    "log"
    "os"
)

func main() {
    rpc := deltachat.NewRpc()
    defer rpc.Stop()
    rpc.Start()

    bot := deltachat.NewBotFromAccountManager(deltachat.NewAccountManager(rpc))
    bot.OnNewMsg(func(msg *deltachat.Message) {
        snapshot, _ := msg.Snapshot()
        chat := snapshot["chat"].(*deltachat.Chat)
        chat.SendText(snapshot["text"].(string))
    })

    if !bot.IsConfigured() {
        log.Println("Bot not configured, configuring...")
        err := bot.Configure(os.Args[1], os.Args[2])
        if err != nil {
            log.Fatalln(err)
        }
    }

    addr, _ := bot.GetConfig("addr")
    log.Println("Listening at:", addr)
    bot.Run()
}
```

Save that example as `echobot.go` then run:

```sh
go run ./examples/echobot.go bot@example.com PASSWORD
```

Check the [examples folder](https://github.com/deltachat/deltachat-rpc-client-go/tree/master/examples)
for more examples.
