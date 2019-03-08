package servicenow

import (
	"fmt"
	"os"
	"sync"

	"github.com/byuoitav/common/jsonhttp"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/structs"
)

var token = os.Getenv("SERVICENOW_WSO2_TOKEN")

func init() {
	mutexMap = make(map[string]*sync.Mutex)
	mutexMapMutex = sync.Mutex{}
}

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

var mutexMap map[string]*sync.Mutex
var mutexMapMutex sync.Mutex

func SyncServiceNowWithRoomIssue(RoomIssue structs.RoomIssue) (string, error) {
	key := RoomIssue.RoomID //no more severity

	mutexMapMutex.Lock()
	mutie, ok := mutexMap[key]

	if !ok {
		mutie = &sync.Mutex{}
		mutexMap[key] = mutie
	}

	mutexMapMutex.Unlock()

	mutie.Lock()
	defer mutie.Unlock()

	//fot now
	return "", nil

	// if structs.ContainsAllTags(RoomIssue.AlertSeverities, structs.Warning) {
	// 	repairResponse, err := SyncRepairWithRoomIssue(RoomIssue)

	// 	if err != nil {
	// 		return "", err
	// 	}

	// 	return repairResponse.Number, nil
	// }

	// incidentReponse, err := SyncIncidentWithRoomIssue(RoomIssue)

	// if err != nil {
	// 	return "", err
	// }

	// return incidentReponse.Number, nil
}

func QueryAllUsers() (structs.QueriedUsers, error) {
	weburl := fmt.Sprint("https://api.byu.edu:443/domains/servicenow/tableapi/v1/table/sys_user?sysparm_query=active=true^assignment_group=javascript:getMyAssignmentGroups()")
	var output structs.QueriedUsers
	input := ""
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}
	outputJson, _, err := jsonhttp.CreateAndExecuteJSONRequest("queryUsers", "GET", weburl,
		input, headers, 500, &output)
	log.L.Debugf("Output JSON: %s", outputJson)
	log.L.Debugf("Output JSON: %+v", output)
	return output, err
}

func GetResolutionActions() (structs.ResolutionCategories, error) {
	weburl := "https://api.byu.edu:443/domains/servicenow/tableapi/v1/table/u_inc_resolution_cat?sysparm_query=active%3Dtrue%5Eassignment_group%3Djavascript%3AgetMyAssignmentGroups()"
	log.L.Debugf("WebURL: %s", weburl)
	var output structs.ResolutionCategories
	input := ""
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}
	outputJson, _, err := jsonhttp.CreateAndExecuteJSONRequest("querycategory", "GET", weburl,
		input, headers, 20, &output)
	log.L.Debugf("Output JSON: %s", outputJson)
	log.L.Debugf("Output JSON: %+v", output)
	return output, err
}
