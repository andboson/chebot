package repositories

import (
	"fmt"
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
	platform = detectPlatform(message)

	text = message.Text
	id = message.From.ID
	ctx, _ := userContexts[id]

	if ctx != "" {

		switch ctx {
		case "kino":
			sendFilmsReplyMessage(&message, text, platform)
			setUserContext(id, "")
		}

		return
	}

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

func detectPlatform(activity skypeapi.Activity) string {
	var platform = "web"
	if len(activity.Entities) > 0 {
		entity, ok := activity.Entities[0].(map[string]string)
		if ok {
			platformRaw, ok := entity["platform"]
			if ok {
				platform = strings.ToLower(platformRaw)
			}
		}
	}

	return platform
}

func sendFilmsReplyMessage(activity *skypeapi.Activity, location, platform string) {
	log.Printf("activity: %s  ----- %+v", activity.Action, activity)
	name, ok := KinoNamesRu[location]
	if !ok {
		skypeapi.SendReplyMessage(activity, "Не знаю такое место", SkypeToken.AccessToken)
		log.Printf("unknown location: %s", location)
		return
	}
	skypeapi.SendReplyMessage(activity, "Фильмы в "+name, SkypeToken.AccessToken)
	films := GetMovies(location)

	var text string
	for _, film := range films {
		var filmText = " \n "
		filmText += fmt.Sprintf("\r\n **%s** ", film.Title)
		filmText += fmt.Sprintf("\n `%s` \n ", film.TimeBlock)
		//filmText += fmt.Sprintf("![img](%s)", URL_PREFIX+"/"+film.Img)
		filmText += fmt.Sprintf("![img](%s)", "http://www.publicdomainpictures.net/pictures/30000/t2/duck-on-a-rock.jpg")
		text += filmText
	}
	skypeapi.SendReplyMessage(activity, text, SkypeToken.AccessToken)
}

func setUserContext(id string, ctx string) {
	// clear context
	if ctx == "" {
		userContexts[id] = ""
		return
	}

	// check and hold
	_, ok := userContexts[id]
	if !ok {
		userContexts[id] = ctx
		go holdUserContext(id)
		return
	}

	userContextsUpdated[id] <- true
}

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

func sendChoicePlaceReplyMessage(activity *skypeapi.Activity, message, authorizationToken string) error {
	responseActivity := &skypeapi.Activity{
		Type:         activity.Type,
		From:         activity.Recipient,
		Conversation: activity.Conversation,
		Recipient:    activity.From,
		Text:         message,
		InputHint:    "место (lyubava\\plaza)",
		SuggestedActions: skypeapi.SuggestedActions{
			Actions: []skypeapi.CardAction{
				{
					Title: "Любава",
					Type:  "imBack",
					Value: "lyubava",
				},
				{
					Title: "Днепроплаза",
					Type:  "imBack",
					Value: "plaza",
				},
			},
		},
		ReplyToID: activity.ID,
	}
	replyUrl := fmt.Sprintf("%vv3/conversations/%v/activities", activity.ServiceURL, activity.Conversation.ID)

	return skypeapi.SendActivityRequest(responseActivity, replyUrl, authorizationToken)
}
