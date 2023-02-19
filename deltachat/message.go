package deltachat

import "fmt"

// Delta Chat Message.
type Message struct {
	acc *Account
	Id  uint64
}

// Implement Stringer.
func (self *Message) String() string {
	return fmt.Sprintf("Message(Id=%v, acc=%v)", self.Id, self.acc.Id)
}

// Return map of this account configuration parameters.
func (self *Message) Snapshot() (map[string]any, error) {
	var snapshot map[string]any
	err := self.acc.rpc.CallResult(&snapshot, "get_message", self.acc.Id, self.Id)
	snapshot["chat"] = NewChat(self.acc, uint64(snapshot["chatId"].(float64)))
	snapshot["sender"] = NewContact(self.acc, uint64(snapshot["fromId"].(float64)))
	snapshot["message"] = self
	return snapshot, err
}

// Mark the message as seen.
func (self *Message) MarkSeen() error {
	return self.acc.rpc.Call("markseen_msgs", self.acc.Id, []any{self.Id})
}

// Message factory
func NewMessage(acc *Account, id uint64) *Message {
	return &Message{acc, id}
}
