package health

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/byuoitav/common/events"
	"github.com/labstack/echo"
)

func SendSuccessfulStartup(healthCheck func() map[string]string, MicroserviceName string, publish func(events.Event)) error {
	log.Printf("[HealthCheck] will report success in 10 seconds, waiting for listening services to be up")
	time.Sleep(10 * time.Second)
	log.Printf("[HealthCheck] Reporting microsrevice startup complete")

	log.Printf("[HealthCheck] Checking Health...")
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

	log.Printf("[HealthCheck] Reporting...")
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
	publish(BuildEvent(
		events.HEALTH,
		events.STARTUP,
		k, v, name,
	))
}

func BuildEvent(Type events.EventType, Cause events.EventCause, Key string, Value string, Device string) events.Event {

	info := events.EventInfo{
		Type:           Type,
		EventCause:     Cause,
		Device:         Device,
		EventInfoKey:   Key,
		EventInfoValue: Value,
	}

	hostname := os.Getenv("PI_HOSTNAME")
	split := strings.Split(hostname, "-")

	return events.Event{
		Hostname:         hostname,
		Timestamp:        time.Now().Format(time.RFC3339),
		LocalEnvironment: len(os.Getenv("LOCAL_ENVIRONMENT")) > 0,
		Event:            info,
		Building:         split[0],
		Room:             split[1],
	}

}

func HealthCheck(context echo.Context) error {
	return context.JSON(http.StatusOK, "Uh, had a slight weapons malfunction. But, uh, everything's perfectly all right now. We're fine. We're all fine here, now, thank you. How are you?")

}
