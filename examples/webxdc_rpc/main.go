package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/deltachat/deltachat-rpc-client-go/deltachat"
	"github.com/deltachat/deltachat-rpc-client-go/deltachat/transport"
	"github.com/deltachat/deltachat-rpc-client-go/deltachat/xdcrpc"
)

func onEvent(bot *deltachat.Bot, accId deltachat.AccountId, event deltachat.Event) {
	switch ev := event.(type) {
	case deltachat.EventWebxdcStatusUpdate:
		rawUpdate, err := xdcrpc.GetUpdate(bot.Rpc, accId, ev.MsgId, ev.StatusUpdateSerial)
		if err != nil {
			return
		}
		if !xdcrpc.IsFromSelf(rawUpdate) {
			response := xdcrpc.HandleMessage(&API{}, rawUpdate)
			if response != nil {
				sendPayload(bot.Rpc, accId, ev.MsgId, response)
			}
		}
	}
}

// Send a WebXDC status update with the given payload
func sendPayload[T any](rpc *deltachat.Rpc, accId deltachat.AccountId, msgId deltachat.MsgId, payload T) error {
	data, err := json.Marshal(xdcrpc.StatusUpdate[T]{Payload: payload})
	if err != nil {
		return err
	}
	return rpc.SendWebxdcStatusUpdate(accId, msgId, string(data), "")
}

func main() {
	trans := transport.NewIOTransport()
	trans.Open()
	defer trans.Close()
	rpc := &deltachat.Rpc{Context: context.Background(), Transport: trans}

	bot := deltachat.NewBot(rpc)
	accId := deltachat.GetAccount(rpc)

	bot.OnUnhandledEvent(onEvent)
	bot.OnNewMsg(func(bot *deltachat.Bot, accId deltachat.AccountId, msgId deltachat.MsgId) {
		msg, _ := bot.Rpc.GetMessage(accId, msgId)
		if msg.FromId > deltachat.ContactLastSpecial {
			// TODO: send the webxdc app here
			return
		}
	})

	if isConf, _ := bot.Rpc.IsConfigured(accId); !isConf {
		log.Println("Bot not configured, configuring...")
		err := bot.Configure(accId, os.Args[1], os.Args[2])
		if err != nil {
			log.Fatalln(err)
		}
	}

	bot.Run()
}
