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
type Rpc struct {
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

// Implement Stringer.
func (self *Rpc) String() string {
	return fmt.Sprintf("Rpc(AccountsDir=%v)", self.AccountsDir)
}

func (self *Rpc) Start() error {
	self.closed = false
	self.cmd = exec.Command("deltachat-rpc-server")
	if self.AccountsDir != "" {
		self.cmd.Env = append(os.Environ(), "DC_ACCOUNTS_PATH="+self.AccountsDir)
	}
	self.cmd.Stderr = self.Stderr
	self.stdin, _ = self.cmd.StdinPipe()
	stdout, _ := self.cmd.StdoutPipe()
	if err := self.cmd.Start(); err != nil {
		return err
	}

	self.ctx = context.Background()
	self.events = make(map[uint64]chan map[string]any)
	options := jrpc2.ClientOptions{OnNotify: self._onNotify}
	self.client = jrpc2.NewClient(channel.Line(stdout, self.stdin), &options)
	return nil
}

func (self *Rpc) Stop() {
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

func (self *Rpc) _initEventChannel(accountId uint64) {
	self.eventsMutex.Lock()
	if _, ok := self.events[accountId]; !ok {
		self.events[accountId] = make(chan map[string]any)
	}
	self.eventsMutex.Unlock()
}

func (self *Rpc) _onNotify(req *jrpc2.Request) {
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

func (self *Rpc) WaitForEvent(accountId uint64) map[string]any {
	self._initEventChannel(accountId)
	v, _ := <-self.events[accountId]
	return v
}

func (self *Rpc) Call(method string, params ...any) error {
	_, err := self.client.Call(self.ctx, method, params)
	return err
}

func (self *Rpc) CallResult(result any, method string, params ...any) error {
	return self.client.CallResult(self.ctx, method, params, &result)
}

func NewRpc() *Rpc {
	return &Rpc{Stderr: os.Stderr}
}
