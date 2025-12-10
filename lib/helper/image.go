package helper

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

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

	for i := 0; i < 5; i++ {
		resp, err := http.Get(imageURL)
		if err != nil {
			fmt.Printf("Gagal HTTP GET dari %s (Percobaan %d): %v\n", imageURL, i+1, err)
			time.Sleep(200 * time.Millisecond)
			continue
		}

		if resp.StatusCode != 200 {
			fmt.Printf("API mengembalikan status non-200: %d (Percobaan %d)\n", resp.StatusCode, i+1)
			resp.Body.Close()
			time.Sleep(200 * time.Millisecond)
			continue
		}

		mediaData, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Println("Gagal Read Body:", err)
			sendFailure("Aku gagal membaca data medianya.")
			return
		}
		
		contentType := http.DetectContentType(mediaData)
		if !strings.Contains(contentType, "image") && !strings.Contains(contentType, "gif") {
			fmt.Printf("Konten bukan gambar/gif (%s), mencoba lagi...\n", contentType)
			time.Sleep(200 * time.Millisecond)
			continue
		}

		sendMediaWithUpload(ctx, client, msg, mediaData, contentType, caption)
		return
	}

	sendFailure("Hmph, aku gagal mendapatkan media yang valid setelah beberapa kali mencoba.")
}

func sendMediaWithUpload(ctx context.Context, client *whatsmeow.Client, msg *events.Message, mediaData []byte, contentType, caption string) {
	var uploaded whatsmeow.UploadResponse
	var err error
	var message *waE2E.Message

	if strings.Contains(contentType, "gif") {
		uploaded, err = client.Upload(ctx, mediaData, whatsmeow.MediaVideo)
		if err != nil {
			fmt.Println("Gagal Upload Fallback GIF:", err)
			return
		}
		
		isGif := true
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
				GifPlayback:   &isGif,
				Caption:       &caption,
			},
		}
	} else {
		uploaded, err = client.Upload(ctx, mediaData, whatsmeow.MediaImage)
		if err != nil {
			fmt.Println("Gagal Upload Fallback Gambar:", err)
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
