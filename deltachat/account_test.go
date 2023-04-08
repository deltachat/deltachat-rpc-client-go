package deltachat

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAccount_String(t *testing.T) {
	t.Parallel()
	manager := acfactory.NewAcManager()
	defer manager.Rpc.Stop()

	acc, err := manager.AddAccount()
	assert.Nil(t, err)

	assert.NotEmpty(t, acc.String())
}

func TestAccount_GetEventChannel(t *testing.T) {
	t.Parallel()
	manager := acfactory.NewAcManager()
	defer manager.Rpc.Stop()

	acc, err := manager.AddAccount()
	assert.Nil(t, err)

	assert.NotNil(t, acc.GetEventChannel())
}

func TestAccount_Select(t *testing.T) {
	t.Parallel()
	manager := acfactory.NewAcManager()
	defer manager.Rpc.Stop()

	acc, err := manager.AddAccount()
	assert.Nil(t, err)

	assert.Nil(t, acc.Select())
}

func TestAccount_StartIO(t *testing.T) {
	t.Parallel()
	manager := acfactory.NewAcManager()
	defer manager.Rpc.Stop()

	acc, err := manager.AddAccount()
	assert.Nil(t, err)

	assert.Nil(t, acc.StartIO())
}

func TestAccount_StopIO(t *testing.T) {
	t.Parallel()
	manager := acfactory.NewAcManager()
	defer manager.Rpc.Stop()

	acc, err := manager.AddAccount()
	assert.Nil(t, err)

	assert.Nil(t, acc.StopIO())
}

func TestAccount_Connectivity(t *testing.T) {
	t.Parallel()
	manager := acfactory.NewAcManager()
	defer manager.Rpc.Stop()

	acc, err := manager.AddAccount()
	assert.Nil(t, err)

	conn, err := acc.Connectivity()
	assert.Nil(t, err)
	assert.True(t, conn > 0)

	_, err = acc.ConnectivityHtml()
	assert.NotNil(t, err)

	acc = acfactory.GetOnlineAccount()

	html, err := acc.ConnectivityHtml()
	assert.Nil(t, err)
	assert.NotEmpty(t, html)
}

func TestAccount_Info(t *testing.T) {
	t.Parallel()
	manager := acfactory.NewAcManager()
	defer manager.Rpc.Stop()

	acc, err := manager.AddAccount()
	assert.Nil(t, err)

	info, err := acc.Info()
	assert.Nil(t, err)
	assert.NotEmpty(t, info["sqlite_version"])
}

func TestAccount_Size(t *testing.T) {
	t.Parallel()
	manager := acfactory.NewAcManager()
	defer manager.Rpc.Stop()

	acc, err := manager.AddAccount()
	assert.Nil(t, err)

	size, err := acc.Size()
	assert.Nil(t, err)
	assert.NotEqual(t, size, 0)
}

func TestAccount_IsConfigured(t *testing.T) {
	t.Parallel()
	acc := acfactory.GetUnconfiguredAccount()
	defer acc.Manager.Rpc.Stop()

	configured, err := acc.IsConfigured()
	assert.Nil(t, err)
	assert.False(t, configured)

	assert.Nil(t, acc.Configure())

	configured, err = acc.IsConfigured()
	assert.Nil(t, err)
	assert.True(t, configured)
}

func TestAccount_SetAndGetConfig(t *testing.T) {
	t.Parallel()
	manager := acfactory.NewAcManager()
	defer manager.Rpc.Stop()

	acc, err := manager.AddAccount()
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

	assert.Nil(t, acc.SetConfig("selfavatar", acfactory.GetTestImage()))
	WaitForEvent(acc, eventSelfavatarChanged)
}

func TestAccount_Avatar(t *testing.T) {
	t.Parallel()
	manager := acfactory.NewAcManager()
	defer manager.Rpc.Stop()

	acc, err := manager.AddAccount()
	assert.Nil(t, err)

	err = acc.SetAvatar("invalid.jpg")
	assert.Contains(t, err.Error(), "failed to open file")

	avatar, err := acc.Avatar()
	assert.Nil(t, err)
	assert.Equal(t, avatar, "")
}

