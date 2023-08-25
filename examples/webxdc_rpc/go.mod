module rpcbot

go 1.19

require github.com/deltachat/deltachat-rpc-client-go v0.17.1-0.20230731132031-99c0b7b46920

require (
	github.com/creachadair/jrpc2 v1.1.0 // indirect
	github.com/creachadair/mds v0.1.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
)

// this is needed only for tests, don't add it in your project's go.mod
replace github.com/deltachat/deltachat-rpc-client-go => ../../
