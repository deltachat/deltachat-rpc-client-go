package deltachat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContact_String(t *testing.T) {
	t.Parallel()
	acc := acfactory.GetOnlineAccount()
	defer acc.Manager.Rpc.Stop()

	contact, err := acc.CreateContact("test@example.com", "test")
	assert.Nil(t, err)
	assert.NotNil(t, contact)

	assert.NotEmpty(t, contact.String())
}

func TestContact_Block(t *testing.T) {
	t.Parallel()
	acc := acfactory.GetOnlineAccount()
	defer acc.Manager.Rpc.Stop()

	contact, err := acc.CreateContact("test@example.com", "test")
	assert.Nil(t, err)
	assert.NotNil(t, contact)

	assert.Nil(t, contact.Block())
	snapshot, _ := contact.Snapshot()
	assert.True(t, snapshot.IsBlocked)

	assert.Nil(t, contact.Unblock())
	snapshot, _ = contact.Snapshot()
	assert.False(t, snapshot.IsBlocked)
}

func TestContact_Delete(t *testing.T) {
	t.Parallel()
	acc := acfactory.GetOnlineAccount()
	defer acc.Manager.Rpc.Stop()

	contact, err := acc.CreateContact("test@example.com", "test")
	assert.Nil(t, err)
	assert.NotNil(t, contact)

	assert.Nil(t, contact.Delete())
}

func TestContact_SetName(t *testing.T) {
	t.Parallel()
	acc := acfactory.GetOnlineAccount()
	defer acc.Manager.Rpc.Stop()

	contact, err := acc.CreateContact("test@example.com", "test")
	assert.Nil(t, err)
	assert.NotNil(t, contact)

	assert.Nil(t, contact.SetName("new name"))
	snapshot, _ := contact.Snapshot()
	assert.Equal(t, snapshot.Name, "new name")
}

func TestContact_EncryptionInfo(t *testing.T) {
	t.Parallel()
	acc := acfactory.GetOnlineAccount()
	defer acc.Manager.Rpc.Stop()

	contact, err := acc.CreateContact("test@example.com", "test")
	assert.Nil(t, err)
	assert.NotNil(t, contact)

	_, err = contact.EncryptionInfo()
	assert.Nil(t, err)
}

func TestContact_CreateChat(t *testing.T) {
	t.Parallel()
	acc := acfactory.GetOnlineAccount()
	defer acc.Manager.Rpc.Stop()

	contact, err := acc.CreateContact("test@example.com", "test")
	assert.Nil(t, err)
	assert.NotNil(t, contact)

	_, err = contact.CreateChat()
	assert.Nil(t, err)
}
