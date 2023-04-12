package deltachat

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestContact_String(t *testing.T) {
	t.Parallel()
	acc := acfactory.OnlineAccount()
	defer acfactory.StopRpc(acc)

	contact, err := acc.CreateContact("null@localhost", "test")
	assert.Nil(t, err)
	assert.NotNil(t, contact)

	assert.NotEmpty(t, contact.String())
}

func TestContact_Block(t *testing.T) {
	t.Parallel()
	acc := acfactory.OnlineAccount()
	defer acfactory.StopRpc(acc)

	contact, err := acc.CreateContact("null@localhost", "test")
	assert.Nil(t, err)
	assert.NotNil(t, contact)

	assert.Nil(t, contact.Block())
	snapshot, err := contact.Snapshot()
	assert.Nil(t, err)
	assert.True(t, snapshot.IsBlocked)

	assert.Nil(t, contact.Unblock())
	snapshot, err = contact.Snapshot()
	assert.Nil(t, err)
	assert.False(t, snapshot.IsBlocked)
}

func TestContact_Delete(t *testing.T) {
	t.Parallel()
	acc := acfactory.OnlineAccount()
	defer acfactory.StopRpc(acc)

	contact, err := acc.CreateContact("null@localhost", "test")
	assert.Nil(t, err)
	assert.NotNil(t, contact)

	assert.Nil(t, contact.Delete())
}

func TestContact_SetName(t *testing.T) {
	t.Parallel()
	acc := acfactory.OnlineAccount()
	defer acfactory.StopRpc(acc)

	contact, err := acc.CreateContact("null@localhost", "test")
	assert.Nil(t, err)
	assert.NotNil(t, contact)

	assert.Nil(t, contact.SetName("new name"))
	snapshot, err := contact.Snapshot()
	assert.Nil(t, err)
	assert.Equal(t, snapshot.Name, "new name")
}

func TestContact_EncryptionInfo(t *testing.T) {
	t.Parallel()
	acc := acfactory.OnlineAccount()
	defer acfactory.StopRpc(acc)

	contact, err := acc.CreateContact("null@localhost", "test")
	assert.Nil(t, err)
	assert.NotNil(t, contact)

	_, err = contact.EncryptionInfo()
	assert.Nil(t, err)
}

func TestContact_CreateChat(t *testing.T) {
	t.Parallel()
	acc := acfactory.OnlineAccount()
	defer acfactory.StopRpc(acc)

	contact, err := acc.CreateContact("null@localhost", "test")
	assert.Nil(t, err)
	assert.NotNil(t, contact)

	_, err = contact.CreateChat()
	assert.Nil(t, err)
}

func TestContact_Snapshot(t *testing.T) {
	t.Parallel()
	acc1 := acfactory.OnlineAccount()
	defer acfactory.StopRpc(acc1)
	acc2 := acfactory.OnlineAccount()
	defer acfactory.StopRpc(acc2)

	addr1, err := acc1.GetConfig("configured_addr")
	assert.Nil(t, err)
	contact1, err := acc2.CreateContact(addr1, "")
	assert.Nil(t, err)

	snapshot, err := contact1.Snapshot()
	assert.Nil(t, err)
	assert.Equal(t, Timestamp{time.Unix(0, 0)}, snapshot.LastSeen)

	chat1, err := acfactory.CreateChat(acc1, acc2)
	assert.Nil(t, err)
	_, err = chat1.SendText("hi")
	assert.Nil(t, err)
	msgSnapshot, err := acfactory.NextMsg(acc2)
	assert.Nil(t, err)
	assert.Equal(t, "hi", msgSnapshot.Text)

	snapshot, err = contact1.Snapshot()
	assert.Nil(t, err)
	assert.NotEqual(t, Timestamp{time.Unix(0, 0)}, snapshot.LastSeen)
}
