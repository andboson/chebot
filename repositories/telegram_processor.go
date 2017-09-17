package repositories

import (
	"github.com/andboson/chebot/models"
	"github.com/labstack/gommon/log"
	"gopkg.in/telegram-bot-api.v4"
	"strconv"
	"strings"
	"sync"
)

const no_understand = "Не понимаю :\\"

var KinoNamesRu = map[string]string{
	"lyubava": "Любаве",
	"plaza":   "Днепроплазе",
}

var TeleLastMsgID map[int]int
var TeleBot *tgbotapi.BotAPI
var mu sync.Mutex

func init() {
	mu.Lock()
	defer mu.Unlock()
	TeleLastMsgID = make(map[int]int)
}

type TelegramProcessor struct {
	Update    tgbotapi.Update
	Uid       int
	chatId    int64
	text      string
	messageId int
}

func NewTelegramProcessor(update tgbotapi.Update) TelegramProcessor {
	var chatId int64
	var text string
	var messageId int
	var id int
	var s = TelegramProcessor{
		Update: update,
	}

	if update.CallbackQuery != nil {
		id = s.Update.CallbackQuery.From.ID
		chatId = s.Update.CallbackQuery.Message.Chat.ID
		messageId = s.Update.CallbackQuery.Message.MessageID
		text = s.Update.CallbackQuery.Data
	} else if s.Update.Message != nil {
		id = s.Update.Message.From.ID
		chatId = s.Update.Message.Chat.ID
		messageId = s.Update.Message.MessageID
		text = s.Update.Message.Text
	} else {
		log.Printf("[tlgrm] unable to init proc %#v", s.Update.Message)

		return s
	}

	s.chatId = chatId
	s.messageId = messageId
	s.text = text
	s.Uid = id

	return s
}


func (s TelegramProcessor) NoResults() {
	helpText := no_understand
	msg := tgbotapi.NewMessage(s.chatId, helpText)
	_, err := TeleBot.Send(msg)

	if err != nil {
		log.Printf("[tlgrm] error help messaging: %s", err)
	}

}

func (s TelegramProcessor) ShowHelp() {
	helpText := "Доступные команды:  \r\n # " + strings.Join(models.CmdList, "\r\n # ")
	msg := tgbotapi.NewMessage(s.chatId, helpText)
	_, err := TeleBot.Send(msg)

	if err != nil {
		log.Printf("[tlgrm] error help messaging: %s", err)
	}

}

func (s TelegramProcessor) ShowKinoPlaces() {
	msg := tgbotapi.NewMessage(s.chatId, "Выберите кинотеатр")
	butt := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Lyubava", "lyubava"),
		tgbotapi.NewInlineKeyboardButtonData("Dniproplaza", "plaza"),
	)
	keyb := tgbotapi.NewInlineKeyboardMarkup(butt)
	msg.ReplyMarkup = &keyb
	_, err := TeleBot.Send(msg)
	if err != nil {
		log.Printf("[tlgrm] error messaging: %s", err)
	}
	TeleLastMsgID[s.Uid] = 0
}

func (s TelegramProcessor) ShowFilms(location string) {
	msgId := processKinoRequest(s.text, s.chatId, s.Uid)
	TeleLastMsgID[s.Uid] = msgId
}

func (s TelegramProcessor) ShowTaxiList() {
	processTaxiRequest("", s.chatId)
}

func (s TelegramProcessor) GetText() string {

	return s.text
}

func (s TelegramProcessor) GetUid() string {

	return strconv.Itoa(s.Uid)
}
