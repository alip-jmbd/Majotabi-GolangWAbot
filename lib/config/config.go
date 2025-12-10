package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	OwnerNumber   string   `json:"owner_number"`
	OwnerLID      string   `json:"owner_lid"`
	BotName       string   `json:"bot_name"`
	Prefix        string   `json:"prefix"`
	Thumbnail     string   `json:"thumbnail"`
	GeminiAPIKeys []string `json:"gemini_api_keys"`
}

var Current Config
var IsPublic bool = true
var IsElainaActive bool = false

func LoadConfig() error {
	file, err := os.ReadFile("config.json")
	if err != nil {
		return err
	}
	return json.Unmarshal(file, &Current)
}
