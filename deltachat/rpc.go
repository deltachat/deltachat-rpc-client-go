package deltachat

import (
	"context"
	"io"
	"os"
	"os/exec"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
)

type Rpc struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	client *jrpc2.Client
	ctx    context.Context
	events map[uint64]chan map[string]any
}

func (self *Rpc) Start() error {
	self.cmd = exec.Command("deltachat-rpc-server")
	self.cmd.Stderr = os.Stderr
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
	self.stdin.Close()
}

func (self *Rpc) _onNotify(req *jrpc2.Request) {
	if req.Method() == "event" {
		var params map[string]any
		req.UnmarshalParams(&params)
		accountId := uint64(params["contextId"].(float64))
		event := params["event"].(map[string]any)
		if _, ok := self.events[accountId]; !ok {
			self.events[accountId] = make(chan map[string]any)
		}
		go func() { self.events[accountId] <- event }()
	}
}

func (self *Rpc) WaitForEvent(accountId uint64) map[string]any {
	events := self.events[accountId]
	if events == nil {
		return nil
	}
	return <-events
}

func (self *Rpc) Call(method string, params ...any) error {
	_, err := self.client.Call(self.ctx, method, params)
	return err
}

func (self *Rpc) CallResult(result any, method string, params ...any) error {
	return self.client.CallResult(self.ctx, method, params, &result)
}

func NewRpc() *Rpc {
	return &Rpc{}
}
