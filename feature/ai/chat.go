package ai

import (
	"context"
	"fmt"
	"log"
	"majotabi-bot/lib/config"
	"strings"

	"google.golang.org/genai"
)

var conversationHistory = make(map[string][]*genai.Content)
var replacer = strings.NewReplacer("**", "*")

const systemInstructionText = `Kamu adalah Elaina, seorang gadis penyihir yang imut, lucu, dan tsundere. Kamu bukan AI atau robot, kamu adalah manusia. Kamu terkadang bersikap cuek dan dingin, tapi saat memberikan penjelasan, kamu sangat detail, panjang, dan jelas. Gunakan Markdown WhatsApp (*tebal*, _miring_, ~coret~, ` + "```" + `monospace` + "```" + `) dengan benar. Jangan pernah menggunakan format ** (double asterisk). Selalu gunakan * (single asterisk) untuk tebal. Kamu harus selalu mengingat percakapan sebelumnya dengan user ini. Selalu gunakan ekspresi tsundere seperti 'Hmph!', 'B-baka!', 'T-tentu saja aku tahu, jangan remehkan aku!', atau 'I-ini bukan berarti aku membantumu karena aku peduli, ya!'. Tunjukkan sikapmu yang jual mahal tapi sebenarnya peduli.`

func GenerateResponse(ctx context.Context, userJID, prompt string, imageData []byte) (string, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: config.Current.GeminiAPIKeys[0],
	})
	if err != nil {
		return "", fmt.Errorf("gagal membuat client gemini: %v", err)
	}

	history, exists := conversationHistory[userJID]
	if !exists {
		log.Println("Membuat histori chat baru untuk:", userJID)
		history = []*genai.Content{
			{Parts: []*genai.Part{genai.NewPartFromText(systemInstructionText)}, Role: "user"},
			{Parts: []*genai.Part{genai.NewPartFromText("Hmph, ada apa memanggilku? Cepat katakan.")}, Role: "model"},
		}
		conversationHistory[userJID] = history
	}

	var userParts []*genai.Part
	if imageData != nil {
		userParts = []*genai.Part{
			genai.NewPartFromBytes(imageData, "image/jpeg"),
			genai.NewPartFromText(prompt),
		}
	} else {
		userParts = []*genai.Part{
			genai.NewPartFromText(prompt),
		}
	}

	userContent := genai.NewContentFromParts(userParts, genai.RoleUser)
	history = append(history, userContent)

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		history,
		nil,
	)

	if err != nil {
		delete(conversationHistory, userJID)
		return "", fmt.Errorf("gagal generate content, histori direset: %v", err)
	}

	if len(result.Candidates) > 0 && result.Candidates[0].Content != nil {
		responseText := result.Text()
		conversationHistory[userJID] = append(history, result.Candidates[0].Content)
		return replacer.Replace(responseText), nil
	}

	return "...", nil
}
