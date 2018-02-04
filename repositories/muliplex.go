package repositories

import (
	"path/filepath"
	"github.com/andrewstuart/goq"
	"regexp"
	"strings"
	"os"
	"encoding/gob"
	"log"
	"github.com/kardianos/osext"
	"net/http"
	"time"
)


const (
	LUBAVA_URL = "https://multiplex.ua/cinema/cherkasy/lyubava"
	PLAZA_URL  = "https://multiplex.ua/cinema/cherkasy/dniproplaza"
	URL_PREFIX = "https://multiplex.ua/"
	RECACHE_TIME_HOURS = 6
	films_cache_file = "films.cache"
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
	FilmTds []Film `goquery:"div.film"`
}

type Film struct {
	Img       string `goquery:"div.poster,[style]"`
	Link      string `goquery:"a,[href]"`
	Title     string `goquery:".info a"`
	TimeBlock string `goquery:".info .sessions,text"`
}

func GetMovies(cinema string, force bool) []Film {
	url, ok := KinoUrls[cinema]
	if !ok {
		log.Printf("Cinema not found!")
		return []Film{}
	}

	//need to check in cache
	if force != true {
		filmsC, errC := loadCacheFilmsFromFile(cinema)
		if len(filmsC) > 1 {
			return filmsC
		}
		if errC != nil {
			log.Printf("err decode films cache", cinema, errC)
		}
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

	re := regexp.MustCompile(`\s+|\n+`)
	for idx, flm := range ex.FilmTds {
		img := strings.TrimRight(flm.Img, "')")
		img = strings.TrimLeft(img, "background-image: url('")
		times := re.ReplaceAllLiteralString(flm.TimeBlock, " ")
		ex.FilmTds[idx].Img = img
		ex.FilmTds[idx].TimeBlock = times
	}

	//filter doubles
	var uniq = map[string]string{}
	var result []Film
	for _, flm := range ex.FilmTds  {
		if _, ok := uniq[flm.Title]; ok {
			continue
		} else {
			uniq[flm.Title] = flm.Title
		}
		result = append(result, flm)
	}

	// caching
	saveFilmsCacheToFile(result, cinema)
	time.AfterFunc( RECACHE_TIME_HOURS * time.Hour, func() {
		GetMovies(cinema, force)
	})

	return result
}

func loadCacheFilmsFromFile(cinema string) ([]Film, error) {
	var err error
	var result []Film

	curDir, _ := osext.ExecutableFolder()
	file, err := os.Open(filepath.Join(curDir, cinema + "." + films_cache_file))
	if err == nil {
		log.Printf("trying to load films cache from file...")
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(&result)
	} else {
		log.Printf("trying to load films cache from file error: %s", err)
	}
	file.Close()

	return result, err
}

func saveFilmsCacheToFile(collect []Film, cinema string) {
	curDir, _ := osext.ExecutableFolder()
	file, err := os.Create(filepath.Join(curDir, cinema + "." + films_cache_file))
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(collect)
	} else {
		log.Printf("create films cache file error: %s", err)
	}
	file.Close()
}
