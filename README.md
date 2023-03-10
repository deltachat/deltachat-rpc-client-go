# Delta Chat API for Go

![Latest release](https://img.shields.io/github/v/tag/deltachat/deltachat-rpc-client-go?label=release)
![Go version](https://img.shields.io/github/go-mod/go-version/deltachat/deltachat-rpc-client-go)
[![CI](https://github.com/deltachat/deltachat-rpc-client-go/actions/workflows/ci.yml/badge.svg)](https://github.com/deltachat/deltachat-rpc-client-go/actions/workflows/ci.yml)
![Coverage](https://img.shields.io/badge/Coverage-85.9%25-brightgreen)

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

```go
package main

import (
    "github.com/deltachat/deltachat-rpc-client-go/deltachat"
    "log"
    "os"
)

func main() {
    rpc := deltachat.NewRpcIO()
    defer rpc.Stop()
    rpc.Start()

    bot := deltachat.NewBotFromAccountManager(&deltachat.AccountManager{rpc})
    bot.OnNewMsg(func(msg *deltachat.Message) {
        snapshot, _ := msg.Snapshot()
        chat := deltachat.Chat{bot.Account, snapshot.ChatId}
        chat.SendText(snapshot.Text)
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
    bot.RunForever()
}
```

Save the previous code snippet as `echobot.go` then run:

```sh
go run ./echobot.go bot@example.com PASSWORD
```

Check the [examples folder](https://github.com/deltachat/deltachat-rpc-client-go/tree/master/examples)
for more examples.
