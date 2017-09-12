package repositories

import (
	"github.com/andboson/skypeapi"
	"github.com/labstack/gommon/log"
	"time"
	"github.com/andboson/chebot/models"
	"strings"
)

var SkypeToken skypeapi.TokenResponse

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

type SkypeProcessor struct {
	Message *skypeapi.Activity
}

func (s SkypeProcessor) ShowHelp()  {
	helpText := "Доступные команды:  \r\n " + strings.Join(models.CmdList, "\r\n")
	err := skypeapi.SendReplyMessage(s.Message, helpText, SkypeToken.AccessToken)
	if err != nil {
		log.Printf("[skype] error messaging: %s", err)
	}

}

func (s SkypeProcessor) ShowKinoPlaces()  {
	err := sendChoicePlaceReplyMessage(s.Message, "Выберите  кинотеатр", SkypeToken.AccessToken)
	if err != nil {
		log.Printf("[skype] error messaging: %s", err)
	}
}

func (s SkypeProcessor) ShowFilms(location string)  {
	sendFilmsReplyMessage(s.Message, location)

}

func (s SkypeProcessor) ShowTaxiList()  {
	err := SendTaxiList(s.Message)
	if err != nil {
		log.Printf("[skype] taxi err messaging %s", err)
	}
}

func (s SkypeProcessor) GetText() string {

	return s.Message.Text
}

func (s SkypeProcessor) GetUid() string {

	return s.Message.From.ID
}