package main_feature

import (
	"context"
	"majotabi-bot/lib/config"
	"majotabi-bot/lib/helper"
	"majotabi-bot/lib/registry"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
)

func init() {
	registry.Register("elaina", ElainaHandler)
}

func ElainaHandler(ctx context.Context, client *whatsmeow.Client, msg *events.Message) {
	txt := ""
	if msg.Message.Conversation != nil {
		txt = *msg.Message.Conversation
	} else if msg.Message.ExtendedTextMessage != nil {
		txt = *msg.Message.ExtendedTextMessage.Text
	}
	args := strings.Split(txt, " ")

	var resp string
	if len(args) > 1 {
		switch strings.ToLower(args[1]) {
		case "on":
			config.IsElainaActive = true
			resp = "Baiklah, aku akan menemanimu. Tapi jangan berharap banyak, ya! 츤데레"
		case "off":
			config.IsElainaActive = false
			resp = "Hmph, akhirnya aku bisa istirahat."
		default:
			resp = "Gunakan `!elaina on` atau `!elaina off`."
		}
	} else {
		resp = "Gunakan `!elaina on` atau `!elaina off`."
	}
	
	client.SendMessage(ctx, msg.Info.Chat, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text:        &resp,
			ContextInfo: helper.GetContext(msg),
		},
	})
}
