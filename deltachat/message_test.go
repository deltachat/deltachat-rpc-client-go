package deltachat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
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
}

func TestMsgSnapshot_ParseMemberAddedRemoved(t *testing.T) {
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
		event := waitForEvent(acc1, EVENT_MSGS_CHANGED, chat1.Id)
		if event.MsgId == msg.Id {
			break
		}
	}

	// add new member
	assert.Nil(t, chat1.AddContact(contact2))
	// acc1 side
	event := waitForEvent(acc1, EVENT_MSGS_CHANGED, chat1.Id)
	assert.NotEmpty(t, event.MsgId)
	msg = &Message{acc1, event.MsgId}
	snapshot, err := msg.Snapshot()
	assert.Nil(t, err)
	assert.Equal(t, SYSMSG_TYPE_MEMBER_ADDED_TO_GROUP, snapshot.SystemMessageType)
	actor, target, err := snapshot.ParseMemberAdded()
	assert.Nil(t, err)
	assert.Equal(t, acc1.Me(), actor)
	assert.Equal(t, contact2, target)
	// acc2 side
	snapshot, err = acfactory.GetNextMsg(acc2)
	assert.Nil(t, err)
	assert.Equal(t, SYSMSG_TYPE_MEMBER_ADDED_TO_GROUP, snapshot.SystemMessageType)
	actor, target, err = snapshot.ParseMemberAdded()
	assert.Nil(t, err)
	assert.Equal(t, contact1, actor)
	assert.Equal(t, acc2.Me(), target)

	// remove new member
	assert.Nil(t, chat1.RemoveContact(contact3acc1))
	// acc1 side
	event = waitForEvent(acc1, EVENT_MSGS_CHANGED, chat1.Id)
	assert.NotEmpty(t, event.MsgId)
	msg = &Message{acc1, event.MsgId}
	snapshot, err = msg.Snapshot()
	assert.Nil(t, err)
	assert.Equal(t, SYSMSG_TYPE_MEMBER_REMOVED_FROM_GROUP, snapshot.SystemMessageType)
	actor, target, err = snapshot.ParseMemberRemoved()
	assert.Nil(t, err)
	assert.Equal(t, acc1.Me(), actor)
	assert.Equal(t, contact3acc1, target)
	// acc2 side
	snapshot, err = acfactory.GetNextMsg(acc2)
	assert.Nil(t, err)
	assert.Equal(t, SYSMSG_TYPE_MEMBER_REMOVED_FROM_GROUP, snapshot.SystemMessageType)
	actor, target, err = snapshot.ParseMemberRemoved()
	assert.Nil(t, err)
	assert.Equal(t, contact1, actor)
	assert.Equal(t, contact3acc2, target)

	// leave
	assert.Nil(t, chat1.Leave())
	// acc1 side
	event = waitForEvent(acc1, EVENT_MSGS_CHANGED, chat1.Id)
	assert.NotEmpty(t, event.MsgId)
	msg = &Message{acc1, event.MsgId}
	snapshot, err = msg.Snapshot()
	assert.Nil(t, err)
	assert.Equal(t, SYSMSG_TYPE_MEMBER_REMOVED_FROM_GROUP, snapshot.SystemMessageType)
	actor, target, err = snapshot.ParseMemberRemoved()
	assert.Nil(t, err)
	assert.Equal(t, acc1.Me(), actor)
	assert.Equal(t, acc1.Me(), target)
	// acc2 side
	snapshot, err = acfactory.GetNextMsg(acc2)
	assert.Nil(t, err)
	assert.Equal(t, SYSMSG_TYPE_MEMBER_REMOVED_FROM_GROUP, snapshot.SystemMessageType)
	actor, target, err = snapshot.ParseMemberRemoved()
	assert.Nil(t, err)
	assert.Equal(t, contact1, actor)
	assert.Equal(t, contact1, target)
}
