package deltachat

import "fmt"

// Message data provided to Chat.SendMsg()
type MsgData struct {
	Text               string      `json:"text,omitempty"`
	Html               string      `json:"html,omitempty"`
	ViewType           string      `json:"viewtype,omitempty"`
	File               string      `json:"file,omitempty"`
	Location           *[2]float64 `json:"location,omitempty"`
	OverrideSenderName string      `json:"overrideSenderName,omitempty"`
	QuotedMessageId    uint64      `json:"quotedMessageId,omitempty"`
}

// Message quote. Only the Text property is warrantied to be present, all other fields are optional.
type MsgQuote struct {
	Text               string
	MessageId          uint64
	AuthorDisplayName  string
	AuthorDisplayColor string
	OverrideSenderName string
	Image              string
	IsForwarded        bool
	ViewType           string
}

// Message search result.
type MsgSearchResult struct {
	Id                 uint64
	AuthorProfileImage string
	AuthorName         string
	AuthorColor        string
	ChatName           string
	Message            string
	Timestamp          int64
}

// Delta Chat Message.
type Message struct {
	Account *Account
	Id      uint64
}

// Implement Stringer.
func (self *Message) String() string {
	return fmt.Sprintf("Message(Id=%v, Account=%v)", self.Id, self.Account.Id)
}

// Return map of this account configuration parameters.
func (self *Message) Snapshot() (*MsgSnapshot, error) {
	var snapshot MsgSnapshot
	err := self.rpc().CallResult(&snapshot, "get_message", self.Account.Id, self.Id)
	if err != nil {
		return nil, err
	}
	snapshot.Account = self.Account
	return &snapshot, err
}

// Get the HTML part of this message.
func (self *Message) Html() (string, error) {
	var html string
	err := self.rpc().CallResult(&html, "get_message_html", self.Account.Id, self.Id)
	return html, err
}

// Get an informational text for a single message.
func (self *Message) Info() (string, error) {
	var info string
	err := self.rpc().CallResult(&info, "get_message_info", self.Account.Id, self.Id)
	return info, err
}

// Delete message.
func (self *Message) Delete() error {
	return self.rpc().Call("delete_messages", self.Account.Id, []uint64{self.Id})
}

// Asks the core to start downloading a message fully.
func (self *Message) Download() error {
	return self.rpc().Call("download_full_message", self.Account.Id, self.Id)
}

// Mark the message as seen.
func (self *Message) MarkSeen() error {
	return self.rpc().Call("markseen_msgs", self.Account.Id, []uint64{self.Id})
}

// Send a reaction to this message.
func (self *Message) SendReaction(reaction ...string) error {
	err := self.rpc().Call("send_reaction", self.Account.Id, self.Id, reaction)
	return err
}

// Continue the AutoCrypt key transfer process.
func (self *Message) ContinueAutocryptKeyTransfer(setupCode string) error {
	return self.rpc().Call("continue_autocrypt_key_transfer", self.Account.Id, self.Id, setupCode)
}

// Send status update for the webxdc instance of this message.
func (self *Message) SendStatusUpdate(update, description string) error {
	return self.rpc().Call("send_webxdc_status_update", self.Account.Id, self.Id, update, description)
}

// Get the status updates of this webxdc message as a JSON string.
func (self *Message) StatusUpdates(lastKnownSerial uint) (string, error) {
	var data string
	err := self.rpc().CallResult(&data, "get_webxdc_status_updates", self.Account.Id, self.Id, lastKnownSerial)
	return data, err
}

// Get info from this webxdc message.
func (self *Message) WebxdcInfo() (*WebxdcMsgInfo, error) {
	var info *WebxdcMsgInfo
	err := self.rpc().CallResult(info, "get_webxdc_info", self.Account.Id, self.Id)
	return info, err
}

func (self *Message) rpc() Rpc {
	return self.Account.Manager.Rpc
}
