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

func (rpc *Rpc) Start() error {
	rpc.cmd = exec.Command("deltachat-rpc-server")
	rpc.cmd.Stderr = os.Stderr
	rpc.stdin, _ = rpc.cmd.StdinPipe()
	stdout, _ := rpc.cmd.StdoutPipe()
	if err := rpc.cmd.Start(); err != nil {
		return err
	}

	rpc.ctx = context.Background()
	rpc.events = make(map[uint64]chan map[string]any)
	options := jrpc2.ClientOptions{OnNotify: rpc._onNotify}
	rpc.client = jrpc2.NewClient(channel.Line(stdout, rpc.stdin), &options)
	return nil
}

func (rpc *Rpc) Stop() {
	rpc.stdin.Close()
}

func (rpc *Rpc) _onNotify(req *jrpc2.Request) {
	if req.Method() == "event" {
		var params map[string]any
		req.UnmarshalParams(&params)
		accountId := uint64(params["contextId"].(float64))
		event := params["event"].(map[string]any)
		if _, ok := rpc.events[accountId]; !ok {
			rpc.events[accountId] = make(chan map[string]any)
		}
		go func() { rpc.events[accountId] <- event }()
	}
}

func (rpc *Rpc) WaitForEvent(accountId uint64) map[string]any {
	events := rpc.events[accountId]
	if events == nil {
		return nil
	}
	return <-events
}

func (rpc *Rpc) Call(method string, params ...any) error {
	_, err := rpc.client.Call(rpc.ctx, method, params)
	return err
}

func (rpc *Rpc) CallResult(result any, method string, params ...any) error {
	return rpc.client.CallResult(rpc.ctx, method, params, &result)
}

func NewRpc() *Rpc {
	return &Rpc{}
}
