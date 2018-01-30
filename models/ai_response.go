package models

type AiResponse struct {
	Speech      string      `json:"speech,omitempty"`
	DisplayText string      `json:"displayText,omitempty"`
	Data        interface{} `json:"data,omitempty"`
	Source      string      `json:"source,omitempty"`
	ContextOut  []Context   `json:"contextOut,omitempty"`
	Messages    []Message   `json:"messages,omitempty"`
}

type Data struct {
	Google struct {
		ExpectUserResponse bool         `json:"expectUserResponse"`
		RichResponse       RichResponse `json:"richResponse"`
		SystemIntent       SystemIntent `json:"systemIntent"`
		ExpectedInputs     []ExpectedInput `json:"expectedInputs,omitempty"`
	} `json:"google"`
}

type RichResponse struct {
	Items             []map[string]interface{} `json:"items"`
	Suggestions       []Suggestion             `json:"suggestions"`
	LinkOutSuggestion []interface{}            `json:"linkOutSuggestion"`
}

type ExpectedInput struct {
	PossibleIntents []PossibleIntent `json:"possibleIntents"`
}

type PossibleIntent struct {
	Intent 		   string       `json:"intent"`
	InputValueData InputValueData       `json:"inputValueData,omitempty"`
}

type InputValueData struct {
	Type           string         `json:"@type"`
}

type SuggestionResponse struct {
	Suggestion
}

type Item struct {
	SimpleResponse `json:"simple_response,omitempty"`
	BasicCard      `json:"basic_card,omitempty"`
}

type SimpleResponse struct {
	TextToSpeech string `json:"textToSpeech,omitempty"`
	DisplayText  string `json:"displayText"`
	Ssml         string `json:"ssml,omitempty"`
}

type BasicCard struct {
	Title         string `json:"title"`
	Subtitle      string `json:"subtitle"`
	FormattedText string `json:"formattedText"`
	Image         Image  `json:"image"`
}

type Image struct {
	URL               string `json:"url"`
	AccessibilityText string `json:"accessibilityText"`
}

////// coursel
type SystemIntent struct {
	Intent string       `json:"intent"` //actions.intent.OPTION
	Data   CouruselData `json:"data,omitempty"`
}

type CouruselData struct {
	Type           string         `json:"@type"` //type.googleapis.com/google.actions.v2.OptionValueSpec
	CarouselSelect CarouselSelect `json:"carouselSelect"`
}

type CarouselSelect struct {
	Items []CouruselItems `json:"items"`
}

type CouruselItems struct {
	OptionInfo  OptionInfo `json:"optionInfo"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Image       Image      `json:"image"`
}

type OptionInfo struct {
	Key      string   `json:"key"`
	Synonyms []string `json:"synonyms"`
}

//{
//"basicCard": models.BasicCard{
//Title:         "Card title",
//Subtitle:      "subtitle",
//FormattedText: "text text \n text",
//Image: models.Image{
//AccessibilityText: "32121",
//URL:               "https://www.gstatic.com/mobilesdk/170329_assistant/assistant_color_96dp.png",
//}},
//},
