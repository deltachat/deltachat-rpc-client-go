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

// Message snapshot.
type MsgSnapshot struct {
	Id                    uint64
	ChatId                uint64
	FromId                uint64
	Quote                 *MsgQuote
	ParentId              uint64
	Text                  string
	HasLocation           bool
	HasHtml               bool
	ViewType              string
	State                 int
	Error                 string
	Timestamp             int
	SortTimestamp         int
	ReceivedTimestamp     int
	HasDeviatingTimestamp bool
	Subject               string
	ShowPadlock           bool
	IsSetupmessage        bool
	IsInfo                bool
	IsForwarded           bool
	IsBot                 bool
	SystemMessageType     string
	Duration              int
	DimensionsHeight      int
	DimensionsWidth       int
	VideochatType         int
	VideochatUrl          string
	OverrideSenderName    string
	Sender                *ContactSnapshot
	SetupCodeBegin        string
	File                  string
	FileMime              string
	FileBytes             uint64
	FileName              string
	WebxdcInfo            *WebxdcMsgInfo
	DownloadState         string
	Reactions             *Reactions
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
	return &snapshot, err
}

// Mark the message as seen.
func (self *Message) MarkSeen() error {
	return self.rpc().Call("markseen_msgs", self.Account.Id, []any{self.Id})
}

func (self *Message) rpc() Rpc {
	return self.Account.Manager.Rpc
}
