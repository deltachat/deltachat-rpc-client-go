package deltachat

// Delta Chat accounts manager. This is the root of the API.
type AccountManager struct {
	rpc *Rpc
}

// Create a new account.
func (man *AccountManager) AddAccount() (*Account, error) {
	var id uint64
	err := man.rpc.CallResult(&id, "add_account")
	return newAccount(man.rpc, id), err
}

// Return a list of all available accounts.
func (man *AccountManager) Accounts() ([]*Account, error) {
	var ids []uint64
	err := man.rpc.CallResult(&ids, "get_all_account_ids")
	var accounts []*Account
	if err == nil {
		accounts = make([]*Account, len(ids))
		for i := range ids {
			accounts[i] = newAccount(man.rpc, ids[i])
		}
	}
	return accounts, err
}

// Start the I/O of all accounts.
func (man *AccountManager) StartIO() error {
	return man.rpc.Call("start_io_for_all_accounts")
}

// Stop the I/O of all accounts.
func (man *AccountManager) StopIO() error {
	return man.rpc.Call("stop_io_for_all_accounts")
}

// Indicate that the network likely has come back or just that the network conditions might have changed.
func (man *AccountManager) MaybeNetwork() error {
	return man.rpc.Call("maybe_network")
}

// Get information about the Delta Chat core in this system.
func (man *AccountManager) GetSystemInfo() (map[string]any, error) {
	var info map[string]any
	return info, man.rpc.CallResult(&info, "get_system_info")
}

// Set stock translation strings.
func (man *AccountManager) SetTranslations(translations map[string]string) error {
	return man.rpc.Call("set_stock_strings", translations)
}

// AccountManager factory
func NewAccountManager(rpc *Rpc) *AccountManager {
	return &AccountManager{rpc}
}
