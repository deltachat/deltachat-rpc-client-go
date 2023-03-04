package deltachat

// Delta Chat accounts manager. This is the root of the API.
type AccountManager struct {
	Rpc Rpc
}

// Implement Stringer.
func (self *AccountManager) String() string {
	return "AccountManager(Rpc=" + self.Rpc.String() + ")"
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
	if err == nil {
		accounts = make([]*Account, len(ids))
		for i := range ids {
			accounts[i] = &Account{self, ids[i]}
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

// Get the current connectivity, i.e. whether the device is connected to the IMAP server.
// One of:
// - DC_CONNECTIVITY_NOT_CONNECTED (1000-1999): Show e.g. the string "Not connected" or a red dot
// - DC_CONNECTIVITY_CONNECTING (2000-2999): Show e.g. the string "Connecting…" or a yellow dot
// - DC_CONNECTIVITY_WORKING (3000-3999): Show e.g. the string "Getting new messages" or a spinning wheel
// - DC_CONNECTIVITY_CONNECTED (>=4000): Show e.g. the string "Connected" or a green dot
func (self *AccountManager) Connectivity() (uint, error) {
	var info uint
	err := self.Rpc.CallResult(&info, "get_connectivity")
	return info, err
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
