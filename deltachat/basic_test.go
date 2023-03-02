package deltachat

import "testing"
import "log"

func logEvent(event *Event) {
	log.Printf("%v: %v", event.Type, event.Msg)
}

func TestBasic(t *testing.T) {
	rpc := NewRpcIO()
	defer rpc.Stop()
	rpc.Start()

	manager := &AccountManager{rpc}
	sysinfo, _ := manager.SystemInfo()
	if len(sysinfo["deltachat_core_version"]) == 0 {
		t.Error("invalid deltachat_core_version")
	}

	bot := NewBotFromAccountManager(manager)
	bot.On(EVENT_INFO, logEvent)
	bot.On(EVENT_WARNING, logEvent)
	bot.On(EVENT_ERROR, logEvent)
	bot.OnNewMsg(func(msg *Message) {
		snapshot, _ := msg.Snapshot()
		chat := Chat{bot.Account, snapshot.ChatId}
		chat.SendText(snapshot.Text)
	})

	if bot.IsConfigured() {
		t.Error("bot.IsConfigured() returning true, expected false")
	}
}
