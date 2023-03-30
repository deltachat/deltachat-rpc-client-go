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

// Mark all messages in this chat as noticed.
func (self *Chat) MarkNoticed() error {
	return self.rpc().Call("marknoticed_chat", self.Account.Id, self.Id)
}

// Set mute duration of this chat.
// duration value can be:
//
//	0 - Chat is not muted.
//
// -1 - Chat is muted until the user unmutes the chat.
//
//	t - Chat is muted for a limited period of time.
func (self *Chat) SetMuteDuration(duration int64) error {
	var data any
	switch duration {
	case -1:
		data = "Forever"
	case 0:
		data = "NotMuted"
	default:
		data = map[string]int64{"Until": duration}
	}
	return self.rpc().Call("set_chat_mute_duration", self.Account.Id, self.Id, data)
}

// Set name of this chat.
func (self *Chat) SetName(name string) error {
	return self.rpc().Call("set_chat_name", self.Account.Id, self.Id, name)
}

// Set profile image of this chat.
func (self *Chat) SetImage(path string) error {
	return self.rpc().Call("set_chat_profile_image", self.Account.Id, self.Id, path)
}

// Pin this chat.
func (self *Chat) Pin() error {
	return self.rpc().Call("set_chat_visibility", self.Account.Id, self.Id, CHAT_VISIBILITY_PINNED)
}

// Unpin this chat.
func (self *Chat) Unpin() error {
	return self.rpc().Call("set_chat_visibility", self.Account.Id, self.Id, CHAT_VISIBILITY_NORMAL)
}

// Archive this chat.
func (self *Chat) Archive() error {
	return self.rpc().Call("set_chat_visibility", self.Account.Id, self.Id, CHAT_VISIBILITY_ARCHIVED)
}

// Unarchive this chat.a
func (self *Chat) Unarchive() error {
	return self.rpc().Call("set_chat_visibility", self.Account.Id, self.Id, CHAT_VISIBILITY_NORMAL)
}

// Remove profile image of this chat.
func (self *Chat) RemoveImage() error {
	return self.rpc().Call("set_chat_profile_image", self.Account.Id, self.Id, nil)
}

// Add contact to this group.
func (self *Chat) AddContact(contact *Contact) error {
	return self.rpc().Call("add_contact_to_chat", self.Account.Id, self.Id, contact.Id)
}

// Remove contact from this group.
func (self *Chat) RemoveContact(contact *Contact) error {
	return self.rpc().Call("remove_contact_from_chat", self.Account.Id, self.Id, contact.Id)
}

// Get the list of contacts in this chat.
func (self *Chat) Contacts() ([]*Contact, error) {
	var contacts []*Contact
	var ids []uint64
	err := self.rpc().CallResult(&ids, "get_chat_contacts", self.Account.Id, self.Id)
	if err != nil {
		return contacts, err
	}
	contacts = make([]*Contact, len(ids))
	for i := range ids {
		contacts[i] = &Contact{self.Account, ids[i]}
	}
	return contacts, nil
}

// Set ephemeral timer of this chat.
func (self *Chat) SetEphemeralTimer(timer uint) error {
	return self.rpc().Call("set_chat_ephemeral_timer", self.Account.Id, self.Id, timer)
}

// Get ephemeral timer of this chat.
func (self *Chat) EphemeralTimer() (uint, error) {
	var timer uint
	err := self.rpc().CallResult(&timer, "get_chat_ephemeral_timer", self.Account.Id, self.Id)
	return timer, err
}

// Get Join-Group QR code text and SVG data.
func (self *Chat) QrCode() (string, string, error) {
	var data [2]string
	err := self.rpc().CallResult(&data, "get_chat_securejoin_qr_code_svg", self.Account.Id, self.Id)
	return data[0], data[1], err
}

// Get encryption info for this chat.
// Get a multi-line encryption info, containing encryption preferences of all members.
// Can be used to find out why messages sent to group are not encrypted.
//
// returns Multi-line text
func (self *Chat) EncryptionInfo() (string, error) {
	var data string
	err := self.rpc().CallResult(&data, "get_chat_encryption_info", self.Account.Id, self.Id)
	return data, err
}

// Get the list of messages in this chat.
func (self *Chat) Messages(infoOnly, addDaymarker bool) ([]*Message, error) {
	var msgs []*Message
	var ids []uint64
	err := self.rpc().CallResult(&ids, "get_message_ids", self.Account.Id, self.Id, infoOnly, addDaymarker)
	if err != nil {
		return msgs, err
	}
	msgs = make([]*Message, len(ids))
	for i := range ids {
		msgs[i] = &Message{self.Account, ids[i]}
	}
	return msgs, nil
}

// Get the number of fresh messages in this chat.
func (self *Chat) FreshMsgCount() (uint, error) {
	var count uint
	err := self.rpc().CallResult(&count, "get_fresh_msg_cnt", self.Account.Id, self.Id)
	return count, err
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

// Send a video chat invitation.
func (self *Chat) SendVideoChatInvitation() (*Message, error) {
	var id uint64
	err := self.rpc().CallResult(&id, "send_videochat_invitation", self.Account.Id, self.Id)
	if err != nil {
		return nil, err
	}
	return &Message{self.Account, id}, nil
}

// Get first unread message in this chat.
func (self *Chat) FirstUnreadMsg() (*Message, error) {
	var id uint64
	err := self.rpc().CallResult(&id, "get_first_unread_message_of_chat", self.Account.Id, self.Id)
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

// Forward a list of messages to this chat.
func (self *Chat) DeleteMsgs(messages []*Message) error {
	ids := make([]uint64, len(messages))
	for i := range messages {
		ids[i] = messages[i].Id
	}
	return self.rpc().Call("forward_messages", self.Account.Id, ids, self.Id)
}

func (self *Chat) rpc() Rpc {
	return self.Account.Manager.Rpc
}
