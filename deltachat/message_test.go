package deltachat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	t.Parallel()
	acc := acfactory.OnlineAccount()
	defer acfactory.StopRpc(acc)

	chat, err := acc.Me().CreateChat()
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
	acfactory.WaitForEventInChat(msg.Account, EventReactionsChanged{}, chat.Id)

	assert.Nil(t, msg.Delete())
}

func TestMessage_WebxdcInfo(t *testing.T) {
	t.Parallel()
	acc := acfactory.OnlineAccount()
	defer acfactory.StopRpc(acc)

	chat, err := acc.Me().CreateChat()
	assert.Nil(t, err)

	msg, err := chat.SendText("test")
	assert.Nil(t, err)
	_, err = msg.WebxdcInfo()
	assert.NotNil(t, err)

	msg, err = chat.SendMsg(MsgData{Text: "testing webxdc", File: acfactory.TestWebxdc()})
	assert.Nil(t, err)

	info, err := msg.WebxdcInfo()
	assert.Nil(t, err)
	assert.NotEmpty(t, info.Name)
}

func TestMessage_SendMsg(t *testing.T) {
	t.Parallel()
	acc := acfactory.OnlineAccount()
	defer acfactory.StopRpc(acc)

	chat, err := acc.Me().CreateChat()
	assert.Nil(t, err)

	_, err = chat.SendMsg(MsgData{Location: &[2]float64{1, 1}})
	assert.Nil(t, err)

	acfactory.WaitForEvent(acc, EventLocationChanged{})

	acc2 := acfactory.OnlineAccount()
	defer acfactory.StopRpc(acc2)
	chat, _ = acfactory.CreateChat(acc2, acc)
	assert.Nil(t, acc.SetConfig("delete_server_after", "1"))
	_, err = chat.SendMsg(MsgData{Text: "test"})
	assert.Nil(t, err)
	acfactory.WaitForEvent(acc, EventImapMessageDeleted{})
}

func TestMessage_StatusUpdates(t *testing.T) {
	t.Parallel()
	acc1 := acfactory.OnlineAccount()
	defer acfactory.StopRpc(acc1)
	acc2 := acfactory.OnlineAccount()
	defer acfactory.StopRpc(acc2)

	chat1, err := acfactory.CreateChat(acc1, acc2)
	assert.Nil(t, err)

	msg, err := chat1.SendMsg(MsgData{Text: "testing webxdc", File: acfactory.TestWebxdc()})
	assert.Nil(t, err)
	snapshot, err := acfactory.NextMsg(acc2)
	assert.Nil(t, err)

	assert.Nil(t, msg.SendStatusUpdate(`{"payload": "test payload"}`, "update 1"))
	acfactory.WaitForEvent(acc2, EventWebxdcStatusUpdate{})
	assert.Nil(t, msg.Delete())
	acfactory.WaitForEvent(acc1, EventWebxdcInstanceDeleted{})

	msg = &Message{acc2, snapshot.Id}
	updates, err := msg.StatusUpdates(0)
	assert.Nil(t, err)
	assert.Contains(t, updates, "test payload")
}

func TestMessage_ContinueAutocryptKeyTransfer(t *testing.T) {
	t.Parallel()
	acc1 := acfactory.OnlineAccount()
	defer acfactory.StopRpc(acc1)
	acc2 := acfactory.UnconfiguredAccount()
	defer acfactory.StopRpc(acc2)

	addr, err := acc1.GetConfig("configured_addr")
	assert.Nil(t, err)
	assert.NotEmpty(t, addr)
	assert.Nil(t, acc2.SetConfig("addr", addr))
	password, err := acc1.GetConfig("configured_mail_pw")
	assert.Nil(t, err)
	assert.NotEmpty(t, password)
	assert.Nil(t, acc2.SetConfig("mail_pw", password))
	assert.Nil(t, acc2.Configure())

	code, err := acc1.InitiateAutocryptKeyTransfer()
	assert.Nil(t, err)
	assert.NotEmpty(t, code)

	selfchat, err := acc2.Me().CreateChat()
	assert.Nil(t, err)
	event := acfactory.WaitForEventInChat(acc2, EventMsgsChanged{}, selfchat.Id).(EventMsgsChanged)
	assert.NotEmpty(t, event.MsgId)
	msg := &Message{acc2, event.MsgId}
	assert.Nil(t, msg.ContinueAutocryptKeyTransfer(code))
}

