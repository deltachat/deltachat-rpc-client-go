package deltachat

type eventType string

const (
	eventInfo                       eventType = "Info"
	eventSmtpConnected              eventType = "SmtpConnected"
	eventImapConnected              eventType = "ImapConnected"
	eventSmtpMessageSent            eventType = "SmtpMessageSent"
	eventImapMessageDeleted         eventType = "ImapMessageDeleted"
	eventImapMessageMoved           eventType = "ImapMessageMoved"
	eventImapInboxIdle              eventType = "ImapInboxIdle"
	eventNewBlobFile                eventType = "NewBlobFile"
	eventDeletedBlobFile            eventType = "DeletedBlobFile"
	eventWarning                    eventType = "Warning"
	eventError                      eventType = "Error"
	eventErrorSelfNotInGroup        eventType = "ErrorSelfNotInGroup"
	eventMsgsChanged                eventType = "MsgsChanged"
	eventReactionsChanged           eventType = "ReactionsChanged"
	eventIncomingMsg                eventType = "IncomingMsg"
	eventIncomingMsgBunch           eventType = "IncomingMsgBunch"
	eventMsgsNoticed                eventType = "MsgsNoticed"
	eventMsgDelivered               eventType = "MsgDelivered"
	eventMsgFailed                  eventType = "MsgFailed"
	eventMsgRead                    eventType = "MsgRead"
	eventChatModified               eventType = "ChatModified"
	eventChatEphemeralTimerModified eventType = "ChatEphemeralTimerModified"
	eventContactsChanged            eventType = "ContactsChanged"
	eventLocationChanged            eventType = "LocationChanged"
	eventConfigureProgress          eventType = "ConfigureProgress"
	eventImexProgress               eventType = "ImexProgress"
	eventImexFileWritten            eventType = "ImexFileWritten"
	eventSecurejoinInviterProgress  eventType = "SecurejoinInviterProgress"
	eventSecurejoinJoinerProgress   eventType = "SecurejoinJoinerProgress"
	eventConnectivityChanged        eventType = "ConnectivityChanged"
	eventSelfavatarChanged          eventType = "SelfavatarChanged"
	eventWebxdcStatusUpdate         eventType = "WebxdcStatusUpdate"
	eventWebxdcInstanceDeleted      eventType = "WebxdcInstanceDeleted"
)

// Delta Chat core Event
type Event interface {
	eventType() eventType
}

// The library-user may write an informational string to the log.
//
// This event should *not* be reported to the end-user using a popup or something like
// that.
type EventInfo struct {
	Msg string
}

func (self EventInfo) eventType() eventType {
	return eventInfo
}

// Emitted when SMTP connection is established and login was successful.
type EventSmtpConnected struct {
	Msg string
}

func (self EventSmtpConnected) eventType() eventType {
	return eventSmtpConnected
}

// Emitted when IMAP connection is established and login was successful.
type EventImapConnected struct {
	Msg string
}

func (self EventImapConnected) eventType() eventType {
	return eventImapConnected
}

// Emitted when a message was successfully sent to the SMTP server.
type EventSmtpMessageSent struct {
	Msg string
}

func (self EventSmtpMessageSent) eventType() eventType {
	return eventSmtpMessageSent
}

// Emitted when an IMAP message has been marked as deleted
type EventImapMessageDeleted struct {
	Msg string
}

func (self EventImapMessageDeleted) eventType() eventType {
	return eventImapMessageDeleted
}

// Emitted when an IMAP message has been moved
type EventImapMessageMoved struct {
	Msg string
}

func (self EventImapMessageMoved) eventType() eventType {
	return eventImapMessageMoved
}

// Emitted before going into IDLE on the Inbox folder.
type EventImapInboxIdle struct{}

func (self EventImapInboxIdle) eventType() eventType {
	return eventImapInboxIdle
}

// Emitted when an new file in the $BLOBDIR was created
type EventNewBlobFile struct {
	File string
}

func (self EventNewBlobFile) eventType() eventType {
	return eventNewBlobFile
}

// Emitted when an file in the $BLOBDIR was deleted
type EventDeletedBlobFile struct {
	File string
}

func (self EventDeletedBlobFile) eventType() eventType {
	return eventDeletedBlobFile
}

