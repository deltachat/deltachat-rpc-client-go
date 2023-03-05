package deltachat

import (
	"fmt"
	"sync"
)

type EventHandler func(event *Event)
type NewMsgHandler func(msg *Message)

// Delta Chat bot that listen to events of a single account.
type Bot struct {
	Account         *Account
	newMsgHandler   NewMsgHandler
	handlerMap      map[string]EventHandler
	handlerMapMutex sync.RWMutex
	quitChan        chan struct{}
	running         bool
}

// Create a new Bot that will process events from the given account
func NewBot(account *Account) *Bot {
	return &Bot{Account: account, handlerMap: make(map[string]EventHandler), quitChan: make(chan struct{})}
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
	self.Account.UpdateConfig(
		map[string]string{
			"bot":     "1",
			"addr":    addr,
			"mail_pw": password,
		},
	)
	return self.Account.Configure()
}

// Return true if the bot's account is configured, false otherwise.
func (self *Bot) IsConfigured() bool {
	configured, _ := self.Account.IsConfigured()
	return configured
}

// Tweak several account configuration values in a batch.
func (self *Bot) UpdateConfig(config map[string]string) error {
	return self.Account.UpdateConfig(config)
}

// Set configuration value.
func (self *Bot) SetConfig(key string, value string) error {
	return self.Account.SetConfig(key, value)
}

// Get configuration value.
func (self *Bot) GetConfig(key string) (string, error) {
	return self.Account.GetConfig(key)
}

// This bot's contact object.
func (self *Bot) Me() *Contact {
	return self.Account.Me()
}

// Process events until Stop() is called.
func (self *Bot) Run() {
	self.running = true
	defer func() { self.running = false }()

	if self.IsConfigured() {
		self.Account.StartIO()
		self.processMessages() // Process old messages.
	}

	eventChan := self.Account.GetEventChannel()
	for {
		select {
		case <-self.quitChan:
			return
		case event, _ := <-eventChan:
			if event == nil {
				return
			}
			self.onEvent(event)
			if event.Type == EVENT_INCOMING_MSG {
				self.processMessages()
			}
		}
	}
}

// Stop processing events.
func (self *Bot) Stop() {
	if self.running {
		self.quitChan <- struct{}{}
	}
}

func (self *Bot) onEvent(event *Event) {
	self.handlerMapMutex.RLock()
	handler, ok := self.handlerMap[event.Type]
	self.handlerMapMutex.RUnlock()
	if ok {
		handler(event)
	}
}

func (self *Bot) processMessages() {
	msgs, err := self.Account.FreshMsgsInArrivalOrder()
	if err != nil {
		return
	}
	for _, msg := range msgs {
		if self.newMsgHandler != nil {
			self.newMsgHandler(msg)
		}
		msg.MarkSeen()
	}
}
