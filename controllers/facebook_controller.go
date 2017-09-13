package controllers

import (
	"github.com/andboson/chebot/models"
	"github.com/labstack/gommon/log"
	"github.com/maciekmm/messenger-platform-go-sdk"
	"net/http"
	"github.com/andboson/chebot/repositories"
)

var fbmess *messenger.Messenger

type FaceBookCheck struct {
	HubMode      string `json:"hub.mode"`
	HubChallenge string `json:"hub.challenge"`
	HubToken     string `json:"hub.verify_token"`
}

func InitFb() {
	fbmess = &messenger.Messenger{
		VerifyToken: models.Conf.FbVerifyToken,
		AppSecret:   models.Conf.FbAppSecret,
		AccessToken: models.Conf.FbPageToken,
	}
	fbmess.MessageReceived = MessageReceived
	fbmess.Postback = MessagePostback
	go func() {
		http.HandleFunc("/facebook.hook", fbmess.Handler)
		log.Fatal(http.ListenAndServe(":1324", nil))
	}()
}

func MessagePostback(event messenger.Event, opts messenger.MessageOpts, payload messenger.Postback) {
	var proc = repositories.FbProcesssor{}
	proc.Messenger = fbmess
	proc.Payload = payload
	proc.Opts = opts

	repositories.ProcessMessage(proc)
}

func MessageReceived(event messenger.Event, opts messenger.MessageOpts, msg messenger.ReceivedMessage) {
	var proc = repositories.FbProcesssor{}
	proc.Messenger = fbmess
	proc.Msg = msg
	proc.Opts = opts

	repositories.ProcessMessage(proc)
}
