package main

import (
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

func runEchoBot(bot *deltachat.Bot) {
	sysinfo, _ := bot.Rpc.GetSystemInfo()
	log.Println("Running deltachat core", sysinfo["deltachat_core_version"])

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

func main() {
	trans := transport.NewProcessTransport()
	trans.Open()
	defer trans.Close()
	rpc := &deltachat.Rpc{Transport: trans}
	runEchoBot(deltachat.NewBot(rpc, 0))
}
