package health

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/byuoitav/common/log"

	"github.com/byuoitav/common/v2/events"
	"github.com/labstack/echo"
)

func SendSuccessfulStartup(healthCheck func() map[string]string, MicroserviceName string, publish func(events.Event)) error {
	log.L.Infof("[HealthCheck] will report success in 10 seconds, waiting for listening services to be up")
	time.Sleep(10 * time.Second)
	log.L.Infof("[HealthCheck] Reporting microsrevice startup complete")

	log.L.Infof("[HealthCheck] Checking Health...")
	statusReport := healthCheck()
	allSuccess := true
	for _, v := range statusReport {
		if v != "ok" {
			allSuccess = false
		}
	}

	report := make(map[string]interface{})
	if allSuccess {
		report["success"] = "ok"
	} else {
		report["success"] = "errors"
	}
	report["report"] = statusReport
	report["Microservice"] = MicroserviceName

	log.L.Infof("[HealthCheck] Reporting...")
	for k, v := range statusReport {
		publishEvent(publish, k, v, MicroserviceName)
	}

	if allSuccess {
		publishEvent(publish, "ready", "true", MicroserviceName)
	} else {
		publishEvent(publish, "ready", "false", MicroserviceName)
	}
	return nil
}

func publishEvent(publish func(events.Event), k string, v string, name string) {
	publish(BuildEvent(k, v, name))
}

func BuildEvent(Key string, Value string, Device string) events.Event {
	hostname := os.Getenv("SYSTEM_ID")
	split := strings.Split(hostname, "-")
	room := fmt.Sprintf("%s-%s", split[0], split[1])

	roomInfo := events.GenerateBasicRoomInfo(room)

	deviceInfo := events.GenerateBasicDeviceInfo(Device)

	e := events.Event{
		GeneratingSystem: hostname,
		Timestamp:        time.Now(),
		AffectedRoom:     roomInfo,
		TargetDevice:     deviceInfo,
		Key:              Key,
		Value:            Value,
	}

	return e
}

func HealthCheck(context echo.Context) error {
	return context.JSON(http.StatusOK, "Uh, had a slight weapons malfunction. But, uh, everything's perfectly all right now. We're fine. We're all fine here, now, thank you. How are you?")
}
