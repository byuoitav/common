package servicenow

import (
	"fmt"
	"strings"

	"github.com/byuoitav/common/jsonhttp"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/structs"
)

var (
	incidentWebURL               = "https://api.byu.edu:443/domains/servicenow/tableapi/v1/table/incident"
	incidentModifyWebURL         = "https://api.byu.edu:443/domains/servicenow/apiTable/v1/table/incident"
	incidentAssignmentGroup      = "OIT-AV Support"
	incidentService              = "TEC ROOM"
	incidentResolutionService    = "TEC Room"
	incidentClosureCode          = "Resolved"
	incidentDefaultRequestor     = "AV Metrics Web Service"
	incidentWorkStatus           = "Very Low"
	incidentSensitivity          = "Very Low"
	incidentSeverity             = "Very Low"
	incidentReach                = "Very Low"
	incidentDefaultContactNumber = "801-422-7671"
)

//CreateIncident will create a new incident for a Room Issue
func CreateIncident(RoomIssue structs.RoomIssue) (structs.IncidentResponse, error) {

	alertTypes := getRoomIssueAlertTypeList(RoomIssue)

	shortDescription := fmt.Sprintf("%s is alerting with %v Alerts of type %s.", RoomIssue.RoomID, len(RoomIssue.Alerts), alertTypes)

	internalNotes := ""

	if len(RoomIssue.Notes) > 0 {
		internalNotes += "\n--------Room Notes-------\n"
		internalNotes += RoomIssue.Notes + "\n"
	}

	for _, alert := range RoomIssue.Alerts {
		if len(alert.Message) > 0 {
			internalNotes += fmt.Sprintf("\n--------%s Notes-------\n", alert.DeviceID)
			internalNotes += alert.Message + "\n"
		}
	}

	internalNotes = strings.TrimSpace(internalNotes)

	workLog := ""
	resolutionClosureCode := ""
	resolutionService := ""
	resolutionAction := ""

	if RoomIssue.Resolved {
		workLog += "\n-------Resolution Info-------\n"
		workLog += RoomIssue.ResolutionInfo.Code + "\n"
		workLog += RoomIssue.ResolutionInfo.Notes + "\n"

		resolutionClosureCode = incidentClosureCode
		resolutionService = incidentResolutionService
		resolutionAction = RoomIssue.ResolutionInfo.Code
	}

	workLog = strings.TrimSpace(workLog)

	roomIDreplaced := strings.Replace(RoomIssue.RoomID, "-", " ", -1)

	requester := ""

	for _, alert := range RoomIssue.Alerts {
		if len(alert.Requester) > 0 {
			requester = alert.Requester
		}
	}

	if len(requester) == 0 {
		requester = incidentDefaultRequestor
	}

	input := structs.IncidentRequest{
		Service:       incidentService,
		CallerID:      requester,
		ContactNumber: incidentDefaultContactNumber,

		AssignmentGroup: incidentAssignmentGroup,
		Room:            roomIDreplaced,

		ShortDescription: shortDescription,

		Severity:    incidentSeverity,
		Reach:       incidentReach,
		WorkStatus:  incidentWorkStatus,
		Sensitivity: incidentSensitivity,

		InternalNotes: internalNotes,
		WorkLog:       workLog,

		ClosureCode:       resolutionClosureCode,
		ResolutionService: resolutionService,
		ResolutionAction:  resolutionAction,
	}

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	var output structs.IncidentResponseWrapper

	outputJSON, _, err := jsonhttp.CreateAndExecuteJSONRequest("CreateRequest", "POST", incidentWebURL,
		input, headers, 20, &output)
	log.L.Debugf("Output JSON: %s", outputJSON)
	log.L.Debugf("Output JSON: %+v", output)

	return output.Result, err
}

// //we need to be able to access the sysID of the incident Ticket
// //TO DO: takes incident ID and string for internal notes
// func ModifyIncident(Alert structs.Alert) (structs.ReceiveIncident, error) {
// 	IncidentNumber := Alert.IncidentID
// 	SysID, _ := GetSysID(IncidentNumber)
// 	weburl := fmt.Sprintf("https://api.byu.edu/domains/servicenow/incident/v1.1/incident/%s?sysparm_display_value=true", SysID)
// 	log.L.Debugf("WebURL: %s", weburl)
// 	var internalNotes string
// 	var state string

