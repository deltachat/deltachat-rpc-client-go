module echobot

go 1.19

require (
	github.com/deltachat/deltachat-rpc-client-go v0.17.1-0.20230503140011-50d66d05c271
	github.com/stretchr/testify v1.8.2
)

require (
	github.com/creachadair/jrpc2 v1.0.0 // indirect
	github.com/creachadair/mds v0.0.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// this is needed only for tests, don't add it in your project's go.mod
replace github.com/deltachat/deltachat-rpc-client-go => ../../
