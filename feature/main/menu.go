package main_feature

import (
	"context"
	"fmt"
	"io"
	"majotabi-bot/lib/config"
	"majotabi-bot/lib/helper"
	"majotabi-bot/lib/registry"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
)

func init() {
	registry.Register("menu", MenuHandler)
	registry.Register("allmenu", MenuHandler)
}

func MenuHandler(ctx context.Context, client *whatsmeow.Client, msg *events.Message) {
	txt := ""
	if msg.Message.Conversation != nil {
		txt = *msg.Message.Conversation
	} else if msg.Message.ExtendedTextMessage != nil {
		txt = *msg.Message.ExtendedTextMessage.Text
	} else if msg.Message.ImageMessage != nil && msg.Message.ImageMessage.Caption != nil {
		txt = *msg.Message.ImageMessage.Caption
	}

	args := strings.Split(txt, " ")
	command := strings.TrimPrefix(args[0], config.Current.Prefix)
	param := ""
	if len(args) > 1 {
		param = args[1]
	}

	greeting, emoji := getGreeting()
	header := fmt.Sprintf("*%s,* @%s %s\n_Watashi wa Elaina - Go !_\n", greeting, msg.Info.Sender.User, emoji)
	
	var menuText string
	var subHeader string

	if command == "allmenu" {
		subHeader = "Berikut adalah *SEMUA* menu ku !\n"
		menuText = generateTreeAll()
	} else if param != "" {
		subHeader = fmt.Sprintf("Berikut adalah menu *%s* !\n", strings.ToUpper(param))
		menuText = generateTreeCategory(param)
	} else {
		subHeader = "Berikut adalah semua kategori *MENU* ku !\n"
		menuText = generateCategoryList()
	}

	footer := "\n_Majotabi - Go by Lipp Majotabi ☘️_"
	finalText := header + subHeader + menuText + footer

	imgData, err := downloadImage(config.Current.Thumbnail)
	if err != nil {
		client.SendMessage(ctx, msg.Info.Chat, &waE2E.Message{
			ExtendedTextMessage: &waE2E.ExtendedTextMessage{
				Text:        &finalText,
				ContextInfo: helper.GetContext(msg),
			},
		})
		return
	}

	uploaded, err := client.Upload(ctx, imgData, whatsmeow.MediaImage)
	if err != nil {
		client.SendMessage(ctx, msg.Info.Chat, &waE2E.Message{
			ExtendedTextMessage: &waE2E.ExtendedTextMessage{
				Text:        &finalText,
				ContextInfo: helper.GetContext(msg),
			},
		})
		return
	}

	mimeType := "image/jpeg"
	client.SendMessage(ctx, msg.Info.Chat, &waE2E.Message{
		ImageMessage: &waE2E.ImageMessage{
			Caption:       &finalText,
			URL:           &uploaded.URL,
			DirectPath:    &uploaded.DirectPath,
			MediaKey:      uploaded.MediaKey,
			Mimetype:      &mimeType,
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    &uploaded.FileLength,
			ContextInfo:   helper.GetContext(msg),
		},
	})
}

func getGreeting() (string, string) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	hour := time.Now().In(loc).Hour()

	if hour >= 3 && hour < 11 {
		return "Ohayouu", "🌅"
	} else if hour >= 11 && hour < 15 {
		return "Konnichiwa", "☀️"
	} else if hour >= 15 && hour < 18 {
		return "Konnichiwa", "🌇"
	} else {
		return "Konbanwa", "🌙"
	}
}

func generateCategoryList() string {
	entries, err := os.ReadDir("feature")
	if err != nil {
		return "Error reading features."
	}

	res := "\n`📁 Feature`\n"
	var dirs []string
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, titleCase(e.Name()))
		}
	}

	for i, d := range dirs {
		if i == len(dirs)-1 {
			res += fmt.Sprintf("└── %s\n", d)
		} else {
			res += fmt.Sprintf("├── %s\n", d)
		}
	}
	res += "\n💡 Tutorial: Ketik *" + config.Current.Prefix + "menu main*"
	return res
}

func generateTreeCategory(category string) string {
	path := filepath.Join("feature", category)
	entries, err := os.ReadDir(path)
	if err != nil {
		return "⚠️ Kategori tidak ditemukan."
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".go") {
			files = append(files, strings.TrimSuffix(e.Name(), ".go"))
		}
	}

	res := fmt.Sprintf("\n`📁 %s`\n", titleCase(category))
	for i, f := range files {
		prefix := config.Current.Prefix
		if i == len(files)-1 {
			res += fmt.Sprintf("└── %s%s\n", prefix, f)
		} else {
			res += fmt.Sprintf("├── %s%s\n", prefix, f)
		}
	}
	return res
}

func generateTreeAll() string {
	res := ""
	
	dirs, _ := os.ReadDir("feature")
	for _, d := range dirs {
		if d.IsDir() {
			res += fmt.Sprintf("\n`📁 %s`\n", titleCase(d.Name()))
			
			files, _ := os.ReadDir(filepath.Join("feature", d.Name()))
			var cleanFiles []string
			for _, f := range files {
				if !f.IsDir() && strings.HasSuffix(f.Name(), ".go") {
					cleanFiles = append(cleanFiles, strings.TrimSuffix(f.Name(), ".go"))
				}
			}

			for i, f := range cleanFiles {
				prefix := config.Current.Prefix
				if i == len(cleanFiles)-1 {
					res += fmt.Sprintf("└── %s%s\n", prefix, f)
				} else {
					res += fmt.Sprintf("├── %s%s\n", prefix, f)
				}
			}
		}
	}
	return res
}

func titleCase(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.ToUpper(str[:1]) + str[1:]
}

func downloadImage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
