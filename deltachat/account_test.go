package deltachat

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAccount_String(t *testing.T) {
	acc, err := server.AccountManager().AddAccount()
	assert.Nil(t, err)

	assert.NotEmpty(t, acc.String())
}

func TestAccount_GetEventChannel(t *testing.T) {
	acc, err := server.AccountManager().AddAccount()
	assert.Nil(t, err)

	assert.NotNil(t, acc.GetEventChannel())
}

func TestAccount_Select(t *testing.T) {
	acc, err := server.AccountManager().AddAccount()
	assert.Nil(t, err)

	assert.Nil(t, acc.Select())
}

func TestAccount_StartIO(t *testing.T) {
	acc, err := server.AccountManager().AddAccount()
	assert.Nil(t, err)

	assert.Nil(t, acc.StartIO())
}

func TestAccount_StopIO(t *testing.T) {
	acc, err := server.AccountManager().AddAccount()
	assert.Nil(t, err)

	assert.Nil(t, acc.StopIO())
}

func TestAccount_Connectivity(t *testing.T) {
	acc, err := server.AccountManager().AddAccount()
	assert.Nil(t, err)

	conn, err := acc.Connectivity()
	assert.Nil(t, err)
	assert.True(t, conn > 0)
}

func TestAccount_Info(t *testing.T) {
	acc, err := server.AccountManager().AddAccount()
	assert.Nil(t, err)

	info, err := acc.Info()
	assert.Nil(t, err)
	assert.NotEmpty(t, info["sqlite_version"])
}

func TestAccount_Size(t *testing.T) {
	acc, err := server.AccountManager().AddAccount()
	assert.Nil(t, err)

	size, err := acc.Size()
	assert.Nil(t, err)
	assert.NotEqual(t, size, 0)
}

func TestAccount_IsConfigured(t *testing.T) {
	acc := server.GetUnconfiguredAccount()

	configured, err := acc.IsConfigured()
	assert.Nil(t, err)
	assert.False(t, configured)

	assert.Nil(t, acc.Configure())

	configured, err = acc.IsConfigured()
	assert.Nil(t, err)
	assert.True(t, configured)
}

func TestAccount_SetAndGetConfig(t *testing.T) {
	acc, err := server.AccountManager().AddAccount()
	assert.Nil(t, err)

	assert.Nil(t, acc.SetConfig("displayname", "test name"))
	name, err := acc.GetConfig("displayname")
	assert.Nil(t, err)
	assert.Equal(t, name, "test name")

	err = acc.UpdateConfig(map[string]string{
		"displayname": "new name",
		"selfstatus":  "test status",
	})
	assert.Nil(t, err)
	name, err = acc.GetConfig("displayname")
	assert.Nil(t, err)
	assert.Equal(t, name, "new name")
}

func TestAccount_Avatar(t *testing.T) {
	acc, err := server.AccountManager().AddAccount()
	assert.Nil(t, err)

	err = acc.SetAvatar("invalid.jpg")
	assert.Contains(t, err.Error(), "failed to open file")

	avatar, err := acc.Avatar()
	assert.Nil(t, err)
	assert.Equal(t, avatar, "")
}

func TestAccount_Remove(t *testing.T) {
	acc, err := server.AccountManager().AddAccount()
	assert.Nil(t, err)

	assert.Nil(t, acc.Remove())
}

func TestAccount_Configure(t *testing.T) {
	acc := server.GetUnconfiguredAccount()
	assert.Nil(t, acc.Configure())
}

func TestAccount_Contacts(t *testing.T) {
	acc, err := server.GetOnlineAccount()
	assert.Nil(t, err)

	contacts, err := acc.Contacts()
	assert.Nil(t, err)
	assert.Empty(t, contacts)

	contact, err := acc.CreateContact("test@example.com", "test")
	assert.Nil(t, err)
	assert.NotNil(t, contact)

	contacts, err = acc.Contacts()
	assert.Nil(t, err)
	assert.Contains(t, contacts, contact)

	contacts, err = acc.QueryContacts("unknown", 0)
	assert.Nil(t, err)
	assert.Empty(t, contacts)
	contacts, err = acc.QueryContacts("test", 0)
	assert.Nil(t, err)
	assert.Contains(t, contacts, contact)
}

