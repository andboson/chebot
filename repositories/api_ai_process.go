package repositories

import (
	"github.com/andboson/chebot/models"
	"github.com/essentialkaos/translit"
	"strconv"
)



func GetMovieListResponse(films []Film, cinema string, isVoice bool) models.Data {
	var empty string
	if len(films) == 0 {
		empty = " is empty"
	}

	name := KinoNames[cinema]
	data := models.Data{}
	data.Google.ExpectUserResponse = false
	data.Google.RichResponse = models.RichResponse{
		Items: []map[string]interface{}{
			{
				"simpleResponse": models.SimpleResponse{
					DisplayText:  name + "films list" + empty,
					TextToSpeech: name + "films list" + empty,
				},
			},
		},
		Suggestions: []models.Suggestion{
			{
				Title: "Lubava",
			},
			{
				Title: "Dniproplaza",
			},
		},
	}

	var items []models.CouruselItems
	var speechFilms = "<speak>"

	for idx, film := range films {
		//SpeechAndStore(film.Title)
		name := UseRHVoice(film.Title)
		speech := models.Conf.Mp3HttpPath + film.Title + ".mp3"
		speech = models.Conf.Mp3HttpPath + name + ".mp3"
		speechFilms += "<p><s>film"+ strconv.Itoa(idx) +": " +
		"<audio src=\""+ speech + "\">" +
			translit.EncodeToISO9B(film.Title) + "</audio></s></p><break time=\"400ms\"/>"

		var item = models.CouruselItems{
			Title:       film.Title,
			Description: film.TimeBlock,
			Image: models.Image{
				URL_PREFIX + film.Img,
				film.Title,
			},
			OptionInfo: models.OptionInfo{
				Key:      film.Title,
				Synonyms: []string{},
			},
		}
		items = append(items, item)
	}


	if isVoice {
		simpleTitle := map[string]interface{}{
			"simpleResponse": models.SimpleResponse{
				DisplayText:  " ",
				Ssml: speechFilms + "</speak>",
			},
		}
		data.Google.RichResponse.Items = append(data.Google.RichResponse.Items, simpleTitle)
	}

	data.Google.SystemIntent = &models.SystemIntent{
		Intent: "actions.intent.OPTION",
		Data: models.CouruselData{
			Type: "type.googleapis.com/google.actions.v2.OptionValueSpec",
			CarouselSelect: models.CarouselSelect{
				Items: items,
			},
		},
	}



	return data
}

func GetTaxiResponse() models.Data {
	data := models.Data{}
	data.Google.ExpectUserResponse = false
	data.Google.RichResponse = models.RichResponse{
		Items: []map[string]interface{}{
			{
				"simpleResponse": models.SimpleResponse{
					DisplayText:  "",
					TextToSpeech: "taxi list",
				},
			},
		},
	}

	return  data
}
