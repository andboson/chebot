package tests

import (
	"testing"
	"github.com/andboson/chebot/repositories"
	"github.com/stretchr/testify/assert"
	"log"
)

func TestParse(t *testing.T) {
	films := repositories.GetMovies("lyubava")
	assert.True(t, len(films) > 0)

	//data := repositories.GetMovieListResponse(films, "lyubava")
	//log.Printf("--data~:  %+v", data)
	log.Printf("--data~: \r\n  %+v", films[0])
}

