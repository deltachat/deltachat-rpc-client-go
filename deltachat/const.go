package deltachat

const (
	//Special contact ids
	CONTACT_SELF         uint64 = 1
	CONTACT_INFO         uint64 = 2
	CONTACT_DEVICE       uint64 = 5
	CONTACT_LAST_SPECIAL uint64 = 9

	// Chatlist Flags
	CHATLIST_FLAG_ARCHIVED_ONLY    uint = 0x01
	CHATLIST_FLAG_NO_SPECIALS      uint = 0x02
	CHATLIST_FLAG_ADD_ALLDONE_HINT uint = 0x04
	CHATLIST_FLAG_FOR_FORWARDING   uint = 0x08

	// Contact Flags
	CONTACT_FLAG_VERIFIED_ONLY = 0x01
	CONTACT_FLAG_ADD_SELF      = 0x02

	//Chat types
	CHAT_TYPE_UNDEFINED   uint = 0
	CHAT_TYPE_SINGLE      uint = 100
	CHAT_TYPE_GROUP       uint = 120
	CHAT_TYPE_MAILINGLIST uint = 140
	CHAT_TYPE_BROADCAST   uint = 160

	// Chat visibility types
	CHAT_VISIBILITY_NORMAL   = "Normal"
	CHAT_VISIBILITY_ARCHIVED = "Archived"
	CHAT_VISIBILITY_PINNED   = "Pinned"

	//Message download states
	DOWNLOAD_STATE_DONE        = "Done"
	DOWNLOAD_STATE_AVAILABLE   = "Available"
	DOWNLOAD_STATE_FAILURE     = "Failure"
	DOWNLOAD_STATE_IN_PROGRESS = "InProgress"

	//Message view types
	MSG_TYPE_UNKNOWN              = "Unknown"
	MSG_TYPE_TEXT                 = "Text"
	MSG_TYPE_IMAGE                = "Image"
	MSG_TYPE_GIF                  = "Gif"
	MSG_TYPE_STICKER              = "Sticker"
	MSG_TYPE_AUDIO                = "Audio"
	MSG_TYPE_VOICE                = "Voice"
	MSG_TYPE_VIDEO                = "Video"
	MSG_TYPE_FILE                 = "File"
	MSG_TYPE_VIDEOCHAT_INVITATION = "VideochatInvitation"
	MSG_TYPE_WEBXDC               = "Webxdc"

	//System message types
	SYSMSG_TYPE_UNKNOWN                    = "Unknown"
	SYSMSG_TYPE_GROUP_NAME_CHANGED         = "GroupNameChanged"
	SYSMSG_TYPE_GROUP_IMAGE_CHANGED        = "GroupImageChanged"
	SYSMSG_TYPE_MEMBER_ADDED_TO_GROUP      = "MemberAddedToGroup"
	SYSMSG_TYPE_MEMBER_REMOVED_FROM_GROUP  = "MemberRemovedFromGroup"
	SYSMSG_TYPE_AUTOCRYPT_SETUP_MESSAGE    = "AutocryptSetupMessage"
	SYSMSG_TYPE_SECUREJOIN_MESSAGE         = "SecurejoinMessage"
	SYSMSG_TYPE_LOCATION_STREAMING_ENABLED = "LocationStreamingEnabled"
	SYSMSG_TYPE_LOCATION_ONLY              = "LocationOnly"
	SYSMSG_TYPE_CHAT_PROTECTION_ENABLED    = "ChatProtectionEnabled"
	SYSMSG_TYPE_CHAT_PROTECTION_DISABLED   = "ChatProtectionDisabled"
	SYSMSG_TYPE_WEBXDC_STATUS_UPDATE       = "WebxdcStatusUpdate"
	SYSMSG_TYPE_EPHEMERAL_TIMER_CHANGED    = "EphemeralTimerChanged"
	SYSMSG_TYPE_MULTI_DEVICE_SYNC          = "MultiDeviceSync"
	SYSMSG_TYPE_WEBXDC_INFO_MESSAGE        = "WebxdcInfoMessage"

	// Event types
	EVENT_INFO                          = "Info"
	EVENT_SMTP_CONNECTED                = "SmtpConnected"
	EVENT_IMAP_CONNECTED                = "ImapConnected"
	EVENT_SMTP_MESSAGE_SENT             = "SmtpMessageSent"
	EVENT_IMAP_MESSAGE_DELETED          = "ImapMessageDeleted"
	EVENT_IMAP_MESSAGE_MOVED            = "ImapMessageMoved"
	EVENT_IMAP_INBOX_IDLE               = "ImapInboxIdle"
	EVENT_NEW_BLOB_FILE                 = "NewBlobFile"
	EVENT_DELETED_BLOB_FILE             = "DeletedBlobFile"
	EVENT_WARNING                       = "Warning"
	EVENT_ERROR                         = "Error"
	EVENT_ERROR_SELF_NOT_IN_GROUP       = "ErrorSelfNotInGroup"
	EVENT_MSGS_CHANGED                  = "MsgsChanged"
	EVENT_REACTIONS_CHANGED             = "ReactionsChanged"
	EVENT_INCOMING_MSG                  = "IncomingMsg"
	EVENT_INCOMING_MSG_BUNCH            = "IncomingMsgBunch"
	EVENT_MSGS_NOTICED                  = "MsgsNoticed"
	EVENT_MSG_DELIVERED                 = "MsgDelivered"
	EVENT_MSG_FAILED                    = "MsgFailed"
	EVENT_MSG_READ                      = "MsgRead"
	EVENT_CHAT_MODIFIED                 = "ChatModified"
	EVENT_CHAT_EPHEMERAL_TIMER_MODIFIED = "ChatEphemeralTimerModified"
	EVENT_CONTACTS_CHANGED              = "ContactsChanged"
	EVENT_LOCATION_CHANGED              = "LocationChanged"
	EVENT_CONFIGURE_PROGRESS            = "ConfigureProgress"
	EVENT_IMEX_PROGRESS                 = "ImexProgress"
	EVENT_IMEX_FILE_WRITTEN             = "ImexFileWritten"
	EVENT_SECUREJOIN_INVITER_PROGRESS   = "SecurejoinInviterProgress"
	EVENT_SECUREJOIN_JOINER_PROGRESS    = "SecurejoinJoinerProgress"
	EVENT_CONNECTIVITY_CHANGED          = "ConnectivityChanged"
	EVENT_SELFAVATAR_CHANGED            = "SelfavatarChanged"
	EVENT_WEBXDC_STATUS_UPDATE          = "WebxdcStatusUpdate"
	EVENT_WEBXDC_INSTANCE_DELETED       = "WebxdcInstanceDeleted"
)
