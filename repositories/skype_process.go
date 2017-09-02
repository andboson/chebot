package repositories

import (
	"github.com/andboson/chebot/models"
	"github.com/andboson/skypeapi"
	"log"
	"strings"
	"time"
)

var userContextsUpdated map[string]chan bool
var userContexts map[string]string
var SkypeToken skypeapi.TokenResponse

func init() {
	mu.Lock()
	mu.Unlock()
	userContexts = make(map[string]string)
	userContextsUpdated = make(map[string]chan bool)
}

func InitSkype() {
	token, err := skypeapi.RequestAccessToken(models.Conf.SkypeAppID, models.Conf.SkypePassword)
	if err != nil {
		log.Printf("[---- SKYPE AUTH ERROR ----]  %s", err)
	}
	mu.Lock()
	defer mu.Unlock()
	SkypeToken = token

	log.Printf("[skype] Authorized")
	time.AfterFunc(time.Duration(token.ExpiresIn)*time.Second, InitSkype)
}

func ProcessSkypeMessage(message skypeapi.Activity) {
	var id string
	var text string
	var platform string
	var err error
	platform = detectPlatform(message)
	text = message.Text
	id = message.From.ID
	ctx, _ := userContexts[id]

	// help
	if text == "/?" {
		helpText := "Доступные команды:  \r\n  1. `kino`  - Фильмы в кинотеатрах"
		err = skypeapi.SendReplyMessage(&message, helpText, SkypeToken.AccessToken)
		if err != nil {
			log.Printf("[skype] error messaging: %s", err)
		}

		return
	}

	// process text with context
	if ctx != "" {
		switch ctx {
		case "kino":
			sendFilmsReplyMessage(&message, text, platform)
			setUserContext(id, "")
		}

		return
	}

	// catch commands if empty context
	if ctx == "" && (strings.ToLower(text) == "kino" || strings.ToLower(text) == "films") {
		setUserContext(id, "kino")
		var prompt = " (lyubava\\plaza)"
		if platform != "web" {
			prompt = ""
		}
		err := sendChoicePlaceReplyMessage(&message, "Выберите  кинотеатр"+prompt, SkypeToken.AccessToken)
		if err != nil {
			log.Printf("[skype] error messaging: %s", err)
		}
	}
}

// detect sender platform
func detectPlatform(activity skypeapi.Activity) string {
	var platform = "web"
	if len(activity.Entities) > 0 {
		entity, ok := activity.Entities[0].(map[string]interface{})
		if ok {
			platformRaw, ok := entity["platform"]
			if ok {
				platform = strings.ToLower(platformRaw.(string))
			}
		}
	}

	return platform
}

// set context
func setUserContext(id string, ctx string) {
	// clear context
	if ctx == "" {
		userContexts[id] = ""
		return
	}
	userContextsUpdated[id] = make(chan bool)

	// check and hold
	holded, ok := userContexts[id]
	if !ok || holded == "" {
		userContexts[id] = ctx
		go holdUserContext(id)
	}

	userContextsUpdated[id] <- true
}

// hold, update, forget context
func holdUserContext(id string) {
	defer func() {
		delete(userContexts, id)
	}()
	forgetContext := time.After(DEFAULT_CONTEXT_LIFETIME * time.Second)

	for {
		select {
		case <-userContextsUpdated[id]:
			forgetContext = time.After(DEFAULT_CONTEXT_LIFETIME * time.Second)
		case <-forgetContext:
			userContexts[id] = ""
			return
		}
	}
}
