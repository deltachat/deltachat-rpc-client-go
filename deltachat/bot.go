package deltachat

import (
	"fmt"
	"sync"
)

type EventHandler func(event map[string]any)
type NewMsgHandler func(msg *Message)

// Delta Chat bot that listen to events of a single account.
type Bot struct {
	Account         *Account
	newMsgHandler   NewMsgHandler
	handlerMap      map[string]EventHandler
	handlerMapMutex sync.RWMutex
}

// Implement Stringer.
func (self *Bot) String() string {
	return fmt.Sprintf("Bot(Account=%v)", self.Account.Id)
}

// Set an EventHandler for the given event type. Calling On() several times
// with the same event type will override the previously set EventHandler.
func (self *Bot) On(event string, handler EventHandler) {
	self.handlerMapMutex.Lock()
	self.handlerMap[event] = handler
	self.handlerMapMutex.Unlock()
}

// Set the NewMsgHandler for this bot.
func (self *Bot) OnNewMsg(handler NewMsgHandler) {
	self.newMsgHandler = handler
}

// Configure the bot's account.
func (self *Bot) Configure(addr string, password string) error {
	self.Account.SetConfig("bot", "1")
	self.Account.SetConfig("addr", addr)
	self.Account.SetConfig("mail_pw", password)
	return self.Account.Configure()
}

// Return true if the bot's account is configured, false otherwise.
func (self *Bot) IsConfigured() bool {
	configured, _ := self.Account.IsConfigured()
	return configured
}

// Set configuration value.
func (self *Bot) SetConfig(key string, value string) error {
	return self.Account.SetConfig(key, value)
}

// Get configuration value.
func (self *Bot) GetConfig(key string) (string, error) {
	return self.Account.GetConfig(key)
}

// Process events forever.
func (self *Bot) Run() {
	if self.IsConfigured() {
		self.Account.StartIO()
		self._processMessages() // Process old messages.
	}
	for {
		event := self.Account.WaitForEvent()
		if event == nil {
			break
		}
		self._onEvent(event)
		if event["type"].(string) == EVENT_INCOMING_MSG {
			self._processMessages()
		}
	}
}

func (self *Bot) _onEvent(event map[string]any) {
	eventType := event["type"].(string)
	self.handlerMapMutex.RLock()
	handler, ok := self.handlerMap[eventType]
	self.handlerMapMutex.RUnlock()
	if ok {
		handler(event)
	}
}

func (self *Bot) _processMessages() {
	msgs, _ := self.Account.FreshMsgsInArrivalOrder()
	for _, msg := range msgs {
		if self.newMsgHandler != nil {
			self.newMsgHandler(msg)
		}
		msg.MarkSeen()
	}
}

// Create a new Bot that will process events from the given account
func NewBot(account *Account) *Bot {
	return &Bot{Account: account, handlerMap: make(map[string]EventHandler)}
}

// Helper function to create a new Bot from the given AccountManager.
// The first available account will be used, a new account will be created if none exists.
func NewBotFromAccountManager(manager *AccountManager) *Bot {
	accounts, _ := manager.Accounts()
	var acc *Account
	if len(accounts) == 0 {
		acc, _ = manager.AddAccount()
	} else {
		acc = accounts[0]
	}
	return NewBot(acc)
}
