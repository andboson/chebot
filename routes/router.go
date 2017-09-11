package routes

import (
	"github.com/andboson/chebot/controllers"
	"github.com/labstack/echo"
	"net/http"
	"github.com/labstack/gommon/log"
)

func Router() *echo.Echo {
	e := echo.New()
	e.POST("/ai.get_movies", controllers.GetMovies)
	e.POST("/skype.hook", controllers.SkypeHook)

	go func() {
		http.HandleFunc("/facebook.hook", controllers.FbMess.Handler)
		log.Fatal(http.ListenAndServe(":1324", nil))
	}()

	return e
}