func TestAccount_Remove(t *testing.T) {
	t.Parallel()
	manager := acfactory.NewAcManager()
	defer manager.Rpc.Stop()

	acc, err := manager.AddAccount()
	assert.Nil(t, err)

	assert.Nil(t, acc.Remove())
}

func TestAccount_Configure(t *testing.T) {
	t.Parallel()
	acc := acfactory.GetUnconfiguredAccount()
	defer acc.Manager.Rpc.Stop()
	assert.Nil(t, acc.Configure())
}

func TestAccount_Contacts(t *testing.T) {
	t.Parallel()
	acc := acfactory.GetOnlineAccount()
	defer acc.Manager.Rpc.Stop()

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
	t.Parallel()
	acc := acfactory.GetOnlineAccount()
	defer acc.Manager.Rpc.Stop()

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
	t.Parallel()
	acc := acfactory.GetOnlineAccount()
	defer acc.Manager.Rpc.Stop()

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
	t.Parallel()
	acc := acfactory.GetOnlineAccount()
	defer acc.Manager.Rpc.Stop()

	assert.NotNil(t, acc.Me())
}

func TestAccount_CreateBroadcastList(t *testing.T) {
	t.Parallel()
	acc := acfactory.GetOnlineAccount()
	defer acc.Manager.Rpc.Stop()

	chat, err := acc.CreateBroadcastList()
	assert.Nil(t, err)
	assert.NotNil(t, chat)
}

func TestAccount_CreateGroup(t *testing.T) {
	t.Parallel()
	acc := acfactory.GetOnlineAccount()
	defer acc.Manager.Rpc.Stop()

	chat, err := acc.CreateGroup("test group", true)
	assert.Nil(t, err)
	assert.NotNil(t, chat)
}

func TestAccount_QrCode(t *testing.T) {
	t.Parallel()
	acc := acfactory.GetOnlineAccount()
	defer acc.Manager.Rpc.Stop()

	qrdata, svg, err := acc.QrCode()
	assert.Nil(t, err)
	assert.NotEmpty(t, qrdata)
	assert.NotEmpty(t, svg)

	acc2 := acfactory.GetOnlineAccount()
	defer acc2.Manager.Rpc.Stop()
	acfactory.IntroduceEachOther(acc, acc2)
	chat2, err := acc2.SecureJoin(qrdata)
	assert.Nil(t, err)
	assert.NotNil(t, chat2)

	WaitForEvent(acc, eventSecurejoinInviterProgress)
	WaitForEvent(acc2, eventSecurejoinJoinerProgress)
}

func TestAccount_ImportSelfKeys(t *testing.T) {
	t.Parallel()
	acc := acfactory.GetOnlineAccount()
	defer acc.Manager.Rpc.Stop()

	dir, err := os.MkdirTemp("", "")
	assert.Nil(t, err)
	defer os.RemoveAll(dir)
	assert.Nil(t, acc.ExportSelfKeys(dir))
	assert.Nil(t, acc.ImportSelfKeys(dir))
}

func TestAccount_ImportBackup(t *testing.T) {
	t.Parallel()
	acc := acfactory.GetOnlineAccount()
	defer acc.Manager.Rpc.Stop()

	dir, err := os.MkdirTemp("", "")
	assert.Nil(t, err)
	defer os.RemoveAll(dir)
	assert.Nil(t, acc.ExportBackup(dir, ""))
	files, err := os.ReadDir(dir)
	assert.Nil(t, err)
	assert.Equal(t, len(files), 1)

	t.Skip("skipping ImportBackup due to bug in deltachat-rpc-server")
	backup := filepath.Join(dir, files[0].Name())
	assert.FileExists(t, backup)
	manager := acfactory.NewAcManager()
	defer manager.Rpc.Stop()
	acc2, err := manager.AddAccount()
	assert.Nil(t, err)
	assert.Nil(t, acc2.ImportBackup(backup, ""))
}

