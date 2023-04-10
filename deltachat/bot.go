package deltachat

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// ErrRpcRunning is returned by Bot.Run() if the Bot is already running
var ErrBotRunning = errors.New("bot is already running")

type EventHandler func(event Event)
type NewMsgHandler func(msg *Message)

// Delta Chat bot that listen to events of a single account.
type Bot struct {
	Account         *Account
	newMsgHandler   NewMsgHandler
	handlerMap      map[eventType]EventHandler
	handlerMapMutex sync.RWMutex
	ctx             context.Context
	close           context.CancelFunc
}

// Create a new Bot that will process events from the given account
func NewBot(account *Account) *Bot {
	return &Bot{Account: account, handlerMap: make(map[eventType]EventHandler)}
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
func (self *Bot) On(event Event, handler EventHandler) {
	self.handlerMapMutex.Lock()
	self.handlerMap[event.eventType()] = handler
	self.handlerMapMutex.Unlock()
}

// Remove EventHandler for the given event type.
func (self *Bot) RemoveEventHandler(event Event) {
	self.handlerMapMutex.Lock()
	delete(self.handlerMap, event.eventType())
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
func (self *Bot) Run() error {
	if self.ctx != nil && self.ctx.Err() == nil {
		return ErrRpcRunning
	}
	self.ctx, self.close = context.WithCancel(context.Background())
	defer func() { self.close() }()

	if self.IsConfigured() {
		self.Account.StartIO()
		self.processMessages() // Process old messages.
	}

	for {
		event := self.Account.GetNextEvent(self.ctx)
		if event == nil {
			break
		}
		self.onEvent(event)
		if event.eventType() == eventTypeIncomingMsg {
			self.processMessages()
		}
	}
	return nil
}

// Stop processing events.
func (self *Bot) Stop() {
	if self.ctx == nil {
		return
	}
	select {
	case <-self.ctx.Done():
	default:
		self.close()
	}
}

func (self *Bot) onEvent(event Event) {
	self.handlerMapMutex.RLock()
	handler, ok := self.handlerMap[event.eventType()]
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
