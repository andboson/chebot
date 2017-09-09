package controllers

import (
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"net/http"
)

const token = "f1c9830e3258c325be00bc6f7cd5324a"

type FaceBookCheck struct {
	HubMode      string `json:"hub.mode"`
	HubChallenge string `json:"hub.challenge"`
	HubToken     string `json:"hub.verify_token"`
}

func FacebookHook(c echo.Context) error {
	var request FaceBookCheck
	err := c.Bind(&request)
	if err != nil {
		log.Printf("--- skype decode msg error!: %+v  >>%s", c.Request(), err)
	}

	log.Printf("user --  %#v", request)

	return c.String(http.StatusOK, request.HubChallenge)
}
