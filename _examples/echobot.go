package main

import (
	"github.com/deltachat/deltachat-rpc-client-go/deltachat"
	"log"
	"os"
)

func logEvent(event map[string]any) {
	log.Printf("%v: %v", event["type"], event["msg"])
}

func main() {
	rpc := deltachat.NewRpcIO()
	defer rpc.Stop()
	rpc.Start()

	manager := &deltachat.AccountManager{rpc}
	sysinfo, _ := manager.SystemInfo()
	log.Println("Running deltachat core", sysinfo["deltachat_core_version"])

	bot := deltachat.NewBotFromAccountManager(manager)
	bot.On(deltachat.EVENT_INFO, logEvent)
	bot.On(deltachat.EVENT_WARNING, logEvent)
	bot.On(deltachat.EVENT_ERROR, logEvent)
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
