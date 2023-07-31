package deltachat

import (
	"fmt"
	"sync"

	"github.com/deltachat/deltachat-rpc-client-go/deltachat/option"
)

type EventHandler func(bot *Bot, event Event)
type NewMsgHandler func(bot *Bot, msgId MsgId)

// BotAlreadyStartedErr is returned by Bot.Run() if the Bot is already running
type BotAlreadyStartedErr struct{}

func (self *BotAlreadyStartedErr) Error() string {
	return "bot was already started"
}

// Delta Chat bot that listen to events of a single account.
type Bot struct {
	Rpc             *Rpc
	AccountId       AccountId
	newMsgHandler   NewMsgHandler
	handlerMap      map[eventType]EventHandler
	handlerMapMutex sync.RWMutex
	startedMutex    sync.Mutex
	started         bool
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

// Process events until Stop() is called. Can only be executed once, any subsequent call,
// even after the bot is stopped, will return an error.
func (self *Bot) Run() error {
	self.startedMutex.Lock()
	if self.started {
		self.startedMutex.Unlock()
		return &BotAlreadyStartedErr{}
	}
	self.started = true
	self.startedMutex.Unlock()

	if self.IsConfigured() {
		self.Rpc.StartIo(self.AccountId) //nolint:errcheck
		self.processMessages()           // Process old messages.
	}

	eventChan := make(chan Event)
	go func() {
		for {
			accId, event, err := self.Rpc.GetNextEvent()
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
			return nil
		}
		self.onEvent(event)
		if event.eventType() == eventTypeIncomingMsg {
			self.processMessages()
		}
	}
}

// Close the Rpc's Transport and stop processing events.
// The Rpc and bot should not be used anymore after calling Stop().
func (self *Bot) Stop() {
	self.Rpc.Transport.Close()
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
