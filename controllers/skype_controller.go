package controllers

import (
	"github.com/andboson/skypeapi"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"net/http"
	"github.com/andboson/chebot/repositories"
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

	repositories.ProcessMessage(proc)
	if request.Text == "taxi add" || request.Text == "taxi clear"{
		repositories.ProcessSkypeTaxiManage(request)
	}

	resp := map[string]string{
		"status": "success",
	}

	log.Printf("user --", request.From.Name, request.From.ID)

	return c.JSON(http.StatusOK, resp)
}
