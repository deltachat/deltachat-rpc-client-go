// Package acfactory facilitates unit testing Delta Chat clients/bots.
//
// This package must be used in conjunction with a test mail server service, for example:
// https://github.com/deltachat/mail-server-tester
//
// Typical usage is as follows:
//
//	import (
//		"testing"

//		"github.com/deltachat/deltachat-rpc-client-go/acfactory"
//	)

//	func TestMain(m *testing.M) {
//		cfg := map[string]string{
//			"mail_server":   "localhost",
//			"send_server":   "localhost",
//			"mail_port":     "3143",
//			"send_port":     "3025",
//			"mail_security": "3",
//			"send_security": "3",
//		}
//		acfactory.TearUp(cfg)
//		defer acfactory.TearDown()
//		m.Run()
//	}

package acfactory

import (
	"os"

	"github.com/deltachat/deltachat-rpc-client-go/deltachat"
)

var acf *deltachat.AcFactory

// Prepare the AcFactory, defaultAcConfig is the default settings to apply to
// the new accounts created with UnconfiguredAccount(), OnlineAccount(), OnlineBot()
// and RunningBot().
//
// If the test mail server has not standard configuration, you should set the custom configuration
// here.
func TearUp(defaultAcConfig map[string]string) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		panic(err)
	}
	acf = &deltachat.AcFactory{}
	acf.TearUp(defaultAcConfig, dir, os.Getenv("TEST_DEBUG") == "1")
}

// Do cleanup, removing temporary directories and files created by the configured test accounts.
// Usually TearDown() is called with defer immediately after the call to TearUp().
func TearDown() {
	acf.TearDown()
}

// Stop the Rpc of the given Account, Bot or AccountManager.
func StopRpc(accountOrBot any) {
	acf.StopRpc(accountOrBot)
}

// MkdirTemp creates a new temporary directory. The directory is automatically removed on TearDown().
func MkdirTemp() string {
	return acf.MkdirTemp()
}

// Create a new AccountManager.
func NewAcManager() *deltachat.AccountManager {
	return acf.NewAcManager()
}

// Get a new Account that is not yet configured, but it is ready to be configured
// calling Account.Configure().
func UnconfiguredAccount() *deltachat.Account {
	return acf.UnconfiguredAccount()
}

// Get a new account configured and with I/O already started.
func OnlineAccount() *deltachat.Account {
	return acf.OnlineAccount()
}

// Get a new bot configured and already listening to new events/messages.
// It is ensured that Bot.IsRunning() is true for the returned bot.
func OnlineBot() *deltachat.Bot {
	return acf.OnlineBot()
}

// Wait for the next incoming message in the given account.
func NextMsg(account *deltachat.Account) (*deltachat.MsgSnapshot, error) {
	return acf.NextMsg(account)
}

// Introduce two accounts to each other creating a 1:1 chat between them and exchanging messages.
func IntroduceEachOther(account1, account2 *deltachat.Account) {
	acf.IntroduceEachOther(account1, account2)
}

// Create a 1:1 chat with acc2 in the chatlist of acc1.
func CreateChat(acc1, acc2 *deltachat.Account) (*deltachat.Chat, error) {
	return acf.CreateChat(acc1, acc2)
}

// Get a path to an image file that can be used for testing.
func TestImage() string {
	return acf.TestImage()
}

// Get a path to a Webxdc file that can be used for testing.
func TestWebxdc() string {
	return acf.TestWebxdc()
}

// Wait for an event of the same type as the given event, the event must belong to the chat
// with the given ChatId.
func WaitForEventInChat(account *deltachat.Account, event deltachat.Event, chatId deltachat.ChatId) deltachat.Event {
	return acf.WaitForEventInChat(account, event, chatId)
}

// Wait for an event of the same type as the given event.
func WaitForEvent(account *deltachat.Account, event deltachat.Event) deltachat.Event {
	return acf.WaitForEvent(account, event)
}
