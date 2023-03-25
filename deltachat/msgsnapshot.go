package deltachat

import (
	"fmt"
	"regexp"
	"strings"
)

// Message snapshot.
type MsgSnapshot struct {
	Account *Account

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

// Extract metadata from system message with type SYSMSG_TYPE_MEMBER_ADDED_TO_GROUP.
func (self *MsgSnapshot) ParseMemberAdded() (actor *Contact, target *Contact, err error) {
	action, actor, target, err := self.parseMemberAddRemove()
	if err != nil {
		return nil, nil, err
	}
	if action == "added" {
		return actor, target, nil
	}
	return nil, nil, fmt.Errorf("System message does not match")
}

// Extract metadata from system message with type SYSMSG_TYPE_MEMBER_REMOVED_FROM_GROUP.
func (self *MsgSnapshot) ParseMemberRemoved() (actor *Contact, target *Contact, err error) {
	action, actor, target, err := self.parseMemberAddRemove()
	if err != nil {
		return nil, nil, err
	}
	if action == "removed" {
		return actor, target, nil
	}
	return nil, nil, fmt.Errorf("System message does not match")
}

func (self *MsgSnapshot) parseMemberAddRemove() (string, *Contact, *Contact, error) {
	text := strings.ToLower(self.Text)
	actor := &Contact{self.Account, self.FromId}

	regex := regexp.MustCompile(`^member (.+) (removed|added) by .+\.$`)
	match := regex.FindStringSubmatch(text)
	if len(match) > 0 {
		target, err := self.extractContact(match[1])
		if err != nil {
			return "", nil, nil, err
		}
		return match[2], actor, target, nil
	}

	regex = regexp.MustCompile(`^you (removed|added) member (.+)\.$`)
	match = regex.FindStringSubmatch(text)
	if len(match) > 0 {
		target, err := self.extractContact(match[2])
		if err != nil {
			return "", nil, nil, err
		}
		return match[1], actor, target, nil
	}

	regex = regexp.MustCompile(`^group left by .+\.$`)
	match = regex.FindStringSubmatch(text)
	if len(match) > 0 {
		return "removed", actor, actor, nil
	}

	regex = regexp.MustCompile(`^you left the group\.$`)
	if regex.MatchString(text) {
		return "removed", actor, actor, nil
	}

	return "", nil, nil, fmt.Errorf("System message does not match")
}

func (self *MsgSnapshot) extractContact(text string) (*Contact, error) {
	regex := regexp.MustCompile(`^.*\((.+@.+)\)$`)
	match := regex.FindStringSubmatch(text)
	if len(match) > 0 {
		text = match[1]
	}
	return self.Account.GetContactByAddr(text)
}
