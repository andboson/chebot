package repositories

import (
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"sync"
	"time"
)

const (
	DEFAULT_CONTEXT_LIFETIME = 10
	CONTEXT_KINO             = "kino"
)

var KinoNamesRu = map[string]string{
	"lyubava": "Любаве",
	"plaza":   "Днипроплазе",
}

var TeleMsgIn map[int]chan tgbotapi.Update
var TeleUserMessagesFuncStop map[int]chan bool
var TeleBot *tgbotapi.BotAPI

func init() {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	TeleMsgIn = make(map[int]chan tgbotapi.Update)
	TeleUserMessagesFuncStop = make(map[int]chan bool)
}

// return channel for user messages
func GetOrNewUserChannel(update tgbotapi.Update) chan tgbotapi.Update {
	var id int
	if update.Message == nil {
		//log.Fatalf("[-]   %+v", update.CallbackQuery.From.ID)
		id = update.CallbackQuery.From.ID
	} else {
		id = update.Message.From.ID
	}
	_, ok := TeleMsgIn[id]
	if !ok {
		TeleMsgIn[id] = make(chan tgbotapi.Update)
	}
	CheckOrCreateUserMessagesFunc(id)

	return TeleMsgIn[id]
}

// check is user messages func started
// start if not or brake and start if not sure
func CheckOrCreateUserMessagesFunc(id int) {
	_, ok := TeleUserMessagesFuncStop[id]
	if !ok {
		TeleUserMessagesFuncStop[id] = make(chan bool)
	} else {
		return
	}

	go UserMessagesFunc(id)
}

// routine for current user messages
func UserMessagesFunc(id int) {
	var context string
	defer delete(TeleUserMessagesFuncStop, id)
	forgetContext := time.After(DEFAULT_CONTEXT_LIFETIME * time.Second)

	for {
		select {
		case update := <-TeleMsgIn[id]:
			context = processUserMessage(update, context)
			forgetContext = time.After(DEFAULT_CONTEXT_LIFETIME * time.Second)
		case stop := <-TeleUserMessagesFuncStop[id]:
			if stop {
				context = ""
				return
			}
		case <-forgetContext:
			context = ""
			return
		}
	}
}

func GoodbyeMsg(chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "До скорого!")
	TeleBot.Send(msg)
}

func processUserMessage(update tgbotapi.Update, ctx string) string {
	var chatId int64
	var reply, text string
	var messageId int

	if update.Message == nil {
		chatId = update.CallbackQuery.Message.Chat.ID
		messageId = update.CallbackQuery.Message.MessageID
		text = update.CallbackQuery.Data
	} else {
		chatId = update.Message.Chat.ID
		messageId = update.Message.MessageID
		text = update.Message.Text
	}

	//log.Printf("--text:  %s", text, messageId, ctx)

	// process context
	if text != "" && ctx != "" {
		switch ctx {
		case CONTEXT_KINO:
			processKinoRequest(text, ctx, chatId, messageId)

		}
		return ctx
	}

	msg := tgbotapi.NewMessage(chatId, "")

	// catch /kino command
	if (text == "/kino" || text == "/films" ) && ctx != CONTEXT_KINO {
		reply = "Выберите кинотеатр"

		butt := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Lyubava", "lyubava"),
			tgbotapi.NewInlineKeyboardButtonData("Dniproplaza", "plaza"),
		)
		keyb := tgbotapi.NewInlineKeyboardMarkup(butt)
		msg.ReplyMarkup = &keyb
		ctx = CONTEXT_KINO
	}

	// new user
	if update.Message != nil && update.Message.NewChatMember != nil && update.Message.NewChatMember.UserName != "" {
		// В чат вошел новый пользователь
		reply = fmt.Sprintf(`Добро пожаловать @%s!.`,
			update.Message.NewChatMember.UserName)
		msg.ReplyMarkup = nil
	}

	if reply != "" {
		msg.Text = reply
		TeleBot.Send(msg)
	}

	return ctx
}

func processKinoRequest(location string, ctx string, chatId int64, messageId int) {
	name, ok := KinoNamesRu[location]
	if !ok {
		log.Printf("unknown loction: %s", location)
		return
	}

	msg := tgbotapi.NewMessage(chatId, "Фильмы в кинотеатре в "+name+":")
	TeleBot.Send(msg)

	films := GetMovies(location)
	msg = tgbotapi.NewMessage(chatId, "")
	msgPhoto := tgbotapi.NewPhotoShare(chatId, "")
	for _, film := range films {
		msg.Text = fmt.Sprintf("*%s*",film.Title) //URL_PREFIX + "/" + film.Link)
		msg.ParseMode = tgbotapi.ModeMarkdown
		TeleBot.Send(msg)

		msgPhoto = tgbotapi.NewPhotoShare(chatId, URL_PREFIX+film.Img)
		msgPhoto.Caption = film.TimeBlock
		TeleBot.Send(msgPhoto)
	}

	return
}
