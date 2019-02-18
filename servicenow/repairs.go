package servicenow

import (
	"fmt"
	"time"

	"github.com/byuoitav/common/jsonhttp"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/structs"
)

var (
	repairWebURL             = "https://api.byu.edu:443/domains/servicenow/tableapi/v1/table/u_maint_repair"
	repairAssignmentGroup    = "OIT-AV Support"
	repairRequestOrigination = "On-Site"
	repairEquipmentReturn    = "On-Site"
	repairService            = "TEC Room"
)

//CreateRepair is to create a new repair ticket from an alert
func CreateRepair(Alert structs.Alert) (structs.Repair, error) {

	shortDescription := fmt.Sprintf("%s is alerting: %s.", Alert.DeviceID, Alert.Type)

	description := fmt.Sprintf("%s is alerting: %s.", Alert.DeviceID, Alert.Message)

	internalNotes := fmt.Sprintf("%s\nData: %v", Alert.Message, Alert.Data)

	year, month, day := time.Now().Date()

	requestdate := fmt.Sprintf("%v-%v-%v", year, int(month), day)

	input := structs.Repair{
		Service:            repairService,
		Building:           Alert.BuildingID,
		Room:               Alert.RoomID,
		AssignmentGroup:    repairAssignmentGroup,
		ShortDescription:   shortDescription,
		Description:        description,
		InternalNotes:      internalNotes,
		CallerId:           Alert.Requester,
		RequestDate:        requestdate,
		DateNeeded:         "ASAP",
		RequestOrigination: repairRequestOrigination,
		EquipmentReturn:    repairEquipmentReturn,
	}
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	var output structs.RepairWrapper

	outputJSON, _, err := jsonhttp.CreateAndExecuteJSONRequest("CreateRequest", "POST", repairWebURL,
		input, headers, 20, &output)

	log.L.Debugf("Output JSON: %s", outputJSON)
	log.L.Debugf("Output JSON: %+v", output)

	return output.Result[0], err
}

//CreateRepairAsClosedChildClone is used to create a sub-ticket if we close just a portion of the alerts alerting on a room
func CreateRepairAsClosedChildClone(Alert structs.Alert) (structs.Repair, error) {

	room := fmt.Sprintf("%s %s", Alert.BuildingID, Alert.RoomID)

	shortDescription := fmt.Sprintf("%s in room %s-%s has the following alert: %s.", Alert.DeviceID, Alert.BuildingID, Alert.RoomID, Alert.Message)

	description := fmt.Sprintf("%s in room %s-%s has the following alert: %s.", Alert.DeviceID, Alert.BuildingID, Alert.RoomID, Alert.Message)

	internalNotes := fmt.Sprintf("%v", Alert.Data)

	year, month, day := time.Now().Date()

	requestdate := fmt.Sprintf("%v-%v-%v", year, int(month), day)

	input := structs.Repair{
		Service:            repairService,
		Room:               room,
		AssignmentGroup:    repairAssignmentGroup,
		ShortDescription:   shortDescription,
		Description:        description,
		InternalNotes:      internalNotes,
		CallerId:           Alert.Requester,
		RequestDate:        requestdate,
		DateNeeded:         "ASAP",
		RequestOrigination: repairRequestOrigination,
		EquipmentReturn:    repairEquipmentReturn,
	}
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	var output structs.RepairWrapper

	outputJSON, _, err := jsonhttp.CreateAndExecuteJSONRequest("CreateRequest", "POST", repairWebURL,
		input, headers, 20, &output)

	log.L.Debugf("Output JSON: %s", outputJSON)
	log.L.Debugf("Output JSON: %+v", output)

	return output.Result[0], err
}

//ModifyRepair updates the notes on a repair ticket
func ModifyRepair(Alert structs.Alert) (structs.RecieveRepair, error) {
	RepairNum := Alert.IncidentID

	SysID, _ := GetRepairID(RepairNum)

	weburl := fmt.Sprintf("%s/%s?sysparm_display_value=true", repairWebURL, SysID)

	var internalNotes string
	var state string

	//if you want to pull info from the alert
	if Alert.HelpSentAt.IsZero() == false && Alert.HelpArrivedAt.IsZero() == true {
		internalNotes = "Help was was sent at: " + fmt.Sprintf("%s", Alert.HelpSentAt)
		state = "Assigned"
	}

	if Alert.HelpSentAt.IsZero() == false && Alert.HelpArrivedAt.IsZero() == false {
		internalNotes += "\n" + " Help arrived at: " + fmt.Sprintf("%s", Alert.HelpArrivedAt)
		state = "Work In Progress"
	}

	input := structs.Repair{
		State:         state,
		InternalNotes: internalNotes,
	}

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	var output structs.ReceiveRepairWrapper

	outputJSON, _, err := jsonhttp.CreateAndExecuteJSONRequest("ModifyIncident", "PUT", weburl,
		input, headers, 20, &output)

	log.L.Debugf("Output JSON: %s", outputJSON)
	log.L.Debugf("Output JSON: %+v", output)

	return output.Result, err
}

