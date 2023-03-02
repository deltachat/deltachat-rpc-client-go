package deltachat

import "fmt"

// Message quote. Only the Text property is warrantied to be present, all other fields are optional.
type MessageQuote struct {
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
type MessageSnapshot struct {
	Id       uint64
	ChatId   uint64
	FromId   uint64
	Quote    MessageQuote
	ParentId uint64

	Text        string
	HasLocation bool
	HasHtml     bool
	ViewType    string
	State       int
	Error       string

	Timestamp             int
	SortTimestamp         int
	ReceivedTimestamp     int
	HasDeviatingTimestamp bool

	Subject        string
	ShowPadlock    bool
	IsSetupmessage bool
	IsInfo         bool
	IsForwarded    bool

	IsBot bool

	SystemMessageType string

	Duration         int
	DimensionsHeight int
	DimensionsWidth  int

	VideochatType int
	VideochatUrl  string

	OverrideSenderName string
	Sender             ContactSnapshot

	SetupCodeBegin string

	File      string
	FileMime  string
	FileBytes uint64
	FileName  string

	WebxdcInfo WebxdcMessageInfo

	DownloadState string

	Reactions Reactions
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
func (self *Message) Snapshot() (*MessageSnapshot, error) {
	var snapshot MessageSnapshot
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
