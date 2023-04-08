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

	toEvent(&_Event{Type: eventTypeImapMessageDeleted, Msg: "test"})
	toEvent(&_Event{Type: eventTypeImapInboxIdle})
	toEvent(&_Event{Type: eventTypeNewBlobFile, File: "test.jpg"})
	toEvent(&_Event{Type: eventTypeDeletedBlobFile, File: "test.jpg"})
	toEvent(&_Event{Type: eventTypeError, Msg: "test"})
	toEvent(&_Event{Type: eventTypeMsgFailed, ChatId: 0, MsgId: 0})
}
