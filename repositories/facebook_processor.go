package repositories

import (
	"github.com/maciekmm/messenger-platform-go-sdk"
	"time"
	"strings"
	"github.com/andboson/chebot/models"
	"github.com/labstack/gommon/log"
	"github.com/maciekmm/messenger-platform-go-sdk/template"
)

type FbProcesssor struct {
	Messenger *messenger.Messenger
	Opts messenger.MessageOpts
	Payload messenger.Postback
	Msg messenger.ReceivedMessage
}


func (s FbProcesssor) ShowHelp() {
	helpText := "Доступные команды:  \r\n # " + strings.Join(models.CmdList, "\r\n # ")
	_, err := s.Messenger.SendSimpleMessage(s.Opts.Sender.ID, helpText)
	if err != nil {
		log.Printf("[fb] error messaging: %s", err)
	}

}

func (s FbProcesssor) ShowKinoPlaces() {
	s.Messenger.SendAction(messenger.Recipient{ID: s.Opts.Sender.ID}, messenger.SenderActionTypingOn)
	time.Sleep(100 * time.Millisecond)
	s.Messenger.SendAction(messenger.Recipient{ID: s.Opts.Sender.ID}, messenger.SenderActionTypingOff)
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
	mq.Text("Выберите кинотеатр")
	mq.RecipientID(s.Opts.Sender.ID)
	_, err := s.Messenger.SendMessage(mq)

	if err != nil {
		log.Printf("[fb] error messaging: %s", err)
	}
}

func (s FbProcesssor) ShowFilms(location string) {
	//sendFilmsReplyMessage(s.Message, location)
	tpl1 := template.GenericTemplate{
		Title: "Выберите кинотеатр",
		Subtitle: "2232323",
		ItemURL: "https://cherkassy.multiplex.ua",
		ImageURL: "https://cherkassy.multiplex.ua/Images/Upload/origin.%D0%92%D0%B0%D0%BB%D0%B5%D1%80%D1%96%D0%B0%D0%BD%20%D1%82%D0%B0%20%D0%BC%D1%96%D1%81%D1%82%D0%BE%20%D1%82%D0%B8%D1%81%D1%8F%D1%87%D1%96%20%D0%BF%D0%BB%D0%B0%D0%BD%D0%B5%D1%82%203%D0%94.jpg",
	}
	tpl2 := template.GenericTemplate{
		Title: "Выберите 222",
		Subtitle: "2232323 333",
		ItemURL: "https://cherkassy.multiplex.ua",
		ImageURL: "https://cherkassy.multiplex.ua/Images/Upload/origin.%D0%92%D0%B0%D0%BB%D0%B5%D1%80%D1%96%D0%B0%D0%BD%20%D1%82%D0%B0%20%D0%BC%D1%96%D1%81%D1%82%D0%BE%20%D1%82%D0%B8%D1%81%D1%8F%D1%87%D1%96%20%D0%BF%D0%BB%D0%B0%D0%BD%D0%B5%D1%82%203%D0%94.jpg",
	}

	mq := messenger.MessageQuery{}
	mq.Template(tpl1)
	mq.Template(tpl2)
	mq.Text("jndtn")
	mq.RecipientID(s.Opts.Sender.ID)
	resp, err2 := s.Messenger.SendMessage(mq)
	log.Printf("[fb postback] %#v", resp, err2)
}

func (s FbProcesssor) ShowTaxiList() {

}

func (s FbProcesssor) GetText() string {
	var text string
	text =  s.Msg.Text

	if s.Payload.Payload != "" {
		text = s.Payload.Payload
	}

	return text
}

func (s FbProcesssor) GetUid() string {

	return s.Opts.Sender.ID
}