// The library-user should write a warning string to the log.
//
// This event should *not* be reported to the end-user using a popup or something like
// that.
type EventWarning struct {
	Msg string
}

func (self EventWarning) eventType() eventType {
	return eventWarning
}

// The library-user should report an error to the end-user.
//
// As most things are asynchronous, things may go wrong at any time and the user
// should not be disturbed by a dialog or so.  Instead, use a bubble or so.
//
// However, for ongoing processes (eg. Account.Configure())
// or for functions that are expected to fail (eg. Message.AutocryptContinueKeyTransfer())
// it might be better to delay showing these events until the function has really
// failed (returned false). It should be sufficient to report only the *last* error
// in a messasge box then.
type EventError struct {
	Msg string
}

func (self EventError) eventType() eventType {
	return eventError
}

// An action cannot be performed because the user is not in the group.
// Reported eg. after a call to
// Chat.SetName(), Chat.SetImage(),
// Chat.AddContact(), Chat.RemoveContact(),
// and messages sending functions.
type EventErrorSelfNotInGroup struct {
	Msg string
}

func (self EventErrorSelfNotInGroup) eventType() eventType {
	return eventErrorSelfNotInGroup
}

// Messages or chats changed.  One or more messages or chats changed for various
// reasons in the database:
// - Messages sent, received or removed
// - Chats created, deleted or archived
// - A draft has been set
//
// ChatId is set if only a single chat is affected by the changes, otherwise 0.
// MsgId is set if only a single message is affected by the changes, otherwise 0.
type EventMsgsChanged struct {
	ChatId ChatId
	MsgId  MsgId
}

func (self EventMsgsChanged) eventType() eventType {
	return eventMsgsChanged
}

// Reactions for the message changed.
type EventReactionsChanged struct {
	ChatId    ChatId
	MsgId     MsgId
	ContactId ContactId
}

func (self EventReactionsChanged) eventType() eventType {
	return eventReactionsChanged
}

// There is a fresh message. Typically, the user will show an notification
// when receiving this message.
//
// There is no extra EventMsgsChanged event send together with this event.
type EventIncomingMsg struct {
	ChatId ChatId
	MsgId  MsgId
}

func (self EventIncomingMsg) eventType() eventType {
	return eventIncomingMsg
}

// Downloading a bunch of messages just finished. This is an experimental
// event to allow the UI to only show one notification per message bunch,
// instead of cluttering the user with many notifications.
//
// msg_ids contains the message ids.
type EventIncomingMsgBunch struct {
	MsgIds []MsgId
}

func (self EventIncomingMsgBunch) eventType() eventType {
	return eventIncomingMsgBunch
}

// Messages were seen or noticed.
// chat id is always set.
type EventMsgsNoticed struct {
	ChatId ChatId
}

func (self EventMsgsNoticed) eventType() eventType {
	return eventMsgsNoticed
}

// A single message is sent successfully. State changed from  MsgStateOutPending to
// MsgStateOutDelivered.
type EventMsgDelivered struct {
	ChatId ChatId
	MsgId  MsgId
}

func (self EventMsgDelivered) eventType() eventType {
	return eventMsgDelivered
}

// A single message could not be sent. State changed from MsgStateOutPending or MsgStateOutDelivered to
// MsgStateOutFailed.
type EventMsgFailed struct {
	ChatId ChatId
	MsgId  MsgId
}

func (self EventMsgFailed) eventType() eventType {
	return eventMsgFailed
}

// A single message is read by the receiver. State changed from MsgStateOutDelivered to
// MsgStateOutMdnRcvd.
type EventMsgRead struct {
	ChatId ChatId
	MsgId  MsgId
}

func (self EventMsgRead) eventType() eventType {
	return eventMsgRead
}

// Chat changed.  The name or the image of a chat group was changed or members were added or removed.
// Or the verify state of a chat has changed.
// See Chat.SetName(), Chat.SetImage(), Chat.AddContact()
// and Chat.RemoveContact().
//
// This event does not include ephemeral timer modification, which
// is a separate event.
type EventChatModified struct {
	ChatId ChatId
}

func (self EventChatModified) eventType() eventType {
	return eventChatModified
}

