package tests

import (
	"testing"
	"github.com/andboson/chebot/repositories"
	"github.com/andboson/chebot/models"
	"os/exec"
	"fmt"
	"github.com/labstack/gommon/log"
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


func TestSpeechRhVoice(t *testing.T) {
	name := repositories.GetMD5Hash("tes")
	log.Printf("----- %s", name)

	text := "привет всем, friends"
	cmd := fmt.Sprintf("echo \"%s\" | RHVoice-test -p Irina", text)
	out, err := exec.Command("bash","-c",cmd).Output()
	if err != nil {
		t.Errorf("Failed to execute command: %s", cmd)
	}
	log.Printf(">>> %s", out)
	if err != nil {
		t.Fatalf("unable to get speech file", err)
	}
}
