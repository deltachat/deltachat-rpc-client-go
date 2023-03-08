package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	acc, err := server.GetOnlineAccount()
	assert.Nil(t, err)

	contact, err := acc.CreateContact("test@example.com", "test")
	assert.Nil(t, err)
	assert.NotNil(t, contact)

	chat, err := contact.CreateChat()
	assert.Nil(t, err)

	msg, err := chat.SendText("test")
	assert.Nil(t, err)

	assert.NotEmpty(t, msg.String())

	snapshot, err := msg.Snapshot()
	assert.Nil(t, err)
	assert.Equal(t, snapshot.Text, "test")

	_, err = msg.Html()
	assert.Nil(t, err)

	_, err = msg.Info()
	assert.Nil(t, err)

	assert.NotNil(t, msg.Download())

	assert.Nil(t, msg.MarkSeen())

	assert.Nil(t, msg.SendReaction(":)"))
}
