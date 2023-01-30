package deltachat

// Delta Chat Message.
type Message struct {
	acc *Account
	rpc *Rpc
	Id  uint64
}

// Return map of this account configuration parameters.
func (msg Message) GetSnapshot() (map[string]any, error) {
	var snapshot map[string]any
	err := msg.rpc.CallResult(&snapshot, "get_message", msg.acc.Id, msg.Id)
	snapshot["chat"] = NewChat(msg.acc, uint64(snapshot["chatId"].(float64)))
	snapshot["sender"] = NewContact(msg.acc, uint64(snapshot["fromId"].(float64)))
	snapshot["message"] = msg
	return snapshot, err
}

// Mark the message as seen.
func (msg Message) MarkSeen() error {
	return msg.rpc.Call("markseen_msgs", msg.acc.Id, []any{msg.Id})
}

// Message factory
func newMessage(acc *Account, id uint64) Message {
	return Message{acc, acc.rpc, id}
}
