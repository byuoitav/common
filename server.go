package common

import (
	"net/http"

	"github.com/byuoitav/common/log"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// NewRouter returns a echo router with default routes/middleware to be used by all microservices.
func NewRouter() *echo.Echo {
	router := echo.New()

	router.Pre(middleware.RemoveTrailingSlash())
	router.Use(middleware.CORS())

	router.GET("/health", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, "Did you ever hear the tragedy of Darth Plagueis The Wise?")
	})

	router.PUT("/log-level/:level", log.SetLogLevel)
	router.GET("/log-level/:level", log.GetLogLevel)

	return router
}
