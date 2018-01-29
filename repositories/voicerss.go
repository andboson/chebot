package repositories

import (
	"net/url"
	"net/http"
	"github.com/andboson/chebot/models"
	"io/ioutil"
	"os"
	"time"
)

func init()  {
	models.InitConfig()
	ClearOldFiles()
}

func GetSpeechFileContent(text string) ([]byte, error){
	resp, err := http.Get("http://api.voicerss.org/?key=" + models.Conf.VoiceRssApiKey +
		"&hl=ru-ru&src=" + url.QueryEscape(text))
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

func StoreFile(text string, content []byte) error{
	err := ioutil.WriteFile( models.Conf.VoiceMp3sFolder +"/" + text + ".mp3", content, 0755)

	return err
}

func SpeechAndStore(text string) error {
	content, err := GetSpeechFileContent(text)
	if err != nil {
		return err
	}

	return StoreFile(text, content)
}

func ClearOldFiles()  {
	os.RemoveAll(models.Conf.VoiceMp3sFolder +"/*")

	time.AfterFunc( 1 * 24 * time.Hour, ClearOldFiles)
}