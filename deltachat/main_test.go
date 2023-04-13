package deltachat

import (
	"os"
	"testing"
)

var acfactory *AcFactory

func tearUp(acf *AcFactory) {
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

	acf.TearUp(cfg, dir, os.Getenv("TEST_DEBUG") == "1")
}

func TestMain(m *testing.M) {
	acfactory = &AcFactory{}
	tearUp(acfactory)
	defer acfactory.TearDown()
	m.Run()
}