// 	//if you want to pull info from the alert
// 	if Alert.HelpSentAt.IsZero() == false && Alert.HelpArrivedAt.IsZero() == true {
// 		internalNotes = "Help was was sent at: " + fmt.Sprintf("%s", Alert.HelpSentAt)
// 		state = "Assigned"
// 	}

// 	if Alert.HelpSentAt.IsZero() == false && Alert.HelpArrivedAt.IsZero() == false {
// 		internalNotes += "\n" + " Help arrived at: " + fmt.Sprintf("%s", Alert.HelpArrivedAt)
// 		state = "Work In Progress"
// 	}

// 	input := structs.Incident{
// 		State:         state,
// 		InternalNotes: internalNotes,
// 		Description:   "This is a description, want to see what happens",
// 	}
// 	headers := map[string]string{
// 		"Authorization": "Bearer " + token,
// 		"Content-Type":  "application/json",
// 	}

// 	var output structs.ReceiveIncidentWrapper
// 	outputJson, _, err := jsonhttp.CreateAndExecuteJSONRequest("ModifyIncident", "PUT", weburl,
// 		input, headers, 20, &output)
// 	log.L.Debugf("Output JSON: %s", outputJson)
// 	log.L.Debugf("Output JSON: %+v", output)
// 	return output.Result, err

// }

// func GetResolutionActions() (structs.ResolutionCategories, error) {
// 	weburl := "https://api.byu.edu:443/domains/servicenow/tableapi/v1/table/u_inc_resolution_cat?sysparm_query=active%3Dtrue%5Eassignment_group%3Djavascript%3AgetMyAssignmentGroups()"
// 	log.L.Debugf("WebURL: %s", weburl)
// 	var output structs.ResolutionCategories
// 	input := ""
// 	headers := map[string]string{
// 		"Authorization": "Bearer " + token,
// 		"Content-Type":  "application/json",
// 	}
// 	outputJson, _, err := jsonhttp.CreateAndExecuteJSONRequest("querycategory", "GET", weburl,
// 		input, headers, 20, &output)
// 	log.L.Debugf("Output JSON: %s", outputJson)
// 	log.L.Debugf("Output JSON: %+v", output)
// 	return output, err
// }

// func CloseIncident(Alert structs.Alert) (structs.ReceiveIncident, error) {
// 	IncidentID := Alert.IncidentID
// 	SysID, _ := GetSysID(IncidentID)
// 	weburl := fmt.Sprintf("https://api.byu.edu/domains/servicenow/incident/v1.1/incident/%s?sysparm_display_value=true", SysID)
// 	log.L.Debugf("WebURL: %s", weburl)
// 	state := "Closed"
// 	closurecode := "Resolved"
// 	resolutionservice := "TEC Room"

// 	input := structs.Incident{
// 		State:             state,
// 		ClosureCode:       closurecode,
// 		ResolutionService: resolutionservice,
// 		ResolutionAction:  Alert.ResolutionInfo.Code,
// 		WorkLog:           Alert.ResolutionInfo.Notes,
// 	}
// 	headers := map[string]string{
// 		"Authorization": "Bearer " + token,
// 		"Content-Type":  "application/json",
// 	}

// 	var output structs.ReceiveIncidentWrapper
// 	outputJson, _, err := jsonhttp.CreateAndExecuteJSONRequest("ModifyIncident", "PUT", weburl,
// 		input, headers, 20, &output)
// 	log.L.Debugf("Output JSON: %s", outputJson)
// 	log.L.Debugf("Output JSON: %+v", output)
// 	return output.Result, err
// }

