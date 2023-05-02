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
	mu          sync.Mutex
}

func NewRpcIO() *RpcIO {
	return &RpcIO{Cmd: deltachatRpcServerBin, Stderr: os.Stderr}
}

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

	self.client = jrpc2.NewClient(channel.Line(stdout, self.stdin), nil)
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
}

func (self *RpcIO) Call(method string, params ...any) error {
	_, err := self.client.Call(self.ctx, method, params)
	return err
}

func (self *RpcIO) CallResult(result any, method string, params ...any) error {
	return self.client.CallResult(self.ctx, method, params, &result)
}
