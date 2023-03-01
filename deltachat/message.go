package deltachat

import "fmt"

// Delta Chat Message.
type Message struct {
	Account *Account
	Id      uint64
}

// Implement Stringer.
func (self *Message) String() string {
	return fmt.Sprintf("Message(Id=%v, Account=%v)", self.Id, self.Account.Id)
}

// Return map of this account configuration parameters.
func (self *Message) Snapshot() (map[string]any, error) {
	var snapshot map[string]any
	err := self.rpc().CallResult(&snapshot, "get_message", self.Account.Id, self.Id)
	snapshot["chat"] = &Chat{self.Account, uint64(snapshot["chatId"].(float64))}
	snapshot["sender"] = &Contact{self.Account, uint64(snapshot["fromId"].(float64))}
	snapshot["message"] = self
	return snapshot, err
}

// Mark the message as seen.
func (self *Message) MarkSeen() error {
	return self.rpc().Call("markseen_msgs", self.Account.Id, []any{self.Id})
}

func (self *Message) rpc() Rpc {
	return self.Account.Manager.Rpc
}
