package deltachat

import (
	"fmt"
	"regexp"
	"strings"
)

// Message snapshot.
type MsgSnapshot struct {
	Account *Account

	Id                    MsgId
	ChatId                ChatId
	FromId                ContactId
	Quote                 *MsgQuote
	ParentId              MsgId
	Text                  string
	HasLocation           bool
	HasHtml               bool
	ViewType              MsgType
	State                 MsgState
	Error                 string
	Timestamp             Timestamp
	SortTimestamp         Timestamp
	ReceivedTimestamp     Timestamp
	HasDeviatingTimestamp bool
	Subject               string
	ShowPadlock           bool
	IsSetupmessage        bool
	IsInfo                bool
	IsForwarded           bool
	IsBot                 bool
	SystemMessageType     SysmsgType
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
	DownloadState         DownloadState
	Reactions             *Reactions
}

// Extract metadata from system message with type SysmsgTypeMemberAddedToGroup.
func (self *MsgSnapshot) ParseMemberAdded() (actor *Contact, target *Contact, err error) {
	return self.parseMemberAddRemove("added")
}

// Extract metadata from system message with type SysmsgTypeMemberRemovedFromGroup.
func (self *MsgSnapshot) ParseMemberRemoved() (actor *Contact, target *Contact, err error) {
	return self.parseMemberAddRemove("removed")
}

func (self *MsgSnapshot) parseMemberAddRemove(action string) (actor *Contact, target *Contact, err error) {
	text := strings.ToLower(self.Text)
	actor = &Contact{self.Account, self.FromId}

	regex := regexp.MustCompile(`^member (.+) ` + action + ` by .+\.$`)
	match := regex.FindStringSubmatch(text)
	if len(match) > 0 {
		target, err := self.extractContact(match[1])
		if err != nil {
			return nil, nil, err
		}
		return actor, target, nil
	}

	regex = regexp.MustCompile(`^you ` + action + ` member (.+)\.$`)
	match = regex.FindStringSubmatch(text)
	if len(match) > 0 {
		target, err := self.extractContact(match[1])
		if err != nil {
			return nil, nil, err
		}
		return actor, target, nil
	}

	if action == "removed" {
		regex = regexp.MustCompile(`^group left by .+\.$`)
		match = regex.FindStringSubmatch(text)
		if len(match) > 0 {
			return actor, actor, nil
		}

		regex = regexp.MustCompile(`^you left the group\.$`)
		if regex.MatchString(text) {
			return actor, actor, nil
		}
	}

	return nil, nil, fmt.Errorf("System message does not match")
}

func (self *MsgSnapshot) extractContact(text string) (*Contact, error) {
	regex := regexp.MustCompile(`^.*\((.+@.+)\)$`)
	match := regex.FindStringSubmatch(text)
	if len(match) > 0 {
		text = match[1]
	}
	return self.Account.GetContactByAddr(text)
}
