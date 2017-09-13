package repositories

import (
	"fmt"
	"github.com/andboson/chebot/models"
	"github.com/labstack/gommon/log"
	"github.com/maciekmm/messenger-platform-go-sdk"
	"github.com/maciekmm/messenger-platform-go-sdk/template"
	"strings"
	"time"
)

type FbProcesssor struct {
	Messenger *messenger.Messenger
	Opts      messenger.MessageOpts
	Payload   messenger.Postback
	Msg       messenger.ReceivedMessage
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
	name, ok := KinoNamesRu[location]
	url, _ := KinoUrls[location]
	if !ok {
		s.Messenger.SendSimpleMessage(s.Opts.Sender.ID, "Не знаю такое место")
		log.Printf("[fb] unknown location: %s", location)
		return
	}
	films := GetMovies(location)
	name = fmt.Sprintf("[%s](%s)", name, url)
	mq := messenger.MessageQuery{}
	for idx, film := range films {
		filmTpl := template.GenericTemplate{
			Title:    film.Title,
			Subtitle: film.TimeBlock,
			ItemURL:  URL_PREFIX + film.Link,
			ImageURL: URL_PREFIX + film.Img,
		}
		mq.Template(filmTpl)
	    if idx == 8 {
break
}
	}

	mq.Text("Фильмы в " + name)
	mq.RecipientID(s.Opts.Sender.ID)
	_, err := s.Messenger.SendMessage(mq)
	if err != nil {
		log.Printf("[fb] error messaging films: %s", err)
	}
}

func (s FbProcesssor) ShowTaxiList() {
	taxiList := LoadTaxi()
	text := fmt.Sprintf("Список такси: (%d)", len(taxiList))
	for number, firm := range taxiList {
		text += " \r\n " + number + " - " + firm
	}

	_, err := s.Messenger.SendSimpleMessage(s.Opts.Sender.ID, text)
	if err != nil {
		log.Printf("[fb] error messaging: %s", err)
	}
}

func (s FbProcesssor) GetText() string {
	var text string
	text = s.Msg.Text

	if s.Payload.Payload != "" {
		text = s.Payload.Payload
	}

	return text
}

func (s FbProcesssor) GetUid() string {

	return s.Opts.Sender.ID
}
