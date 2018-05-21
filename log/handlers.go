package log

import (
	"net/http"

	"github.com/labstack/echo"
)

//SetLogLevel depends on a :level parameter in the endpoint
func SetLogLevel(c echo.Context) error {
	lvl := c.Param("level")

	L.Info("Setting log level to %v")
	err := SetLevel(lvl)
	if err != nil {
		return c.Json(http.StatusBadRequest, err.Error())
	}

	L.Info("Log level set to %v", lvl)
	return c.JSON(http.StatusOK, "ok")
}
