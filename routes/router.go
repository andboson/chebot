package routes

import (
	"github.com/andboson/chebot/controllers"
	"github.com/labstack/echo"
)

func Router() *echo.Echo {
	e := echo.New()
	e.POST("/ai.get_movies", controllers.GetMovies)
	e.POST("/skype.hook", controllers.SkypeHook)
	e.POST("/facebook.hook", controllers.FacebookHook)

	return e
}
