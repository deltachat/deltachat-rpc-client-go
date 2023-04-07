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
		var event Event
		switch params.Event.Type {
		case eventInfo:
			event = EventInfo{Msg: params.Event.Msg}
		case eventSmtpConnected:
			event = EventSmtpConnected{Msg: params.Event.Msg}
		case eventImapConnected:
			event = EventImapConnected{Msg: params.Event.Msg}
		case eventSmtpMessageSent:
			event = EventSmtpMessageSent{Msg: params.Event.Msg}
		case eventImapMessageDeleted:
			event = EventImapMessageDeleted{Msg: params.Event.Msg}
		case eventImapMessageMoved:
			event = EventImapMessageMoved{Msg: params.Event.Msg}
		case eventImapInboxIdle:
			event = EventImapInboxIdle{}
		case eventNewBlobFile:
			event = EventNewBlobFile{File: params.Event.File}
		case eventDeletedBlobFile:
			event = EventDeletedBlobFile{File: params.Event.File}
		case eventWarning:
			event = EventWarning{Msg: params.Event.Msg}
		case eventError:
			event = EventError{Msg: params.Event.Msg}
		case eventErrorSelfNotInGroup:
			event = EventErrorSelfNotInGroup{Msg: params.Event.Msg}
		case eventMsgsChanged:
			event = EventMsgsChanged{ChatId: params.Event.ChatId, MsgId: params.Event.MsgId}
		case eventReactionsChanged:
			event = EventReactionsChanged{
				ChatId:    params.Event.ChatId,
				MsgId:     params.Event.MsgId,
				ContactId: params.Event.ContactId,
			}
		case eventIncomingMsg:
			event = EventIncomingMsg{ChatId: params.Event.ChatId, MsgId: params.Event.MsgId}
		case eventIncomingMsgBunch:
			event = EventIncomingMsgBunch{MsgIds: params.Event.MsgIds}
		case eventMsgsNoticed:
			event = EventMsgsNoticed{ChatId: params.Event.ChatId}
		case eventMsgDelivered:
			event = EventMsgDelivered{ChatId: params.Event.ChatId, MsgId: params.Event.MsgId}
		case eventMsgFailed:
			event = EventMsgFailed{ChatId: params.Event.ChatId, MsgId: params.Event.MsgId}
		case eventMsgRead:
			event = EventMsgRead{ChatId: params.Event.ChatId, MsgId: params.Event.MsgId}
		case eventChatModified:
			event = EventChatModified{ChatId: params.Event.ChatId}
		case eventChatEphemeralTimerModified:
			event = EventChatEphemeralTimerModified{
				ChatId: params.Event.ChatId,
				Timer:  params.Event.Timer,
			}
		case eventContactsChanged:
			event = EventContactsChanged{ContactId: params.Event.ContactId}
		case eventLocationChanged:
			event = EventLocationChanged{ContactId: params.Event.ContactId}
		case eventConfigureProgress:
			event = EventConfigureProgress{Progress: params.Event.Progress, Comment: params.Event.Comment}
		case eventImexProgress:
			event = EventImexProgress{Progress: params.Event.Progress}
		case eventImexFileWritten:
			event = EventImexFileWritten{Path: params.Event.Path}
		case eventSecurejoinInviterProgress:
			event = EventSecurejoinInviterProgress{
				ContactId: params.Event.ContactId,
				Progress:  params.Event.Progress,
			}
		case eventSecurejoinJoinerProgress:
			event = EventSecurejoinJoinerProgress{
				ContactId: params.Event.ContactId,
				Progress:  params.Event.Progress,
			}
		case eventConnectivityChanged:
			event = EventConnectivityChanged{}
		case eventSelfavatarChanged:
			event = EventSelfavatarChanged{}
		case eventWebxdcStatusUpdate:
			event = EventWebxdcStatusUpdate{
				MsgId:              params.Event.MsgId,
				StatusUpdateSerial: params.Event.StatusUpdateSerial,
			}
		case eventWebxdcInstanceDeleted:
			event = EventWebxdcInstanceDeleted{MsgId: params.Event.MsgId}
		}
		if !self.closed {
			go func() { channel <- event }()
		}
	}
}
