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

// Delta Chat core RPC
type Rpc interface {
	Start() error
	Stop()
	WaitForEvent(accountId uint64) map[string]any
	Call(method string, params ...any) error
	CallResult(result any, method string, params ...any) error
}

// Delta Chat core RPC working over IO
type RpcIO struct {
	cmd         *exec.Cmd
	stdin       io.WriteCloser
	Stderr      *os.File
	AccountsDir string
	client      *jrpc2.Client
	ctx         context.Context
	events      map[uint64]chan map[string]any
	eventsMutex sync.Mutex
	closed      bool
}

func NewRpcIO() *RpcIO {
	return &RpcIO{Stderr: os.Stderr}
}

// Implement Stringer.
func (self *RpcIO) String() string {
	return fmt.Sprintf("Rpc(AccountsDir=%v)", self.AccountsDir)
}

func (self *RpcIO) Start() error {
	self.closed = false
	self.cmd = exec.Command("deltachat-rpc-server")
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
	self.events = make(map[uint64]chan map[string]any)
	options := jrpc2.ClientOptions{OnNotify: self._onNotify}
	self.client = jrpc2.NewClient(channel.Line(stdout, self.stdin), &options)
	return nil
}

func (self *RpcIO) Stop() {
	self.eventsMutex.Lock()
	if !self.closed {
		self.stdin.Close()
		self.cmd.Process.Wait()
		for _, value := range self.events {
			close(value)
		}
		self.closed = true
	}
	self.eventsMutex.Unlock()
}

func (self *RpcIO) WaitForEvent(accountId uint64) map[string]any {
	self._initEventChannel(accountId)
	v, _ := <-self.events[accountId]
	return v
}

func (self *RpcIO) Call(method string, params ...any) error {
	_, err := self.client.Call(self.ctx, method, params)
	return err
}

func (self *RpcIO) CallResult(result any, method string, params ...any) error {
	return self.client.CallResult(self.ctx, method, params, &result)
}

func (self *RpcIO) _initEventChannel(accountId uint64) {
	self.eventsMutex.Lock()
	if _, ok := self.events[accountId]; !ok {
		self.events[accountId] = make(chan map[string]any, 5)
	}
	self.eventsMutex.Unlock()
}

func (self *RpcIO) _onNotify(req *jrpc2.Request) {
	if req.Method() == "event" {
		var params map[string]any
		req.UnmarshalParams(&params)
		accountId := uint64(params["contextId"].(float64))
		event := params["event"].(map[string]any)
		self._initEventChannel(accountId)
		if !self.closed {
			go func() { self.events[accountId] <- event }()
		}
	}
}
