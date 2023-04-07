package deltachat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	t.Parallel()
	acc := acfactory.GetOnlineAccount()
	defer acc.Manager.Rpc.Stop()

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

	assert.Nil(t, msg.Delete())
}

func TestMessage_WebxdcInfo(t *testing.T) {
	t.Parallel()
	acc := acfactory.GetOnlineAccount()
	defer acc.Manager.Rpc.Stop()

	chat, err := acc.CreateChat(acc)
	assert.Nil(t, err)

	msg, err := chat.SendText("test")
	assert.Nil(t, err)
	info, err := msg.WebxdcInfo()
	assert.NotNil(t, err)

	msg, err = chat.SendMsg(MsgData{Text: "testing webxdc", File: acfactory.GetTestWebxdc()})
	assert.Nil(t, err)

	info, err = msg.WebxdcInfo()
	assert.Nil(t, err)
	assert.NotEmpty(t, info.Name)
}

func TestMessage_StatusUpdates(t *testing.T) {
	t.Parallel()
	acc1 := acfactory.GetOnlineAccount()
	defer acc1.Manager.Rpc.Stop()
	acc2 := acfactory.GetOnlineAccount()
	defer acc2.Manager.Rpc.Stop()

	chat1, err := acc1.CreateChat(acc2)
	assert.Nil(t, err)

	msg, err := chat1.SendMsg(MsgData{Text: "testing webxdc", File: acfactory.GetTestWebxdc()})
	assert.Nil(t, err)
	snapshot, err := acfactory.GetNextMsg(acc2)
	assert.Nil(t, err)

	assert.Nil(t, msg.SendStatusUpdate(`{"payload": "test payload"}`, "update 1"))
	WaitForEvent(acc2, eventWebxdcStatusUpdate)

	msg = &Message{acc2, snapshot.Id}
	updates, err := msg.StatusUpdates(0)
	assert.Nil(t, err)
	assert.Contains(t, updates, "test payload")
}

func TestMessage_ContinueAutocryptKeyTransfer(t *testing.T) {
	t.Parallel()
	acc1 := acfactory.GetOnlineAccount()
	defer acc1.Manager.Rpc.Stop()
	acc2 := acfactory.GetUnconfiguredAccount()
	defer acc2.Manager.Rpc.Stop()

	addr, err := acc1.GetConfig("configured_addr")
	assert.Nil(t, err)
	assert.NotEmpty(t, addr)
	acc2.SetConfig("addr", addr)
	password, err := acc1.GetConfig("configured_mail_pw")
	assert.Nil(t, err)
	assert.NotEmpty(t, password)
	acc2.SetConfig("mail_pw", password)
	assert.Nil(t, acc2.Configure())

	code, err := acc1.InitiateAutocryptKeyTransfer()
	assert.Nil(t, err)
	assert.NotEmpty(t, code)

	selfchat, err := acc2.Me().CreateChat()
	assert.Nil(t, err)
	event := waitForEvent(acc2, eventMsgsChanged, selfchat.Id).(EventMsgsChanged)
	assert.NotEmpty(t, event.MsgId)
	msg := &Message{acc2, event.MsgId}
	assert.Nil(t, msg.ContinueAutocryptKeyTransfer(code))
}

func TestMsgSnapshot_ParseMemberAddedRemoved(t *testing.T) {
	t.Parallel()
	acc1 := acfactory.GetOnlineAccount()
	defer acc1.Manager.Rpc.Stop()
	addr1, err := acc1.GetConfig("configured_addr")
	assert.Nil(t, err)
	acc2 := acfactory.GetOnlineAccount()
	defer acc2.Manager.Rpc.Stop()
	addr2, err := acc2.GetConfig("configured_addr")
	assert.Nil(t, err)

	acfactory.IntroduceEachOther(acc1, acc2)
	contact1, err := acc2.CreateContact(addr1, "")
	assert.Nil(t, err)
	contact2, err := acc1.CreateContact(addr2, "")
	assert.Nil(t, err)
	contact3acc1, err := acc1.CreateContact("test@example.com", "")
	assert.Nil(t, err)
	contact3acc2, err := acc2.CreateContact("test@example.com", "")
	assert.Nil(t, err)

	chat1, err := acc1.CreateGroup("test group", false)
	assert.Nil(t, err)
	assert.Nil(t, chat1.AddContact(contact3acc1))

	// promote group
	msg, _ := chat1.SendText("test")
	for {
		event := waitForEvent(acc1, eventMsgsChanged, chat1.Id).(EventMsgsChanged)
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
	event := waitForEvent(acc1, eventMsgsChanged, chat1.Id).(EventMsgsChanged)
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
	snapshot, err = acfactory.GetNextMsg(acc2)
	assert.Nil(t, err)
	assert.Equal(t, SysmsgMemberAddedToGroup, snapshot.SystemMessageType)
	actor, target, err = snapshot.ParseMemberAdded()
	assert.Nil(t, err)
	assert.Equal(t, contact1, actor)
	assert.Equal(t, acc2.Me(), target)

	// remove new member
	assert.Nil(t, chat1.RemoveContact(contact3acc1))
	// acc1 side
	event = waitForEvent(acc1, eventMsgsChanged, chat1.Id).(EventMsgsChanged)
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
	snapshot, err = acfactory.GetNextMsg(acc2)
	assert.Nil(t, err)
	assert.Equal(t, SysmsgMemberRemovedFromGroup, snapshot.SystemMessageType)
	actor, target, err = snapshot.ParseMemberRemoved()
	assert.Nil(t, err)
	assert.Equal(t, contact1, actor)
	assert.Equal(t, contact3acc2, target)

	// leave
	assert.Nil(t, chat1.Leave())
	// acc1 side
	event = waitForEvent(acc1, eventMsgsChanged, chat1.Id).(EventMsgsChanged)
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
	snapshot, err = acfactory.GetNextMsg(acc2)
	assert.Nil(t, err)
	assert.Equal(t, SysmsgMemberRemovedFromGroup, snapshot.SystemMessageType)
	actor, target, err = snapshot.ParseMemberRemoved()
	assert.Nil(t, err)
	assert.Equal(t, contact1, actor)
	assert.Equal(t, contact1, target)
}
