package main

import (
	"github.com/deltachat/deltachat-rpc-client-go/deltachat"
	"log"
	"os"
)

// Get the first available account or create a new one if none exists.
// A real client would need to provide account selection.
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

// Dummy function that just prints some events, here your client's UI would process the event
func handleEvent(event map[string]any) {
	switch evtype := event["type"].(string); evtype {
	case deltachat.EVENT_INFO:
		log.Println("INFO:", event["msg"])
	case deltachat.EVENT_WARNING:
		log.Println("WARNING:", event["msg"])
	case deltachat.EVENT_ERROR:
		log.Println("ERROR:", event["msg"])
	case deltachat.EVENT_INCOMING_MSG:
		log.Println("Got new message!")
	}
}

func main() {
	rpc := deltachat.NewRpc()
	rpc.Stderr = nil // disable printing logs from core RPC, do this if your client is a TUI
	defer rpc.Stop()
	rpc.Start() // start communication with Delta Chat core

	acc := getAccount(deltachat.NewAccountManager(rpc))

	if configured, _ := acc.IsConfigured(); configured {
		acc.StartIO()
	} else {
		log.Println("Account not configured, configuring...")
		acc.SetConfig("addr", os.Args[1])
		acc.SetConfig("mail_pw", os.Args[2])
		if err := acc.Configure(); err != nil {
			log.Fatalln(err)
		}
	}

	for {
		handleEvent(acc.WaitForEvent())
	}
}