//CloseRepair will just close the repair with the specified resolution info
func CloseRepair(Alert structs.Alert) (structs.RecieveRepair, error) {

	RepairNum := Alert.IncidentID

	SysID, _ := GetRepairID(RepairNum)

	weburl := fmt.Sprintf("%s/%s?sysparm_display_value=true", repairWebURL, SysID)

	state := "Closed"

	input := structs.Repair{
		State:          state,
		WorkLog:        Alert.ResolutionInfo.Notes,
		ResolutionCode: Alert.ResolutionInfo.Code,
	}
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	var output structs.ReceiveRepairWrapper

	outputJSON, _, err := jsonhttp.CreateAndExecuteJSONRequest("ModifyIncident", "PUT", weburl,
		input, headers, 20, &output)

	log.L.Debugf("Output JSON: %s", outputJSON)
	log.L.Debugf("Output JSON: %+v", output)

	return output.Result, err
}

// we don't want this because we want to limit to a specific room and group
// func QueryRepairsByGroup(GroupName string) (structs.QueriedRepairs, error) {
// 	weburl := fmt.Sprintf("https://api.byu.edu/domains/servicenow/incident/v1.1/u_maint_repair?active=true&assignment_group=%s&sysparm_display_value=true", GroupName)
// 	log.L.Debugf("WebURL: %s", weburl)
// 	var output structs.QueriedRepairs
// 	input := ""
// 	headers := map[string]string{
// 		"Authorization": "Bearer " + token,
// 		"Content-Type":  "application/json",
// 	}
// 	outputJSON, _, err := jsonhttp.CreateAndExecuteJSONRequest("querycategory", "GET", weburl,
// 		input, headers, 20, &output)
// 	log.L.Debugf("Output JSON: %s", outputJSON)
// 	log.L.Debugf("Output JSON: %+v", output)
// 	return output, err
// }

// we don't want this because we want to limit to a specific room and group
// func QueryRepairsByRoom(BuildingID string, RoomID string) (structs.QueriedRepairs, error) {
// 	weburl := fmt.Sprintf("https://api.byu.edu/domains/servicenow/incident/v1.1/u_maint_repair?active=true&sysparm_display_value=true&u_room=%s+%s", BuildingID, RoomID)
// 	log.L.Debugf("WebURL: %s", weburl)
// 	var output structs.QueriedRepairs
// 	input := ""
// 	headers := map[string]string{
// 		"Authorization": "Bearer " + token,
// 		"Content-Type":  "application/json",
// 	}
// 	outputJSON, _, err := jsonhttp.CreateAndExecuteJSONRequest("querycategory", "GET", weburl,
// 		input, headers, 20, &output)
// 	log.L.Debugf("Output JSON: %s", outputJSON)
// 	log.L.Debugf("Output JSON: %+v", output)
// 	return output, err
// }

//QueryRepairsByRoomAndGroupName gets a list of repairs by room assigned to specified group
func QueryRepairsByRoomAndGroupName(BuildingID string, RoomID string, GroupName string) (structs.QueriedRepairs, error) {
	weburl := fmt.Sprintf("%s?active=true&sysparm_display_value=true&u_room=%s+%s&assignment_group=%s", repairWebURL, BuildingID, RoomID, GroupName)

	var output structs.QueriedRepairs

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	outputJSON, _, err := jsonhttp.CreateAndExecuteJSONRequest("querycategory", "GET", weburl,
		"", headers, 20, &output)

	log.L.Debugf("Output JSON: %s", outputJSON)
	log.L.Debugf("Output JSON: %+v", output)

	return output, err
}

//GetRepairID Get ticket by RPR# return the sysID
func GetRepairID(RepairNum string) (string, error) {
	weburl := fmt.Sprintf("%s?sysparm_query=number=%s", repairWebURL, RepairNum)

	var output structs.IncidentWrapper

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	outputJSON, _, err := jsonhttp.CreateAndExecuteJSONRequest("queryUsers", "GET", weburl,
		"", headers, 200, &output)

	log.L.Debugf("Output JSON: %s", outputJSON)
	log.L.Debugf("Output JSON: %+v", output)

	SysID := output.Result[0].SysId

	log.L.Debugf("Output sysID: %+v", SysID)

	return SysID, err
}
