package main // replace with your package name

import (
	"testing"

	"github.com/deltachat/deltachat-rpc-client-go/acfactory"
	"github.com/stretchr/testify/assert"
)

func TestEchoBot(t *testing.T) {
	bot := acfactory.OnlineBot()
	defer acfactory.StopRpc(bot) // do this for every account/bot to release resources soon in your tests

	user := acfactory.OnlineAccount()
	defer acfactory.StopRpc(user)

	go runEchoBot(bot) // this is the function we are testing

	chatWithBot, err := acfactory.CreateChat(user, bot.Account)
	assert.Nil(t, err)

	chatWithBot.SendText("hi")
	msg, err := acfactory.NextMsg(user)
	assert.Nil(t, err)
	assert.Equal(t, "hi", msg.Text) // check that bot echoes back the "hi" message from user
}