// //query all incidents for a given assignment group
// func QueryIncidentsByGroup(GroupName string) (structs.QueriedIncidents, error) {
// 	weburl := fmt.Sprintf("https://api.byu.edu/domains/servicenow/incident/v1.1/incident?active=true&assignment_group=%s&sysparm_display_value=true", GroupName)
// 	log.L.Debugf("WebURL: %s", weburl)
// 	var output structs.QueriedIncidents
// 	input := ""
// 	headers := map[string]string{
// 		"Authorization": "Bearer " + token,
// 		"Content-Type":  "application/json",
// 	}
// 	outputJson, _, err := jsonhttp.CreateAndExecuteJSONRequest("querycategory", "GET", weburl,
// 		input, headers, 20, &output)
// 	log.L.Debugf("Output JSON: %s", outputJson)
// 	log.L.Debugf("Output JSON: %+v", output)
// 	return output, err
// }

// //query all incidents by room number
// func QueryIncidentsByRoom(BuildingID string, RoomID string) (structs.QueriedIncidents, error) {
// 	weburl := fmt.Sprintf("https://api.byu.edu/domains/servicenow/incident/v1.1/incident?active=true&sysparm_display_value=true&u_room=%s+%s", BuildingID, RoomID)
// 	log.L.Debugf("WebURL: %s", weburl)
// 	var output structs.QueriedIncidents
// 	input := ""
// 	headers := map[string]string{
// 		"Authorization": "Bearer " + token,
// 		"Content-Type":  "application/json",
// 	}
// 	outputJson, _, err := jsonhttp.CreateAndExecuteJSONRequest("querycategory", "GET", weburl,
// 		input, headers, 20, &output)
// 	log.L.Debugf("Output JSON: %s", outputJson)
// 	log.L.Debugf("Output JSON: %+v", output)
// 	return output, err
// }

// //query all incidents by room number and group
// func QueryIncidentsByRoomAndGroupName(BuildingID string, RoomID string, GroupName string) (structs.QueriedIncidents, error) {
// 	weburl := fmt.Sprintf("https://api.byu.edu/domains/servicenow/incident/v1.1/incident?active=true&sysparm_display_value=true&u_room=%s+%s&assignment_group=%s", BuildingID, RoomID, GroupName)
// 	log.L.Debugf("WebURL: %s", weburl)
// 	var output structs.QueriedIncidents
// 	input := ""
// 	headers := map[string]string{
// 		"Authorization": "Bearer " + token,
// 		"Content-Type":  "application/json",
// 	}
// 	outputJson, _, err := jsonhttp.CreateAndExecuteJSONRequest("querycategory", "GET", weburl,
// 		input, headers, 20, &output)
// 	log.L.Debugf("Output JSON: %s", outputJson)
// 	log.L.Debugf("Output JSON: %+v", output)
// 	return output, err
// }

// //TODO query all of the users in the system (Net_id)

// func QueryAllUsers() (structs.QueriedUsers, error) {
// 	weburl := fmt.Sprint("https://api.byu.edu:443/domains/servicenow/tableapi/v1/table/sys_user?sysparm_query=active=true^assignment_group=javascript:getMyAssignmentGroups()")
// 	var output structs.QueriedUsers
// 	input := ""
// 	headers := map[string]string{
// 		"Authorization": "Bearer " + token,
// 		"Content-Type":  "application/json",
// 	}
// 	outputJson, _, err := jsonhttp.CreateAndExecuteJSONRequest("queryUsers", "GET", weburl,
// 		input, headers, 500, &output)
// 	log.L.Debugf("Output JSON: %s", outputJson)
// 	log.L.Debugf("Output JSON: %+v", output)
// 	return output, err
// }

// //Get ticket by INC# return the sysID
// func GetSysID(IncidentNumber string) (string, error) {
// 	weburl := fmt.Sprintf("https://api.byu.edu:443/domains/servicenow/tableapi/v1/table/incident?sysparm_query=number=%s", IncidentNumber)
// 	var output structs.IncidentWrapper
// 	input := ""
// 	headers := map[string]string{
// 		"Authorization": "Bearer " + token,
// 		"Content-Type":  "application/json",
// 	}
// 	outputJson, _, err := jsonhttp.CreateAndExecuteJSONRequest("queryUsers", "GET", weburl,
// 		input, headers, 200, &output)
// 	log.L.Debugf("Output JSON: %s", outputJson)
// 	log.L.Debugf("Output JSON: %+v", output)
// 	SysID := output.Result[0].SysId
// 	log.L.Debugf("Output sysID: %+v", SysID)
// 	return SysID, err
// }
