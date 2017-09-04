package tests

import (
	"github.com/andboson/chebot/repositories"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestAddTaxi(t *testing.T) {
	os.Remove(repositories.TAXI_LIST_FILE)
	err := repositories.AddTaxi("number", "firm")
	assert.Nil(t, err)

	repositories.AddTaxi("number2", "firm2")

	taxis := repositories.LoadTaxi()
	assert.Equal(t, 2, len(taxis))

	log.Printf(">>>>> %#v", taxis)

}