func TestAccount_GetContactByAddr(t *testing.T) {
	acc, err := server.GetOnlineAccount()
	assert.Nil(t, err)

	contact, err := acc.CreateContact("test@example.com", "test")
	assert.Nil(t, err)
	assert.NotNil(t, contact)

	contact2, err := acc.GetContactByAddr("unknown@example.com")
	assert.Nil(t, err)
	assert.Nil(t, contact2)

	contact2, err = acc.GetContactByAddr("test@example.com")
	assert.Nil(t, err)
	assert.NotNil(t, contact2)
	assert.Equal(t, contact, contact2)
}

func TestAccount_BlockedContacts(t *testing.T) {
	acc, err := server.GetOnlineAccount()
	assert.Nil(t, err)

	contact, err := acc.CreateContact("test@example.com", "test")
	assert.Nil(t, err)
	assert.NotNil(t, contact)

	blocked, err := acc.BlockedContacts()
	assert.Nil(t, err)
	assert.Empty(t, blocked)

	assert.Nil(t, contact.Block())

	blocked, err = acc.BlockedContacts()
	assert.Nil(t, err)
	assert.NotEmpty(t, blocked)
}

func TestAccount_Me(t *testing.T) {
	acc, err := server.GetOnlineAccount()
	assert.Nil(t, err)

	assert.NotNil(t, acc.Me())
}

func TestAccount_CreateBroadcastList(t *testing.T) {
	acc, err := server.GetOnlineAccount()
	assert.Nil(t, err)

	chat, err := acc.CreateBroadcastList()
	assert.Nil(t, err)
	assert.NotNil(t, chat)
}

func TestAccount_CreateGroup(t *testing.T) {
	acc, err := server.GetOnlineAccount()
	assert.Nil(t, err)

	chat, err := acc.CreateGroup("test group", true)
	assert.Nil(t, err)
	assert.NotNil(t, chat)
}

func TestAccount_QrCode(t *testing.T) {
	acc, err := server.GetOnlineAccount()
	assert.Nil(t, err)

	qrdata, svg, err := acc.QrCode()
	assert.Nil(t, err)
	assert.NotEmpty(t, qrdata)
	assert.NotEmpty(t, svg)

	acc2, err := server.GetOnlineAccount()
	assert.Nil(t, err)
	chat2, err := acc2.SecureJoin(qrdata)
	assert.Nil(t, err)
	assert.NotNil(t, chat2)

	event := acc.WaitForEvent(EVENT_SECUREJOIN_INVITER_PROGRESS)
	assert.NotNil(t, event)

	event = acc2.WaitForEvent(EVENT_SECUREJOIN_JOINER_PROGRESS)
	assert.NotNil(t, event)
}

func TestAccount_ImportSelfKeys(t *testing.T) {
	acc, err := server.GetOnlineAccount()
	assert.Nil(t, err)

	dir, err := os.MkdirTemp("", "")
	assert.Nil(t, err)
	defer os.RemoveAll(dir)
	assert.Nil(t, acc.ExportSelfKeys(dir))
	assert.Nil(t, acc.ImportSelfKeys(dir))
}

func TestAccount_ImportBackup(t *testing.T) {
	acc, err := server.GetOnlineAccount()
	assert.Nil(t, err)

	dir, err := os.MkdirTemp("", "")
	assert.Nil(t, err)
	defer os.RemoveAll(dir)
	assert.Nil(t, acc.ExportBackup(dir, ""))
	files, err := os.ReadDir(dir)
	assert.Nil(t, err)
	assert.Equal(t, len(files), 1)

	t.Skip("skipping ImportBackup due to bug")
	backup := filepath.Join(dir, files[0].Name())
	assert.FileExists(t, backup)
	acc2, err := server.AccountManager().AddAccount()
	assert.Nil(t, err)
	assert.Nil(t, acc2.ImportBackup(backup, ""))
}