// Chat ephemeral timer changed.
type EventChatEphemeralTimerModified struct {
	ChatId ChatId
	Timer  int
}

func (self EventChatEphemeralTimerModified) eventType() eventType {
	return eventChatEphemeralTimerModified
}

// Contact(s) created, renamed, blocked or deleted.
type EventContactsChanged struct {
	// The id of contact that has changed, or zero if several contacts have changed.
	ContactId ContactId
}

func (self EventContactsChanged) eventType() eventType {
	return eventContactsChanged
}

// Location of one or more contact has changed.
type EventLocationChanged struct {
	// The id of contact for which the location has changed, or zero if the locations of several contacts have been changed.
	ContactId ContactId
}

func (self EventLocationChanged) eventType() eventType {
	return eventLocationChanged
}

// Inform about the configuration progress started by Account.Configure().
type EventConfigureProgress struct {
	// Progress.
	// 0=error, 1-999=progress in permille, 1000=success and done
	Progress uint

	// Optional progress comment or error, something to display to the user.
	Comment string
}

func (self EventConfigureProgress) eventType() eventType {
	return eventConfigureProgress
}

// Inform about the import/export progress.
type EventImexProgress struct {
	// Progress.
	// (usize) 0=error, 1-999=progress in permille, 1000=success and done
	Progress uint
}

func (self EventImexProgress) eventType() eventType {
	return eventImexProgress
}

// A file has been exported.
// This event may be sent after a call to Account.ExportBackup() or Account.ExportSelfKeys().
//
// A typical purpose for a handler of this event may be to make the file public to some system
// services.
type EventImexFileWritten struct {
	Path string
}

func (self EventImexFileWritten) eventType() eventType {
	return eventImexFileWritten
}

// Progress information of a secure-join handshake from the view of the inviter
// (Alice, the person who shows the QR code).
//
// These events are typically sent after a joiner has scanned the QR code
// generated by Account.QrCode() or Chat.QrCode().
type EventSecurejoinInviterProgress struct {
	// ID of the contact that wants to join.
	ContactId ContactId

	// Progress as:
	// 300=vg-/vc-request received, typically shown as "bob@addr joins".
	// 600=vg-/vc-request-with-auth received, vg-member-added/vc-contact-confirm sent, typically shown as "bob@addr verified".
	// 800=vg-member-added-received received, shown as "bob@addr securely joined GROUP", only sent for the verified-group-protocol.
	// 1000=Protocol finished for this contact.
	Progress uint
}

func (self EventSecurejoinInviterProgress) eventType() eventType {
	return eventSecurejoinInviterProgress
}

// Progress information of a secure-join handshake from the view of the joiner
// (Bob, the person who scans the QR code).
// The events are typically sent while Account.SecureJoin(), which
// may take some time, is executed.
type EventSecurejoinJoinerProgress struct {
	// ID of the inviting contact.
	ContactId ContactId

	// Progress as:
	// 400=vg-/vc-request-with-auth sent, typically shown as "alice@addr verified, introducing myself."
	// (Bob has verified alice and waits until Alice does the same for him)
	Progress uint
}

func (self EventSecurejoinJoinerProgress) eventType() eventType {
	return eventSecurejoinJoinerProgress
}

// The connectivity to the server changed.
// This means that you should refresh the connectivity view
// and possibly the connectivtiy HTML; see Account.Connectivity() and
// Account.ConnectivityHtml() for details.
type EventConnectivityChanged struct{}

func (self EventConnectivityChanged) eventType() eventType {
	return eventConnectivityChanged
}

// The user's avatar changed.
type EventSelfavatarChanged struct{}

func (self EventSelfavatarChanged) eventType() eventType {
	return eventSelfavatarChanged
}

// Webxdc status update received.
type EventWebxdcStatusUpdate struct {
	MsgId              MsgId
	StatusUpdateSerial uint
}

func (self EventWebxdcStatusUpdate) eventType() eventType {
	return eventWebxdcStatusUpdate
}

// Inform that a message containing a webxdc instance has been deleted
type EventWebxdcInstanceDeleted struct {
	MsgId MsgId
}

func (self EventWebxdcInstanceDeleted) eventType() eventType {
	return eventWebxdcInstanceDeleted
}
