package deltachat

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// AcFactory facilitates unit testing Delta Chat clients/bots.
// It must be used in conjunction with a test mail server service, for example:
// https://github.com/deltachat/mail-server-tester
type AcFactory struct {
	tempDir     string
	acCfg       map[string]string
	debug       bool
	serial      int64
	startTime   int64
	serialMutex sync.Mutex
	tearUp      bool
}

// Prepare the AcFactory, defaultAcConfig is the default settings to apply to
// the new accounts created with UnconfiguredAccount(), OnlineAccount(), OnlineBot()
// and RunningBot().
//
// If the test mail server has not standard configuration, you should set the custom configuration
// here.
func (self *AcFactory) TearUp(defaultAcConfig map[string]string, tempDir string, debug bool) {
	self.acCfg = defaultAcConfig
	self.tempDir = tempDir
	self.debug = debug
	self.startTime = time.Now().Unix()
	self.tearUp = true
}

// Do cleanup, removing temporary directories and files created by the configured test accounts.
// Usually TearDown() is called with defer immediately after the creation of the AcFactory instance.
func (self *AcFactory) TearDown() {
	self.ensureTearUp()
	os.RemoveAll(self.tempDir)
}

// Stop the Rpc of the given Account, Bot or AccountManager.
func (self *AcFactory) StopRpc(accountOrBot any) {
	switch obj := accountOrBot.(type) {
	case *Bot:
		obj.Account.Manager.Rpc.Stop()
	case *Account:
		obj.Manager.Rpc.Stop()
	case *AccountManager:
		obj.Rpc.Stop()
	default:
		panic("invalid type provided to StopRpc()")
	}
}

// MkdirTemp creates a new temporary directory. The directory is automatically removed on TearDown().
func (self *AcFactory) MkdirTemp() string {
	dir, err := os.MkdirTemp(self.tempDir, "")
	if err != nil {
		panic(err)
	}
	return dir
}

// Create a new AccountManager.
func (self *AcFactory) NewAcManager() *AccountManager {
	self.ensureTearUp()
	rpc := NewRpcIO()
	if !self.debug {
		rpc.Stderr = nil
	}
	dir := self.MkdirTemp()
	rpc.AccountsDir = filepath.Join(dir, "accounts")
	err := rpc.Start()
	if err != nil {
		panic(err)
	}
	return &AccountManager{rpc}
}

// Get a new Account that is not yet configured, but it is ready to be configured
// calling Account.Configure().
func (self *AcFactory) UnconfiguredAccount() *Account {
	account, err := self.NewAcManager().AddAccount()
	if err != nil {
		panic(err)
	}
	self.serialMutex.Lock()
	self.serial++
	serial := self.serial
	self.serialMutex.Unlock()

	if len(self.acCfg) != 0 {
		err = account.UpdateConfig(self.acCfg)
		if err != nil {
			panic(err)
		}
	}
	err = account.UpdateConfig(map[string]string{
		"addr":    fmt.Sprintf("acc%v.%v@localhost", serial, self.startTime),
		"mail_pw": fmt.Sprintf("password%v", serial),
	})
	if err != nil {
		panic(err)
	}
	return account
}

// Get a new account configured and with I/O already started.
func (self *AcFactory) OnlineAccount() *Account {
	account := self.UnconfiguredAccount()
	err := account.Configure()
	if err != nil {
		panic(err)
	}
	return account
}

// Get a new bot configured and with its account I/O already started. The bot is not running yet.
func (self *AcFactory) OnlineBot() *Bot {
	account := self.UnconfiguredAccount()
	addr, _ := account.GetConfig("addr")
	pass, _ := account.GetConfig("mail_pw")
	bot := NewBot(account)
	err := bot.Configure(addr, pass)
	if err != nil {
		panic(err)
	}
	return bot
}

// Get a new bot configured and already listening to new events/messages.
// It is ensured that Bot.IsRunning() is true for the returned bot.
func (self *AcFactory) RunningBot() *Bot {
	bot := self.OnlineBot()
	var err error
	go func() { err = bot.Run() }()
	for {
		if bot.IsRunning() {
			break
		}
		if err != nil {
			panic(err)
		}
	}
	return bot
}

// Wait for the next incoming message in the given account.
func (self *AcFactory) NextMsg(account *Account) (*MsgSnapshot, error) {
	event := self.WaitForEvent(account, EventIncomingMsg{}).(EventIncomingMsg)
	msg := Message{account, event.MsgId}
	return msg.Snapshot()
}

