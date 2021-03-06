package repositories

import (
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
)

func processTaxiRequest(ctx string, chatId int64) int {
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
