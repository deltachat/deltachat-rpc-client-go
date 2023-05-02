package deltachat

// Delta Chat RPC client.
type Rpc interface {
	// Start the communication with the RPC server.
	Start() error
	// Stop the communication with the RPC server.
	Stop()
	// Request the RPC server to call a function that does not have a return value.
	Call(method string, params ...any) error
	// Request the RPC server to call a function that does have a return value.
	CallResult(result any, method string, params ...any) error
	// String representation of the Rpc instance.
	String() string
}
