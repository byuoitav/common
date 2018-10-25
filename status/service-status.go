package status

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"time"

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

// Status represents the microservice's health status
type Status struct {
	Name       string                 `json:"name"`
	Bin        string                 `json:"bin"`
	StatusCode string                 `json:"statuscode"`
	Version    string                 `json:"version"`
	Uptime     string                 `json:"uptime"`
	Info       map[string]interface{} `json:"info"`
}

var startTime time.Time

func init() {
	startTime = time.Now()
}

// NewStatus retuns an empty, initalized status struct
func NewStatus() Status {
	return Status{
		Info: make(map[string]interface{}),
	}
}

// DefaultStatusHandler can be used as a default mstatus handler
func DefaultStatusHandler(ctx echo.Context) error {
	log.L.Debugf("Status request from %v", ctx.Request().RemoteAddr)

	var err error
	status := NewStatus()

	status.Bin = os.Args[0]
	status.Uptime = GetProgramUptime().String()

	status.Version, err = GetMicroserviceVersion()
	if err != nil {
		status.StatusCode = Sick
		status.Info["error"] = "failed to open version.txt"
		return ctx.JSON(http.StatusInternalServerError, status)
	}

	status.StatusCode = Healthy
	status.Info[""] = "used default status handler"
	return ctx.JSON(http.StatusOK, status)
}

// DatabaseStatusHandler validates that the microservice can talk to the database.
func DatabaseStatusHandler(ctx echo.Context) error {
	log.L.Infof("Status request from %v", ctx.Request().RemoteAddr)

	var err error
	status := NewStatus()

	status.Bin = os.Args[0]
	status.Uptime = GetProgramUptime().String()

	status.Version, err = GetMicroserviceVersion()
	if err != nil {
		status.Info["error"] = "failed to open version.txt"
		status.StatusCode = Sick

		return ctx.JSON(http.StatusInternalServerError, status)
	}

	// Test a database retrieval to assess the status.
	vals, err := db.GetDB().GetAllBuildings()
	if len(vals) == 0 || err != nil {
		status.StatusCode = Dead
		status.Info["error"] = fmt.Sprintf("unable to access database: %s", err)
	} else {
		status.StatusCode = Healthy
		status.Info[""] = "Connected to database"
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

// GetProgramUptime returns how long the program has been running
func GetProgramUptime() time.Duration {
	return time.Since(startTime).Truncate(time.Second)
}
