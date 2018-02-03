package controllers

import (
	"github.com/andboson/chebot/models"
	"github.com/andboson/chebot/repositories"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"net/http"
)

func GetAiResponse(c echo.Context) error {
	var request models.AiRequest
	err := c.Bind(&request)
	if err != nil {
		log.Printf("[---] request error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "request decoding error",
			"error":   err.Error(),
		})
	}

	token := c.Request().Header.Get("x-secret-header")
	if token != models.Conf.IncomingGoogleToken {
		log.Printf("[---] token error", token)
		return c.JSON(http.StatusForbidden, map[string]string{
			"message": "forbidden",
		})
	}

	var resp = new(models.AiResponse)
	resp.Speech = "films"
	resp.Source = "bot"

	var data models.Data
	var context = ""
	var isVoice = false
	// get context form contexts
	for _, ctx := range request.Result.Contexts {
		if ctx.Name == "google_assistant_input_type_voice" {
			isVoice = true
		}
		if _, ok := repositories.AvailContexts[ctx.Name]; ok {
			context = ctx.Name
		}
	}
	// get context from query
	if _, ok := repositories.AvailContexts[request.Result.ResolvedQuery]; ok {
		context = request.Result.ResolvedQuery
	}

	log.Printf("--- %+v  --context: %s == %+v", request, context, data)

	switch context {
	case repositories.CONTEXT_KINO:
		log.Printf("--- 0", context)

		films := repositories.GetMovies(request.Result.Parameters.Cinema, false)
		data = repositories.GetMovieListResponse(films, request.Result.Parameters.Cinema, isVoice)
	case repositories.CONTEXT_TAXI:
		log.Printf("--- 1", context)

		data = repositories.GetTaxiResponse()
	}

	resp.Data = data

	return c.JSON(http.StatusOK, resp)
}