func TestAccount_GetBackup(t *testing.T) {
	var err error
	acc, err := server.GetOnlineAccount()
	assert.Nil(t, err)

	go func() { assert.Nil(t, acc.ProvideBackup()) }()
	var qrData string
	qrData, err = acc.GetBackupQr()
	for err != nil {
		time.Sleep(time.Millisecond * 200)
		qrData, err = acc.GetBackupQr()
	}
	assert.NotNil(t, qrData)

	qrSvg, err := acc.GetBackupQrSvg()
	assert.Nil(t, err)
	assert.NotNil(t, qrSvg)

	acc2, err := server.AccountManager().AddAccount()
	assert.Nil(t, err)
	assert.Nil(t, acc2.GetBackup(qrData))
}

func TestAccount_InitiateAutocryptKeyTransfer(t *testing.T) {
	acc, err := server.GetOnlineAccount()
	assert.Nil(t, err)

	code, err := acc.InitiateAutocryptKeyTransfer()
	assert.Nil(t, err)
	assert.NotEmpty(t, code)
}

func TestAccount_FreshMsgs(t *testing.T) {
	acc, err := server.GetOnlineAccount()
	assert.Nil(t, err)
	acc2, err := server.GetOnlineAccount()
	assert.Nil(t, err)

	chat2, err := acc2.CreateChat(acc)
	assert.Nil(t, err)
	chat2.SendText("hi")
	msg, err := server.GetNextMsg(acc)
	assert.Nil(t, err)
	assert.Equal(t, msg.Text, "hi")

	msgs, err := acc.FreshMsgs()
	assert.Nil(t, err)
	assert.NotEmpty(t, msgs)

	msgs, err = acc.FreshMsgsInArrivalOrder()
	assert.Nil(t, err)
	assert.NotEmpty(t, msgs)

	assert.Nil(t, acc.MarkSeenMsgs(msgs))

	msgs, err = acc.FreshMsgs()
	assert.Nil(t, err)
	assert.Empty(t, msgs)

	msgs, err = acc.FreshMsgsInArrivalOrder()
	assert.Nil(t, err)
	assert.Empty(t, msgs)
}

func TestAccount_DeleteMsgs(t *testing.T) {
	acc, err := server.GetOnlineAccount()
	assert.Nil(t, err)
	acc2, err := server.GetOnlineAccount()
	assert.Nil(t, err)

	chat2, err := acc2.CreateChat(acc)
	assert.Nil(t, err)
	chat2.SendText("hi")
	msg, err := server.GetNextMsg(acc)
	assert.Nil(t, err)
	assert.Equal(t, msg.Text, "hi")

	msgs, err := acc.FreshMsgs()
	assert.Nil(t, err)
	assert.NotEmpty(t, msgs)

	assert.Nil(t, acc.DeleteMsgs(msgs))

	msgs, err = acc.FreshMsgs()
	assert.Nil(t, err)
	assert.Empty(t, msgs)
}

func TestAccount_ChatListItems(t *testing.T) {
	acc, err := server.GetOnlineAccount()
	assert.Nil(t, err)

	contact, err := acc.CreateContact("test@example.com", "test")
	assert.Nil(t, err)
	_, err = contact.CreateChat()
	assert.Nil(t, err)

	chatitems, err := acc.ChatListItems()
	assert.Nil(t, err)
	assert.NotEmpty(t, chatitems)

	chatitems, err = acc.QueryChatListItems("unknown", nil, 0)
	assert.Nil(t, err)
	assert.Empty(t, chatitems)
}

func TestAccount_ChatListEntries(t *testing.T) {
	acc, err := server.GetOnlineAccount()
	assert.Nil(t, err)

	contact, err := acc.CreateContact("test@example.com", "test")
	assert.Nil(t, err)
	_, err = contact.CreateChat()
	assert.Nil(t, err)

	chats, err := acc.ChatListEntries()
	assert.Nil(t, err)
	assert.NotEmpty(t, chats)

	chats, err = acc.QueryChatListEntries("unknown", nil, 0)
	assert.Nil(t, err)
	assert.Empty(t, chats)
}

func TestAccount_AddDeviceMsg(t *testing.T) {
	acc, err := server.GetOnlineAccount()
	assert.Nil(t, err)

	message, err := acc.AddDeviceMsg("test", "new message")
	assert.Nil(t, err)
	assert.NotNil(t, message)
	msg, err := message.Snapshot()
	assert.Nil(t, err)
	assert.Equal(t, msg.Text, "new message")
}
