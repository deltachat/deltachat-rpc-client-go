package deltachat

import "fmt"

// Full chat snapshot.
type FullChatSnapshot struct {
	Id                  uint64
	Name                string
	IsProtected         bool
	ProfileImage        string
	Archived            bool
	ChatType            uint
	IsUnpromoted        bool
	IsSelfTalk          bool
	Contacts            []*ContactSnapshot
	ContactIds          []uint64
	Color               string
	FreshMessageCounter uint
	IsContactRequest    bool
	IsDeviceChat        bool
	SelfInGroup         bool
	IsMuted             bool
	EphemeralTimer      uint
	CanSend             bool
	WasSeenRecently     bool
	MailingListAddress  string
}

// Cheaper version of FullChatSnapshot.
type BasicChatSnapshot struct {
	Id               uint64
	Name             string
	IsProtected      bool
	ProfileImage     string
	Archived         bool
	ChatType         uint
	IsUnpromoted     bool
	IsSelfTalk       bool
	Color            string
	IsContactRequest bool
	IsDeviceChat     bool
	IsMuted          bool
}

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

// Add contact to this group.
func (self *Chat) AddContact(contact *Contact) error {
	return self.rpc().Call("add_contact_to_chat", self.Account.Id, self.Id, contact.Id)
}

// Remove contact from this group.
func (self *Chat) RemoveContact(contact *Contact) error {
	return self.rpc().Call("remove_contact_from_chat", self.Account.Id, self.Id, contact.Id)
}

// Get Join-Group QR code text and SVG data.
func (self *Chat) QrCode() ([2]string, error) {
	var data [2]string
	err := self.rpc().CallResult(&data, "get_chat_securejoin_qr_code_svg", self.Account.Id, self.Id)
	return data, err
}

// Send a message and return the resulting Message instance.
func (self *Chat) SendMsg(msgData MsgData) (*Message, error) {
	var id uint64
	err := self.rpc().CallResult(&id, "send_msg", self.Account.Id, self.Id, msgData)
	if err != nil {
		return nil, err
	}
	return &Message{self.Account, id}, nil
}

// Send a text message and return the resulting Message instance.
func (self *Chat) SendText(text string) (*Message, error) {
	var id uint64
	err := self.rpc().CallResult(&id, "misc_send_text_message", self.Account.Id, self.Id, text)
	if err != nil {
		return nil, err
	}
	return &Message{self.Account, id}, nil
}

// Get a chat snapshot with basic info about this chat.
func (self *Chat) BasicSnapshot() (*BasicChatSnapshot, error) {
	var result BasicChatSnapshot
	err := self.rpc().CallResult(&result, "get_basic_chat_info", self.Account.Id, self.Id)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Get a full snapshot of this chat.
func (self *Chat) FullSnapshot() (*FullChatSnapshot, error) {
	var result FullChatSnapshot
	err := self.rpc().CallResult(&result, "get_full_chat_by_id", self.Account.Id, self.Id)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (self *Chat) rpc() Rpc {
	return self.Account.Manager.Rpc
}
