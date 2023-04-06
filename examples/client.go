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
func handleEvent(acc *deltachat.Account, event *deltachat.Event) {
	switch event.Type {
	case deltachat.EventInfo:
		log.Println("INFO:", event.Msg)
	case deltachat.EventWarning:
		log.Println("WARNING:", event.Msg)
	case deltachat.EventError:
		log.Println("ERROR:", event.Msg)
	case deltachat.EventIncomingMsg:
		msg := deltachat.Message{acc, event.MsgId}
		snapshot, _ := msg.Snapshot()
		log.Printf("Got new message from %v: %v", snapshot.Sender.DisplayName, snapshot.Text)
	}
}

func main() {
	rpc := deltachat.NewRpcIO()
	rpc.Stderr = nil // disable printing logs from core RPC, do this if your client is a TUI
	defer rpc.Stop()
	rpc.Start() // start communication with Delta Chat core

	acc := getAccount(&deltachat.AccountManager{rpc})

	if configured, _ := acc.IsConfigured(); configured {
		acc.StartIO()
	} else {
		log.Println("Account not configured, configuring...")
		acc.UpdateConfig(
			map[string]string{
				"addr":    os.Args[1],
				"mail_pw": os.Args[2],
			},
		)
		if err := acc.Configure(); err != nil {
			log.Fatalln(err)
		}
	}

	addr, _ := acc.GetConfig("addr")
	log.Println("Using account:", addr)

	eventChan := acc.GetEventChannel()
	for {
		event, ok := <-eventChan
		if !ok {
			break
		}
		handleEvent(acc, event)
	}
}
