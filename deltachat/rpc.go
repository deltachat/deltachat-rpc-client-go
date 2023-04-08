package deltachat

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
)

type _Event struct {
	Type               eventType
	Msg                string
	File               string
	ChatId             ChatId
	MsgId              MsgId
	ContactId          ContactId
	MsgIds             []MsgId
	Timer              int
	Progress           uint
	Comment            string
	Path               string
	StatusUpdateSerial uint
}

type _Params struct {
	ContextId uint64
	Event     *_Event
}

// Delta Chat core RPC
type Rpc interface {
	Start() error
	Stop()
	GetEventChannel(accountId AccountId) <-chan Event
	Call(method string, params ...any) error
	CallResult(result any, method string, params ...any) error
	String() string
}

// Delta Chat core RPC working over IO
type RpcIO struct {
	Stderr      io.Writer
	AccountsDir string
	Cmd         string
	cmd         *exec.Cmd
	stdin       io.WriteCloser
	client      *jrpc2.Client
	ctx         context.Context
	events      map[AccountId]chan Event
	eventsMutex sync.Mutex
	closed      bool
}

func NewRpcIO() *RpcIO {
	return &RpcIO{Cmd: "deltachat-rpc-server", Stderr: os.Stderr, closed: true}
}

// Implement Stringer.
func (self *RpcIO) String() string {
	return fmt.Sprintf("Rpc(AccountsDir=%#v)", self.AccountsDir)
}

func (self *RpcIO) Start() error {
	if !self.closed {
		return fmt.Errorf("Rpc is already running")
	}
	self.closed = false
	self.cmd = exec.Command(self.Cmd)
	if self.AccountsDir != "" {
		self.cmd.Env = append(os.Environ(), "DC_ACCOUNTS_PATH="+self.AccountsDir)
	}
	self.cmd.Stderr = self.Stderr
	self.stdin, _ = self.cmd.StdinPipe()
	stdout, _ := self.cmd.StdoutPipe()
	if err := self.cmd.Start(); err != nil {
		self.closed = true
		return err
	}

	self.ctx = context.Background()
	self.events = make(map[AccountId]chan Event)
	options := jrpc2.ClientOptions{OnNotify: self.onNotify}
	self.client = jrpc2.NewClient(channel.Line(stdout, self.stdin), &options)
	return nil
}

func (self *RpcIO) Stop() {
	self.eventsMutex.Lock()
	defer self.eventsMutex.Unlock()
	if !self.closed {
		self.closed = true
		self.stdin.Close()
		self.cmd.Process.Wait()
		for _, channel := range self.events {
		loop:
			for {
				select {
				case <-channel:
					continue
				default:
					break loop
				}
			}
			close(channel)
		}
	}
}

func (self *RpcIO) GetEventChannel(accountId AccountId) <-chan Event {
	return self.getEventChannel(accountId)
}

func (self *RpcIO) Call(method string, params ...any) error {
	_, err := self.client.Call(self.ctx, method, params)
	return err
}

func (self *RpcIO) CallResult(result any, method string, params ...any) error {
	return self.client.CallResult(self.ctx, method, params, &result)
}

func (self *RpcIO) getEventChannel(accountId AccountId) chan Event {
	self.eventsMutex.Lock()
	defer self.eventsMutex.Unlock()
	channel, ok := self.events[accountId]
	if !ok {
		channel = make(chan Event, 10)
		self.events[accountId] = channel
	}
	return channel
}

func (self *RpcIO) onNotify(req *jrpc2.Request) {
	if req.Method() == "event" {
		var params _Params
		req.UnmarshalParams(&params)
		channel := self.getEventChannel(AccountId(params.ContextId))
		event := toEvent(params.Event)
		if !self.closed {
			go func() { channel <- event }()
		}
	}
}

func toEvent(ev *_Event) Event {
	var event Event
	switch ev.Type {
	case eventInfo:
		event = EventInfo{Msg: ev.Msg}
	case eventSmtpConnected:
		event = EventSmtpConnected{Msg: ev.Msg}
	case eventImapConnected:
		event = EventImapConnected{Msg: ev.Msg}
	case eventSmtpMessageSent:
		event = EventSmtpMessageSent{Msg: ev.Msg}
	case eventImapMessageDeleted:
		event = EventImapMessageDeleted{Msg: ev.Msg}
	case eventImapMessageMoved:
		event = EventImapMessageMoved{Msg: ev.Msg}
	case eventImapInboxIdle:
		event = EventImapInboxIdle{}
	case eventNewBlobFile:
		event = EventNewBlobFile{File: ev.File}
	case eventDeletedBlobFile:
		event = EventDeletedBlobFile{File: ev.File}
	case eventWarning:
		event = EventWarning{Msg: ev.Msg}
	case eventError:
		event = EventError{Msg: ev.Msg}
	case eventErrorSelfNotInGroup:
		event = EventErrorSelfNotInGroup{Msg: ev.Msg}
	case eventMsgsChanged:
		event = EventMsgsChanged{ChatId: ev.ChatId, MsgId: ev.MsgId}
	case eventReactionsChanged:
		event = EventReactionsChanged{
			ChatId:    ev.ChatId,
			MsgId:     ev.MsgId,
			ContactId: ev.ContactId,
		}
	case eventIncomingMsg:
		event = EventIncomingMsg{ChatId: ev.ChatId, MsgId: ev.MsgId}
	case eventIncomingMsgBunch:
		event = EventIncomingMsgBunch{MsgIds: ev.MsgIds}
	case eventMsgsNoticed:
		event = EventMsgsNoticed{ChatId: ev.ChatId}
	case eventMsgDelivered:
		event = EventMsgDelivered{ChatId: ev.ChatId, MsgId: ev.MsgId}
	case eventMsgFailed:
		event = EventMsgFailed{ChatId: ev.ChatId, MsgId: ev.MsgId}
	case eventMsgRead:
		event = EventMsgRead{ChatId: ev.ChatId, MsgId: ev.MsgId}
	case eventChatModified:
		event = EventChatModified{ChatId: ev.ChatId}
	case eventChatEphemeralTimerModified:
		event = EventChatEphemeralTimerModified{
			ChatId: ev.ChatId,
			Timer:  ev.Timer,
		}
	case eventContactsChanged:
		event = EventContactsChanged{ContactId: ev.ContactId}
	case eventLocationChanged:
		event = EventLocationChanged{ContactId: ev.ContactId}
	case eventConfigureProgress:
		event = EventConfigureProgress{Progress: ev.Progress, Comment: ev.Comment}
	case eventImexProgress:
		event = EventImexProgress{Progress: ev.Progress}
	case eventImexFileWritten:
		event = EventImexFileWritten{Path: ev.Path}
	case eventSecurejoinInviterProgress:
		event = EventSecurejoinInviterProgress{
			ContactId: ev.ContactId,
			Progress:  ev.Progress,
		}
	case eventSecurejoinJoinerProgress:
		event = EventSecurejoinJoinerProgress{
			ContactId: ev.ContactId,
			Progress:  ev.Progress,
		}
	case eventConnectivityChanged:
		event = EventConnectivityChanged{}
	case eventSelfavatarChanged:
		event = EventSelfavatarChanged{}
	case eventWebxdcStatusUpdate:
		event = EventWebxdcStatusUpdate{
			MsgId:              ev.MsgId,
			StatusUpdateSerial: ev.StatusUpdateSerial,
		}
	case eventWebxdcInstanceDeleted:
		event = EventWebxdcInstanceDeleted{MsgId: ev.MsgId}
	}
	return event
}
