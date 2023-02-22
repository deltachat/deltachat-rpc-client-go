package deltachat

import "fmt"

// Delta Chat account.
type Account struct {
	Manager *AccountManager
	Id      uint64
	rpc     *Rpc
}

// Implement Stringer.
func (self *Account) String() string {
	return fmt.Sprintf("Account(Id=%v)", self.Id)
}

// Wait until the next event and return it.
// Stop waiting and return nil when the Rpc connection is closed.
func (self *Account) WaitForEvent() map[string]any {
	return self.rpc.WaitForEvent(self.Id)
}

// Remove the account.
func (self *Account) Remove() error {
	return self.rpc.Call("remove_account", self.Id)
}

// Start the account I/O.
func (self *Account) StartIO() error {
	return self.rpc.Call("start_io", self.Id)
}

// Stop the account I/O.
func (self *Account) StopIO() error {
	return self.rpc.Call("stop_io", self.Id)
}

// Return map of this account configuration parameters.
func (self *Account) Info() (map[string]string, error) {
	var info map[string]string
	return info, self.rpc.CallResult(&info, "get_info", self.Id)
}

// Get the combined filesize of an account in bytes.
func (self *Account) Size() (int, error) {
	var size int
	return size, self.rpc.CallResult(&size, "get_account_file_size", self.Id)
}

// Return true if this account is configured, false otherwise.
func (self *Account) IsConfigured() (bool, error) {
	var configured bool
	return configured, self.rpc.CallResult(&configured, "is_configured", self.Id)
}

// Set configuration value.
func (self *Account) SetConfig(key string, value string) error {
	return self.rpc.Call("set_config", self.Id, key, value)
}

// Get configuration value.
func (self *Account) GetConfig(key string) (string, error) {
	var value string
	return value, self.rpc.CallResult(&value, "get_config", self.Id, key)
}

// Set self avatar. Passing nil will discard the currently set avatar.
func (self *Account) SetAvatar(path string) error {
	return self.SetConfig("selfavatar", path)
}

// Get self avatar path.
func (self *Account) Avatar() (string, error) {
	return self.GetConfig("selfavatar")
}

// Configure an account.
func (self *Account) Configure() error {
	return self.rpc.Call("configure", self.Id)
}

// Create a new Contact or return an existing one.
// If there already is a Contact with that e-mail address, it is unblocked and its display
// name is updated if specified.
func (self *Account) CreateContact(addr string, name string) (*Contact, error) {
	var id uint64
	err := self.rpc.CallResult(&id, "create_contact", self.Id, addr, name)
	return NewContact(self, id), err
}

// Check if an e-mail address belongs to a known and unblocked contact.
func (self *Account) GetContactByAddr(addr string) (*Contact, error) {
	var id uint64
	err := self.rpc.CallResult(&id, "lookup_contact_id_by_addr", self.Id, addr)
	if id > 0 {
		return NewContact(self, id), err
	}
	return nil, err
}

// Return a list with snapshots of all blocked contacts.
func (self *Account) GetBlockedContacts() ([]map[string]any, error) {
	var contacts []map[string]any
	return contacts, self.rpc.CallResult(&contacts, "get_blocked_contacts", self.Id)
}

// This account's identity as a Contact.
func (self *Account) SelfContact() *Contact {
	return NewContact(self, CONTACT_SELF)
}

// Create a new group chat.
// After creation, the group has only self-contact as member and is in unpromoted state.
func (self *Account) CreateGroup(name string, protected bool) (*Chat, error) {
	var id uint64
	err := self.rpc.CallResult(&id, "create_group_chat", self.Id, name, protected)
	return NewChat(self, id), err
}

// Continue a Setup-Contact or Verified-Group-Invite protocol started on another device.
func (self *Account) SecureJoin(qrdata string) (*Chat, error) {
	var id uint64
	err := self.rpc.CallResult(&id, "secure_join", self.Id, qrdata)
	return NewChat(self, id), err
}

// Get Setup-Contact QR Code text and SVG data.
func (self *Account) QrCode() ([2]string, error) {
	var data [2]string
	return data, self.rpc.CallResult(&data, "get_chat_securejoin_qr_code_svg", self.Id)
}

// Mark the given set of messages as seen.
func (self *Account) MarkSeenMsgs(messages []*Message) error {
	ids := make([]uint64, len(messages))
	for i := range messages {
		ids[i] = messages[i].Id
	}
	return self.rpc.Call("markseen_msgs", self.Id, ids)
}

// Delete the given set of messages (local and remote).
func (self *Account) DeleteMsgs(messages []*Message) error {
	ids := make([]uint64, len(messages))
	for i := range messages {
		ids[i] = messages[i].Id
	}
	return self.rpc.Call("delete_messages", self.Id, ids)
}

// Return the list of fresh messages, newest messages first.
// This call is intended for displaying notifications.
func (self *Account) FreshMsgs() ([]*Message, error) {
	var ids []uint64
	err := self.rpc.CallResult(&ids, "get_fresh_msgs", self.Id)
	var msgs []*Message
	if err == nil {
		msgs = make([]*Message, len(ids))
		for i := range ids {
			msgs[i] = NewMessage(self, ids[i])
		}
	}
	return msgs, err
}

// Return fresh messages list sorted in the order of their arrival, with ascending IDs.
func (self *Account) FreshMsgsInArrivalOrder() ([]*Message, error) {
	var ids []uint64
	err := self.rpc.CallResult(&ids, "get_fresh_msgs", self.Id)
	var msgs []*Message
	if err == nil {
		msgs = make([]*Message, len(ids))
		for i := range ids {
			msgs[i] = NewMessage(self, ids[i])
		}
	}
	return msgs, err
}

// Account factory
func NewAccount(manager *AccountManager, id uint64) *Account {
	return &Account{manager, manager.Rpc, id}
}
