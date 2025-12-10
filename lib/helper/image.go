package helper

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
)

func SendImageFromURL(ctx context.Context, client *whatsmeow.Client, msg *events.Message, imageURL string, caption string) {
	sendFailure := func(text string) {
		client.SendMessage(ctx, msg.Info.Chat, &waE2E.Message{
			ExtendedTextMessage: &waE2E.ExtendedTextMessage{
				Text:        &text,
				ContextInfo: GetContext(msg),
			},
		})
	}

	resp, err := http.Get(imageURL)
	if err != nil {
		fmt.Printf("Gagal HTTP GET dari %s: %v\n", imageURL, err)
		sendFailure("Hmph, aku tidak bisa mengambil gambar dari URL itu.")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("API mengembalikan status non-200: %d\n", resp.StatusCode)
		sendFailure("Hmph, sepertinya API-nya sedang bermasalah.")
		return
	}

	mediaData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Gagal Read Body:", err)
		sendFailure("Aku gagal membaca data medianya.")
		return
	}
	
	contentType := http.DetectContentType(mediaData)
	var message *waE2E.Message

	if strings.Contains(contentType, "gif") {
		uploaded, err := client.Upload(ctx, mediaData, whatsmeow.MediaVideo)
		if err != nil {
			fmt.Println("Gagal Upload GIF:", err)
			sendFailure("Aku gagal mengunggah GIF-nya.")
			return
		}
		
		isGifFlag := true
		mimetype := "video/mp4"
		message = &waE2E.Message{
			VideoMessage: &waE2E.VideoMessage{
				URL:           &uploaded.URL,
				DirectPath:    &uploaded.DirectPath,
				MediaKey:      uploaded.MediaKey,
				Mimetype:      &mimetype,
				FileEncSHA256: uploaded.FileEncSHA256,
				FileSHA256:    uploaded.FileSHA256,
				FileLength:    &uploaded.FileLength,
				ContextInfo:   GetContext(msg),
				GifPlayback:   &isGifFlag,
				Caption:       &caption,
			},
		}
	} else {
		uploaded, err := client.Upload(ctx, mediaData, whatsmeow.MediaImage)
		if err != nil {
			fmt.Println("Gagal Upload Gambar:", err)
			sendFailure("Aku gagal mengunggah gambarnya.")
			return
		}
		
		mimetype := "image/jpeg"
		message = &waE2E.Message{
			ImageMessage: &waE2E.ImageMessage{
				URL:           &uploaded.URL,
				DirectPath:    &uploaded.DirectPath,
				MediaKey:      uploaded.MediaKey,
				Mimetype:      &mimetype,
				FileEncSHA256: uploaded.FileEncSHA256,
				FileSHA256:    uploaded.FileSHA256,
				FileLength:    &uploaded.FileLength,
				ContextInfo:   GetContext(msg),
				Caption:       &caption,
			},
		}
	}

	client.SendMessage(ctx, msg.Info.Chat, message)
}
