package tests

import (
	"testing"
)

var server *EmailServer

func TestMain(m *testing.M) {
	var err error
	server, err = NewEmailServer()
	if err != nil {
		panic(err)
	}
	defer server.Stop()
	server.Start()
	m.Run()
}
