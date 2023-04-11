package deltachat

// BotRunningErr is returned by Bot.Run() if the Bot is already running
type BotRunningErr struct{}

func (self *BotRunningErr) Error() string {
	return "bot is already running"
}

// RpcRunningErr is returned by Rpc.Start() if the Rpc is already running
type RpcRunningErr struct{}

func (self *RpcRunningErr) Error() string {
	return "RPC is already running"
}
