package repositories

import (
	"fmt"
	"github.com/andboson/skypeapi"
	"log"
)

func sendFilmsReplyMessage(activity *skypeapi.Activity, location string) {
	name, ok := KinoNamesRu[location]
	url, _ := KinoUrls[location]
	if !ok {
		skypeapi.SendReplyMessage(activity, "Не знаю такое место", SkypeToken.AccessToken)
		log.Printf("unknown location: %s", location)
		return
	}
	films := GetMovies(location, false)
	name = fmt.Sprintf("[%s](%s)", name, url)

	err := sendReplyMessageRich(activity, "Фильмы в "+name, SkypeToken.AccessToken, films)
	if err != nil {
		log.Printf("[skype] error films msg: %s", err)
	}

}

func sendChoicePlaceReplyMessage(activity *skypeapi.Activity, message, authorizationToken string) error {
	responseActivity := &skypeapi.Activity{
		Type:             activity.Type,
		From:             activity.Recipient,
		Conversation:     activity.Conversation,
		Recipient:        activity.From,
		Text:             message,
		InputHint:        "место (lyubava\\plaza)",
		AttachmentLayout: "carousel",
		Attachments: []skypeapi.Attachment{
			{
				ContentType: "application/vnd.microsoft.card.hero",
				Content: skypeapi.AttachmentContent{
					Title: "Любава",
					Text:  "нажмите, чтобы выбрать",
					Tap: &skypeapi.CardAction{
						Title: "Любава",
						Type:  "imBack",
						Value: "lyubava",
					},
				},
			},
			{
				ContentType: "application/vnd.microsoft.card.hero",
				Content: skypeapi.AttachmentContent{
					Title: "Днепроплаза",
					Text:  "нажмите, чтобы выбрать",
					Tap: &skypeapi.CardAction{
						Title: "Днепроплаза",
						Type:  "imBack",
						Value: "plaza",
					},
				},
			},
		},
		ReplyToID: activity.ID,
	}
	replyUrl := fmt.Sprintf("%vv3/conversations/%v/activities", activity.ServiceURL, activity.Conversation.ID)

	return skypeapi.SendActivityRequest(responseActivity, replyUrl, authorizationToken)
}

func sendReplyMessageRich(activity *skypeapi.Activity, message, authorizationToken string, films []Film) error {
	var attchmts []skypeapi.Attachment
	i := 0
	for _, film := range films {
		i++
		var att = skypeapi.Attachment{
			ContentType: "application/vnd.microsoft.card.hero",
			Content: skypeapi.AttachmentContent{
				Title: film.Title,
				Text:  film.TimeBlock,
				Tap: &skypeapi.CardAction{
					Type:  "openUrl",
					Value: URL_PREFIX + film.Link,
				},
				Images: []skypeapi.CardImage{
					{
						URL: URL_PREFIX + film.Img,
						Alt: film.Title,
					},
				},
			},
		}
		if i == 9 {
			break
		}

		attchmts = append(attchmts, att)
	}

	responseActivity := &skypeapi.Activity{
		Type:             activity.Type,
		AttachmentLayout: "carousel",
		From:             activity.Recipient,
		Conversation:     activity.Conversation,
		Recipient:        activity.From,
		Text:             message,
		Attachments:      attchmts,
		ReplyToID:        activity.ID,
	}
	replyUrl := fmt.Sprintf("%v/v3/conversations/%v/activities/%v", activity.ServiceURL, activity.Conversation.ID, activity.ID)
	//log.Printf("[skype] ---- %#v", films, replyUrl)

	return skypeapi.SendActivityRequest(responseActivity, replyUrl, authorizationToken)
}
