package log

import (
	"net/http"

	"github.com/labstack/echo"
)

//SetLogLevel depends on a :level parameter in the endpoint
func SetLogLevel(c echo.Context) error {
	lvl := c.Param("level")

	L.Infof("Setting log level to %s", lvl)
	err := SetLevel(lvl)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	L.Infof("Log level set to %s", lvl)
	return c.JSON(http.StatusOK, "ok")
}

// GetLogLevel returns the current log level
func GetLogLevel(c echo.Context) error {

	L.Info("Getting log level.")
	lvl, err := GetLevel()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	L.Infof("Log level is %s", lvl)

	m := make(map[string]string)
	m["log-level"] = lvl

	return c.JSON(http.StatusOK, m)
}
