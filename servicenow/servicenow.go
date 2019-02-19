package servicenow

import (
	"os"

	"github.com/byuoitav/common/structs"
)

var token = os.Getenv("SERVICENOW_WSO2_TOKEN")

func getRoomIssueAlertTypeList(RoomIssue structs.RoomIssue) []structs.AlertType {
	var output []structs.AlertType

	for _, alert := range RoomIssue.Alerts {
		add := true
		for _, ele := range output {
			if ele == alert.Type {
				add = false
				break
			}
		}

		if add {
			output = append(output, alert.Type)
		}
	}

	return output
}
