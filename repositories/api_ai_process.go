package repositories

import (
	"github.com/andboson/chebot/models"
	"github.com/andrewstuart/goq"
	"log"
	"net/http"
)

const (
	LUBAVA_URL = "http://cherkassy.multiplex.ua/Poster.aspx?id=16"
	PLAZA_URL  = "http://cherkassy.multiplex.ua/Poster.aspx?id=10"
	URL_PREFIX = "https://cherkassy.multiplex.ua"
)

var KinoUrls = map[string]string{
	"lyubava": LUBAVA_URL,
	"plaza":   PLAZA_URL,
}

var KinoNames = map[string]string{
	"lyubava": "Lyubava ",
	"plaza":   "Dniproplaza ",
}


type Films struct {
	FilmTds []Film `goquery:"td.afisha_td_bottom"`
}

type Film struct {
	Img       string `goquery:"a img,[src]"`
	Link      string `goquery:"a,[href]"`
	Title     string `goquery:".afisha_film"`
	TimeBlock string `goquery:".afisha_time_block"`
}

func GetMovies(cinema string) []Film {
	url, ok := KinoUrls[cinema]
	if !ok {
		log.Fatal("Cinema not found!")
	}

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var ex Films

	err = goq.NewDecoder(res.Body).Decode(&ex)
	if err != nil {
		log.Fatal(err)
	}

	return ex.FilmTds
}

func GetMovieListResponse(films []Film, cinema string) models.Data {
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
	for _, film := range films {
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

	data.Google.SystemIntent = models.SystemIntent{
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
