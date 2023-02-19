package deltachat

import "fmt"

// Delta Chat Contact.
type Contact struct {
	acc *Account
	Id  uint64
}

// Implement Stringer.
func (self *Contact) String() string {
	return fmt.Sprintf("Contact(Id=%v, acc=%v)", self.Id, self.acc.Id)
}

// Block contact.
func (self *Contact) Block() error {
	return self.acc.rpc.Call("block_contact", self.acc.Id, self.Id)
}

// Unblock contact.
func (self *Contact) Unblock() error {
	return self.acc.rpc.Call("unblock_contact", self.acc.Id, self.Id)
}

// Delete contact.
func (self *Contact) Delete() error {
	return self.acc.rpc.Call("delete_contact", self.acc.Id, self.Id)
}

// Return a map with a snapshot of all contact properties.
func (self *Contact) Snapshot() (map[string]any, error) {
	var data map[string]any
	err := self.acc.rpc.CallResult(&data, "get_contact", self.acc.Id, self.Id)
	return data, err
}

// Create or get an existing 1:1 chat for this contact.
func (self *Contact) CreateChat() (*Chat, error) {
	var id uint64
	err := self.acc.rpc.CallResult(&id, "create_chat_by_contact_id", self.acc.Id, self.Id)
	return NewChat(self.acc, id), err
}

// Contact factory
func NewContact(acc *Account, id uint64) *Contact {
	return &Contact{acc, id}
}
