package status

import (
	"bufio"
	"net/http"
	"os"
	"time"

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
	Name       string                 `json:"name,omitempty"`
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

//NewBaseStatus initializes a status stuct with the default stuff.
func NewBaseStatus() Status {
	var err error
	status := NewStatus()

	status.Bin = os.Args[0]
	status.Uptime = GetProgramUptime().String()

	status.Version, err = GetMicroserviceVersion()
	if err != nil {
		status.StatusCode = Sick
		status.Info["error"] = "failed to open version.txt"
	}

	status.StatusCode = Healthy
	return status
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
