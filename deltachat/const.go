package deltachat

type ChatListFlag uint

type ContactFlag uint

type ChatVisibility string

type DownloadState string

type MsgType string

type SysmsgType string

type EventType string

const (
	//Special contact ids
	ContactSelf        ContactId = 1
	ContactInfo        ContactId = 2
	ContactDevice      ContactId = 5
	ContactLastSpecial ContactId = 9

	// Chatlist Flags
	ChatListFlagArchivedOnly   ChatListFlag = 0x01
	ChatListFlagNoSpecials     ChatListFlag = 0x02
	ChatListFlagAddAlldoneHint ChatListFlag = 0x04
	ChatListFlagForForwarding  ChatListFlag = 0x08

	// Contact Flags
	ContactFlagVerifiedOnly ContactFlag = 0x01
	ContactFlagAddSelf      ContactFlag = 0x02

	//Chat types
	ChatUndefined   ChatType = 0
	ChatSingle      ChatType = 100
	ChatGroup       ChatType = 120
	ChatMailinglist ChatType = 140
	ChatBroadcast   ChatType = 160

	// Chat visibility types
	ChatVisibilityNormal   ChatVisibility = "Normal"
	ChatVisibilityArchived ChatVisibility = "Archived"
	ChatVisibilityPinned   ChatVisibility = "Pinned"

	//Message download states
	DownloadDone       DownloadState = "Done"
	DownloadAvailable  DownloadState = "Available"
	DownloadFailure    DownloadState = "Failure"
	DownloadInProgress DownloadState = "InProgress"

	//Message view types
	MsgUnknown             MsgType = "Unknown"
	MsgText                MsgType = "Text"
	MsgImage               MsgType = "Image"
	MsgGif                 MsgType = "Gif"
	MsgSticker             MsgType = "Sticker"
	MsgAudio               MsgType = "Audio"
	MsgVoice               MsgType = "Voice"
	MsgVideo               MsgType = "Video"
	MsgFile                MsgType = "File"
	MsgVideochatInvitation MsgType = "VideochatInvitation"
	MsgWebxdc              MsgType = "Webxdc"

	//System message types
	SysmsgUnknown                  SysmsgType = "Unknown"
	SysmsgGroupNameChanged         SysmsgType = "GroupNameChanged"
	SysmsgGroupImageChanged        SysmsgType = "GroupImageChanged"
	SysmsgMemberAddedToGroup       SysmsgType = "MemberAddedToGroup"
	SysmsgMemberRemovedFromGroup   SysmsgType = "MemberRemovedFromGroup"
	SysmsgAutocryptSetupMessage    SysmsgType = "AutocryptSetupMessage"
	SysmsgSecurejoinMessage        SysmsgType = "SecurejoinMessage"
	SysmsgLocationStreamingEnabled SysmsgType = "LocationStreamingEnabled"
	SysmsgLocationOnly             SysmsgType = "LocationOnly"
	SysmsgChatProtectionEnabled    SysmsgType = "ChatProtectionEnabled"
	SysmsgChatProtectionDisabled   SysmsgType = "ChatProtectionDisabled"
	SysmsgWebxdcStatusUpdate       SysmsgType = "WebxdcStatusUpdate"
	SysmsgEphemeralTimerChanged    SysmsgType = "EphemeralTimerChanged"
	SysmsgMultiDeviceSync          SysmsgType = "MultiDeviceSync"
	SysmsgWebxdcInfoMessage        SysmsgType = "WebxdcInfoMessage"

	// Event types
	EventInfo                       EventType = "Info"
	EventSmtpConnected              EventType = "SmtpConnected"
	EventImapConnected              EventType = "ImapConnected"
	EventSmtpMessageSent            EventType = "SmtpMessageSent"
	EventImapMessageDeleted         EventType = "ImapMessageDeleted"
	EventImapMessageMoved           EventType = "ImapMessageMoved"
	EventImapInboxIdle              EventType = "ImapInboxIdle"
	EventNewBlobFile                EventType = "NewBlobFile"
	EventDeletedBlobFile            EventType = "DeletedBlobFile"
	EventWarning                    EventType = "Warning"
	EventError                      EventType = "Error"
	EventErrorSelfNotInGroup        EventType = "ErrorSelfNotInGroup"
	EventMsgsChanged                EventType = "MsgsChanged"
	EventReactionsChanged           EventType = "ReactionsChanged"
	EventIncomingMsg                EventType = "IncomingMsg"
	EventIncomingMsgBunch           EventType = "IncomingMsgBunch"
	EventMsgsNoticed                EventType = "MsgsNoticed"
	EventMsgDelivered               EventType = "MsgDelivered"
	EventMsgFailed                  EventType = "MsgFailed"
	EventMsgRead                    EventType = "MsgRead"
	EventChatModified               EventType = "ChatModified"
	EventChatEphemeralTimerModified EventType = "ChatEphemeralTimerModified"
	EventContactsChanged            EventType = "ContactsChanged"
	EventLocationChanged            EventType = "LocationChanged"
	EventConfigureProgress          EventType = "ConfigureProgress"
	EventImexProgress               EventType = "ImexProgress"
	EventImexFileWritten            EventType = "ImexFileWritten"
	EventSecurejoinInviterProgress  EventType = "SecurejoinInviterProgress"
	EventSecurejoinJoinerProgress   EventType = "SecurejoinJoinerProgress"
	EventConnectivityChanged        EventType = "ConnectivityChanged"
	EventSelfavatarChanged          EventType = "SelfavatarChanged"
	EventWebxdcStatusUpdate         EventType = "WebxdcStatusUpdate"
	EventWebxdcInstanceDeleted      EventType = "WebxdcInstanceDeleted"
)
