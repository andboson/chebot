package models

import (
	"encoding/json"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"sync"
)

const ConfigFile = "config.json"

var Conf *Config

type Config struct {
	IncomingGoogleToken string `json:"incoming_google_token"`
	TelegramToken       string `json:"telegram_token"`
	SkypePassword       string `json:"skype_pass"`
	SkypeAppID          string `json:"skype_id"`
	FbPageToken         string `json:"fb_page_token"`
	FbAppSecret         string `json:"fb_app_secret"`
	FbVerifyToken       string `json:"fb_verify_token"`
}

var CmdList = []string{"1. kino  - Фильмы в кинотеатрах", "2. taxi - Список такси"}

func InitConfig() {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	Conf = LoadConfig()
}

func LoadConfig() *Config {
	var conf Config
	result, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		log.Fatalf("Unable to load config file: %s with error: %s", ConfigFile, err)
	}
	err = json.Unmarshal(result, &conf)
	if err != nil {
		log.Fatalf("Unable to unmarshall config file: %s with error: %s", ConfigFile, err)
	}

	return &conf
}
