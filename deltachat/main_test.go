package deltachat

import (
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
	serialMutex sync.Mutex
	tempDir     string
}

func InitAcFactory() {
	if acfactory == nil {
		dir, err := os.MkdirTemp("", "")
		if err != nil {
			panic(err)
		}
		acfactory = &AcFactory{tempDir: dir, serial: time.Now().Unix()}
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

	account.UpdateConfig(map[string]string{
		"mail_server":   "localhost",
		"send_server":   "localhost",
		"mail_port":     "3143",
		"send_port":     "3025",
		"mail_security": "3",
		"send_security": "3",
		"addr":          fmt.Sprintf("acc%v@localhost", serial),
		"mail_pw":       fmt.Sprintf("password%v", serial),
	})
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
	event := WaitForEvent(account, EVENT_INCOMING_MSG)
	msg := Message{account, event.MsgId}
	return msg.Snapshot()
}

func (self *AcFactory) IntroduceEachOther(account1, account2 *Account) {
	chat, _ := account1.CreateChat(account2)
	chat.SendText("hi")
	waitForEvent(account1, EVENT_MSGS_CHANGED, chat.Id)
	snapshot, _ := self.GetNextMsg(account2)
	if snapshot.Text != "hi" {
		panic("unexpected message: " + snapshot.Text)
	}

	chat = &Chat{account2, snapshot.ChatId}
	chat.Accept()
	chat.SendText("hello")
	waitForEvent(account2, EVENT_MSGS_CHANGED, chat.Id)
	snapshot, _ = self.GetNextMsg(account1)
	if snapshot.Text != "hello" {
		panic("unexpected message: " + snapshot.Text)
	}
}

func TestMain(m *testing.M) {
	InitAcFactory()
	defer acfactory.TearDown()
	m.Run()
}

func waitForEvent(account *Account, eventType string, chatId ChatId) *Event {
	for {
		event := WaitForEvent(account, eventType)
		if event.ChatId == chatId {
			return event
		}
	}
}

func WaitForEvent(account *Account, eventType string) *Event {
	eventChan := account.GetEventChannel()
	debug := os.Getenv("TEST_DEBUG") == "1"
	for {
		event := <-eventChan
		if debug {
			fmt.Printf("Waiting for event %v, got: %v\n", eventType, event.Type)
		}
		if event.Type == eventType {
			return event
		}
	}
}
