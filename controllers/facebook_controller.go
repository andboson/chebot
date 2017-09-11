package controllers

import (
	"github.com/andboson/chebot/models"
	"github.com/maciekmm/messenger-platform-go-sdk"
	"fmt"
	"github.com/labstack/gommon/log"
)

var FbMess *messenger.Messenger

type FaceBookCheck struct {
	HubMode      string `json:"hub.mode"`
	HubChallenge string `json:"hub.challenge"`
	HubToken     string `json:"hub.verify_token"`
}

func InitFb() {
	FbMess = &messenger.Messenger{
		VerifyToken: models.Conf.FbVerifyToken,
		AppSecret:   models.Conf.FbAppSecret,
		AccessToken: models.Conf.FbPageToken,
	}
	FbMess.MessageReceived = MessageReceived
}

func MessageReceived(event messenger.Event, opts messenger.MessageOpts, msg messenger.ReceivedMessage) {
	profile, err := FbMess.GetProfile(opts.Sender.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := FbMess.SendSimpleMessage(opts.Sender.ID, fmt.Sprintf("Hello, %s %s, %s", profile.FirstName, profile.LastName, msg.Text))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", resp)
	log.Printf("[fb] %#v", event)
}
