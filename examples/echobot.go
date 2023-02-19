package main

import (
	"github.com/deltachat/deltachat-rpc-client-go/deltachat"
	"log"
	"os"
)

func getAccount(manager *deltachat.AccountManager) *deltachat.Account {
	accounts, _ := manager.Accounts()
	var acc *deltachat.Account
	if len(accounts) == 0 {
		acc, _ = manager.AddAccount()
	} else {
		acc = accounts[0]
	}
	return acc
}

func configure(acc *deltachat.Account) {
	acc.SetConfig("bot", "1")
	if configured, _ := acc.IsConfigured(); configured {
		log.Println("Account is already configured.")
		acc.StartIO()
	} else {
		log.Println("Account not configured, configuring...")
		acc.SetConfig("addr", os.Args[1])
		acc.SetConfig("mail_pw", os.Args[2])
		err := acc.Configure()
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Account configured.")
	}
}

func processMessages(acc *deltachat.Account) {
	msgs, _ := acc.FreshMsgsInArrivalOrder()
	for _, msg := range msgs {
		snapshot, _ := msg.Snapshot()
		if !snapshot["isInfo"].(bool) {
			chat := snapshot["chat"].(deltachat.Chat)
			chat.SendText(snapshot["text"].(string))
		}
		msg.MarkSeen()
	}
}

func main() {
	rpc := deltachat.NewRpc()
	defer rpc.Stop()
	rpc.Start()

	manager := deltachat.NewAccountManager(rpc)
	sysinfo, _ := manager.SystemInfo()
	log.Println("Running deltachat core", sysinfo["deltachat_core_version"])

	acc := getAccount(manager)
	configure(acc)
	addr, _ := acc.GetConfig("addr")
	log.Println("Listening at:", addr)

	// Process old messages.
	processMessages(acc)
	for {
		data := acc.WaitForEvent()
		if data == nil {
			return
		}
		switch evtype := data["type"].(string); evtype {
		case "Info":
			log.Println("INFO:", data["msg"])
		case "Warning":
			log.Println("WARNING:", data["msg"])
		case "Error":
			log.Println("ERROR:", data["msg"])
		case "IncomingMsg":
			processMessages(acc)
		}
	}
}
