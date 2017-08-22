package tests

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"bytes"
	"github.com/andboson/chebot/routes"
	"github.com/labstack/gommon/log"
	"github.com/labstack/echo"
)

func TestCreateUser(t *testing.T) {
	// Setup
	b := bytes.NewBufferString(`{}`)
	r, _ := http.NewRequest("POST", "/ai.get_movies", b)
	r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	w := httptest.NewRecorder()

	router := routes.Router()
	router.ServeHTTP(w, r)

	body := w.Body.String()
	log.Printf("----", body)


}
