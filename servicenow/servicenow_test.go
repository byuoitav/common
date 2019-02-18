package servicenow

import (
	"testing"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/structs"
	"github.com/byuoitav/common/v2/events"
)

func TestServiceNow(t *testing.T) {
	log.SetLevel("debug")

	log.L.Debugf("Testing create of new repair")

	TestAlert := structs.Alert{
		BasicDeviceInfo: events.BasicDeviceInfo{
			BasicRoomInfo: events.BasicRoomInfo{
				BuildingID: "ITB",
				RoomID:     "ITB-1108",
			},
			DeviceID: "ITB-1108-CP1",
		},
		AlertID:    "ITB-1108-CP1^System Communication Error^System^Critical",
		Type:       "System Communication Error",
		Category:   "System",
		Severity:   "Critical",
		MessageLog: []string{"ITB-1108-CP1 has not reported any state since 2019-02-18 12:38:18.639429878 -0700 MST"},
		SystemType: "pi",
		Data:       "Data of the event goes here",
	}

	log.L.Debugf("Test alert %v", TestAlert)

	repair, err := CreateRepair(TestAlert)

	if err != nil {
		log.L.Debugf("Error creating repair: %v", err)
	} else {
		log.L.Debugf("Repair: %v", repair)
	}
}
