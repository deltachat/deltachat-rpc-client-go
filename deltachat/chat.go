package deltachat

import "fmt"

// Delta Chat Chat.
type Chat struct {
	acc *Account
	Id  uint64
}

// Implement Stringer.
func (self *Chat) String() string {
	return fmt.Sprintf("Chat(Id=%v, acc=%v)", self.Id, self.acc.Id)
}

// Delete this chat and all its messages.
func (self *Chat) Delete() error {
	return self.acc.rpc.Call("delete_chat", self.acc.Id, self.Id)
}

// Block this chat.
func (self *Chat) Block() error {
	return self.acc.rpc.Call("block_chat", self.acc.Id, self.Id)
}

// Accept this contact request chat.
func (self *Chat) Accept() error {
	return self.acc.rpc.Call("accept_chat", self.acc.Id, self.Id)
}

// Leave this group chat.
func (self *Chat) Leave() error {
	return self.acc.rpc.Call("leave_group", self.acc.Id, self.Id)
}

// Set name of this chat.
func (self *Chat) SetName(name string) error {
	return self.acc.rpc.Call("set_chat_name", self.acc.Id, self.Id, name)
}

// Get Join-Group QR code text and SVG data.
func (self *Chat) QrCode() ([2]string, error) {
	var data [2]string
	err := self.acc.rpc.CallResult(&data, "get_chat_securejoin_qr_code_svg", self.acc.Id, self.Id)
	return data, err
}

// Send a text message and return the resulting Message instance.
func (self *Chat) SendText(text string) (*Message, error) {
	var id uint64
	err := self.acc.rpc.CallResult(&id, "misc_send_text_message", self.acc.Id, self.Id, text)
	return NewMessage(self.acc, id), err
}

// Chat factory
func NewChat(acc *Account, id uint64) *Chat {
	return &Chat{acc, id}
}
