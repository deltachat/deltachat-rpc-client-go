package deltachat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvent(t *testing.T) {
	t.Parallel()

	assert.NotEmpty(t, EventInfo{}.eventType())
	assert.NotEmpty(t, EventSmtpConnected{}.eventType())
	assert.NotEmpty(t, EventImapConnected{}.eventType())
	assert.NotEmpty(t, EventSmtpMessageSent{}.eventType())
	assert.NotEmpty(t, EventImapMessageDeleted{}.eventType())
	assert.NotEmpty(t, EventImapMessageMoved{}.eventType())
	assert.NotEmpty(t, EventImapInboxIdle{}.eventType())
	assert.NotEmpty(t, EventNewBlobFile{}.eventType())
	assert.NotEmpty(t, EventDeletedBlobFile{}.eventType())
	assert.NotEmpty(t, EventWarning{}.eventType())
	assert.NotEmpty(t, EventError{}.eventType())
	assert.NotEmpty(t, EventErrorSelfNotInGroup{}.eventType())
	assert.NotEmpty(t, EventInfo{}.eventType())
	assert.NotEmpty(t, EventMsgsChanged{}.eventType())
	assert.NotEmpty(t, EventReactionsChanged{}.eventType())
	assert.NotEmpty(t, EventIncomingMsg{}.eventType())
	assert.NotEmpty(t, EventIncomingMsgBunch{}.eventType())
	assert.NotEmpty(t, EventMsgsNoticed{}.eventType())
	assert.NotEmpty(t, EventMsgDelivered{}.eventType())
	assert.NotEmpty(t, EventMsgFailed{}.eventType())
	assert.NotEmpty(t, EventMsgRead{}.eventType())
	assert.NotEmpty(t, EventChatModified{}.eventType())
	assert.NotEmpty(t, EventChatEphemeralTimerModified{}.eventType())
	assert.NotEmpty(t, EventContactsChanged{}.eventType())
	assert.NotEmpty(t, EventLocationChanged{}.eventType())
	assert.NotEmpty(t, EventConfigureProgress{}.eventType())
	assert.NotEmpty(t, EventImexProgress{}.eventType())
	assert.NotEmpty(t, EventImexFileWritten{}.eventType())
	assert.NotEmpty(t, EventSecurejoinInviterProgress{}.eventType())
	assert.NotEmpty(t, EventConnectivityChanged{}.eventType())
	assert.NotEmpty(t, EventSelfavatarChanged{}.eventType())
	assert.NotEmpty(t, EventWebxdcStatusUpdate{}.eventType())
	assert.NotEmpty(t, EventWebxdcInstanceDeleted{}.eventType())
}
