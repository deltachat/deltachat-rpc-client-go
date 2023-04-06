package deltachat

import "fmt"

type ContactId uint64

// Delta Chat Contact snapshot.
type ContactSnapshot struct {
	Address         string
	Color           string
	AuthName        string
	Status          string
	DisplayName     string
	Id              ContactId
	Name            string
	ProfileImage    string
	NameAndAddr     string
	IsBlocked       bool
	IsVerified      bool
	VerifierAddr    string
	VerifierId      ContactId
	LastSeen        Timestamp
	WasSeenRecently bool
}

// Delta Chat Contact.
type Contact struct {
	Account *Account
	Id      ContactId
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

// Set name of this contact.
func (self *Contact) SetName(name string) error {
	return self.rpc().Call("change_contact_name", self.Account.Id, self.Id, name)
}

// Get encryption info for this contact.
// Get a multi-line encryption info, containing your fingerprint and the
// fingerprint of the contact, used e.g. to compare the fingerprints for a simple out-of-band verification.
func (self *Contact) EncryptionInfo() (string, error) {
	var data string
	err := self.rpc().CallResult(&data, "get_contact_encryption_info", self.Account.Id, self.Id)
	return data, err
}

// Return a map with a snapshot of all contact properties.
func (self *Contact) Snapshot() (*ContactSnapshot, error) {
	var snapshot ContactSnapshot
	err := self.rpc().CallResult(&snapshot, "get_contact", self.Account.Id, self.Id)
	return &snapshot, err
}

// Create or get an existing 1:1 chat for this contact.
func (self *Contact) CreateChat() (*Chat, error) {
	var id ChatId
	err := self.rpc().CallResult(&id, "create_chat_by_contact_id", self.Account.Id, self.Id)
	return &Chat{self.Account, id}, err
}

func (self *Contact) rpc() Rpc {
	return self.Account.Manager.Rpc
}
