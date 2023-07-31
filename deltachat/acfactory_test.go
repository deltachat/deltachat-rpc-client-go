package deltachat

import (
	"testing"
)

func TestAcFactory_TearDown(t *testing.T) {
	t.Parallel()
	acf := &AcFactory{}
	acf.TearUp()
	acf.TearDown()
}

func TestAcFactory_getChatId(t *testing.T) {
	t.Parallel()
	getChatId(EventIncomingMsg{})
	getChatId(EventMsgsNoticed{})
	getChatId(EventMsgDelivered{})
	getChatId(EventMsgFailed{})
	getChatId(EventMsgRead{})
	getChatId(EventChatModified{})
}
