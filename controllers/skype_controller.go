package controllers

import (
	"github.com/andboson/chebot/repositories"
	"github.com/andboson/skypeapi"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"net/http"
)

func SkypeHook(c echo.Context) error {
	var request skypeapi.Activity
	err := c.Bind(&request)
	if err != nil {
		log.Printf("--- skype decode msg error!: %+v  >>%s", c.Request(), err)
	}
	var proc repositories.Processer
	proc = repositories.SkypeProcessor{
		Message: &request,
	}

	var result bool
	result = repositories.ProcessMessage(proc)
	if !result {
		result = repositories.ProcessSkypeTaxiManage(request)
	}

	if !result {
		proc.NoResults()
	}

	resp := map[string]string{
		"status": "success",
	}

	log.Printf("user --", request.From.Name, request.From.ID)

	return c.JSON(http.StatusOK, resp)
}
