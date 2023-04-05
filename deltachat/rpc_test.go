package deltachat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRpc_Start(t *testing.T) {
	manager := acfactory.NewAcManager()
	defer manager.Rpc.Stop()

	assert.NotNil(t, manager.Rpc.Start())
}
