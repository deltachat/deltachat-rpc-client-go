package deltachat

// Delta Chat Contact.
type Contact struct {
	acc *Account
	Id  uint64
}

// Contact factory
func NewContact(acc *Account, id uint64) *Contact {
	return &Contact{acc, id}
}
