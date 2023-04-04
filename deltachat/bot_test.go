package deltachat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBot_NewBotFromAccountManager(t *testing.T) {
	manager := acfactory.NewAcManager()
	defer manager.Rpc.Stop()

	bot := NewBotFromAccountManager(manager)
	assert.NotNil(t, bot)
}

func TestBot_NewBot(t *testing.T) {
	manager := acfactory.NewAcManager()
	defer manager.Rpc.Stop()

	acc, err := manager.AddAccount()
	assert.Nil(t, err)
	bot := NewBot(acc)
	assert.NotNil(t, bot)
}

func TestBot_String(t *testing.T) {
	acc := acfactory.GetUnconfiguredAccount()
	defer acc.Manager.Rpc.Stop()

	bot := NewBot(acc)
	assert.NotEmpty(t, bot.String())
}

func TestBot_OnNewMsg(t *testing.T) {
	bot := acfactory.GetOnlineBot()
	defer bot.Account.Manager.Rpc.Stop()
	defer bot.Stop()

	onCalled := false
	bot.On(EVENT_INFO, func(event *Event) {
		onCalled = true
	})
	bot.OnNewMsg(func(msg *Message) {
		snapshot, _ := msg.Snapshot()
		chat := &Chat{bot.Account, snapshot.ChatId}
		chat.SendText(snapshot.Text)
	})

	acc := acfactory.GetOnlineAccount()
	defer acc.Manager.Rpc.Stop()

	chatWithBot, err := acc.CreateChat(bot.Account)
	assert.Nil(t, err)

	chatWithBot.SendText("test")
	msg, _ := acfactory.GetNextMsg(acc)
	assert.Equal(t, msg.Text, "test")
	assert.True(t, onCalled)
}

func TestBot_IsConfigured(t *testing.T) {
	acc := acfactory.GetUnconfiguredAccount()
	defer acc.Manager.Rpc.Stop()

	bot := NewBot(acc)
	assert.False(t, bot.IsConfigured())

	assert.Nil(t, acc.Configure())

	assert.True(t, bot.IsConfigured())
}

func TestBot_UpdateConfig(t *testing.T) {
	acc := acfactory.GetUnconfiguredAccount()
	defer acc.Manager.Rpc.Stop()

	bot := NewBot(acc)
	assert.Nil(t, bot.UpdateConfig(map[string]string{"selfstatus": "status"}))
}

func TestBot_SetConfig(t *testing.T) {
	acc := acfactory.GetUnconfiguredAccount()
	defer acc.Manager.Rpc.Stop()

	bot := NewBot(acc)
	assert.Nil(t, bot.SetConfig("selfstatus", "testing"))
	val, err := bot.GetConfig("selfstatus")
	assert.Nil(t, err)
	assert.Equal(t, val, "testing")
}

func TestBot_Me(t *testing.T) {
	acc := acfactory.GetUnconfiguredAccount()
	defer acc.Manager.Rpc.Stop()

	bot := NewBot(acc)
	assert.NotNil(t, bot.Me())
}
