package controllers

import (
	"fmt"
	"github.com/andboson/chebot/models"
	"github.com/labstack/gommon/log"
	"github.com/maciekmm/messenger-platform-go-sdk"
	"github.com/maciekmm/messenger-platform-go-sdk/template"
	"net/http"
	"time"
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
	FbMess.Postback = MessagePostback
	go func() {
		http.HandleFunc("/facebook.hook", FbMess.Handler)
		log.Fatal(http.ListenAndServe(":1324", nil))
	}()
}

func MessagePostback(event messenger.Event, opts messenger.MessageOpts, payload messenger.Postback) {
	log.Printf("====---==== %+v----%+v ---- %+v", opts, event, payload)

	tpl1 := template.GenericTemplate{
		Title: "Выберите кинотеатр",
		Subtitle: "2232323",
		ImageURL: "https://cherkassy.multiplex.ua/Images/Upload/origin.%D0%92%D0%B0%D0%BB%D0%B5%D1%80%D1%96%D0%B0%D0%BD%20%D1%82%D0%B0%20%D0%BC%D1%96%D1%81%D1%82%D0%BE%20%D1%82%D0%B8%D1%81%D1%8F%D1%87%D1%96%20%D0%BF%D0%BB%D0%B0%D0%BD%D0%B5%D1%82%203%D0%94.jpg",
	}
	tpl2 := template.GenericTemplate{
		Title: "Выберите 222",
		Subtitle: "2232323 333",
		ImageURL: "https://cherkassy.multiplex.ua/Images/Upload/origin.%D0%92%D0%B0%D0%BB%D0%B5%D1%80%D1%96%D0%B0%D0%BD%20%D1%82%D0%B0%20%D0%BC%D1%96%D1%81%D1%82%D0%BE%20%D1%82%D0%B8%D1%81%D1%8F%D1%87%D1%96%20%D0%BF%D0%BB%D0%B0%D0%BD%D0%B5%D1%82%203%D0%94.jpg",
	}

	mq := messenger.MessageQuery{}
	mq.Template(tpl1)
	mq.Template(tpl2)
	mq.Text("jndtn")
	mq.RecipientID(opts.Sender.ID)
	resp, err2 := FbMess.SendMessage(mq)
	log.Printf("[fb postback] %#v", resp, err2)


}

func MessageReceived(event messenger.Event, opts messenger.MessageOpts, msg messenger.ReceivedMessage) {
	profile, err := FbMess.GetProfile(opts.Sender.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Printf("======== %+v----%+v", opts, event)

	resp, err := FbMess.SendSimpleMessage(opts.Sender.ID, fmt.Sprintf("Hello, %s %s, %s", profile.FirstName, profile.LastName, msg.Text))
	if err != nil {
		fmt.Println(err)
	}

	FbMess.SendAction(messenger.Recipient{ID: opts.Sender.ID}, messenger.SenderActionTypingOn)
	time.Sleep(1 * time.Second)
	FbMess.SendAction(messenger.Recipient{ID: opts.Sender.ID}, messenger.SenderActionTypingOff)
	btns := template.GenericTemplate{
		Title: "Выберите кинотеатр",
		Buttons: []template.Button{
			{
				Title:   "Любава",
				Type:    template.ButtonTypePostback,
				Payload: "lyubava",
			},
			{
				Title:   "Днепроплаза",
				Type:    template.ButtonTypePostback,
				Payload: "plaza",
			},
		},
	}

	mq := messenger.MessageQuery{}
	mq.Template(btns)
	mq.Text("jndtn")
	mq.RecipientID(opts.Sender.ID)
	resp, err2 := FbMess.SendMessage(mq)
	fmt.Printf("%+v", resp)
	log.Printf("[fb] %#v", resp, err2)
}
