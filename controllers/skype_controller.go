package controllers

import (
	"fmt"
	"github.com/andboson/chebot/models"
	"github.com/andboson/skypeapi"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"net/http"
	"sync"
	"time"
)

var SkypeToken skypeapi.TokenResponse

func InitSkype() {
	var mu sync.Mutex
	log.Printf("[token] request")

	token, err := skypeapi.RequestAccessToken(models.Conf.SkypeAppID, models.Conf.SkypePassword)
	if err != nil {
		log.Printf("[---- SKYPE AUTH ERROR ----]  %s", err)
	}
	mu.Lock()
	defer mu.Unlock()
	SkypeToken = token

	time.AfterFunc(time.Duration(token.ExpiresIn)*time.Second, InitSkype)
}

func SkypeHook(c echo.Context) error {

	//var req models.SkypeMessage
	var req skypeapi.Activity
	err := c.Bind(&req)
	if err != nil {
		log.Printf("--- skype decode msg error!: %+v  ---------- %s", c.Request(), err)

	}

	err = sendReplyMessage(&req, "hello!!1", SkypeToken.AccessToken)
	if err != nil {
		log.Printf("[skype] error messaging: %s", err)
	}

	resp := map[string]string{
		"status": "success",
	}

	return c.JSON(http.StatusOK, resp)
}

func sendReplyMessage(activity *skypeapi.Activity, message, authorizationToken string) error {
	responseActivity := &skypeapi.Activity{
		Type:         activity.Type,
		From:         activity.Recipient,
		Conversation: activity.Conversation,
		Recipient:    activity.From,
		InputHint:    "Выберите  место (lyubava\\plaza)",
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
