package tests

import (
	"testing"
	"github.com/andboson/chebot/repositories"
	"github.com/andboson/chebot/models"
)

func init(){
	models.InitConfig()
}


func TestSpeechFileCreate(t *testing.T) {
	text := "привет всем, friends"
	content, err := repositories.GetSpeechFileContent(text)
	if err != nil {
		t.Fatalf("unable to get speech file", err)
	}

	err = repositories.StoreFile(text, content)
	if err != nil {
		t.Fatalf("unable to store file", err)
	}

}
