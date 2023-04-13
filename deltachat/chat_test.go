package deltachat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChat_String(t *testing.T) {
	t.Parallel()
	acc := acfactory.OnlineAccount()
	defer acfactory.StopRpc(acc)

	chat, err := acc.Me().CreateChat()
	assert.Nil(t, err)

	assert.NotEmpty(t, chat.String())
}

func TestChat_Unpin(t *testing.T) {
	t.Parallel()
	acc := acfactory.OnlineAccount()
	defer acfactory.StopRpc(acc)

	chat, err := acc.Me().CreateChat()
	assert.Nil(t, err)

	assert.Nil(t, chat.Pin())
	assert.Nil(t, chat.Unpin())
}

func TestChat_Unarchive(t *testing.T) {
	t.Parallel()
	acc := acfactory.OnlineAccount()
	defer acfactory.StopRpc(acc)

	chat, err := acc.Me().CreateChat()
	assert.Nil(t, err)

	assert.Nil(t, chat.Archive())
	assert.Nil(t, chat.Unarchive())
}

func TestChat_Basics(t *testing.T) {
	t.Parallel()
	acc := acfactory.OnlineAccount()
	defer acfactory.StopRpc(acc)

	chat, err := acc.Me().CreateChat()
	assert.Nil(t, err)

	assert.Nil(t, chat.Accept())

	assert.Nil(t, chat.MarkNoticed())

	_, err = chat.FirstUnreadMsg()
	assert.Nil(t, err)

	_, err = chat.BasicSnapshot()
	assert.Nil(t, err)

	_, err = chat.FullSnapshot()
	assert.Nil(t, err)

	assert.Nil(t, chat.Block())

	assert.Nil(t, chat.Delete())
}

func TestChat_Groups(t *testing.T) {
	t.Parallel()
	acc := acfactory.OnlineAccount()
	defer acfactory.StopRpc(acc)

	contact, err := acc.CreateContact("null@localhost", "test")
	assert.Nil(t, err)
	assert.NotNil(t, contact)

	chat, err := contact.CreateChat()
	assert.Nil(t, err)

	chat, err = acc.CreateGroup("test group", false)
	assert.Nil(t, err)
	assert.NotNil(t, chat)

	assert.Nil(t, chat.SetImage(acfactory.TestImage()))
	assert.Nil(t, chat.RemoveImage())

	assert.Nil(t, chat.SetMuteDuration(-1))
	assert.Nil(t, chat.SetMuteDuration(100))
	assert.Nil(t, chat.SetMuteDuration(0))

	assert.Nil(t, chat.AddContact(contact))

	assert.Nil(t, chat.RemoveContact(contact))

	_, err = chat.Contacts()
	assert.Nil(t, err)

	assert.Nil(t, chat.SetEphemeralTimer(9000))
	acfactory.WaitForEventInChat(acc, EventChatEphemeralTimerModified{}, chat.Id)

	_, err = chat.EphemeralTimer()
	assert.Nil(t, err)

	_, _, err = chat.QrCode()
	assert.Nil(t, err)

	_, err = chat.EncryptionInfo()
	assert.Nil(t, err)

	_, err = chat.SendText("test")
	assert.Nil(t, err)

	msg, err := chat.SendMsg(MsgData{Text: "test message"})
	assert.Nil(t, err)
	assert.NotNil(t, msg)

	msgs, err := chat.Messages(false, false)
	assert.Nil(t, err)
	assert.NotEmpty(t, msgs)

	results, err := chat.SearchMessages("test message")
	assert.Nil(t, err)
	assert.NotEmpty(t, results)
	assert.Equal(t, results[0].Id, msg.Id)

	_, err = chat.FreshMsgCount()
	assert.Nil(t, err)

	url := "https://test.example.com"
	chat.Account.SetConfig("webrtc_instance", url)
	msg, err = chat.SendVideoChatInvitation()
	assert.Nil(t, err)
	msgData, err := msg.Snapshot()
	assert.Nil(t, err)
	assert.Contains(t, msgData.Text, url)

	assert.Nil(t, chat.DeleteMsgs(msgs))
}

func TestChat_SetName(t *testing.T) {
	t.Parallel()
	acc := acfactory.OnlineAccount()
	defer acfactory.StopRpc(acc)

	chat, err := acc.CreateGroup("test group", false)
	assert.Nil(t, err)
	assert.NotNil(t, chat)

	assert.Nil(t, chat.SetName("new name"))
	assert.Nil(t, chat.Leave())
	assert.NotNil(t, chat.SetName("another name"))
	acfactory.WaitForEvent(acc, EventErrorSelfNotInGroup{})
}
