package repositories

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"gopkg.in/telegram-bot-api.v4"
)

func processTaxiRequest(ctx string, chatId int64, messageId, uid int) int {
	var sentId int
	taxiList := LoadTaxi()

	msg := tgbotapi.NewMessage(chatId, "")
	text := fmt.Sprintf("Список такси: (%d)", len(taxiList))
	for number, firm := range taxiList {
		text += fmt.Sprintf("\r\n %s - %s", number, firm)
	}

	msg.Text = text
	msg.ParseMode = tgbotapi.ModeMarkdown
	msgSent, _ := TeleBot.Send(msg)
	sentId = msgSent.MessageID

	return sentId
}
