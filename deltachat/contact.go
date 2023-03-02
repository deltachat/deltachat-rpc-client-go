package deltachat

import "fmt"

// Delta Chat Contact snapshot.
type ContactSnapshot struct {
	Address         string
	Color           string
	AuthName        string
	Status          string
	DisplayName     string
	Id              uint64
	Name            string
	ProfileImage    string
	NameAndAddr     string
	IsBlocked       bool
	IsVerified      bool
	VerifierAddr    string
	VerifierId      uint64
	LastSeen        uint64
	WasSeenRecently bool
}

// Delta Chat Contact.
type Contact struct {
	Account *Account
	Id      uint64
}

// Implement Stringer.
func (self *Contact) String() string {
	return fmt.Sprintf("Contact(Id=%v, Account=%v)", self.Id, self.Account.Id)
}

// Block contact.
func (self *Contact) Block() error {
	return self.rpc().Call("block_contact", self.Account.Id, self.Id)
}

// Unblock contact.
func (self *Contact) Unblock() error {
	return self.rpc().Call("unblock_contact", self.Account.Id, self.Id)
}

// Delete contact.
func (self *Contact) Delete() error {
	return self.rpc().Call("delete_contact", self.Account.Id, self.Id)
}

// Return a map with a snapshot of all contact properties.
func (self *Contact) Snapshot() (*ContactSnapshot, error) {
	var snapshot ContactSnapshot
	err := self.rpc().CallResult(&snapshot, "get_contact", self.Account.Id, self.Id)
	return &snapshot, err
}

// Create or get an existing 1:1 chat for this contact.
func (self *Contact) CreateChat() (*Chat, error) {
	var id uint64
	err := self.rpc().CallResult(&id, "create_chat_by_contact_id", self.Account.Id, self.Id)
	return &Chat{self.Account, id}, err
}

func (self *Contact) rpc() Rpc {
	return self.Account.Manager.Rpc
}
