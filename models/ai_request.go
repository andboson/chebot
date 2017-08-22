package models

import "time"

type AiRequest struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Lang      string    `json:"lang"`
	Result    Result    `json:"result"`
	Status    Status    `json:"status"`
	SessionID string    `json:"sessionId"`
}

type Status struct {
	Code      int    `json:"code"`
	ErrorType string `json:"errorType"`
}

type Result struct {
	Source           string      `json:"source"`
	ResolvedQuery    string      `json:"resolvedQuery"`
	Speech           string      `json:"speech"`
	Action           string      `json:"action"`
	ActionIncomplete bool        `json:"actionIncomplete"`
	Parameters       Parameters  `json:"parameters"`
	Contexts         []Context   `json:"contexts"`
	Metadata         Metadata    `json:"metadata"`
	Fulfillment      Fulfillment `json:"fulfillment"`
	//	Score            int         `json:"score"`
}

type Context struct {
	Name       string     `json:"name"`
	Parameters Parameters `json:"parameters"`
	//	Lifespan   int        `json:"lifespan"`
}

type Parameters struct {
	CinemaOriginal string `json:"cinema.original"`
	Cinema         string `json:"cinema"`
}

type Metadata struct {
	IntentID                  string `json:"intentId"`
	WebhookUsed               string `json:"webhookUsed"`
	WebhookForSlotFillingUsed string `json:"webhookForSlotFillingUsed"`
	IntentName                string `json:"intentName"`
}

type Fulfillment struct {
	Speech   string    `json:"speech"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Type              interface{}  `json:"type"` // simple_response, basic_card
	Platform          string       `json:"platform,omitempty"`
	TextToSpeech      string       `json:"textToSpeech,omitempty"`
	Speech            string       `json:"speech,omitempty"`
	ImageUrl          string       `json:"imageUrl,omitempty"`
	Suggestions       []Suggestion `json:"suggestions"`
	URL               string       `json:"url,omitempty"`
	AccessibilityText string       `json:"accessibility_text,omitempty"`
}

type Suggestion struct {
	Title string `json:"title"`
}
