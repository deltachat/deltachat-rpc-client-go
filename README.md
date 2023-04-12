# Delta Chat API for Go

![Latest release](https://img.shields.io/github/v/tag/deltachat/deltachat-rpc-client-go?label=release)
[![Go Reference](https://pkg.go.dev/badge/github.com/deltachat/deltachat-rpc-client-go.svg)](https://pkg.go.dev/github.com/deltachat/deltachat-rpc-client-go)
[![CI](https://github.com/deltachat/deltachat-rpc-client-go/actions/workflows/ci.yml/badge.svg)](https://github.com/deltachat/deltachat-rpc-client-go/actions/workflows/ci.yml)
![Coverage](https://img.shields.io/badge/Coverage-97.0%25-brightgreen)
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

```go
package main

import (
	"log"
	"os"

	"github.com/deltachat/deltachat-rpc-client-go/deltachat"
)

func logEvent(event deltachat.Event) {
	switch ev := event.(type) {
	case deltachat.EventInfo:
		log.Printf("INFO: %v", ev.Msg)
	case deltachat.EventWarning:
		log.Printf("WARNING: %v", ev.Msg)
	case deltachat.EventError:
		log.Printf("ERROR: %v", ev.Msg)
	}
}

func main() {
	rpc := deltachat.NewRpcIO()
	rpc.Start()
	defer rpc.Stop()

	manager := &deltachat.AccountManager{rpc}
	sysinfo, _ := manager.SystemInfo()
	log.Println("Running deltachat core", sysinfo["deltachat_core_version"])

	bot := deltachat.NewBotFromAccountManager(manager)
	bot.On(deltachat.EventInfo{}, logEvent)
	bot.On(deltachat.EventWarning{}, logEvent)
	bot.On(deltachat.EventError{}, logEvent)
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

	addr, _ := bot.GetConfig("configured_addr")
	log.Println("Listening at:", addr)
	bot.Run()
}
```

Save the previous code snippet as `echobot.go` then run:

```sh
go run ./echobot.go bot@example.com PASSWORD
```

Check the [examples folder](https://github.com/deltachat/deltachat-rpc-client-go/tree/master/examples)
for more examples.

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

```go
package yourpackage

import (
	"testing"

	"github.com/deltachat/deltachat-rpc-client-go/acfactory"
)

func TestMain(m *testing.M) {
	// cfg is the non-standard configuration of our fake mail server
	cfg := map[string]string{
		"mail_server":   "localhost",
		"send_server":   "localhost",
		"mail_port":     "3143",
		"send_port":     "3025",
		"mail_security": "3",
		"send_security": "3",
	}
	acfactory.TearUp(cfg)
	defer acfactory.TearDown()
	m.Run()
}
```

Now in your other test files you can do:

```go
package mypackage

import (
	"testing"

	"github.com/deltachat/deltachat-rpc-client-go/acfactory"
	"github.com/stretchr/testify/assert"
)


func TestSomething(t *testing.T) {
	botAcc := acfactory.OnlineAccount()
	defer botAcc.Manager.Rpc.Stop() // do this for every account to release resources soon in your tests

	user := acfactory.OnlineAccount()
	defer user.Manager.Rpc.Stop()

    bot := MyEchoBot(botAcc) // MyEchoBot is supposedly an echo bot implemented by you
    go bot.Run()
    defer bot.Stop()

    chatWithBot := acfactory.CreateChat(user, botAcc)

    chatWithBot.SendText("hi")
    msg, err = acfactory.NextMsg(user)
    assert.Nil(t, err)
    assert.Equal(t, "hi", msg.Text) // check that bot echoes back the "hi" message from user
}
```

### GitHub action

To run the tests in a GitHub action with the fake mail server service:

```yaml
name: Test

on:
  push:
    branches: [ "master" ]
  pull_request:
    branches: [ "master" ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v3
        with:
          go-version: 1.19

      - name: Run tests
        run: |
          go test -v ./...

    services:
      mail_server:
        image: ghcr.io/deltachat/mail-server-tester:release
        ports:
          - 3025:25
          - 3143:143
          - 3465:465
          - 3993:993
```

## Contributing

Pull requests are welcome! check [CONTRIBUTING.md](https://github.com/deltachat/deltachat-rpc-client-go/blob/master/CONTRIBUTING.md)
