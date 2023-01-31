package deltachat

// Delta Chat account.
type Account struct {
	rpc *Rpc
	Id  uint64
}

// Wait until the next event and return it.
func (acc *Account) WaitForEvent() map[string]any {
	return acc.rpc.WaitForEvent(acc.Id)
}

// Remove the account.
func (acc *Account) Remove() error {
	return acc.rpc.Call("remove_account", acc.Id)
}

// Start the account I/O.
func (acc *Account) StartIO() error {
	return acc.rpc.Call("start_io", acc.Id)
}

// Stop the account I/O.
func (acc *Account) StopIO() error {
	return acc.rpc.Call("stop_io", acc.Id)
}

// Return map of this account configuration parameters.
func (acc *Account) GetInfo() (map[string]string, error) {
	var info map[string]string
	return info, acc.rpc.CallResult(&info, "get_info", acc.Id)
}

// Get the combined filesize of an account in bytes.
func (acc *Account) GetSize() (int, error) {
	var size int
	return size, acc.rpc.CallResult(&size, "get_account_file_size", acc.Id)
}

// Return true if this account is configured, false otherwise.
func (acc *Account) IsConfigured() (bool, error) {
	var configured bool
	return configured, acc.rpc.CallResult(&configured, "is_configured", acc.Id)
}

// Set configuration value.
func (acc *Account) SetConfig(key string, value string) error {
	return acc.rpc.Call("set_config", acc.Id, key, value)
}

// Get configuration value.
func (acc *Account) GetConfig(key string) (string, error) {
	var value string
	return value, acc.rpc.CallResult(&value, "get_config", acc.Id, key)
}

// Set self avatar. Passing nil will discard the currently set avatar.
func (acc *Account) SetAvatar(path string) error {
	return acc.SetConfig("selfavatar", path)
}

// Get self avatar path.
func (acc *Account) GetAvatar() (string, error) {
	return acc.GetConfig("selfavatar")
}

// Configure an account.
func (acc *Account) Configure() error {
	return acc.rpc.Call("configure", acc.Id)
}

// Return fresh messages list sorted in the order of their arrival, with ascending IDs.
func (acc *Account) GetFreshMsgsInArrivalOrder() ([]*Message, error) {
	var ids []uint64
	err := acc.rpc.CallResult(&ids, "get_fresh_msgs", acc.Id)
	var msgs []*Message
	if err == nil {
		msgs = make([]*Message, len(ids))
		for i := range ids {
			msgs[i] = newMessage(acc, ids[i])
		}
	}
	return msgs, err
}

// Account factory
func newAccount(rpc *Rpc, id uint64) *Account {
	return &Account{rpc, id}
}
