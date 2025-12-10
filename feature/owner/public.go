package owner_feature

import (
	"context"
	"majotabi-bot/lib/config"
	"majotabi-bot/lib/helper"
	"majotabi-bot/lib/registry"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
)

func init() {
	registry.Register("public", PublicHandler)
}

func PublicHandler(ctx context.Context, client *whatsmeow.Client, msg *events.Message) {
	config.IsPublic = true
	resp := "🔓 *Mode Public Aktif*\nSemua orang bisa menggunakan bot."
	client.SendMessage(ctx, msg.Info.Chat, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text:        &resp,
			ContextInfo: helper.GetContext(msg),
		},
	})
}
