module echobot

go 1.21

toolchain go1.21.0

require (
	github.com/deltachat/deltachat-rpc-client-go v0.17.1-0.20230731132031-99c0b7b46920
	github.com/stretchr/testify v1.8.2
)

require (
	github.com/creachadair/jrpc2 v1.1.0 // indirect
	github.com/creachadair/mds v0.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// this is needed only for tests, don't add it in your project's go.mod
replace github.com/deltachat/deltachat-rpc-client-go => ../../