func TestMsgSnapshot_ParseMemberAddedRemoved(t *testing.T) {
	t.Parallel()
	acc1 := acfactory.OnlineAccount()
	defer acfactory.StopRpc(acc1)
	addr1, err := acc1.GetConfig("configured_addr")
	assert.Nil(t, err)
	acc2 := acfactory.OnlineAccount()
	defer acfactory.StopRpc(acc2)
	addr2, err := acc2.GetConfig("configured_addr")
	assert.Nil(t, err)

	acfactory.IntroduceEachOther(acc1, acc2)
	contact1, err := acc2.CreateContact(addr1, "")
	assert.Nil(t, err)
	contact2, err := acc1.CreateContact(addr2, "")
	assert.Nil(t, err)
	contact3acc1, err := acc1.CreateContact("null@localhost", "")
	assert.Nil(t, err)
	contact3acc2, err := acc2.CreateContact("null@localhost", "")
	assert.Nil(t, err)

	chat1, err := acc1.CreateGroup("test group", false)
	assert.Nil(t, err)
	assert.Nil(t, chat1.AddContact(contact3acc1))

	// promote group
	msg, _ := chat1.SendText("test")
	for {
		event := acfactory.WaitForEventInChat(acc1, EventMsgsChanged{}, chat1.Id).(EventMsgsChanged)
		if event.MsgId == msg.Id {
			break
		}
	}
	snapshot, err := msg.Snapshot()
	assert.Nil(t, err)
	_, _, err = snapshot.ParseMemberAdded()
	assert.NotNil(t, err)

	// add new member
	assert.Nil(t, chat1.AddContact(contact2))
	// acc1 side
	event := acfactory.WaitForEventInChat(acc1, EventMsgsChanged{}, chat1.Id).(EventMsgsChanged)
	assert.NotEmpty(t, event.MsgId)
	msg = &Message{acc1, event.MsgId}
	snapshot, err = msg.Snapshot()
	assert.Nil(t, err)
	assert.Equal(t, SysmsgMemberAddedToGroup, snapshot.SystemMessageType)
	_, _, err = snapshot.ParseMemberRemoved()
	assert.NotNil(t, err)
	actor, target, err := snapshot.ParseMemberAdded()
	assert.Nil(t, err)
	assert.Equal(t, acc1.Me(), actor)
	assert.Equal(t, contact2, target)
	// acc2 side
	snapshot, err = acfactory.NextMsg(acc2)
	assert.Nil(t, err)
	assert.Equal(t, SysmsgMemberAddedToGroup, snapshot.SystemMessageType)
	actor, target, err = snapshot.ParseMemberAdded()
	assert.Nil(t, err)
	assert.Equal(t, contact1, actor)
	assert.Equal(t, acc2.Me(), target)

	// remove new member
	assert.Nil(t, chat1.RemoveContact(contact3acc1))
	// acc1 side
	event = acfactory.WaitForEventInChat(acc1, EventMsgsChanged{}, chat1.Id).(EventMsgsChanged)
	assert.NotEmpty(t, event.MsgId)
	msg = &Message{acc1, event.MsgId}
	snapshot, err = msg.Snapshot()
	assert.Nil(t, err)
	assert.Equal(t, SysmsgMemberRemovedFromGroup, snapshot.SystemMessageType)
	_, _, err = snapshot.ParseMemberAdded()
	assert.NotNil(t, err)
	actor, target, err = snapshot.ParseMemberRemoved()
	assert.Nil(t, err)
	assert.Equal(t, acc1.Me(), actor)
	assert.Equal(t, contact3acc1, target)
	// acc2 side
	snapshot, err = acfactory.NextMsg(acc2)
	assert.Nil(t, err)
	assert.Equal(t, SysmsgMemberRemovedFromGroup, snapshot.SystemMessageType)
	actor, target, err = snapshot.ParseMemberRemoved()
	assert.Nil(t, err)
	assert.Equal(t, contact1, actor)
	assert.Equal(t, contact3acc2, target)

	// leave
	assert.Nil(t, chat1.Leave())
	// acc1 side
	event = acfactory.WaitForEventInChat(acc1, EventMsgsChanged{}, chat1.Id).(EventMsgsChanged)
	assert.NotEmpty(t, event.MsgId)
	msg = &Message{acc1, event.MsgId}
	snapshot, err = msg.Snapshot()
	assert.Nil(t, err)
	assert.Equal(t, SysmsgMemberRemovedFromGroup, snapshot.SystemMessageType)
	actor, target, err = snapshot.ParseMemberRemoved()
	assert.Nil(t, err)
	assert.Equal(t, acc1.Me(), actor)
	assert.Equal(t, acc1.Me(), target)
	// acc2 side
	snapshot, err = acfactory.NextMsg(acc2)
	assert.Nil(t, err)
	assert.Equal(t, SysmsgMemberRemovedFromGroup, snapshot.SystemMessageType)
	actor, target, err = snapshot.ParseMemberRemoved()
	assert.Nil(t, err)
	assert.Equal(t, contact1, actor)
	assert.Equal(t, contact1, target)
}
