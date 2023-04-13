package deltachat

import (
	"testing"
)

func TestAcFactory_TearDown(t *testing.T) {
	t.Parallel()
	acf := &AcFactory{}
	tearUp(acf)
	acf.TearDown()
}

func TestAcFactory_NewAcManager(t *testing.T) {
	t.Parallel()
	acfactory := &AcFactory{}
	tearUp(acfactory)
	defer acfactory.TearDown()

	acfactory.debug = false
	manager := acfactory.NewAcManager()
	acfactory.StopRpc(manager)
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
