package deltachat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRpc_Start(t *testing.T) {
	t.Parallel()
	manager := acfactory.NewAcManager()
	defer manager.Rpc.Stop()

	assert.NotNil(t, manager.Rpc.Start())
}

func TestRpc_toEvent(t *testing.T) {
	t.Parallel()

	toEvent(&_Event{Type: eventImapMessageDeleted, Msg: "test"})
	toEvent(&_Event{Type: eventImapInboxIdle})
	toEvent(&_Event{Type: eventNewBlobFile, File: "test.jpg"})
	toEvent(&_Event{Type: eventDeletedBlobFile, File: "test.jpg"})
	toEvent(&_Event{Type: eventError, Msg: "test"})
	toEvent(&_Event{Type: eventMsgFailed, ChatId: 0, MsgId: 0})
}
