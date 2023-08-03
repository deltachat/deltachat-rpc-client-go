package main // replace with your package name

import (
	"testing"

	"github.com/deltachat/deltachat-rpc-client-go/deltachat"
)

var acfactory *deltachat.AcFactory

func TestMain(m *testing.M) {
	acfactory = &deltachat.AcFactory{}
	acfactory.TearUp()
	defer acfactory.TearDown()
	m.Run()
}
