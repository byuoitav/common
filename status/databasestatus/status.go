package databasestatus

import (
	"fmt"
	"net/http"
	"os"

	"github.com/byuoitav/common/db"
	"github.com/byuoitav/common/log"
	stat "github.com/byuoitav/common/status"
	"github.com/labstack/echo"
)

// Handler validates that the microservice can talk to the database.
func Handler(ctx echo.Context) error {
	log.L.Infof("Status request from %v", ctx.Request().RemoteAddr)

	var err error
	status := stat.NewStatus()

	status.Bin = os.Args[0]

	status.Version, err = stat.GetMicroserviceVersion()
	if err != nil {
		status.Info["error"] = "failed to open version.txt"
		status.StatusCode = stat.Sick

		return ctx.JSON(http.StatusInternalServerError, status)
	}

	// Test a database retrieval to assess the status.
	vals, err := db.GetDB().GetAllBuildings()
	if len(vals) == 0 || err != nil {
		status.StatusCode = stat.Dead
		status.Info["error"] = fmt.Sprintf("unable to access database: %s", err)
	} else {
		status.StatusCode = stat.Healthy
		status.Info[""] = "Connected to database"
	}

	return ctx.JSON(http.StatusOK, status)
}
