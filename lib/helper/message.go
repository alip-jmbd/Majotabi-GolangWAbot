package helper

import (
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
)

func GetContext(msg *events.Message) *waE2E.ContextInfo {
	senderString := msg.Info.Sender.String()

	ctxInfo := &waE2E.ContextInfo{
		StanzaID:      &msg.Info.ID,
		Participant:   &senderString,
		QuotedMessage: msg.Message,
		MentionedJID:  []string{senderString},
	}

	var expiration uint32
	if msg.Message != nil {
		if ext := msg.Message.ExtendedTextMessage; ext != nil && ext.ContextInfo != nil && ext.ContextInfo.Expiration != nil {
			expiration = *ext.ContextInfo.Expiration
		} else if ephem := msg.Message.EphemeralMessage; ephem != nil && ephem.Message != nil {
			if innerExt := ephem.Message.ExtendedTextMessage; innerExt != nil && innerExt.ContextInfo != nil && innerExt.ContextInfo.Expiration != nil {
				expiration = *innerExt.ContextInfo.Expiration
			}
		}
	}

	if expiration > 0 {
		ctxInfo.Expiration = &expiration
	}

	return ctxInfo
}
