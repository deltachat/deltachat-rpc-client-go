package deltachat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChat(t *testing.T) {
	acc, err := server.GetOnlineAccount()
	assert.Nil(t, err)

	contact, err := acc.CreateContact("test@example.com", "test")
	assert.Nil(t, err)
	assert.NotNil(t, contact)

	chat, err := contact.CreateChat()
	assert.Nil(t, err)

	assert.NotEmpty(t, chat.String())

	assert.Nil(t, chat.Accept())

	assert.Nil(t, chat.Block())

	assert.Nil(t, chat.MarkNoticed())

	assert.Nil(t, chat.Pin())
	assert.Nil(t, chat.Unpin())

	assert.Nil(t, chat.Archive())
	assert.Nil(t, chat.Unarchive())

	assert.Nil(t, chat.Delete())

	chat, err = acc.CreateGroup("test group", false)
	assert.Nil(t, err)
	assert.NotNil(t, chat)

	assert.Nil(t, chat.SetMuteDuration(-1))
	assert.Nil(t, chat.SetMuteDuration(100))
	assert.Nil(t, chat.SetMuteDuration(0))

	assert.Nil(t, chat.SetName("new name"))

	assert.Nil(t, chat.AddContact(contact))

	assert.Nil(t, chat.RemoveContact(contact))

	_, err = chat.Contacts()
	assert.Nil(t, err)

	assert.Nil(t, chat.SetEphemeralTimer(9000))

	_, err = chat.EphemeralTimer()
	assert.Nil(t, err)

	_, _, err = chat.QrCode()
	assert.Nil(t, err)

	_, err = chat.EncryptionInfo()
	assert.Nil(t, err)

	_, err = chat.SendText("test")
	assert.Nil(t, err)

	_, err = chat.SendMsg(MsgData{Text: "test message"})
	assert.Nil(t, err)

	msgs, err := chat.Messages(false, false)
	assert.Nil(t, err)

	_, err = chat.FreshMsgCount()
	assert.Nil(t, err)

	_, err = chat.SendVideoChatInvitation()
	assert.NotNil(t, err)
	chat.Account.SetConfig("webrtc_instance", "https://meet.jit.si")
	_, err = chat.SendVideoChatInvitation()
	assert.Nil(t, err)

	_, err = chat.FirstUnreadMsg()
	assert.Nil(t, err)

	_, err = chat.BasicSnapshot()
	assert.Nil(t, err)

	_, err = chat.FullSnapshot()
	assert.Nil(t, err)

	assert.Nil(t, chat.DeleteMsgs(msgs))
}
