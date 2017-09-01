package repositories

import (
	"fmt"
	"strings"
	"strconv"
	"gopkg.in/telegram-bot-api.v4"
	"log"
)

func processKinoRequest(location string, ctx string, chatId int64, messageId, uid int) int {
	var num int
	var sentId int

	dataParts := strings.Split(location, "|")
	if len(dataParts) > 1 {
		location = dataParts[0]
		num, _ = strconv.Atoi(dataParts[1])
	}

	name, ok := KinoNamesRu[location]
	if !ok {
		log.Printf("unknown location: %s", location)
		return sentId
	}

	if len(dataParts) == 1 {
		TeleLastMsgID[uid] = 0
		msgWhere := tgbotapi.NewMessage(chatId, "Фильмы в кинотеатре в "+name+":")
		TeleBot.Send(msgWhere)
	}

	films := GetMovies(location)

	butt := tgbotapi.NewInlineKeyboardRow()
	for idx, _ := range films {
		var data = fmt.Sprintf("%s|%d", location, idx)
		var selected string
		if idx == num {
			selected = " *"
		}
		var text = strconv.Itoa(idx) + selected
		btn := tgbotapi.NewInlineKeyboardButtonData(text, data)
		butt = append(butt, btn)
	}

	// send film
	keyb := tgbotapi.NewInlineKeyboardMarkup(butt)
	film := films[num]

	lastMsg, ok := TeleLastMsgID[uid]
	if ok && lastMsg != 0 {
		msgEdit := tgbotapi.NewEditMessageText(chatId, lastMsg,"")
		msgEdit.Text = fmt.Sprintf("*%s*\n %s [:](%s)", film.Title, film.TimeBlock, URL_PREFIX+"/"+film.Img)
		msgEdit.ParseMode = tgbotapi.ModeMarkdown
		msgEdit.ReplyMarkup = &keyb
		TeleBot.Send(msgEdit)
		sentId = TeleLastMsgID[uid]

	} else {
		msg := tgbotapi.NewMessage(chatId, "")

		msg.Text = fmt.Sprintf("*%s*\n %s [:](%s)", film.Title, film.TimeBlock, URL_PREFIX+"/"+film.Img)
		msg.ParseMode = tgbotapi.ModeMarkdown
		msg.ReplyMarkup = &keyb

		msgSent, _ := TeleBot.Send(msg)
		sentId = msgSent.MessageID
	}

	return sentId
}

