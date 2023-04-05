package deltachat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBot_NewBotFromAccountManager(t *testing.T) {
	t.Parallel()
	manager := acfactory.NewAcManager()
	defer manager.Rpc.Stop()

	bot := NewBotFromAccountManager(manager)
	assert.NotNil(t, bot)

	bot2 := NewBotFromAccountManager(manager)
	assert.NotNil(t, bot2)
	assert.Equal(t, bot.Account, bot2.Account)
}

func TestBot_NewBot(t *testing.T) {
	t.Parallel()
	manager := acfactory.NewAcManager()
	defer manager.Rpc.Stop()

	acc, err := manager.AddAccount()
	assert.Nil(t, err)
	bot := NewBot(acc)
	assert.NotNil(t, bot)
}

func TestBot_String(t *testing.T) {
	t.Parallel()
	acc := acfactory.GetUnconfiguredAccount()
	defer acc.Manager.Rpc.Stop()

	bot := NewBot(acc)
	assert.NotEmpty(t, bot.String())
}

func TestBot_OnNewMsg(t *testing.T) {
	t.Parallel()
	bot := acfactory.GetOnlineBot()
	defer bot.Account.Manager.Rpc.Stop()
	defer bot.Stop()

	acc1 := acfactory.GetOnlineAccount()
	defer acc1.Manager.Rpc.Stop()
	chatWithBot1, err := acc1.CreateChat(bot.Account)
	assert.Nil(t, err)

	incomingMsg := make(chan *MsgSnapshot)
	bot.On(EVENT_INCOMING_MSG, func(event *Event) {
		msg := &Message{bot.Account, event.MsgId}
		snapshot, _ := msg.Snapshot()
		incomingMsg <- snapshot
	})

	chatWithBot1.SendText("test1")
	msg := <-incomingMsg
	assert.Equal(t, "test1", msg.Text)
	bot.RemoveEventHandler(EVENT_INCOMING_MSG)
	close(incomingMsg)

	acc2 := acfactory.GetOnlineAccount()
	defer acc2.Manager.Rpc.Stop()
	chatWithBot2, err := acc2.CreateChat(bot.Account)
	assert.Nil(t, err)

	bot.OnNewMsg(func(msg *Message) {
		snapshot, _ := msg.Snapshot()
		chat := &Chat{bot.Account, snapshot.ChatId}
		chat.SendText(snapshot.Text)
	})

	chatWithBot2.SendText("test2")
	msg, _ = acfactory.GetNextMsg(acc2)
	assert.Equal(t, "test2", msg.Text)
}

func TestBot_processMessages(t *testing.T) {
	t.Parallel()
	bot := acfactory.GetOnlineBot()
	defer bot.Account.Manager.Rpc.Stop()

	bot.Account.Manager.Rpc.Stop()
	bot.processMessages()
}

func TestBot_Stop(t *testing.T) {
	t.Parallel()
	bot := acfactory.GetOnlineBot()
	defer bot.Account.Manager.Rpc.Stop()

	bot.Stop()
}

func TestBot_IsConfigured(t *testing.T) {
	t.Parallel()
	acc := acfactory.GetUnconfiguredAccount()
	defer acc.Manager.Rpc.Stop()

	bot := NewBot(acc)
	assert.False(t, bot.IsConfigured())

	assert.Nil(t, acc.Configure())

	assert.True(t, bot.IsConfigured())
}

func TestBot_UpdateConfig(t *testing.T) {
	t.Parallel()
	acc := acfactory.GetUnconfiguredAccount()
	defer acc.Manager.Rpc.Stop()

	bot := NewBot(acc)
	assert.Nil(t, bot.UpdateConfig(map[string]string{"selfstatus": "status"}))
}

func TestBot_SetConfig(t *testing.T) {
	t.Parallel()
	acc := acfactory.GetUnconfiguredAccount()
	defer acc.Manager.Rpc.Stop()

	bot := NewBot(acc)
	assert.Nil(t, bot.SetConfig("selfstatus", "testing"))
	val, err := bot.GetConfig("selfstatus")
	assert.Nil(t, err)
	assert.Equal(t, val, "testing")
}

func TestBot_Me(t *testing.T) {
	t.Parallel()
	acc := acfactory.GetUnconfiguredAccount()
	defer acc.Manager.Rpc.Stop()

	bot := NewBot(acc)
	assert.NotNil(t, bot.Me())
}

func TestBot_Run(t *testing.T) {
	t.Parallel()
	acc := acfactory.GetUnconfiguredAccount()
	defer acc.Manager.Rpc.Stop()

	bot := NewBot(acc)
	go func() { bot.Run() }()
	bot.Stop()
}
