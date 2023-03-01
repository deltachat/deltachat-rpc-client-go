package deltachat

import "fmt"

// Delta Chat Chat.
type Chat struct {
	Account *Account
	Id      uint64
}

// Implement Stringer.
func (self *Chat) String() string {
	return fmt.Sprintf("Chat(Id=%v, Account=%v)", self.Id, self.Account.Id)
}

// Delete this chat and all its messages.
func (self *Chat) Delete() error {
	return self.rpc().Call("delete_chat", self.Account.Id, self.Id)
}

// Block this chat.
func (self *Chat) Block() error {
	return self.rpc().Call("block_chat", self.Account.Id, self.Id)
}

// Accept this contact request chat.
func (self *Chat) Accept() error {
	return self.rpc().Call("accept_chat", self.Account.Id, self.Id)
}

// Leave this group chat.
func (self *Chat) Leave() error {
	return self.rpc().Call("leave_group", self.Account.Id, self.Id)
}

// Set name of this chat.
func (self *Chat) SetName(name string) error {
	return self.rpc().Call("set_chat_name", self.Account.Id, self.Id, name)
}

// Get Join-Group QR code text and SVG data.
func (self *Chat) QrCode() ([2]string, error) {
	var data [2]string
	err := self.rpc().CallResult(&data, "get_chat_securejoin_qr_code_svg", self.Account.Id, self.Id)
	return data, err
}

// Send a text message and return the resulting Message instance.
func (self *Chat) SendText(text string) (*Message, error) {
	var id uint64
	err := self.rpc().CallResult(&id, "misc_send_text_message", self.Account.Id, self.Id, text)
	return &Message{self.Account, id}, err
}

func (self *Chat) rpc() Rpc {
	return self.Account.Manager.Rpc
}
