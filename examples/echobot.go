package main

import (
	"context"
	"log"
	"os"

	"github.com/deltachat/deltachat-rpc-client-go/deltachat"
	"github.com/deltachat/deltachat-rpc-client-go/deltachat/transport"
)

func logEvent(bot *deltachat.Bot, event deltachat.Event) {
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
	trans := transport.NewIOTransport()
	trans.Open()
	defer trans.Close()
	rpc := &deltachat.Rpc{Context: context.Background(), Transport: trans}

	sysinfo, _ := rpc.GetSystemInfo()
	log.Println("Running deltachat core", sysinfo["deltachat_core_version"])

	bot := deltachat.NewBot(rpc, 0)
	bot.On(deltachat.EventInfo{}, logEvent)
	bot.On(deltachat.EventWarning{}, logEvent)
	bot.On(deltachat.EventError{}, logEvent)
	bot.OnNewMsg(func(bot *deltachat.Bot, msgId deltachat.MsgId) {
		msg, _ := bot.Rpc.GetMessage(bot.AccountId, msgId)
		if msg.FromId > deltachat.ContactLastSpecial {
			bot.Rpc.MiscSendTextMessage(bot.AccountId, msg.ChatId, msg.Text)
		}
	})

	if !bot.IsConfigured() {
		log.Println("Bot not configured, configuring...")
		err := bot.Configure(os.Args[1], os.Args[2])
		if err != nil {
			log.Fatalln(err)
		}
	}

	addr, _ := bot.GetConfig("configured_addr")
	log.Println("Listening at:", addr.Unwrap())
	bot.Run()
}