// Introduce two accounts to each other creating a 1:1 chat between them and exchanging messages.
func (self *AcFactory) IntroduceEachOther(account1, account2 *Account) {
	chat, err := self.CreateChat(account1, account2)
	if err != nil {
		panic(err)
	}
	_, err = chat.SendText("hi")
	if err != nil {
		panic(err)
	}
	self.WaitForEventInChat(account1, EventMsgsChanged{}, chat.Id)
	snapshot, _ := self.NextMsg(account2)
	if snapshot.Text != "hi" {
		panic("unexpected message: " + snapshot.Text)
	}

	chat = &Chat{account2, snapshot.ChatId}
	err = chat.Accept()
	if err != nil {
		panic(err)
	}
	_, err = chat.SendText("hello")
	if err != nil {
		panic(err)
	}
	self.WaitForEventInChat(account2, EventMsgsChanged{}, chat.Id)
	snapshot, _ = self.NextMsg(account1)
	if snapshot.Text != "hello" {
		panic("unexpected message: " + snapshot.Text)
	}
}

// Create a 1:1 chat with acc2 in the chatlist of acc1.
func (self *AcFactory) CreateChat(acc1, acc2 *Account) (*Chat, error) {
	addr2, err := acc2.GetConfig("configured_addr")
	if err != nil {
		return nil, err
	}

	contact, err := acc1.CreateContact(addr2, "")
	if err != nil {
		fmt.Println("WARNING: Failed to create contact with: ", addr2)
		return nil, err
	}

	chat, err := contact.CreateChat()
	if err != nil {
		return nil, err
	}

	return chat, nil
}

// Get a path to an image file that can be used for testing.
func (self *AcFactory) TestImage() string {
	acc := self.OnlineAccount()
	defer self.StopRpc(acc)
	chat, err := acc.Me().CreateChat()
	if err != nil {
		panic(err)
	}
	chatData, err := chat.BasicSnapshot()
	if err != nil {
		panic(err)
	}
	return chatData.ProfileImage
}

// Get a path to a Webxdc file that can be used for testing.
func (self *AcFactory) TestWebxdc() string {
	self.ensureTearUp()
	dir := self.MkdirTemp()
	path := filepath.Join(dir, "test.xdc")
	zipFile, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer zipFile.Close()
	writer := zip.NewWriter(zipFile)
	defer writer.Close()

	var files = []struct {
		Name, Body string
	}{
		{"index.html", `<html><head><script src="webxdc.js"></script></head><body>test</body></html>`},
		{"manifest.toml", `name = "TestApp"`},
	}
	for _, file := range files {
		f, err := writer.Create(file.Name)
		if err != nil {
			panic(err)
		}
		_, err = f.Write([]byte(file.Body))
		if err != nil {
			panic(err)
		}
	}

	err = writer.Close()
	if err != nil {
		panic(err)
	}

	return path
}

// Wait for an event of the same type as the given event, the event must belong to the chat
// with the given ChatId.
func (self *AcFactory) WaitForEventInChat(account *Account, event Event, chatId ChatId) Event {
	for {
		event = self.WaitForEvent(account, event)
		if getChatId(event) == chatId {
			return event
		}
	}
}

// Wait for an event of the same type as the given event.
func (self *AcFactory) WaitForEvent(account *Account, event Event) Event {
	eventChan := account.GetEventChannel()
	for {
		ev := <-eventChan
		if self.debug {
			fmt.Printf("Waiting for event %v, got: %v\n", event.eventType(), ev.eventType())
		}
		if ev.eventType() == event.eventType() {
			return ev
		}
	}
}

func (self *AcFactory) ensureTearUp() {
	if !self.tearUp {
		panic("TearUp() required")
	}
}

func getChatId(event Event) ChatId {
	var chatId ChatId
	switch ev := event.(type) {
	case EventMsgsChanged:
		chatId = ev.ChatId
	case EventReactionsChanged:
		chatId = ev.ChatId
	case EventIncomingMsg:
		chatId = ev.ChatId
	case EventMsgsNoticed:
		chatId = ev.ChatId
	case EventMsgDelivered:
		chatId = ev.ChatId
	case EventMsgFailed:
		chatId = ev.ChatId
	case EventMsgRead:
		chatId = ev.ChatId
	case EventChatModified:
		chatId = ev.ChatId
	case EventChatEphemeralTimerModified:
		chatId = ev.ChatId
	}
	return chatId
}
