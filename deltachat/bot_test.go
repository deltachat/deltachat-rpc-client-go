package deltachat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBot_NewBotFromAccountManager(t *testing.T) {
	t.Parallel()
	manager := acfactory.NewAcManager()
	defer acfactory.StopRpc(manager)

	bot := NewBotFromAccountManager(manager)
	assert.NotNil(t, bot)

	bot2 := NewBotFromAccountManager(manager)
	assert.NotNil(t, bot2)
	assert.Equal(t, bot.Account, bot2.Account)
}

func TestBot_NewBot(t *testing.T) {
	t.Parallel()
	manager := acfactory.NewAcManager()
	defer acfactory.StopRpc(manager)

	acc, err := manager.AddAccount()
	assert.Nil(t, err)
	bot := NewBot(acc)
	assert.NotNil(t, bot)
}

func TestBot_String(t *testing.T) {
	t.Parallel()
	acc := acfactory.UnconfiguredAccount()
	defer acfactory.StopRpc(acc)

	bot := NewBot(acc)
	assert.NotEmpty(t, bot.String())
}

func TestBot_OnNewMsg(t *testing.T) {
	t.Parallel()
	bot := acfactory.RunningBot()
	defer acfactory.StopRpc(bot)

	acc1 := acfactory.OnlineAccount()
	defer acfactory.StopRpc(acc1)
	chatWithBot1, err := acfactory.CreateChat(acc1, bot.Account)
	assert.Nil(t, err)

	incomingMsg := make(chan *MsgSnapshot)
	bot.On(EventIncomingMsg{}, func(event Event) {
		ev := event.(EventIncomingMsg)
		msg := &Message{bot.Account, ev.MsgId}
		snapshot, _ := msg.Snapshot()
		incomingMsg <- snapshot
	})

	_, err = chatWithBot1.SendText("test1")
	assert.Nil(t, err)
	msg := <-incomingMsg
	assert.Equal(t, "test1", msg.Text)
	bot.RemoveEventHandler(EventIncomingMsg{})
	close(incomingMsg)

	acc2 := acfactory.OnlineAccount()
	defer acfactory.StopRpc(acc2)
	chatWithBot2, err := acfactory.CreateChat(acc2, bot.Account)
	assert.Nil(t, err)

	bot.OnNewMsg(func(msg *Message) {
		snapshot, _ := msg.Snapshot()
		chat := &Chat{bot.Account, snapshot.ChatId}
		_, err := chat.SendText(snapshot.Text)
		assert.Nil(t, err)
	})

	_, err = chatWithBot2.SendText("test2")
	assert.Nil(t, err)
	msg, _ = acfactory.NextMsg(acc2)
	assert.Equal(t, "test2", msg.Text)
}

func TestBot_processMessages(t *testing.T) {
	t.Parallel()
	bot := acfactory.RunningBot()
	acfactory.StopRpc(bot)
	bot.processMessages()
}

func TestBot_Stop(t *testing.T) {
	t.Parallel()
	acc := acfactory.UnconfiguredAccount()
	defer acfactory.StopRpc(acc)

	done := make(chan error)
	bot := NewBot(acc)

	bot.Stop()
	go func() {
		done <- bot.Run()
	}()
	for {
		if bot.IsRunning() {
			break
		}
	}
	bot.Stop()
	assert.Nil(t, <-done)
}

func TestBot_IsConfigured(t *testing.T) {
	t.Parallel()
	acc := acfactory.UnconfiguredAccount()
	defer acfactory.StopRpc(acc)

	bot := NewBot(acc)
	assert.False(t, bot.IsConfigured())

	assert.Nil(t, acc.Configure())

	assert.True(t, bot.IsConfigured())
}

func TestBot_UpdateConfig(t *testing.T) {
	t.Parallel()
	acc := acfactory.UnconfiguredAccount()
	defer acfactory.StopRpc(acc)

	bot := NewBot(acc)
	assert.Nil(t, bot.UpdateConfig(map[string]string{"selfstatus": "status"}))
}

func TestBot_SetConfig(t *testing.T) {
	t.Parallel()
	acc := acfactory.UnconfiguredAccount()
	defer acfactory.StopRpc(acc)

	bot := NewBot(acc)
	assert.Nil(t, bot.SetConfig("selfstatus", "testing"))
	val, err := bot.GetConfig("selfstatus")
	assert.Nil(t, err)
	assert.Equal(t, val, "testing")
}

func TestBot_SetUiConfig(t *testing.T) {
	t.Parallel()
	acc := acfactory.UnconfiguredAccount()
	defer acfactory.StopRpc(acc)

	bot := NewBot(acc)
	assert.Nil(t, bot.SetUiConfig("testkey", "testing"))
	val, err := bot.GetUiConfig("testkey")
	assert.Nil(t, err)
	assert.Equal(t, val, "testing")

	val, err = bot.GetUiConfig("unknown-key")
	assert.Nil(t, err)
	assert.Empty(t, val)
}

func TestBot_Me(t *testing.T) {
	t.Parallel()
	acc := acfactory.UnconfiguredAccount()
	defer acfactory.StopRpc(acc)

	bot := NewBot(acc)
	assert.NotNil(t, bot.Me())
}

func TestBot_Run(t *testing.T) {
	t.Parallel()
	bot := acfactory.RunningBot()
	defer acfactory.StopRpc(bot)

	assert.NotNil(t, bot.Run())
}
