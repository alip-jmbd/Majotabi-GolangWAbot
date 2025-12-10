package registry

import (
	"context"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

type HandlerFunc func(ctx context.Context, client *whatsmeow.Client, msg *events.Message)

var features = make(map[string]HandlerFunc)

func Register(command string, handler HandlerFunc) {
	features[command] = handler
}

func GetHandler(command string) (HandlerFunc, bool) {
	h, ok := features[command]
	return h, ok
}
