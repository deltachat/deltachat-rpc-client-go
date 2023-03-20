package deltachat

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var serial int
var server *EmailServer

type EmailServer struct {
	cmd         *exec.Cmd
	args        []string
	jar         string
	rpc         Rpc
	manager     *AccountManager
	accountsDir string
}

// create a new EmailServer instance. The given arguments will be passed down to the GreenMail process.
func NewEmailServer(arg ...string) (*EmailServer, error) {
	jar := os.Getenv("GREENMAIL_JAR")
	if jar == "" {
		jar = "greenmail-standalone.jar"
	}
	if len(arg) == 0 {
		arg = append(arg, "-Dgreenmail.setup.test.all", "-Dgreenmail.auth.disabled")
	}
	rpc := NewRpcIO()
	dir, _ := os.MkdirTemp("", "")
	rpc.AccountsDir = filepath.Join(dir, "accounts")
	err := rpc.Start()
	if err != nil {
		return nil, err
	}
	server := &EmailServer{jar: jar, args: arg, rpc: rpc, manager: &AccountManager{rpc}}
	server.accountsDir = dir
	err = server.check()
	if err != nil {
		return nil, err
	}
	return server, nil
}

func (self *EmailServer) AccountManager() *AccountManager {
	return self.manager
}

func (self *EmailServer) Start() error {
	args := append(self.args, "-jar", self.jar)
	self.cmd = exec.Command("java", args...)
	stdout, _ := self.cmd.StdoutPipe()
	if err := self.cmd.Start(); err != nil {
		return err
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Starting GreenMail API server") {
			break
		}
	}
	return nil
}

func (self *EmailServer) Stop() {
	self.rpc.Stop()
	self.cmd.Process.Kill()
	self.cmd.Process.Wait()
	os.RemoveAll(self.accountsDir)
}

func (self *EmailServer) GetUnconfiguredAccount() *Account {
	serial++
	account, _ := self.manager.AddAccount()
	account.UpdateConfig(map[string]string{
		"mail_server":             "localhost",
		"send_server":             "localhost",
		"mail_port":               "3143",
		"send_port":               "3025",
		"mail_security":           "3",
		"send_security":           "3",
		"smtp_certificate_checks": "3",
		"imap_certificate_checks": "3",
		"addr":                    fmt.Sprintf("account%v@localhost", serial),
		"mail_pw":                 fmt.Sprintf("password%v", serial),
	})
	return account
}

func (self *EmailServer) GetOnlineBot() (*Bot, error) {
	account := self.GetUnconfiguredAccount()
	addr, _ := account.GetConfig("addr")
	pass, _ := account.GetConfig("mail_pw")
	bot := NewBot(account)
	err := bot.Configure(addr, pass)
	if err != nil {
		return nil, err
	}
	go bot.Run()
	return bot, nil
}

func (self *EmailServer) GetOnlineAccount() (*Account, error) {
	account := self.GetUnconfiguredAccount()
	return account, account.Configure()
}

func (self *EmailServer) GetNextMsg(account *Account) (*MsgSnapshot, error) {
	event := account.WaitForEvent(EVENT_INCOMING_MSG)
	msg := Message{account, event.MsgId}
	return msg.Snapshot()
}

func (self *EmailServer) IntroduceEachOther(account1, account2 *Account) {
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

func (self *EmailServer) check() error {
	cmd := exec.Command("java", "-jar", self.jar)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func TestMain(m *testing.M) {
	var err error
	server, err = NewEmailServer()
	if err != nil {
		panic(err)
	}
	defer server.Stop()
	server.Start()
	m.Run()
}

func waitForEvent(account *Account, eventType string, chatId uint64) *Event {
	for {
		event := account.WaitForEvent(eventType)
		if event.ChatId == chatId {
			return event
		}
	}
}
