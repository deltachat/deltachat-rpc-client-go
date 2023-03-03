package deltachat

// Chat list item snapshot
type ChatListItem struct {
	Id                  uint64
	Name                string
	AvatarPath          string
	Color               string
	LastUpdated         int64
	SummaryText1        string
	SummaryText2        string
	SummaryStatus       uint32
	IsProtected         bool
	IsGroup             bool
	FreshMessageCounter uint
	IsSelfTalk          bool
	IsDeviceTalk        bool
	IsSendingLocation   bool
	IsSelfInGroup       bool
	IsArchived          bool
	IsPinned            bool
	IsMuted             bool
	IsContactRequest    bool
	IsBroadcast         bool
	DmChatContact       uint64
	WasSeenRecently     bool

	// ArchiveLink
	// FreshMessageCounter uint

	// Error
	// Id uint64
	Error string
}
