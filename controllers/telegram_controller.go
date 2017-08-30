package controllers

import (
	"github.com/andboson/chebot/models"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"github.com/andboson/chebot/repositories"
)

func TelegramMessagesHandler() {
	defer recover()
	bot, err := tgbotapi.NewBotAPI(models.Conf.TelegramToken)

	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false
	bot.RemoveWebhook()
	repositories.TeleBot = bot

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Panic(err)
	}
	// В канал updates будут приходить все новые сообщения.
	for update := range updates {

		userChan := repositories.GetOrNewUserChannel(update)

		go func() {
			userChan <- update
		}()
	}
}
