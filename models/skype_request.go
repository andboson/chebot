package models

import "time"

const MESSAGE_TYPE_MESSAGE = "message"
const TEXT_FORMAT_PLAIN = "plain"

type SkypeMessage struct {
	Type         string       `json:"type"`
	ID           *string       `json:"id"`
	Timestamp    time.Time    `json:"timestamp"`
	ServiceURL   string       `json:"serviceUrl"` //"https://webchat.botframework.com/"
	ChannelID    string       `json:"channelId"`  //"webchat"
	From         From         `json:"from"`
	Conversation Conversation `json:"conversation"`
	Recipient    Recipient    `json:"recipient"`
	TextFormat   string       `json:"textFormat"`
	Locale       string       `json:"locale"`
	Text         string       `json:"text"`
	ReplyToId    string       `json:"replyToId"`
	ChannelData  ChannelData  `json:"channelData"`
}

type From struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Conversation struct {
	ID string `json:"id"`
}

type Recipient struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ChannelData struct {
	ClientActivityID string `json:"clientActivityId"`
}
