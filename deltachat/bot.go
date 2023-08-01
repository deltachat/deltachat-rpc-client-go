package deltachat

import (
	"context"
	"fmt"
	"sync"

	"github.com/deltachat/deltachat-rpc-client-go/deltachat/option"
)

type EventHandler func(bot *Bot, event Event)
type NewMsgHandler func(bot *Bot, msgId MsgId)

// BotRunningErr is returned by Bot.Run() if the Bot is already running
type BotRunningErr struct{}

func (self *BotRunningErr) Error() string {
	return "bot is already running"
}

// Delta Chat bot that listen to events of a single account.
type Bot struct {
	Rpc             *Rpc
	AccountId       AccountId
	newMsgHandler   NewMsgHandler
	handlerMap      map[eventType]EventHandler
	handlerMapMutex sync.RWMutex
	ctxMutex        sync.Mutex
	ctx             context.Context
	stop            context.CancelFunc
}

// Create a new Bot that will process events from the given account.
// If the given account id is zero, the first available account will
// be used, or a new account will be created if none exists.
func NewBot(rpc *Rpc, accountId AccountId) *Bot {
	if accountId == 0 {
		accounts, _ := rpc.GetAllAccountIds()
		if len(accounts) == 0 {
			accountId, _ = rpc.AddAccount()
		} else {
			accountId = accounts[0]
		}
	}
	return &Bot{Rpc: rpc, AccountId: accountId, handlerMap: make(map[eventType]EventHandler)}
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
	err := self.Rpc.BatchSetConfig(
		self.AccountId,
		map[string]option.Option[string]{
			"bot":     option.Some("1"),
			"addr":    option.Some(addr),
			"mail_pw": option.Some(password),
		},
	)
	if err != nil {
		return err
	}
	return self.Rpc.Configure(self.AccountId)
}

// Return true if the bot's account is configured, false otherwise.
func (self *Bot) IsConfigured() bool {
	configured, _ := self.Rpc.IsConfigured(self.AccountId)
	return configured
}

// Tweak several account configuration values in a batch.
func (self *Bot) UpdateConfig(config map[string]option.Option[string]) error {
	return self.Rpc.BatchSetConfig(self.AccountId, config)
}

// Set account configuration value.
func (self *Bot) SetConfig(key string, value option.Option[string]) error {
	return self.Rpc.SetConfig(self.AccountId, key, value)
}

// Get account configuration value.
func (self *Bot) GetConfig(key string) (option.Option[string], error) {
	return self.Rpc.GetConfig(self.AccountId, key)
}

// Set UI-specific configuration value in the bot's account.
// This is useful for custom 3rd party settings set by bot programs.
func (self *Bot) SetUiConfig(key string, value option.Option[string]) error {
	return self.SetConfig("ui."+key, value)
}

// Get custom UI-specific configuration value set with SetUiConfig().
func (self *Bot) GetUiConfig(key string) (option.Option[string], error) {
	return self.GetConfig("ui." + key)
}

// Process events until Stop() is called. If the bot is already running, BotRunningErr is returned.
func (self *Bot) Run() error {
	self.ctxMutex.Lock()
	if self.ctx != nil && self.ctx.Err() == nil {
		self.ctxMutex.Unlock()
		return &BotRunningErr{}
	}
	self.ctx, self.stop = context.WithCancel(context.Background())
	self.ctxMutex.Unlock()

	if self.IsConfigured() {
		self.Rpc.StartIo(self.AccountId) //nolint:errcheck
		self.processMessages()           // Process old messages.
	}

	eventChan := make(chan Event)
	go func() {
		for {
			rpc := &Rpc{Context: self.ctx, Transport: self.Rpc.Transport}
			accId, event, err := rpc.GetNextEvent()
			if err != nil {
				close(eventChan)
				break
			}
			if accId == self.AccountId {
				eventChan <- event
			}
		}
	}()

	for {
		event, ok := <-eventChan
		if !ok {
			self.Stop()
			return nil
		}
		self.onEvent(event)
		if event.eventType() == eventTypeIncomingMsg {
			self.processMessages()
		}
	}
}

// Return true if bot is running (Bot.Run() is running) or false otherwise.
func (self *Bot) IsRunning() bool {
	self.ctxMutex.Lock()
	defer self.ctxMutex.Unlock()
	return self.ctx != nil && self.ctx.Err() == nil
}

// Stop processing events.
func (self *Bot) Stop() {
	self.ctxMutex.Lock()
	defer self.ctxMutex.Unlock()
	if self.ctx != nil && self.ctx.Err() == nil {
		self.stop()
	}
}

func (self *Bot) onEvent(event Event) {
	self.handlerMapMutex.RLock()
	handler, ok := self.handlerMap[event.eventType()]
	self.handlerMapMutex.RUnlock()
	if ok {
		handler(self, event)
	}
}

func (self *Bot) processMessages() {
	msgIds, err := self.Rpc.GetNextMsgs(self.AccountId)
	if err != nil {
		return
	}
	for _, msgId := range msgIds {
		self.Rpc.SetConfig(self.AccountId, "last_msg_id", option.Some(fmt.Sprintf("%v", msgId))) //nolint:errcheck
		if self.newMsgHandler != nil {
			self.newMsgHandler(self, msgId)
		}
	}
}
