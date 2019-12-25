package routes

import (
	"github.com/andboson/chebot/controllers"
	"github.com/andboson/chebot/models"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var apiAuthMiddleware = middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
	KeyLookup: "header:api-key",
	Validator: func(token string, context echo.Context) (bool, error) {
		return token == models.Conf.APIKey, nil
	},
})

func Router() *echo.Echo {
	e := echo.New()
	e.POST("/ai.get_response", controllers.GetAiResponse)
	e.POST("/skype.hook", controllers.SkypeHook)
	e.OPTIONS("/skype.hook", func(c echo.Context) error{
		return  c.NoContent(204)
	})
	e.POST("/web/message", controllers.WebMessage, apiAuthMiddleware)

	return e
}
