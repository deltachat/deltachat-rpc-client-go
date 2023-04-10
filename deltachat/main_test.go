package deltachat

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

var acfactory *AcFactory

type AcFactory struct {
	serial      int64
	startTime   int64
	serialMutex sync.Mutex
	tempDir     string
}

func InitAcFactory() {
	if acfactory == nil {
		dir, err := os.MkdirTemp("", "")
		if err != nil {
			panic(err)
		}
		acfactory = &AcFactory{tempDir: dir, startTime: time.Now().Unix()}
	}
}

func (self *AcFactory) TearDown() {
	os.RemoveAll(self.tempDir)
}

func (self *AcFactory) NewAcManager() *AccountManager {
	rpc := NewRpcIO()
	if os.Getenv("TEST_DEBUG") != "1" {
		rpc.Stderr = nil
	}
	dir, err := os.MkdirTemp(self.tempDir, "")
	if err != nil {
		panic(err)
	}
	rpc.AccountsDir = filepath.Join(dir, "accounts")
	err = rpc.Start()
	if err != nil {
		panic(err)
	}
	return &AccountManager{rpc}
}

func (self *AcFactory) GetUnconfiguredAccount() *Account {
	account, err := self.NewAcManager().AddAccount()
	if err != nil {
		panic(err)
	}

	self.serialMutex.Lock()
	self.serial++
	serial := self.serial
	self.serialMutex.Unlock()

	addr := fmt.Sprintf("acc%v.%v@localhost", serial, self.startTime)
	account.UpdateConfig(map[string]string{
		"mail_server":   "localhost",
		"send_server":   "localhost",
		"mail_port":     "3143",
		"send_port":     "3025",
		"mail_security": "3",
		"send_security": "3",
		"addr":          addr,
		"mail_pw":       fmt.Sprintf("password%v", serial),
	})
	if os.Getenv("TEST_DEBUG") == "1" {
		fmt.Println("Created new account: ", addr)
	}
	return account
}

func (self *AcFactory) GetOnlineBot() *Bot {
	account := self.GetUnconfiguredAccount()
	addr, _ := account.GetConfig("addr")
	pass, _ := account.GetConfig("mail_pw")
	bot := NewBot(account)
	err := bot.Configure(addr, pass)
	if err != nil {
		panic(err)
	}
	go bot.Run()
	return bot
}

func (self *AcFactory) GetOnlineAccount() *Account {
	account := self.GetUnconfiguredAccount()
	account.Configure()
	return account
}

func (self *AcFactory) GetNextMsg(account *Account) (*MsgSnapshot, error) {
	event := WaitForEvent(account, EventIncomingMsg{}).(EventIncomingMsg)
	msg := Message{account, event.MsgId}
	return msg.Snapshot()
}

func (self *AcFactory) IntroduceEachOther(account1, account2 *Account) {
	chat, _ := self.CreateChat(account1, account2)
	chat.SendText("hi")
	waitForEvent(account1, EventMsgsChanged{}, chat.Id)
	snapshot, _ := self.GetNextMsg(account2)
	if snapshot.Text != "hi" {
		panic("unexpected message: " + snapshot.Text)
	}

	chat = &Chat{account2, snapshot.ChatId}
	chat.Accept()
	chat.SendText("hello")
	waitForEvent(account2, EventMsgsChanged{}, chat.Id)
	snapshot, _ = self.GetNextMsg(account1)
	if snapshot.Text != "hello" {
		panic("unexpected message: " + snapshot.Text)
	}
}

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

func (self *AcFactory) GetTestImage() string {
	acc := acfactory.GetOnlineAccount()
	defer acc.Manager.Rpc.Stop()
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

func (self *AcFactory) GetTestWebxdc() string {
	dir, err := os.MkdirTemp(self.tempDir, "")
	if err != nil {
		panic(err)
	}

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

func TestMain(m *testing.M) {
	InitAcFactory()
	defer acfactory.TearDown()
	m.Run()
}

func waitForEvent(account *Account, event Event, chatId ChatId) Event {
	for {
		event = WaitForEvent(account, event)
		var chatId2 ChatId
		switch ev := event.(type) {
		case EventMsgsChanged:
			chatId2 = ev.ChatId
		case EventReactionsChanged:
			chatId2 = ev.ChatId
		case EventIncomingMsg:
			chatId2 = ev.ChatId
		case EventMsgsNoticed:
			chatId2 = ev.ChatId
		case EventMsgDelivered:
			chatId2 = ev.ChatId
		case EventMsgFailed:
			chatId2 = ev.ChatId
		case EventMsgRead:
			chatId2 = ev.ChatId
		case EventChatModified:
			chatId2 = ev.ChatId
		case EventChatEphemeralTimerModified:
			chatId2 = ev.ChatId
		}
		if chatId2 == chatId {
			return event
		}
	}
}

func WaitForEvent(account *Account, event Event) Event {
	eventChan := account.GetEventChannel()
	debug := os.Getenv("TEST_DEBUG") == "1"
	for {
		ev := <-eventChan
		if debug {
			fmt.Printf("Waiting for event %v, got: %v\n", event.eventType(), ev.eventType())
		}
		if ev.eventType() == event.eventType() {
			return ev
		}
	}
}
