package main

import (
	"github.com/labstack/echo/middleware"
	"github.com/andboson/chebot/routes"
	"github.com/andboson/chebot/controllers"
	"github.com/andboson/chebot/repositories"
)

func main() {
	e := routes.Router()
	e.Use(middleware.Logger())
	//e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
	//	Level: 5,
	//}))
	e.Use(middleware.Recover())
	go controllers.TelegramMessagesHandler()
	repositories.InitSkype()

	e.Logger.Fatal(e.Start(":1323"))
}
