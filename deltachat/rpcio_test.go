package deltachat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRpcIO_Start(t *testing.T) {
	t.Parallel()
	manager := acfactory.NewAcManager()
	defer acfactory.StopRpc(manager)

	assert.NotNil(t, manager.Rpc.Start())

	manager.Rpc.Stop()
	assert.Nil(t, manager.Rpc.Start())
}

func TestRpcIO_Stop(t *testing.T) {
	t.Parallel()
	rpc := NewRpcIO()
	rpc.Stop()

	manager := acfactory.NewAcManager()
	manager.Rpc.Stop()
	manager.Rpc.Stop()
}
