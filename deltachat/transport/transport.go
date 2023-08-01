package transport

import (
	"context"
)

// Delta Chat RPC client's transport.
type RpcTransport interface {
	// Request the RPC server to call a function that does not have a return value.
	Call(ctx context.Context, method string, params ...any) error
	// Request the RPC server to call a function that does have a return value.
	CallResult(ctx context.Context, result any, method string, params ...any) error
}
