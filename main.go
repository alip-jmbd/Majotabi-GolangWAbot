package main

import (
	"context"
	"fmt"
	"majotabi-bot/feature/ai"
	"majotabi-bot/lib/config"
	"majotabi-bot/lib/database"
	"majotabi-bot/lib/helper"
	"majotabi-bot/lib/registry"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "majotabi-bot/feature/main"
	_ "majotabi-bot/feature/owner"
	_ "majotabi-bot/feature/random"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var botStartTime time.Time
var wibLocation *time.Location

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
)

const elainaArt = `
⣿⣿⣿⣿⣿⣿⡿⠿⠿⠿⢿⡶⠶⣶⣶⣴⣯⠿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠿⣏⣭⣭⣽⣿⣻⣿⣿⣿⣿⣿⣿⣿
⣿⣿⣿⣿⣿⠟⠁⠀⠀⢀⣀⣀⠉⠉⠚⠋⣝⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡵⢿⡛⠛⠛⠉⠉⠉⠩⣼⣿⣿⣿⣿⣿
⣿⣿⣟⡋⠁⠀⢀⣴⣿⣿⣿⠋⠁⠀⠀⠀⠨⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠚⠁⠀⠠⣶⣶⣦⣄⠀⠀⠙⠿⣿⣿⣿
⣿⣿⠟⠁⠀⣴⣿⣿⣿⣿⣟⠀⠣⡉⢨⠆⢐⣜⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣏⠠⡒⠤⡆⠘⣿⣿⣿⣿⣄⠀⠘⢿⣿⣿
⣍⣀⣀⣀⠀⢿⣿⣿⣿⣿⣿⣄⣀⠈⢃⣠⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡀⠐⠔⠃⢰⣿⣿⣿⣿⣿⠆⠀⣀⣈⣙
⣿⣿⣿⣿⣷⣶⣭⣿⣿⢿⡿⠟⣉⣩⣭⣿⣿⣿⣿⣿⠀⣿⠀⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣒⡒⡚⠻⣿⣿⣿⣿⣵⣾⣿⣿
⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿
`

