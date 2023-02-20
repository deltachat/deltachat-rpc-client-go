package deltachat

const (
	//Special contact ids
	CONTACT_SELF         uint64 = 1
	CONTACT_INFO         uint64 = 2
	CONTACT_DEVICE       uint64 = 5
	CONTACT_LAST_SPECIAL uint64 = 9

	// Event types
	EVENT_INFO                          = "Info"
	EVENT_SMTP_CONNECTED                = "SmtpConnected"
	EVENT_IMAP_CONNECTED                = "ImapConnected"
	EVENT_SMTP_MESSAGE_SENT             = "SmtpMessageSent"
	EVENT_IMAP_MESSAGE_DELETED          = "ImapMessageDeleted"
	EVENT_IMAP_MESSAGE_MOVED            = "ImapMessageMoved"
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
