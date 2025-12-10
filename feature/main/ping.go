package main_feature

import (
	"context"
	"fmt"
	"majotabi-bot/lib/helper"
	"majotabi-bot/lib/registry"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
)

func init() {
	registry.Register("ping", PingHandler)
}

func PingHandler(ctx context.Context, client *whatsmeow.Client, msg *events.Message) {
	requestTime := msg.Info.Timestamp
	now := time.Now()
	latency := now.Sub(requestTime)
	
	if latency < 0 {
		latency = -latency
	}

	resp := fmt.Sprintf("Pong! 🏓\nLatency: %s\nFast Response", latency)
	
	client.SendMessage(ctx, msg.Info.Chat, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text:        &resp,
			ContextInfo: helper.GetContext(msg),
		},
	})
}
