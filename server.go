package common

import (
	"net/http"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/status"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// NewRouter returns a echo router with default routes/middleware to be used by all microservices.
func NewRouter() *echo.Echo {
	router := echo.New()

	router.Pre(middleware.RemoveTrailingSlash())
	router.Use(middleware.CORS())

	// return a default health message
	router.GET("/health", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, "Did you ever hear the tragedy of Darth Plagueis The Wise?")
	})

	// return a default mstatus message
	router.GET("/mstatus", status.DefaultMStatusHandler)
	router.GET("/status", status.DefaultMStatusHandler)

	router.PUT("/log-level/:level", log.SetLogLevel)
	router.GET("/log-level", log.GetLogLevel)

	return router
}
