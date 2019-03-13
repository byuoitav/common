package servicenow

import (
	"fmt"
	"strings"

	"github.com/byuoitav/common/jsonhttp"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/structs"
)

var (
	IncidentWebURL               = "https://api.byu.edu:443/domains/servicenow/tableapi/v1/table/incident"
	IncidentModifyWebURL         = "https://api.byu.edu:443/domains/servicenow/apiTable/v1/table/incident"
	IncidentAssignmentGroup      = "OIT-AV Support"
	IncidentService              = "TEC ROOM"
	IncidentResolutionService    = "TEC Room"
	IncidentClosureCode          = "Resolved"
	IncidentDefaultRequestor     = "AV Metrics Web Service"
	IncidentWorkStatus           = "Very Low"
	IncidentSensitivity          = "Very Low"
	IncidentSeverity             = "Very Low"
	IncidentReach                = "Very Low"
	IncidentDefaultContactNumber = "801-422-7671"
	IncidentClosedState          = "Closed"
)

//CreateIncident will create a new incident to the servicenow API
func CreateIncident(input structs.IncidentRequest) (structs.IncidentResponse, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	var output structs.IncidentResponseWrapper

	outputJSON, _, err := jsonhttp.CreateAndExecuteJSONRequest("CreateRequest", "POST", IncidentWebURL,
		input, headers, 20, &output)
	log.L.Debugf("Output JSON: %s", outputJSON)
	log.L.Debugf("Output JSON: %+v", output)

	return output.Result, err
}

//ModifyIncident to close or post notes to an existing incident via servicenow API
func ModifyIncident(input structs.IncidentRequest, sysID string) (structs.IncidentResponse, error) {

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	tries := 0

	weburl := fmt.Sprintf("%s/%s?sysparm_display_value=true", IncidentModifyWebURL, sysID)

	log.L.Debugf("WebURL: %s", weburl)

	for {
		var output structs.IncidentResponseWrapper

		outputJSON, outputResponse, err := jsonhttp.CreateAndExecuteJSONRequest("ModifyIncident", "PUT", weburl,
			input, headers, 20, &output)
		if err != nil {
			log.L.Errorf("Could not create and execute JSON request: %v", err)
			return output.Result, err
		}
		tries++

		log.L.Debugf("%v-%v", outputResponse.StatusCode, outputResponse.Status)

		log.L.Debugf("Output JSON: %s", outputJSON)
		log.L.Debugf("Output JSON: %+v", output)

		if outputResponse.StatusCode/100 == 2 || tries >= 5 {
			return output.Result, err
		}
	}
}

//QueryIncidentsByRoom - query all incidents by room number for the incident assignment group
func QueryIncidentsByRoom(RoomID string) ([]structs.IncidentResponse, error) {
	roomIDreplaced := strings.Replace(RoomID, "-", "+", -1)
	GroupName := IncidentAssignmentGroup

	weburl := fmt.Sprintf("active=true&sysparm_display_value=true&u_room=%s&assignment_group=%s", roomIDreplaced, GroupName)
	weburl = fmt.Sprintf("%s?%s", IncidentWebURL, weburl)

	log.L.Debugf("WebURL: %s", weburl)

	var output structs.QueriedIncidents

	input := ""

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	outputJSON, _, err := jsonhttp.CreateAndExecuteJSONRequest("querycategory", "GET", weburl,
		input, headers, 20, &output)

	if err != nil {
		log.L.Errorf("Could not create and execute JSON request: %v", err)
		return output.Result, err
	}

	log.L.Debugf("Output JSON: %s", outputJSON)
	log.L.Debugf("Output JSON: %+v", output)
	return output.Result, err
}

//GetIncident - Get ticket by INC#
func GetIncident(IncidentNumber string) (structs.IncidentResponse, error) {

	weburl := fmt.Sprintf("%s?sysparm_query=number=%s&sysparm_display_value=true", IncidentWebURL, IncidentNumber)

	var output structs.MultiIncidentResponseWrapper

	input := ""

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	outputJSON, _, err := jsonhttp.CreateAndExecuteJSONRequest("Get Incident By ID", "GET", weburl,
		input, headers, 200, &output)
	if err != nil {
		log.L.Errorf("Problem getting the incident: %v", err.Error())
		return structs.IncidentResponse{}, fmt.Errorf("Problem getting the incident: %v", err.Error())
	}

	log.L.Debugf("Output JSON: %s", outputJSON)
	log.L.Debugf("Output JSON: %+v", output)

	SysID := output.Result[0].SysID

	log.L.Debugf("Output sysID: %+v", SysID)

	return output.Result[0], err
}
