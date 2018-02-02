package repositories

import (
	"net/url"
	"net/http"
	"github.com/andboson/chebot/models"
	"io/ioutil"
	"os"
	"time"
	"fmt"
	"os/exec"
	"crypto/md5"
	"log"
	"encoding/hex"
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
	if _, errF := os.Stat(models.Conf.VoiceMp3sFolder +"/" + text + ".mp3"); errF == nil {
		return nil
	}

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

func UseRHVoice(text string) string {
	//go
	name := GetMD5Hash(text)
	cmd := fmt.Sprintf(" export $(dbus-launch | xargs) && echo \"%s\" | RHVoice-client -s Irina -v 1  -r 0.1  > %s/%s.wav", text, models.Conf.VoiceMp3sFolder, name)
	_, err := exec.Command("bash","-c",cmd).Output()
	if err != nil {
		log.Printf("Failed to execute command: %s err: %s", cmd, err)
	}

	cmd = fmt.Sprintf("ffmpeg -i %s/%s.wav %s/%s.mp3", models.Conf.VoiceMp3sFolder, name,models.Conf.VoiceMp3sFolder, name)
	_, err = exec.Command("bash","-c",cmd).Output()
	if err != nil {
		log.Printf("Failed to execute command encode: %s err: %s", cmd, err)
	}


	return string(name)
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}