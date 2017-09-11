package main

import (
	"github.com/labstack/echo/middleware"
	"github.com/andboson/chebot/routes"
	"github.com/andboson/chebot/controllers"
	"github.com/andboson/chebot/repositories"
	"github.com/andboson/chebot/models"
)

func main() {
	models.InitConfig()
	e := routes.Router()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	go controllers.TelegramMessagesHandler()
	repositories.InitSkype()
	go controllers.InitFb()

	e.Logger.Fatal(e.Start(":1323"))
}
