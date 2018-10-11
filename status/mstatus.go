package status

import (
	"bufio"
	"fmt"
	"net/http"
	"os"

	"github.com/byuoitav/common/db"
	"github.com/byuoitav/common/log"
	"github.com/labstack/echo"
)

const (
	// Healthy should be the response when the microservice is working 100% properly
	Healthy = "healthy"

	// Sick should be the response when the microservice is partially working or healing
	Sick = "sick"

	// Dead should be the response when the microservice is totally dead
	Dead = "dead"

	versionPath = "version.txt"
)

// MStatus represents the microservice's health status
type MStatus struct {
	StatusCode string      `json:"statuscode"`
	Version    string      `json:"version"`
	Info       interface{} `json:"info"`
}

// DefaultMStatusHandler can be used as a default mstatus handler
func DefaultMStatusHandler(ctx echo.Context) error {
	log.L.Infof("MStatus request from %v", ctx.Request().RemoteAddr)

	var status MStatus
	var err error

	status.Version, err = GetMicroserviceVersion()
	if err != nil {
		status.StatusCode = Sick
		status.Info = "failed to open version.txt"
		return ctx.JSON(http.StatusInternalServerError, status)
	}

	status.StatusCode = Healthy
	status.Info = "used default mstatus handler"
	return ctx.JSON(http.StatusOK, status)
}

// DatabaseMStatusHandler validates that the microservice can talk to the database.
func DatabaseMStatusHandler(ctx echo.Context) error {
	log.L.Infof("MStatus request from %v", ctx.Request().RemoteAddr)

	var status MStatus
	var err error

	status.Version, err = GetMicroserviceVersion()
	if err != nil {
		status.Info = "failed to open version.txt"
		status.StatusCode = Sick

		return ctx.JSON(http.StatusInternalServerError, status)
	}

	// Test a database retrieval to assess the status.
	vals, err := db.GetDB().GetAllBuildings()
	if len(vals) == 0 || err != nil {
		status.StatusCode = Dead
		status.Info = fmt.Sprintf("unable to access database: %s", err)
	} else {
		status.StatusCode = Healthy
		status.Info = "Connected to database"
	}

	return ctx.JSON(http.StatusOK, status)
}

// GetMicroserviceVersion returns the version number located in "version.txt"
func GetMicroserviceVersion() (string, error) {
	file, err := os.Open(versionPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan() // only read first line
	if err := scanner.Err(); err != nil {
		return "", err
	}

	version := scanner.Text()
	return version, nil
}
