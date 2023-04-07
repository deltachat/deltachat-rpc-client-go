package deltachat

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
)

var ErrRpcRunning = errors.New("rpc is already running")

// Delta Chat core Event
type Event struct {
	Type               EventType
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
	Event     *Event
}

// Delta Chat core RPC
type Rpc interface {
	Start() error
	Stop()
	GetEventChannel(accountId AccountId) <-chan *Event
	Call(method string, params ...any) error
	CallResult(result any, method string, params ...any) error
	String() string
}

// Delta Chat core RPC working over IO
type RpcIO struct {
	Stderr        io.Writer
	AccountsDir   string
	Cmd           string
	EventBuffer   int
	cmd           *exec.Cmd
	stdin         io.WriteCloser
	client        *jrpc2.Client
	ctx           context.Context
	cancel        context.CancelFunc
	accountEvents map[AccountId]chan *Event
	events        chan _Params
	mu            sync.Mutex
}

var _ Rpc = &RpcIO{}

func NewRpcIO() *RpcIO {
	return &RpcIO{
		Cmd:         "deltachat-rpc-server",
		Stderr:      os.Stderr,
		EventBuffer: 10,
	}
}

// Implement Stringer.
func (self *RpcIO) String() string {
	return fmt.Sprintf("Rpc(AccountsDir=%#v)", self.AccountsDir)
}

func (self *RpcIO) Start() error {
	self.mu.Lock()
	defer self.mu.Unlock()

	if self.ctx != nil && self.ctx.Err() != nil {
		return ErrRpcRunning
	}

	self.ctx, self.cancel = context.WithCancel(context.Background())
	self.cmd = exec.CommandContext(self.ctx, self.Cmd)
	if self.AccountsDir != "" {
		self.cmd.Env = append(os.Environ(), "DC_ACCOUNTS_PATH="+self.AccountsDir)
	}
	self.cmd.Stderr = self.Stderr
	self.stdin, _ = self.cmd.StdinPipe()
	stdout, _ := self.cmd.StdoutPipe()

	self.events = make(chan _Params, self.EventBuffer)
	self.accountEvents = make(map[AccountId]chan *Event)
	options := jrpc2.ClientOptions{OnNotify: self.onNotify}
	self.client = jrpc2.NewClient(channel.Line(stdout, self.stdin), &options)

	go func() {
		for {
			select {
			case <-self.ctx.Done():
				return
			case params := <-self.events:
				channel := self.getEventChannel(AccountId(params.ContextId))

				var sent bool
				for !sent {
					select {
					case <-self.ctx.Done():
						return
					case channel <- params.Event:
						sent = true
					case <-time.After(time.Second * 1):
						if self.Stderr != nil {
							fmt.Fprintf(self.Stderr, "RPC error: account channel is full, retrying! AccountId:%d, EventType:%s\n", params.ContextId, params.Event.Type)
						}
					}
				}
			}
		}
	}()

	if err := self.cmd.Start(); err != nil {
		self.cancel()
		return err
	}

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
	self.cmd.Process.Wait()
	close(self.events)

	for _, channel := range self.accountEvents {
		close(channel)
	}
}

func (self *RpcIO) GetEventChannel(accountId AccountId) <-chan *Event {
	return self.getEventChannel(accountId)
}

func (self *RpcIO) Call(method string, params ...any) error {
	_, err := self.client.Call(self.ctx, method, params)
	return err
}

func (self *RpcIO) CallResult(result any, method string, params ...any) error {
	return self.client.CallResult(self.ctx, method, params, &result)
}

func (self *RpcIO) getEventChannel(accountId AccountId) chan *Event {
	self.mu.Lock()
	defer self.mu.Unlock()

	channel, ok := self.accountEvents[accountId]
	if !ok {
		channel = make(chan *Event, self.EventBuffer)
		self.accountEvents[accountId] = channel
	}
	return channel
}

func (self *RpcIO) onNotify(req *jrpc2.Request) {
	if req.Method() == "event" {
		var params _Params

		req.UnmarshalParams(&params)

		select {
		case <-self.ctx.Done():
			return
		default:
			self.events <- params
		}
	}
}
