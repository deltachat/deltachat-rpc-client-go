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

// Get next event matching the given type.
func (self *Account) WaitForEvent(eventType string) *Event {
	eventChan := self.GetEventChannel()
	for {
		event := <-eventChan
		if event.Type == eventType {
			return event
		}
	}
}

// Remove the account.
func (self *Account) Remove() error {
	return self.rpc().Call("remove_account", self.Id)
}

// Select the account. The selected account will be returned by AccountManager.SelectedAccount()
func (self *Account) Select() error {
	return self.rpc().Call("select_account", self.Id)
}

// Start the account I/O.
func (self *Account) StartIO() error {
	return self.rpc().Call("start_io", self.Id)
}

// Stop the account I/O.
func (self *Account) StopIO() error {
	return self.rpc().Call("stop_io", self.Id)
}

// Get the current connectivity, i.e. whether the device is connected to the IMAP server.
// One of:
// - DC_CONNECTIVITY_NOT_CONNECTED (1000-1999): Show e.g. the string "Not connected" or a red dot
// - DC_CONNECTIVITY_CONNECTING (2000-2999): Show e.g. the string "Connecting…" or a yellow dot
// - DC_CONNECTIVITY_WORKING (3000-3999): Show e.g. the string "Getting new messages" or a spinning wheel
// - DC_CONNECTIVITY_CONNECTED (>=4000): Show e.g. the string "Connected" or a green dot
func (self *Account) Connectivity() (uint, error) {
	var info uint
	err := self.rpc().CallResult(&info, "get_connectivity", self.Id)
	return info, err
}

// Return map of this account configuration parameters.
func (self *Account) Info() (map[string]string, error) {
	var info map[string]string
	err := self.rpc().CallResult(&info, "get_info", self.Id)
	return info, err
}

// Get the combined filesize of an account in bytes.
func (self *Account) Size() (int, error) {
	var size int
	err := self.rpc().CallResult(&size, "get_account_file_size", self.Id)
	return size, err
}

// Return true if this account is configured, false otherwise.
func (self *Account) IsConfigured() (bool, error) {
	var configured bool
	err := self.rpc().CallResult(&configured, "is_configured", self.Id)
	return configured, err
}

// Set configuration value.
func (self *Account) SetConfig(key string, value string) error {
	return self.rpc().Call("set_config", self.Id, key, value)
}

// Get configuration value.
func (self *Account) GetConfig(key string) (string, error) {
	var value string
	err := self.rpc().CallResult(&value, "get_config", self.Id, key)
	return value, err
}

