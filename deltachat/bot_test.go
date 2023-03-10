package deltachat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBot_NewBotFromAccountManager(t *testing.T) {
	bot := NewBotFromAccountManager(server.AccountManager())
	assert.NotNil(t, bot)
}

func TestBot_NewBot(t *testing.T) {
	acc, err := server.AccountManager().AddAccount()
	assert.Nil(t, err)
	bot := NewBot(acc)
	assert.NotNil(t, bot)
}

func TestBot_String(t *testing.T) {
	bot, err := server.GetOnlineBot()
	assert.Nil(t, err)
	defer bot.Stop()

	assert.NotEmpty(t, bot.String())
}

func TestBot_OnNewMsg(t *testing.T) {
	bot, err := server.GetOnlineBot()
	assert.Nil(t, err)
	defer bot.Stop()

	onCalled := false
	bot.On(EVENT_INFO, func(event *Event) {
		onCalled = true
	})
	bot.OnNewMsg(func(msg *Message) {
		snapshot, _ := msg.Snapshot()
		chat := Chat{bot.Account, snapshot.ChatId}
		chat.SendText(snapshot.Text)
	})

	acc, err := server.GetOnlineAccount()
	assert.Nil(t, err)
	chatWithBot, err := acc.CreateChat(bot.Account)
	assert.Nil(t, err)

	chatWithBot.SendText("test")
	msg, err := server.GetNextMsg(acc)
	assert.Equal(t, msg.Text, "test")
	assert.True(t, onCalled)
}

func TestBot_IsConfigured(t *testing.T) {
	acc := server.GetUnconfiguredAccount()
	bot := NewBot(acc)

	assert.False(t, bot.IsConfigured())

	err := acc.Configure()
	assert.Nil(t, err)

	assert.True(t, bot.IsConfigured())
}

func TestBot_UpdateConfig(t *testing.T) {
	bot := NewBot(server.GetUnconfiguredAccount())

	assert.Nil(t, bot.UpdateConfig(map[string]string{"selfstatus": "status"}))
}

func TestBot_SetConfig(t *testing.T) {
	bot := NewBot(server.GetUnconfiguredAccount())

	assert.Nil(t, bot.SetConfig("selfstatus", "testing"))
	val, err := bot.GetConfig("selfstatus")
	assert.Nil(t, err)
	assert.Equal(t, val, "testing")
}

func TestBot_Me(t *testing.T) {
	bot := NewBot(server.GetUnconfiguredAccount())

	assert.NotNil(t, bot.Me())
}
