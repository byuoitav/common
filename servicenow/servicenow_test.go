package servicenow

import (
	"testing"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/structs"
	"github.com/byuoitav/common/v2/events"
)

func TestServiceNow(t *testing.T) {
	log.SetLevel("debug")

	TestAlert := structs.Alert{
		BasicDeviceInfo: events.BasicDeviceInfo{
			BasicRoomInfo: events.BasicRoomInfo{
				BuildingID: "LSB",
				RoomID:     "LSB-5001",
			},
			DeviceID: "LSB-5001CP1",
		},
		AlertID:    "LSB-5001-CP1^System Communication Error^System^Critical",
		Type:       "System Communication Error",
		Category:   "System",
		Severity:   "warning",
		Message:    "LSB-5001-CP1 has not reported any state since 2019-02-20 13:38:18.639429878 -0700 MST",
		SystemType: "pi",
		Data:       "Data of the event goes here",
	}

	TestRoomIssue := structs.RoomIssue{
		BasicRoomInfo: events.BasicRoomInfo{
			BuildingID: "LSB",
			RoomID:     "LSB-5001",
		},
		Alerts:   []structs.Alert{TestAlert},
		Severity: structs.Critical,
		//IncidentID: "RPR0005375",
		// Notes:         "After consulting with Xuther, it seems that we should reboot AGAIN A THIRD TIME!",
		// HelpSentAt:    time.Now(),
		// Responders:    []string{"Joe", "Danny", "John"},
		// HelpArrivedAt: time.Now(),
		// Resolved:      true,
		// ResolutionInfo: structs.ResolutionInfo{
		// 	Code:  "Alert Cleared",
		// 	Notes: "alrerts auto-resolved",
		// },
	}

	log.L.Debugf("Test alert %v", TestAlert)

	id, err := SyncServiceNowWithRoomIssue(TestRoomIssue)

	if err != nil {
		log.L.Infof("Error: %v", err)
	} else {
		log.L.Infof("Success %v", id)
	}
}
