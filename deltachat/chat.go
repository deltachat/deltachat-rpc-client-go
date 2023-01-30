package deltachat

// Delta Chat Chat.
type Chat struct {
	acc *Account
	Id  uint64
}

// Send a text message and return the resulting Message instance.
func (chat Chat) SendText(text string) (Message, error) {
	var id uint64
	err := chat.acc.rpc.CallResult(&id, "misc_send_text_message", chat.acc.Id, chat.Id, text)
	return newMessage(chat.acc, id), err
}

// Chat factory
func NewChat(acc *Account, id uint64) Chat {
	return Chat{acc, id}
}