func TestAccount_ExportBackup(t *testing.T) {
	t.Parallel()
	acc := acfactory.GetOnlineAccount()
	defer acc.Manager.Rpc.Stop()

	dir, err := os.MkdirTemp("", "")
	assert.Nil(t, err)
	defer os.RemoveAll(dir)
	assert.Nil(t, acc.ExportBackup(dir, "test-phrase"))
	files, err := os.ReadDir(dir)
	assert.Nil(t, err)
	assert.Equal(t, len(files), 1)
}

func TestAccount_GetBackup(t *testing.T) {
	t.Parallel()
	acc := acfactory.GetOnlineAccount()
	defer acc.Manager.Rpc.Stop()

	go func() { assert.Nil(t, acc.ProvideBackup()) }()
	var err error
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

	manager := acfactory.NewAcManager()
	defer manager.Rpc.Stop()
	acc2, err := manager.AddAccount()
	assert.Nil(t, err)
	assert.Nil(t, acc2.GetBackup(qrData))
}

func TestAccount_InitiateAutocryptKeyTransfer(t *testing.T) {
	t.Parallel()
	acc := acfactory.GetOnlineAccount()
	defer acc.Manager.Rpc.Stop()

	code, err := acc.InitiateAutocryptKeyTransfer()
	assert.Nil(t, err)
	assert.NotEmpty(t, code)
}

func TestAccount_FreshMsgs(t *testing.T) {
	t.Parallel()
	acc := acfactory.GetOnlineAccount()
	defer acc.Manager.Rpc.Stop()
	acc2 := acfactory.GetOnlineAccount()
	defer acc2.Manager.Rpc.Stop()

	chat2, err := acc2.CreateChat(acc)
	assert.Nil(t, err)
	chat2.SendText("hi")
	msg, err := acfactory.GetNextMsg(acc)
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
	t.Parallel()
	acc := acfactory.GetOnlineAccount()
	defer acc.Manager.Rpc.Stop()
	acc2 := acfactory.GetOnlineAccount()
	defer acc2.Manager.Rpc.Stop()

	chat2, err := acc2.CreateChat(acc)
	assert.Nil(t, err)
	chat2.SendText("hi")
	msg, err := acfactory.GetNextMsg(acc)
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

func TestAccount_SearchMessages(t *testing.T) {
	t.Parallel()
	acc := acfactory.GetOnlineAccount()
	defer acc.Manager.Rpc.Stop()
	acc2 := acfactory.GetOnlineAccount()
	defer acc2.Manager.Rpc.Stop()

	chat2, err := acc2.CreateChat(acc)
	assert.Nil(t, err)
	chat2.SendText("hi")
	msg, err := acfactory.GetNextMsg(acc)
	assert.Nil(t, err)
	assert.Equal(t, msg.Text, "hi")

	results, err := acc.SearchMessages("hi")
	assert.Nil(t, err)
	assert.NotEmpty(t, results)
	assert.Equal(t, msg.Id, results[0].Id)
	assert.Equal(t, msg.Text, results[0].Message)
}

func TestAccount_ChatListItems(t *testing.T) {
	t.Parallel()
	acc := acfactory.GetOnlineAccount()
	defer acc.Manager.Rpc.Stop()

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
	t.Parallel()
	acc := acfactory.GetOnlineAccount()
	defer acc.Manager.Rpc.Stop()

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
	t.Parallel()
	acc := acfactory.GetOnlineAccount()
	defer acc.Manager.Rpc.Stop()

	message, err := acc.AddDeviceMsg("test", "new message")
	assert.Nil(t, err)
	assert.NotNil(t, message)
	msg, err := message.Snapshot()
	assert.Nil(t, err)
	assert.Equal(t, msg.Text, "new message")
}
