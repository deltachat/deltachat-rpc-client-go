package deltachat

import "fmt"

// Delta Chat accounts manager. This is the root of the API.
type AccountManager struct {
	Rpc Rpc
}

// Implement Stringer.
func (self *AccountManager) String() string {
	return fmt.Sprintf("AccountManager(Rpc=%#v)", self.Rpc.String())
}

// Create a new account.
func (self *AccountManager) AddAccount() (*Account, error) {
	var id uint64
	err := self.Rpc.CallResult(&id, "add_account")
	return &Account{self, id}, err
}

// Get the selected account.
func (self *AccountManager) SelectedAccount() (*Account, error) {
	var id uint64
	err := self.Rpc.CallResult(&id, "get_selected_account_id")
	if id == 0 {
		return nil, err
	}
	return &Account{self, id}, err
}

// Return all available accounts.
func (self *AccountManager) Accounts() ([]*Account, error) {
	var ids []uint64
	err := self.Rpc.CallResult(&ids, "get_all_account_ids")
	var accounts []*Account
	if err != nil {
		return accounts, err
	}
	accounts = make([]*Account, len(ids))
	for i := range ids {
		accounts[i] = &Account{self, ids[i]}
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
func (self *AccountManager) SystemInfo() (map[string]string, error) {
	var info map[string]string
	err := self.Rpc.CallResult(&info, "get_system_info")
	return info, err
}

// Set stock translation strings.
func (self *AccountManager) SetTranslations(translations map[uint]string) error {
	return self.Rpc.Call("set_stock_strings", translations)
}
