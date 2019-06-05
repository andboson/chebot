package routes

import (
	"github.com/andboson/chebot/controllers"
	"github.com/labstack/echo"
)

func Router() *echo.Echo {
	e := echo.New()
	e.POST("/ai.get_response", controllers.GetAiResponse)
	e.POST("/skype.hook", controllers.SkypeHook)
	e.OPTIONS("/skype.hook", func(c echo.Context) error{
		return  c.NoContent(204)
	})

	return e
}
