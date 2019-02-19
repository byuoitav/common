package servicenow

import (
	"fmt"
	"strings"
	"time"

	"github.com/byuoitav/common/jsonhttp"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/structs"
)

var (
	repairWebURL             = "https://api.byu.edu:443/domains/servicenow/tableapi/v1/table/u_maint_repair"
	repairModifyWebURL       = "https://api.byu.edu:443/domains/servicenow/apiTable/v1/table/u_maint_repair"
	repairAssignmentGroup    = "OIT-AV Support"
	repairRequestOrigination = "On-Site"
	repairEquipmentReturn    = "On-Site"
	repairService            = "TEC Room"
	repairDefaultRequestor   = "AV Metrics Web Service"
	repairDateNeeded         = "ASAP"
	repairClosedState        = "Closed Complete"
)

func getNotesForRoomIssueForRepair(RoomIssue structs.RoomIssue) string {
	internalNotes := ""

	if RoomIssue.Resolved {
		internalNotes += "\n---------Resolution Info-----------\n"
		internalNotes += fmt.Sprintf("Resolution Code: %s\n", RoomIssue.ResolutionInfo.Code)
		internalNotes += RoomIssue.ResolutionInfo.Notes
	}

	if RoomIssue.HelpSentAt.IsZero() == false {
		internalNotes += fmt.Sprintf("\nHelp %s was sent at: %s\n", RoomIssue.Responders, RoomIssue.HelpSentAt.Format("01/02/2006 3:04 PM"))
	}

	if RoomIssue.HelpSentAt.IsZero() == false {
		internalNotes += fmt.Sprintf("\nHelp arrived at: %s\n", RoomIssue.HelpArrivedAt.Format("01/02/2006 3:04 PM"))
	}

	if len(RoomIssue.NotesLog) > 0 {
		internalNotes += "\n-----Room Notes-----\n"

		for _, note := range RoomIssue.NotesLog {
			internalNotes += note + "\n"
		}
	}

	for _, alert := range RoomIssue.Alerts {
		if len(alert.MessageLog) > 0 {
			internalNotes += fmt.Sprintf("\n-----System Messages for %s-----\n", alert.DeviceID)

			for _, note := range alert.MessageLog {
				internalNotes += note + "\n"
			}
		}

		dataStr := fmt.Sprintf("%s", alert.Data)
		if len(dataStr) > 0 {
			internalNotes += fmt.Sprintf("\n-----Alert Data for for %s-----\n", alert.DeviceID)
			internalNotes += dataStr
		}
	}

	internalNotes = strings.TrimSpace(internalNotes)

	return internalNotes
}

func createRepairRequest(RoomIssue structs.RoomIssue) structs.RepairRequest {

	alertTypes := getRoomIssueAlertTypeList(RoomIssue)

	shortDescription := fmt.Sprintf("%s is alerting with %v Alerts of type %s.", RoomIssue.RoomID, len(RoomIssue.Alerts), alertTypes)

	internalNotes := getNotesForRoomIssueForRepair(RoomIssue)

	year, month, day := time.Now().Date()

	requestdate := fmt.Sprintf("%v-%v-%v", year, int(month), day)

	roomIDreplaced := strings.Replace(RoomIssue.RoomID, "-", " ", -1)

	requester := ""

	for _, alert := range RoomIssue.Alerts {
		if len(alert.Requester) > 0 {
			requester = alert.Requester
		}
	}

	if len(requester) == 0 {
		requester = repairDefaultRequestor
	}

	input := structs.RepairRequest{
		Service:            repairService,
		Building:           RoomIssue.BuildingID,
		Room:               roomIDreplaced,
		AssignmentGroup:    repairAssignmentGroup,
		ShortDescription:   shortDescription,
		InternalNotes:      internalNotes,
		RequestOriginator:  requester,
		RequestDate:        requestdate,
		DateNeeded:         repairDateNeeded,
		RequestOrigination: repairRequestOrigination,
		EquipmentReturn:    repairEquipmentReturn,
	}

	return input
}

//SyncRepairWithRoomIssue will either create or modify the incident for the room issue
func SyncRepairWithRoomIssue(RoomIssue structs.RoomIssue) (structs.RepairResponse, error) {
	if len(RoomIssue.IncidentID) == 0 {
		return CreateRepair(RoomIssue)
	}

	return ModifyRepair(RoomIssue)
}

