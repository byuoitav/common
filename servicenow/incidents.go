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

//SyncIncidentWithRoomIssue will either create or modify the incident for the room issue
func SyncIncidentWithRoomIssue(RoomIssue structs.RoomIssue) (structs.IncidentResponse, error) {
	if len(RoomIssue.IncidentID) == 0 {
		return CreateIncident(RoomIssue)
	}

	return ModifyIncident(RoomIssue)
}

//CreateIncident will create a new incident for a Room Issue
func CreateIncident(RoomIssue structs.RoomIssue) (structs.IncidentResponse, error) {

	alertTypes := getRoomIssueAlertTypeList(RoomIssue)

	shortDescription := fmt.Sprintf("%s is alerting with %v Alerts of type %s.", RoomIssue.RoomID, len(RoomIssue.Alerts), alertTypes)

	internalNotes := ""

	if RoomIssue.HelpSentAt.IsZero() == false {
		internalNotes += fmt.Sprintf("\nHelp was was sent at: %s\n", RoomIssue.HelpSentAt)
	}

	if RoomIssue.HelpArrivedAt.IsZero() == false {
		internalNotes += fmt.Sprintf("\nHelp arrived at: %s\n", RoomIssue.HelpArrivedAt)
	}

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

//ModifyIncident to close or post notes to an existing incident
func ModifyIncident(RoomIssue structs.RoomIssue) (structs.IncidentResponse, error) {

	IncidentNumber := RoomIssue.IncidentID
	ExistingIncident, _ := GetIncident(IncidentNumber)

	weburl := fmt.Sprintf("%s/%s?sysparm_display_value=true", incidentModifyWebURL, ExistingIncident.SysID)

	log.L.Debugf("WebURL: %s", weburl)

	input := structs.IncidentRequest{}

	internalNotes := ""

	log.L.Warnf("Existing Notes: %s", ExistingIncident.InternalNotes)

	if !strings.Contains(ExistingIncident.InternalNotes, "Help was sent at:") {
		if RoomIssue.HelpSentAt.IsZero() == false {
			internalNotes += fmt.Sprintf("\nHelp was was sent at: %s\n", RoomIssue.HelpSentAt)
		}
	}

	if !strings.Contains(ExistingIncident.InternalNotes, "Help arrived at:") {
		if RoomIssue.HelpArrivedAt.IsZero() == false {
			internalNotes += fmt.Sprintf("\nHelp arrived at: %s\n", RoomIssue.HelpArrivedAt)
		}
	}

	if len(RoomIssue.Notes) > 0 {
		if !strings.Contains(ExistingIncident.InternalNotes, RoomIssue.Notes) {
			internalNotes += "\n--------Room Notes-------\n"
			internalNotes += RoomIssue.Notes + "\n"
		}
	}

	for _, alert := range RoomIssue.Alerts {
		if len(alert.Message) > 0 {
			tmpMessage := fmt.Sprintf("\n--------%s Notes-------\n%s\n", alert.DeviceID, alert.Message)

			if !strings.Contains(ExistingIncident.InternalNotes, tmpMessage) {
				internalNotes += tmpMessage
			}
		}
	}

	internalNotes = strings.TrimSpace(internalNotes)
	input.InternalNotes = internalNotes

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
		input.State = "Closed"
	}

	workLog = strings.TrimSpace(workLog)
	input.WorkLog = workLog
	input.ClosureCode = resolutionClosureCode
	input.ResolutionService = resolutionService
	input.ResolutionAction = resolutionAction

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	tries := 0

	for {
		var output structs.IncidentResponseWrapper

		outputJSON, outputResponse, err := jsonhttp.CreateAndExecuteJSONRequest("ModifyIncident", "PUT", weburl,
			input, headers, 20, &output)

		tries++

		log.L.Debugf("%v-%v", outputResponse.StatusCode, outputResponse.Status)

		log.L.Debugf("Output JSON: %s", outputJSON)
		log.L.Debugf("Output JSON: %+v", output)

		if outputResponse.StatusCode/100 == 2 || tries >= 5 {
			return output.Result, err
		}
	}
}

//QueryIncidentsByRoomAndGroupName - query all incidents by room number and group
func QueryIncidentsByRoomAndGroupName(BuildingID string, RoomID string, GroupName string) ([]structs.IncidentResponse, error) {
	weburl := fmt.Sprintf("%s?active=true&sysparm_display_value=true&u_room=%s+%s&assignment_group=%s", incidentWebURL, BuildingID, RoomID, GroupName)

	log.L.Debugf("WebURL: %s", weburl)

	var output structs.QueriedIncidents

	input := ""

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	outputJson, _, err := jsonhttp.CreateAndExecuteJSONRequest("querycategory", "GET", weburl,
		input, headers, 20, &output)

	log.L.Debugf("Output JSON: %s", outputJson)
	log.L.Debugf("Output JSON: %+v", output)
	return output.Result, err
}

//GetIncident - Get ticket by INC#
func GetIncident(IncidentNumber string) (structs.IncidentResponse, error) {

	weburl := fmt.Sprintf("%s?sysparm_query=number=%s", incidentWebURL, IncidentNumber)

	var output structs.MultiIncidentResponseWrapper

	input := ""

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	outputJson, _, err := jsonhttp.CreateAndExecuteJSONRequest("Get Incident By ID", "GET", weburl,
		input, headers, 200, &output)

	log.L.Debugf("Output JSON: %s", outputJson)
	log.L.Debugf("Output JSON: %+v", output)

	SysID := output.Result[0].SysID

	log.L.Debugf("Output sysID: %+v", SysID)

	return output.Result[0], err
}