func main() {
	botStartTime = time.Now()
	wibLocation, _ = time.LoadLocation("Asia/Jakarta")

	if err := config.LoadConfig(); err != nil {
		panic(fmt.Sprintf("Gagal load config: %v", err))
	}

	if err := database.Connect(); err != nil {
		panic(fmt.Sprintf("Gagal connect database: %v", err))
	}

	clientLog := waLog.Stdout("Client", "ERROR", true)
	deviceStore, err := database.Container.GetFirstDevice(context.Background())
	if err != nil {
		panic(err)
	}

	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(func(evt interface{}) {
		EventHandler(client, evt)
	})

	if client.Store.ID == nil {
		fmt.Println("\n===========================================")
		fmt.Println("   MAJOTABI BOT - SETUP LOGIN")
		fmt.Println("===========================================")
		fmt.Println("1. QR Code")
		fmt.Println("2. Pairing Code (Rekomendasi)")
		fmt.Print("Pilih (1/2): ")
		
		var choice int
		_, err := fmt.Scanln(&choice)
		if err != nil {
			choice = 2
		}

		if choice == 2 {
			err = client.Connect()
			if err != nil {
				panic(fmt.Sprintf("Gagal connect: %v", err))
			}
			
			fmt.Println(colorCyan + ">> Menunggu kestabilan koneksi (5 detik)..." + colorReset)
			time.Sleep(5 * time.Second)

			fmt.Print("\n" + colorYellow + "Masukkan Nomor HP (Format: 628xxx): " + colorReset)
			var phone string
			fmt.Scanln(&phone)
			phone = strings.TrimSpace(phone)

			fmt.Println(colorCyan + ">> Meminta Pairing Code..." + colorReset)
			code, err := client.PairPhone(context.Background(), phone, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
			if err != nil {
				fmt.Printf("\n%sERROR: %v%s\n", colorRed, err, colorReset)
				os.Exit(1)
			}
			
			fmt.Println(colorCyan + elainaArt + colorReset)
			fmt.Println("==========================================================================================================")
			fmt.Printf(" %sKODE PAIRING ANDA:%s  %s%s%s \n", colorYellow, colorReset, colorGreen, code, colorReset)
			fmt.Println("==========================================================================================================")
		} else {
			qrChan, _ := client.GetQRChannel(context.Background())
			err = client.Connect()
			if err != nil {
				panic(err)
			}
			for evt := range qrChan {
				if evt.Event == "code" {
					qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				} else {
					fmt.Println("QR Event:", evt.Event)
				}
			}
		}
	} else {
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("\n%sBot %s Berjalan!%s\n%sOwner:%s %s / %s\n", colorGreen, config.Current.BotName, colorReset, colorYellow, colorReset, config.Current.OwnerNumber, config.Current.OwnerLID)
	
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
	database.Container.Close()
	fmt.Println("Bot Berhenti.")
}

var reactions = []string{"🍀", "🍃", "🍂", "🍁", "🍎"}

func EventHandler(client *whatsmeow.Client, evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		if v.Info.Timestamp.Before(botStartTime) {
			return
		}

		msg := v.Message
		if msg == nil {
			return
		}
		
		txt := ""
		if msg.Conversation != nil {
			txt = *msg.Conversation
		} else if msg.ExtendedTextMessage != nil {
			txt = *msg.ExtendedTextMessage.Text
		} else if msg.ImageMessage != nil && msg.ImageMessage.Caption != nil {
			txt = *msg.ImageMessage.Caption
		}

		senderJID := v.Info.Sender
		senderUser := senderJID.User
		
		isOwner := senderUser == config.Current.OwnerNumber || senderUser == config.Current.OwnerLID
		
		isMentioned := false
		if extMsg := msg.GetExtendedTextMessage(); extMsg != nil && extMsg.ContextInfo != nil {
			for _, mentionedJIDString := range extMsg.ContextInfo.MentionedJID {
				if client.Store.ID != nil && mentionedJIDString == client.Store.ID.String() {
					isMentioned = true
					break
				}
			}
		}

		isReplyToBot := false
		if extMsg := msg.GetExtendedTextMessage(); extMsg != nil && extMsg.ContextInfo != nil {
			if extMsg.ContextInfo.GetParticipant() == client.Store.ID.String() {
				isReplyToBot = true
			}
		}

		if config.IsElainaActive && (strings.Contains(strings.ToLower(txt), "elaina") || isMentioned || isReplyToBot) {
			var imageData []byte
			var err error
			if msg.ImageMessage != nil {
				imageData, err = client.Download(context.Background(), msg.ImageMessage)
				if err != nil {
					fmt.Println("Gagal download gambar:", err)
				}
			}
			
			go func() {
				client.SendChatPresence(context.Background(), senderJID, types.ChatPresenceComposing, types.ChatPresenceMediaText)
				defer client.SendChatPresence(context.Background(), senderJID, types.ChatPresencePaused, types.ChatPresenceMediaText)

				response, err := ai.GenerateResponse(context.Background(), senderJID.String(), txt, imageData)
				if err != nil {
					fmt.Println("Error dari AI:", err)
					response = "Hmph, sepertinya aku sedang tidak mood bicara."
				}
				
				client.SendMessage(context.Background(), v.Info.Chat, &waE2E.Message{
					ExtendedTextMessage: &waE2E.ExtendedTextMessage{
						Text:        &response,
						ContextInfo: helper.GetContext(v),
					},
				})
			}()
			return
		}

		if !config.IsPublic && !isOwner {
			return
		}

		chatType := "PC"
		chatColor := colorPurple
		if v.Info.IsGroup {
			chatType = "GC"
			chatColor = colorGreen
		}

		tag := "[USER]"
		tagColor := colorBlue
		if isOwner {
			tag = "[OWNER]"
			tagColor = colorRed
		}
		
		timestamp := v.Info.Timestamp.In(wibLocation).Format("15:04:05")
		if txt != "" {
			fmt.Printf("%s[%s]%s [%s%s%s] %s%s%s %s%s%s: %s\n",
				colorYellow, timestamp, colorReset,
				chatColor, chatType, colorReset,
				tagColor, tag, colorReset,
				colorCyan, senderUser, colorReset,
				txt)
		}

		if strings.HasPrefix(txt, config.Current.Prefix) {
			rand.Seed(time.Now().UnixNano())
			reaction := reactions[rand.Intn(len(reactions))]
			client.SendMessage(context.Background(), v.Info.Chat, client.BuildReaction(v.Info.Chat, v.Info.Sender, v.Info.ID, reaction))
			
			args := strings.Split(txt, " ")
			cmd := strings.TrimPrefix(args[0], config.Current.Prefix)
			
			if (cmd == "self" || cmd == "public" || cmd == "elaina") && !isOwner {
				return
			}
			
			if handler, ok := registry.GetHandler(cmd); ok {
				go func() {
					handler(context.Background(), client, v)
					time.Sleep(1 * time.Second)
					client.SendMessage(context.Background(), v.Info.Chat, client.BuildReaction(v.Info.Chat, v.Info.Sender, v.Info.ID, ""))
				}()
			}
		}
	}
}
