package deltachat

import (
	"testing"

	"github.com/deltachat/deltachat-rpc-client-go/deltachat/option"
	"github.com/stretchr/testify/assert"
)

func TestBot_NewBot(t *testing.T) {
	t.Parallel()
	acfactory.WithRpc(func(rpc *Rpc) {
		accId, err := rpc.AddAccount()
		assert.Nil(t, err)
		bot := NewBot(rpc, accId)
		assert.NotNil(t, bot)
	})

	acfactory.WithRpc(func(rpc *Rpc) {
		bot := NewBot(rpc, 0)
		assert.NotNil(t, bot)

		bot2 := NewBot(rpc, 0)
		assert.NotNil(t, bot2)
		assert.NotEqual(t, bot.AccountId, 0)
		assert.Equal(t, bot.AccountId, bot2.AccountId)
	})
}

func TestBot_On(t *testing.T) {
	t.Parallel()
	acfactory.WithRunningBot(func(bot *Bot) {
		acfactory.WithOnlineAccount(func(accRpc *Rpc, accId AccountId) {
			incomingMsg := make(chan *MsgSnapshot)
			bot.On(EventIncomingMsg{}, func(bot *Bot, event Event) {
				ev := event.(EventIncomingMsg)
				snapshot, _ := bot.Rpc.GetMessage(bot.AccountId, ev.MsgId)
				incomingMsg <- snapshot
			})

			chatWithBot := acfactory.CreateChat(accRpc, accId, bot.Rpc, bot.AccountId)
			_, err := accRpc.MiscSendTextMessage(accId, chatWithBot, "test1")
			assert.Nil(t, err)
			msg := <-incomingMsg
			assert.Equal(t, "test1", msg.Text)
			bot.RemoveEventHandler(EventIncomingMsg{})
			close(incomingMsg)
		})
	})
}

func TestBot_OnNewMsg(t *testing.T) {
	t.Parallel()
	acfactory.WithRunningBot(func(bot *Bot) {
		acfactory.WithOnlineAccount(func(accRpc *Rpc, accId AccountId) {
			bot.OnNewMsg(func(bot *Bot, msgId MsgId) {
				snapshot, _ := bot.Rpc.GetMessage(bot.AccountId, msgId)
				_, err := bot.Rpc.MiscSendTextMessage(bot.AccountId, snapshot.ChatId, snapshot.Text)
				assert.Nil(t, err)
			})

			chatWithBot := acfactory.CreateChat(accRpc, accId, bot.Rpc, bot.AccountId)
			_, err := accRpc.MiscSendTextMessage(accId, chatWithBot, "test2")
			assert.Nil(t, err)
			msg := acfactory.NextMsg(accRpc, accId)
			assert.Equal(t, "test2", msg.Text)
		})
	})
}

func TestBot_processMessages(t *testing.T) {
	t.Parallel()
	acfactory.WithRunningBot(func(bot *Bot) {
		bot.processMessages()
	})
}

func TestBot_Stop(t *testing.T) {
	t.Parallel()
	acfactory.WithUnconfiguredBot(func(bot *Bot) {
		bot.Stop()
		done := make(chan error)
		go func() {
			done <- bot.Run()
		}()
		bot.Stop()
		assert.Nil(t, <-done)
	})
}

func TestBot_IsConfigured(t *testing.T) {
	t.Parallel()
	acfactory.WithUnconfiguredBot(func(bot *Bot) {
		assert.False(t, bot.IsConfigured())
		assert.Nil(t, bot.Rpc.Configure(bot.AccountId))
		assert.True(t, bot.IsConfigured())
	})
}

func TestBot_UpdateConfig(t *testing.T) {
	t.Parallel()
	acfactory.WithUnconfiguredBot(func(bot *Bot) {
		assert.Nil(t, bot.UpdateConfig(map[string]option.Option[string]{"selfstatus": option.Some("status")}))
	})
}

func TestBot_SetConfig(t *testing.T) {
	t.Parallel()
	acfactory.WithUnconfiguredBot(func(bot *Bot) {
		assert.Nil(t, bot.SetConfig("selfstatus", option.Some("testing")))
		val, err := bot.GetConfig("selfstatus")
		assert.Nil(t, err)
		assert.Equal(t, val.Unwrap(), "testing")
	})
}

func TestBot_SetUiConfig(t *testing.T) {
	t.Parallel()
	acfactory.WithUnconfiguredBot(func(bot *Bot) {
		assert.Nil(t, bot.SetUiConfig("testkey", option.Some("testing")))
		val, err := bot.GetUiConfig("testkey")
		assert.Nil(t, err)
		assert.Equal(t, val.Unwrap(), "testing")

		val, err = bot.GetUiConfig("unknown-key")
		assert.Nil(t, err)
		assert.Empty(t, val.UnwrapOr(""))
	})
}
