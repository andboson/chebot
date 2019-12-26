package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/andboson/chebot/repositories"

	"github.com/andboson/skypeapi"
	"github.com/labstack/echo"
)

const skypeServiceUrl = "https://smba.trafficmanager.net/apis"

type webMessageRequest struct {
	Text           string `json:"text" validate:"required"`
	ConversationID string `json:"conversation_id" validate:"required"`
}

func WebMessage(c echo.Context) error {
	var json webMessageRequest
	if err := c.Bind(&json); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	if strings.TrimSpace(json.Text) == "" || strings.TrimSpace(json.ConversationID) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "fill all required fields"})
	}

	activity := skypeapi.Activity{
		Type: "message",
		Text: json.Text,
	}

	replyUrl := fmt.Sprintf("%s/v3/conversations/%s/activities", skypeServiceUrl, json.ConversationID)

	if err := skypeapi.SendActivityRequest(&activity, replyUrl, repositories.SkypeToken.AccessToken); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