// Tweak several configuration values in a batch.
func (self *Account) UpdateConfig(config map[string]string) error {
	return self.rpc().Call("batch_set_config", self.Id, config)
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
func (self *Account) BlockedContacts() ([]ContactSnapshot, error) {
	var contacts []ContactSnapshot
	err := self.rpc().CallResult(&contacts, "get_blocked_contacts", self.Id)
	return contacts, err
}

// Get the contacts list.
func (self *Account) Contacts() ([]*Contact, error) {
	return self.QueryContacts("", 0)
}

// Get the list of contacts matching the given query.
func (self *Account) QueryContacts(query string, listFlags uint) ([]*Contact, error) {
	var ids []uint64
	err := self.rpc().CallResult(&ids, "get_contact_ids", self.Id, listFlags, query)
	var contacts []*Contact
	if err != nil {
		return contacts, err
	}
	contacts = make([]*Contact, len(ids))
	for i := range ids {
		contacts[i] = &Contact{self, ids[i]}
	}
	return contacts, err
}

// This account's identity as a Contact.
func (self *Account) Me() *Contact {
	return &Contact{self, CONTACT_SELF}
}

// Create a 1:1 chat with the given account.
func (self *Account) CreateChat(account *Account) (*Chat, error) {
	addr, err := account.GetConfig("configured_addr")
	if err != nil {
		return nil, err
	}

	contact, err := self.CreateContact(addr, "")
	if err != nil {
		return nil, err
	}

	chat, err := contact.CreateChat()
	if err != nil {
		return nil, err
	}

	return chat, nil
}

// Create a new group chat.
// After creation, the group has only self-contact as member and is in unpromoted state.
func (self *Account) CreateGroup(name string, protected bool) (*Chat, error) {
	var id uint64
	err := self.rpc().CallResult(&id, "create_group_chat", self.Id, name, protected)
	if err != nil {
		return nil, err
	}
	return &Chat{self, id}, err
}

// Create a new broadcast list.
func (self *Account) CreateBroadcastList() (*Chat, error) {
	var id uint64
	err := self.rpc().CallResult(&id, "create_broadcast_list", self.Id)
	if err != nil {
		return nil, err
	}
	return &Chat{self, id}, err
}

// Continue a Setup-Contact or Verified-Group-Invite protocol started on another device.
func (self *Account) SecureJoin(qrdata string) (*Chat, error) {
	var id uint64
	err := self.rpc().CallResult(&id, "secure_join", self.Id, qrdata)
	return &Chat{self, id}, err
}

// Get Setup-Contact QR Code text and SVG data.
func (self *Account) QrCode() (string, string, error) {
	var data [2]string
	err := self.rpc().CallResult(&data, "get_chat_securejoin_qr_code_svg", self.Id, nil)
	return data[0], data[1], err
}

// Export public and private keys to the specified directory.
// Note that the account does not have to be started.
func (self *Account) ExportSelfKeys(dir string) error {
	return self.rpc().Call("export_self_keys", self.Id, dir, nil)
}

// Import private keys found in the specified directory.
func (self *Account) ImportSelfKeys(dir string) error {
	return self.rpc().Call("import_self_keys", self.Id, dir, nil)
}

// Export account backup.
func (self *Account) ExportBackup(dir, passphrase string) error {
	var data any
	if passphrase == "" {
		data = nil
	} else {
		data = passphrase
	}
	return self.rpc().Call("export_backup", self.Id, dir, data)
}

// Import account backup.
func (self *Account) ImportBackup(file, passphrase string) error {
	var data any
	if passphrase == "" {
		data = nil
	} else {
		data = passphrase
	}
	return self.rpc().Call("import_backup", self.Id, file, data)
}

// Start the AutoCrypt key transfer process.
func (self *Account) InitiateAutocryptKeyTransfer() (string, error) {
	var result string
	err := self.rpc().CallResult(&result, "initiate_autocrypt_key_transfer", self.Id)
	return result, err
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
	if err != nil {
		return msgs, err
	}
	msgs = make([]*Message, len(ids))
	for i := range ids {
		msgs[i] = &Message{self, ids[i]}
	}
	return msgs, nil
}

// Return the default chat list items
func (self *Account) ChatListItems() ([]*ChatListItem, error) {
	return self.QueryChatListItems("", nil, 0)
}

// Return chat list items matching the given query.
func (self *Account) QueryChatListItems(query string, contact *Contact, listFlags uint) ([]*ChatListItem, error) {
	var entries [][]uint64
	var query2 any
	if query == "" {
		query2 = nil
	} else {
		query2 = query
	}
	err := self.rpc().CallResult(&entries, "get_chatlist_entries", self.Id, listFlags, query2, contact)
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

// Return the default chat list entries.
func (self *Account) ChatListEntries() ([]*Chat, error) {
	return self.QueryChatListEntries("", nil, 0)
}

// Return chat list entries matching the given query.
func (self *Account) QueryChatListEntries(query string, contact *Contact, listFlags uint) ([]*Chat, error) {
	var entries [][]uint64
	var query2 any
	if query == "" {
		query2 = nil
	} else {
		query2 = query
	}
	err := self.rpc().CallResult(&entries, "get_chatlist_entries", self.Id, listFlags, query2, contact)
	var items []*Chat
	if err != nil {
		return items, err
	}
	items = make([]*Chat, len(entries))
	for i, entry := range entries {
		items[i] = &Chat{self, entry[0]}
	}
	return items, nil
}

// Add a text message in the "Device messages" chat and return the resulting Message instance.
func (self *Account) AddDeviceMsg(label, text string) (*Message, error) {
	var id uint64
	err := self.rpc().CallResult(&id, "add_device_message", self.Id, label, text)
	if err != nil {
		return nil, err
	}
	return &Message{self, id}, nil
}

func (self *Account) rpc() Rpc {
	return self.Manager.Rpc
}
