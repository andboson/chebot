package main

import (
	"github.com/labstack/echo/middleware"
	"github.com/andboson/chebot/routes"
	"github.com/andboson/chebot/controllers"
)

func main() {
	e := routes.Router()
	e.Use(middleware.Logger())
	//e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
	//	Level: 5,
	//}))
	e.Use(middleware.Recover())
	go controllers.TelegramMessagesHandler()

	controllers.InitSkype()
	e.Logger.Fatal(e.Start(":1323"))
}
