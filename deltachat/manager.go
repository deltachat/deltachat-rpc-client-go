package deltachat

import "fmt"

// Delta Chat accounts manager. This is the root of the API.
type AccountManager struct {
	Rpc Rpc
}

// Implement Stringer.
func (self *AccountManager) String() string {
	return fmt.Sprintf("AccountManager(Rpc=%v)", self.Rpc)
}

// Create a new account.
func (self *AccountManager) AddAccount() (*Account, error) {
	var id uint64
	err := self.Rpc.CallResult(&id, "add_account")
	return NewAccount(self, id), err
}

// Return a list of all available accounts.
func (self *AccountManager) Accounts() ([]*Account, error) {
	var ids []uint64
	err := self.Rpc.CallResult(&ids, "get_all_account_ids")
	var accounts []*Account
	if err == nil {
		accounts = make([]*Account, len(ids))
		for i := range ids {
			accounts[i] = NewAccount(self, ids[i])
		}
	}
	return accounts, err
}

// Start the I/O of all accounts.
func (self *AccountManager) StartIO() error {
	return self.Rpc.Call("start_io_for_all_accounts")
}

// Stop the I/O of all accounts.
func (self *AccountManager) StopIO() error {
	return self.Rpc.Call("stop_io_for_all_accounts")
}

// Indicate that the network likely has come back or just that the network conditions might have changed.
func (self *AccountManager) MaybeNetwork() error {
	return self.Rpc.Call("maybe_network")
}

// Get information about the Delta Chat core in this system.
func (self *AccountManager) SystemInfo() (map[string]any, error) {
	var info map[string]any
	return info, self.Rpc.CallResult(&info, "get_system_info")
}

// Set stock translation strings.
func (self *AccountManager) SetTranslations(translations map[string]string) error {
	return self.Rpc.Call("set_stock_strings", translations)
}

// AccountManager factory
func NewAccountManager(rpc Rpc) *AccountManager {
	return &AccountManager{rpc}
}