//CreateRepair is to create a new repair ticket from a new room issue
func CreateRepair(RoomIssue structs.RoomIssue) (structs.RepairResponse, error) {

	input := createRepairRequest(RoomIssue)

	if RoomIssue.Resolved {
		//update notes with resolution info
		resolutionInfo := RoomIssue.ResolutionInfo.Code + "\n" + RoomIssue.ResolutionInfo.Notes
		input.InternalNotes = fmt.Sprintf("-------Resolution Info------\n%s\n%s", resolutionInfo, input.InternalNotes)

		//set state to closed
		input.State = "closed"
	}

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	var output structs.RepairResponseWrapper

	outputJSON, _, err := jsonhttp.CreateAndExecuteJSONRequest("CreateRequest", "POST", repairWebURL,
		input, headers, 20, &output)

	log.L.Debugf("Output JSON: %s", outputJSON)
	log.L.Debugf("Output JSON: %+v", output)

	return output.Result, err
}

//ModifyRepair updates the notes on a repair ticket
func ModifyRepair(RoomIssue structs.RoomIssue) (structs.RepairResponse, error) {
	RepairNum := RoomIssue.IncidentID

	ExistingRepair, _ := GetRepair(RepairNum)

	//we need to sync up the existing notes
	internalNotes := ""
	log.L.Debugf("Existing Notes: %s", ExistingRepair.InternalNotes)

	if !strings.Contains(ExistingRepair.InternalNotes, "Help was sent at:") {
		if RoomIssue.HelpSentAt.IsZero() == false {
			internalNotes += fmt.Sprintf("\nHelp was sent at: %s\n", RoomIssue.HelpSentAt.Format("01/02/2006 3:04 PM"))
		}
	}

	if !strings.Contains(ExistingRepair.InternalNotes, "Help arrived at:") {
		if RoomIssue.HelpArrivedAt.IsZero() == false {
			internalNotes += fmt.Sprintf("\nHelp arrived at: %s\n", RoomIssue.HelpArrivedAt.Format("01/02/2006 3:04 PM"))
		}
	}

	if len(RoomIssue.Notes) > 0 {
		if !strings.Contains(ExistingRepair.InternalNotes, RoomIssue.Notes) {
			internalNotes += "\n--------Room Notes-------\n"
			internalNotes += RoomIssue.Notes + "\n"
		}
	}

	for _, alert := range RoomIssue.Alerts {
		if len(alert.Message) > 0 {
			tmpMessage := fmt.Sprintf("\n--------%s Notes-------\n%s\n", alert.DeviceID, alert.Message)

			if !strings.Contains(ExistingRepair.InternalNotes, tmpMessage) {
				internalNotes += tmpMessage
			}
		}
	}

	if RoomIssue.Resolved {
		internalNotes += "\n-------Resolution Info-------\n"
		internalNotes += RoomIssue.ResolutionInfo.Code + "\n"
		internalNotes += RoomIssue.ResolutionInfo.Notes + "\n"
	}

	internalNotes = strings.TrimSpace(internalNotes)

	weburl := fmt.Sprintf("%s/%s", repairModifyWebURL, ExistingRepair.SysID)

	log.L.Debugf("Web URL %v", weburl)

	input := structs.RepairRequest{
		InternalNotes: internalNotes,
	}

	//check for resolution
	if RoomIssue.Resolved {
		input.State = repairClosedState
	}

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	tries := 0

	for {
		var output structs.RepairResponseWrapper

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

// //CloseRepair will just close the repair with the specified resolution info

//QueryRepairsByRoomAndGroupName gets a list of repairs by room assigned to specified group
func QueryRepairsByRoomAndGroupName(BuildingID string, RoomID string, GroupName string) ([]structs.RepairResponse, error) {
	weburl := fmt.Sprintf("%s?active=true&sysparm_display_value=true&u_room=%s+%s&assignment_group=%s", repairWebURL, BuildingID, RoomID, GroupName)

	var output structs.MultiRepairResponseWrapper

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	outputJSON, _, err := jsonhttp.CreateAndExecuteJSONRequest("querycategory", "GET", weburl,
		"", headers, 20, &output)

	log.L.Debugf("Output JSON: %s", outputJSON)
	log.L.Debugf("Output JSON: %+v", output)

	return output.Result, err
}

//GetRepair Get repair ticket by the number
func GetRepair(RepairNum string) (structs.RepairResponse, error) {
	weburl := fmt.Sprintf("%s?sysparm_query=number=%s&sysparm_display_value=true", repairWebURL, RepairNum)

	var output structs.MultiRepairResponseWrapper

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	outputJSON, _, err := jsonhttp.CreateAndExecuteJSONRequest("Get Repair", "GET", weburl,
		"", headers, 200, &output)

	if err != nil {
		return structs.RepairResponse{}, err
	}

	log.L.Debugf("Output JSON: %s", outputJSON)
	log.L.Debugf("Output JSON: %+v", output)

	SysID := output.Result[0].SysID

	log.L.Debugf("Output sysID: %+v", SysID)

	return output.Result[0], nil
}
