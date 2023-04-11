package deltachat

import (
	"os"
	"testing"
)

var acfactory *AcFactory

func TestMain(m *testing.M) {
	cfg := map[string]string{
		"mail_server":   "localhost",
		"send_server":   "localhost",
		"mail_port":     "3143",
		"send_port":     "3025",
		"mail_security": "3",
		"send_security": "3",
	}
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		panic(err)
	}
	acfactory = &AcFactory{}
	acfactory.TearUp(cfg, dir, os.Getenv("TEST_DEBUG") == "1")
	defer acfactory.TearDown()
	m.Run()
}
