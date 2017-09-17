package repositories

import (
	"fmt"
	"github.com/andboson/skypeapi"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"os"
	"strings"
)

const TAXI_LIST_FILE = "taxi_numbers.txt"

func ProcessSkypeTaxiManage(message skypeapi.Activity) bool {
	var text string
	var err error
	text = message.Text

	log.Printf("[skype] text: %s :", text)
	text = strings.Replace(text, "CherkassyBot", "", -1)
	text = strings.TrimSpace(text)

	if strings.Contains(text, "taxi add") {
		err = AddTaxiToList(&message, text)
		if err != nil {
			log.Printf("[skype] taxi add err messaging %s", err)
		}

		return true
	}

	if strings.Contains(text, "beer") || strings.Contains(text, "Shvets") {
		skypeapi.SendReplyMessage(&message, "(beer)", SkypeToken.AccessToken)

		return true
	}

	if strings.Contains(text, "taxi clear") {
		err = ClearTaxi(&message, text)
		if err != nil {
			log.Printf("[skype] taxi err clearing %s", err)
		}

		return true
	}

	return false
}

func ClearTaxi(activity *skypeapi.Activity, text string) error {
	err := os.Remove(TAXI_LIST_FILE)
	if err == nil {
		skypeapi.SendReplyMessage(activity, "Done!", SkypeToken.AccessToken)
	}

	return err
}

func AddTaxiToList(activity *skypeapi.Activity, text string) error {
	var err error
	rawTaxi := strings.Trim(text, "taxi add")
	taxiArr := strings.Split(rawTaxi, "=")
	if len(taxiArr) == 2 {
		err = AddTaxi(strings.TrimSpace(taxiArr[0]), strings.TrimSpace(taxiArr[1]))
		if err == nil {
			err = SendTaxiList(activity)
		}
	}

	return err
}

func SendTaxiList(activity *skypeapi.Activity) error {
	taxiList := LoadTaxi()
	var attchmts []skypeapi.Attachment
	var err error

	var btns []skypeapi.CardAction
	for number, firm := range taxiList {
		btn := skypeapi.CardAction{
			Title: firm + " - " + number,
			Type:  "imBack",
			Value: number,
		}

		btns = append(btns, btn)
	}

	var att = skypeapi.Attachment{
		ContentType: "application/vnd.microsoft.card.hero",
		Content: skypeapi.AttachmentContent{
			Title:   "Номера такси " + fmt.Sprintf("(%d)", len(taxiList)),
			Buttons: btns,
		},
	}
	attchmts = append(attchmts, att)
	responseActivity := &skypeapi.Activity{
		Type:             activity.Type,
		AttachmentLayout: "carousel",
		From:             activity.Recipient,
		Conversation:     activity.Conversation,
		Recipient:        activity.From,
		InputHint:        "select number",
		Text:             "Номера такси " + fmt.Sprintf("(%d)", len(taxiList)),
		Attachments:      attchmts,
		ReplyToID:        activity.ID,
	}
	replyUrl := fmt.Sprintf("%v/v3/conversations/%v/activities/%v", activity.ServiceURL, activity.Conversation.ID, activity.ID)
	err = skypeapi.SendActivityRequest(responseActivity, replyUrl, SkypeToken.AccessToken)

	return err
}

func LoadTaxi() map[string]string {
	var taxiList = make(map[string]string)
	content, err := ioutil.ReadFile(TAXI_LIST_FILE)
	if err != nil {
		log.Printf("[taxi] err loading file %s", err)
	}

	contentString := strings.Trim(string(content), "\n")
	contentString = strings.Trim(contentString, "\r")
	lines := strings.Split(contentString, "\r\n")

	for _, line := range lines {
		var taxiArr = strings.Split(line, "|")
		if len(taxiArr) == 2 {
			taxiList[taxiArr[0]] = taxiArr[1]
		}
	}

	return taxiList
}

func AddTaxi(number, firm string) error {
	var err error
	file, err := os.OpenFile(TAXI_LIST_FILE, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0775)
	if err != nil {
		log.Printf("[taxi] err open file %s", err)
		return err
	}
	defer file.Close()
	line := fmt.Sprintf("%s|%s\r\n", number, firm)
	_, err = file.WriteString(line)

	return err
}
