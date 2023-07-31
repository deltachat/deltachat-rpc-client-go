package transport

// Delta Chat RPC client's transport.
type RpcTransport interface {
	// Open the communication channel with the RPC server.
	Open() error
	// Close the communication channel with the RPC server.
	Close()
	// Request the RPC server to call a function that does not have a return value.
	Call(method string, params ...any) error
	// Request the RPC server to call a function that does have a return value.
	CallResult(result any, method string, params ...any) error
}
