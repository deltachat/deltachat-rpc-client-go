package deltachat

import (
	"fmt"
	"sort"
)

// Delta Chat account.
type Account struct {
	Manager *AccountManager
	Id      uint64
}

// Implement Stringer.
func (self *Account) String() string {
	return fmt.Sprintf("Account(Id=%v)", self.Id)
}

// Get this account's event channel.
func (self *Account) GetEventChannel() <-chan *Event {
	return self.rpc().GetEventChannel(self.Id)
}

// Remove the account.
func (self *Account) Remove() error {
	return self.rpc().Call("remove_account", self.Id)
}

// Start the account I/O.
func (self *Account) StartIO() error {
	return self.rpc().Call("start_io", self.Id)
}

// Stop the account I/O.
func (self *Account) StopIO() error {
	return self.rpc().Call("stop_io", self.Id)
}

// Return map of this account configuration parameters.
func (self *Account) Info() (map[string]string, error) {
	var info map[string]string
	return info, self.rpc().CallResult(&info, "get_info", self.Id)
}

// Get the combined filesize of an account in bytes.
func (self *Account) Size() (int, error) {
	var size int
	return size, self.rpc().CallResult(&size, "get_account_file_size", self.Id)
}

// Return true if this account is configured, false otherwise.
func (self *Account) IsConfigured() (bool, error) {
	var configured bool
	return configured, self.rpc().CallResult(&configured, "is_configured", self.Id)
}

// Set configuration value.
func (self *Account) SetConfig(key string, value string) error {
	return self.rpc().Call("set_config", self.Id, key, value)
}

// Get configuration value.
func (self *Account) GetConfig(key string) (string, error) {
	var value string
	return value, self.rpc().CallResult(&value, "get_config", self.Id, key)
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
	return self.rpc().Call("configure", self.Id)
}

// Create a new Contact or return an existing one.
// If there already is a Contact with that e-mail address, it is unblocked and its display
// name is updated if specified.
func (self *Account) CreateContact(addr string, name string) (*Contact, error) {
	var id uint64
	err := self.rpc().CallResult(&id, "create_contact", self.Id, addr, name)
	return &Contact{self, id}, err
}

// Check if an e-mail address belongs to a known and unblocked contact.
func (self *Account) GetContactByAddr(addr string) (*Contact, error) {
	var id uint64
	err := self.rpc().CallResult(&id, "lookup_contact_id_by_addr", self.Id, addr)
	if id > 0 {
		return &Contact{self, id}, err
	}
	return nil, err
}

// Return a list with snapshots of all blocked contacts.
func (self *Account) GetBlockedContacts() ([]ContactSnapshot, error) {
	var contacts []ContactSnapshot
	return contacts, self.rpc().CallResult(&contacts, "get_blocked_contacts", self.Id)
}

// This account's identity as a Contact.
func (self *Account) SelfContact() *Contact {
	return &Contact{self, CONTACT_SELF}
}

// Create a new group chat.
// After creation, the group has only self-contact as member and is in unpromoted state.
func (self *Account) CreateGroup(name string, protected bool) (*Chat, error) {
	var id uint64
	err := self.rpc().CallResult(&id, "create_group_chat", self.Id, name, protected)
	return &Chat{self, id}, err
}

// Continue a Setup-Contact or Verified-Group-Invite protocol started on another device.
func (self *Account) SecureJoin(qrdata string) (*Chat, error) {
	var id uint64
	err := self.rpc().CallResult(&id, "secure_join", self.Id, qrdata)
	return &Chat{self, id}, err
}

// Get Setup-Contact QR Code text and SVG data.
func (self *Account) QrCode() ([2]string, error) {
	var data [2]string
	return data, self.rpc().CallResult(&data, "get_chat_securejoin_qr_code_svg", self.Id)
}

// Mark the given set of messages as seen.
func (self *Account) MarkSeenMsgs(messages []*Message) error {
	ids := make([]uint64, len(messages))
	for i := range messages {
		ids[i] = messages[i].Id
	}
	return self.rpc().Call("markseen_msgs", self.Id, ids)
}

// Delete the given set of messages (local and remote).
func (self *Account) DeleteMsgs(messages []*Message) error {
	ids := make([]uint64, len(messages))
	for i := range messages {
		ids[i] = messages[i].Id
	}
	return self.rpc().Call("delete_messages", self.Id, ids)
}

// Return the list of fresh messages, newest messages first.
// This call is intended for displaying notifications.
func (self *Account) FreshMsgs() ([]*Message, error) {
	var msgs []*Message
	var ids []uint64
	err := self.rpc().CallResult(&ids, "get_fresh_msgs", self.Id)
	if err != nil {
		return msgs, err
	}
	msgs = make([]*Message, len(ids))
	for i := range ids {
		msgs[i] = &Message{self, ids[i]}
	}
	return msgs, nil
}

// Return fresh messages list sorted in the order of their arrival, with ascending IDs.
func (self *Account) FreshMsgsInArrivalOrder() ([]*Message, error) {
	var msgs []*Message
	var ids []uint64
	err := self.rpc().CallResult(&ids, "get_fresh_msgs", self.Id)
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	if err == nil {
		return msgs, err
	}
	msgs = make([]*Message, len(ids))
	for i := range ids {
		msgs[i] = &Message{self, ids[i]}
	}
	return msgs, nil
}

// Return chat list items.
func (self *Account) Chatlist() ([]*ChatListItem, error) {
	var entries [][]uint64
	err := self.rpc().CallResult(&entries, "get_chatlist_entries", self.Id, 0, nil, nil)
	var items []*ChatListItem
	if err != nil {
		return items, err
	}
	var itemsMap map[uint64]*ChatListItem
	err = self.rpc().CallResult(&itemsMap, "get_chatlist_items_by_entries", self.Id, entries)
	if err != nil {
		return items, err
	}
	items = make([]*ChatListItem, len(entries))
	for i, entry := range entries {
		items[i] = itemsMap[entry[0]]
	}
	return items, err
}

// Get the contact list.
func (self *Account) Contactlist() ([]*Contact, error) {
	var ids []uint64
	err := self.rpc().CallResult(&ids, "get_contact_ids", self.Id, 0, nil)
	var contacts []*Contact
	if err == nil {
		contacts = make([]*Contact, len(ids))
		for i := range ids {
			contacts[i] = &Contact{self, ids[i]}
		}
	}
	return contacts, err
}

func (self *Account) rpc() Rpc {
	return self.Manager.Rpc
}
