package main

import (
	"github.com/labstack/echo/middleware"
	"github.com/andboson/chebot/routes"
	"github.com/andboson/chebot/controllers"
	"github.com/andboson/chebot/repositories"
	"github.com/andboson/chebot/models"
)

func main() {
	e := routes.Router()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	models.InitConfig()
	go controllers.TelegramMessagesHandler()
	repositories.InitSkype()

	e.Logger.Fatal(e.Start(":1323"))
}
