package random_feature

import (
	"context"
	"fmt"
	"majotabi-bot/lib/helper"
	"majotabi-bot/lib/registry"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func init() {
	registry.Register("neko", NekoHandler)
}

func NekoHandler(ctx context.Context, client *whatsmeow.Client, msg *events.Message) {
	apiURL := "https://api.nefyu.my.id/api/waifu-sfw/neko"
	caption := fmt.Sprintf("Random Neko ✨\n`%s`", apiURL)
	helper.SendImageFromURL(ctx, client, msg, apiURL, caption)
}
