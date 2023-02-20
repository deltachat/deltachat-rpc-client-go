package deltachat

import "testing"
import "log"

func TestBasic(t *testing.T) {
	rpc := NewRpc()
	defer rpc.Stop()
	rpc.Start()

	manager := NewAccountManager(rpc)
	sysinfo, _ := manager.SystemInfo()
	if len(sysinfo["deltachat_core_version"].(string)) == 0 {
		t.Error("invalid deltachat_core_version")
	}

	bot := NewBotFromAccountManager(manager)
	bot.On(EVENT_INFO, func(event map[string]any) {
		log.Printf("%v: %v", event["type"], event["msg"])
	})
	bot.On(EVENT_WARNING, func(event map[string]any) {
		log.Printf("%v: %v", event["type"], event["msg"])
	})
	bot.On(EVENT_ERROR, func(event map[string]any) {
		log.Printf("%v: %v", event["type"], event["msg"])
	})
	bot.OnNewMsg(func(msg *Message) {
		snapshot, _ := msg.Snapshot()
		chat := snapshot["chat"].(*Chat)
		chat.SendText(snapshot["text"].(string))
	})

	if bot.IsConfigured() {
		t.Error("bot.IsConfigured() returning true, expected false")
	}
}
