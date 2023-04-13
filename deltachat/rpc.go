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
	cancel      context.CancelFunc
	events      map[AccountId]chan Event
	mu          sync.Mutex
}

func NewRpcIO() *RpcIO {
	return &RpcIO{Cmd: deltachatRpcServerBin, Stderr: os.Stderr}
}

// Implement Stringer.
func (self *RpcIO) String() string {
	return fmt.Sprintf("Rpc(AccountsDir=%#v)", self.AccountsDir)
}

func (self *RpcIO) Start() error {
	self.mu.Lock()
	defer self.mu.Unlock()

	if self.ctx != nil && self.ctx.Err() == nil {
		return &RpcRunningErr{}
	}

	self.ctx, self.cancel = context.WithCancel(context.Background())
	self.cmd = exec.CommandContext(self.ctx, self.Cmd)
	if self.AccountsDir != "" {
		self.cmd.Env = append(os.Environ(), "DC_ACCOUNTS_PATH="+self.AccountsDir)
	}
	self.cmd.Stderr = self.Stderr
	self.stdin, _ = self.cmd.StdinPipe()
	stdout, _ := self.cmd.StdoutPipe()
	if err := self.cmd.Start(); err != nil {
		self.cancel()
		return err
	}

	self.events = make(map[AccountId]chan Event)
	options := jrpc2.ClientOptions{OnNotify: self.onNotify}
	self.client = jrpc2.NewClient(channel.Line(stdout, self.stdin), &options)
	return nil
}

func (self *RpcIO) Stop() {
	self.mu.Lock()
	defer self.mu.Unlock()

	if self.ctx == nil {
		return
	}
	select {
	case <-self.ctx.Done():
		return
	default:
	}

	self.stdin.Close()
	self.cancel()
	self.cmd.Process.Wait() //nolint:errcheck
	for _, channel := range self.events {
		close(channel)
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
	self.mu.Lock()
	defer self.mu.Unlock()

	channel, ok := self.events[accountId]
	if !ok {
		channel = make(chan Event, 1000)
		self.events[accountId] = channel
	}
	return channel
}

func (self *RpcIO) onNotify(req *jrpc2.Request) {
	if req.Method() == "event" {
		var params _Params
		err := req.UnmarshalParams(&params)
		if err != nil {
			return
		}
		channel := self.getEventChannel(AccountId(params.ContextId))
		event := toEvent(params.Event)
		select {
		case <-self.ctx.Done():
			return
		default:
		}
		select {
		case channel <- event:
		default:
		}
	}
}

func toEvent(ev *_Event) Event {
	var event Event
	switch ev.Type {
	case eventTypeInfo:
		event = EventInfo{Msg: ev.Msg}
	case eventTypeSmtpConnected:
		event = EventSmtpConnected{Msg: ev.Msg}
	case eventTypeImapConnected:
		event = EventImapConnected{Msg: ev.Msg}
	case eventTypeSmtpMessageSent:
		event = EventSmtpMessageSent{Msg: ev.Msg}
	case eventTypeImapMessageDeleted:
		event = EventImapMessageDeleted{Msg: ev.Msg}
	case eventTypeImapMessageMoved:
		event = EventImapMessageMoved{Msg: ev.Msg}
	case eventTypeImapInboxIdle:
		event = EventImapInboxIdle{}
	case eventTypeNewBlobFile:
		event = EventNewBlobFile{File: ev.File}
	case eventTypeDeletedBlobFile:
		event = EventDeletedBlobFile{File: ev.File}
	case eventTypeWarning:
		event = EventWarning{Msg: ev.Msg}
	case eventTypeError:
		event = EventError{Msg: ev.Msg}
	case eventTypeErrorSelfNotInGroup:
		event = EventErrorSelfNotInGroup{Msg: ev.Msg}
	case eventTypeMsgsChanged:
		event = EventMsgsChanged{ChatId: ev.ChatId, MsgId: ev.MsgId}
	case eventTypeReactionsChanged:
		event = EventReactionsChanged{
			ChatId:    ev.ChatId,
			MsgId:     ev.MsgId,
			ContactId: ev.ContactId,
		}
	case eventTypeIncomingMsg:
		event = EventIncomingMsg{ChatId: ev.ChatId, MsgId: ev.MsgId}
	case eventTypeIncomingMsgBunch:
		event = EventIncomingMsgBunch{MsgIds: ev.MsgIds}
	case eventTypeMsgsNoticed:
		event = EventMsgsNoticed{ChatId: ev.ChatId}
	case eventTypeMsgDelivered:
		event = EventMsgDelivered{ChatId: ev.ChatId, MsgId: ev.MsgId}
	case eventTypeMsgFailed:
		event = EventMsgFailed{ChatId: ev.ChatId, MsgId: ev.MsgId}
	case eventTypeMsgRead:
		event = EventMsgRead{ChatId: ev.ChatId, MsgId: ev.MsgId}
	case eventTypeChatModified:
		event = EventChatModified{ChatId: ev.ChatId}
	case eventTypeChatEphemeralTimerModified:
		event = EventChatEphemeralTimerModified{
			ChatId: ev.ChatId,
			Timer:  ev.Timer,
		}
	case eventTypeContactsChanged:
		event = EventContactsChanged{ContactId: ev.ContactId}
	case eventTypeLocationChanged:
		event = EventLocationChanged{ContactId: ev.ContactId}
	case eventTypeConfigureProgress:
		event = EventConfigureProgress{Progress: ev.Progress, Comment: ev.Comment}
	case eventTypeImexProgress:
		event = EventImexProgress{Progress: ev.Progress}
	case eventTypeImexFileWritten:
		event = EventImexFileWritten{Path: ev.Path}
	case eventTypeSecurejoinInviterProgress:
		event = EventSecurejoinInviterProgress{
			ContactId: ev.ContactId,
			Progress:  ev.Progress,
		}
	case eventTypeSecurejoinJoinerProgress:
		event = EventSecurejoinJoinerProgress{
			ContactId: ev.ContactId,
			Progress:  ev.Progress,
		}
	case eventTypeConnectivityChanged:
		event = EventConnectivityChanged{}
	case eventTypeSelfavatarChanged:
		event = EventSelfavatarChanged{}
	case eventTypeWebxdcStatusUpdate:
		event = EventWebxdcStatusUpdate{
			MsgId:              ev.MsgId,
			StatusUpdateSerial: ev.StatusUpdateSerial,
		}
	case eventTypeWebxdcInstanceDeleted:
		event = EventWebxdcInstanceDeleted{MsgId: ev.MsgId}
	}
	return event
}
