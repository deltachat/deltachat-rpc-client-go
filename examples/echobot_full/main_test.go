package main // replace with your package name

import (
	"log"
	"testing"

	"github.com/deltachat/deltachat-rpc-client-go/deltachat"
)

var acfactory *deltachat.AcFactory

func TestMain(m *testing.M) {
	acfactory = &deltachat.AcFactory{}
	acfactory.TearUp()
	defer acfactory.TearDown()
	acfactory.WithRpc(func(rpc *deltachat.Rpc) {
		sysinfo, _ := rpc.GetSystemInfo()
		log.Println("Running deltachat core", sysinfo["deltachat_core_version"])
	})
	m.Run()
}
