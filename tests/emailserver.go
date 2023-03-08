package tests

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/deltachat/deltachat-rpc-client-go/deltachat"
)

var serial int

type EmailServer struct {
	cmd         *exec.Cmd
	args        []string
	jar         string
	rpc         deltachat.Rpc
	manager     *deltachat.AccountManager
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
	rpc := deltachat.NewRpcIO()
	dir, _ := os.MkdirTemp("", "")
	rpc.AccountsDir = filepath.Join(dir, "accounts")
	rpc.Start()
	server := &EmailServer{jar: jar, args: arg, rpc: rpc, manager: &deltachat.AccountManager{rpc}}
	server.accountsDir = dir
	return server, server.check()
}

func (self *EmailServer) AccountManager() *deltachat.AccountManager {
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

func (self *EmailServer) GetUnconfiguredAccount() *deltachat.Account {
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

func (self *EmailServer) GetOnlineBot() (*deltachat.Bot, error) {
	account := self.GetUnconfiguredAccount()
	addr, _ := account.GetConfig("addr")
	pass, _ := account.GetConfig("mail_pw")
	bot := deltachat.NewBot(account)
	err := bot.Configure(addr, pass)
	if err != nil {
		return nil, err
	}
	go bot.Run()
	return bot, nil
}

func (self *EmailServer) GetOnlineAccount() (*deltachat.Account, error) {
	account := self.GetUnconfiguredAccount()
	return account, account.Configure()
}

func (self *EmailServer) GetNextMsg(account *deltachat.Account) (*deltachat.MsgSnapshot, error) {
	event := account.WaitForEvent(deltachat.EVENT_INCOMING_MSG)
	msg := deltachat.Message{account, event.MsgId}
	return msg.Snapshot()
}

func (self *EmailServer) check() error {
	cmd := exec.Command("java", "-jar", self.jar)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
