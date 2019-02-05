package servicenow

import (
	"fmt"
	"time"

	"github.com/byuoitav/common/jsonhttp"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/structs"
)

func CreateRepair(Alert structs.Alert) (structs.Repair, error) {

	weburl := "https://api.byu.edu:443/domains/servicenow/tableapi/v1/table/u_maint_repair"
	room := fmt.Sprintf("%s %s", Alert.BuildingID, Alert.RoomID)
	assignmentGroup := "AV-Support"
	shortDescription := fmt.Sprintf("%s in room %s-%s has the following alert: %s.", Alert.DeviceID, Alert.BuildingID, Alert.RoomID, Alert.Message)
	description := fmt.Sprintf("%s in room %s-%s has the following alert: %s.", Alert.DeviceID, Alert.BuildingID, Alert.RoomID, Alert.Message)
	internalNotes := fmt.Sprintf("%v", Alert.Data)
	requestOrigination := "On-Site"
	equiptmentReturn := "On-Site"
	year, month, day := time.Now().Date()
	requestdate := fmt.Sprintf("%v-%v-%v", year, int(month), day)
	input := structs.Repair{
		Service:            "TEC Room",
		Room:               room,
		AssignmentGroup:    assignmentGroup,
		ShortDescription:   shortDescription,
		Description:        description,
		InternalNotes:      internalNotes,
		CallerId:           "mjsmith3",
		RequestDate:        requestdate,
		DateNeeded:         "ASAP",
		RequestOrigination: requestOrigination,
		EquiptmentReturn:   equiptmentReturn,
	}
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	var output structs.RepairWrapper
	outputJson, _, err := jsonhttp.CreateAndExecuteJSONRequest("CreateRequest", "POST", weburl,
		input, headers, 20, &output)
	log.L.Debugf("Output JSON: %s", outputJson)
	log.L.Debugf("Output JSON: %+v", output)
	return output.Result[0], err
}

func ModifyRepair(Alert structs.Alert) (structs.RecieveRepair, error) {
	RepairNum := Alert.IncidentID
	SysID, _ := GetRepairID(RepairNum)
	weburl := fmt.Sprintf("https://api.byu.edu:443/domains/servicenow/apiTable/v1/table/u_maint_repair/%s?sysparm_display_value=true", SysID)
	log.L.Debugf("WebURL: %s", weburl)
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
		Description:   "This is a description, want to see what happens",
	}
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	var output structs.ReceiveRepairWrapper
	outputJson, _, err := jsonhttp.CreateAndExecuteJSONRequest("ModifyIncident", "PUT", weburl,
		input, headers, 20, &output)
	log.L.Debugf("Output JSON: %s", outputJson)
	log.L.Debugf("Output JSON: %+v", output)
	return output.Result, err
}

func CloseRepair(Alert structs.Alert) (structs.RecieveRepair, error) {
	RepairNum := Alert.IncidentID
	SysID, _ := GetRepairID(RepairNum)
	weburl := fmt.Sprintf("https://api.byu.edu:443/domains/servicenow/apiTable/v1/table/u_maint_repair/%s?sysparm_display_value=true", SysID)
	log.L.Debugf("WebURL: %s", weburl)
	state := "Closed"

	input := structs.Repair{
		State:   state,
		WorkLog: Alert.ResolutionInfo.Notes,
	}
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	var output structs.ReceiveRepairWrapper
	outputJson, _, err := jsonhttp.CreateAndExecuteJSONRequest("ModifyIncident", "PUT", weburl,
		input, headers, 20, &output)
	log.L.Debugf("Output JSON: %s", outputJson)
	log.L.Debugf("Output JSON: %+v", output)
	return output.Result, err
}

func QueryRepairsByGroup(GroupName string) (structs.QueriedRepairs, error) {
	weburl := fmt.Sprintf("https://api.byu.edu/domains/servicenow/incident/v1.1/u_maint_repair?active=true&assignment_group=%s&sysparm_display_value=true", GroupName)
	log.L.Debugf("WebURL: %s", weburl)
	var output structs.QueriedRepairs
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

func QueryRepairsByRoom(BuildingID string, RoomID string) (structs.QueriedRepairs, error) {
	weburl := fmt.Sprintf("https://api.byu.edu/domains/servicenow/incident/v1.1/u_maint_repair?active=true&sysparm_display_value=true&u_room=%s+%s", BuildingID, RoomID)
	log.L.Debugf("WebURL: %s", weburl)
	var output structs.QueriedRepairs
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

//Get ticket by RPR# return the sysID
func GetRepairID(RepairNum string) (string, error) {
	weburl := fmt.Sprintf("https://api.byu.edu:443/domains/servicenow/tableapi/v1/table/u_maint_repair?sysparm_query=number=%s", RepairNum)
	var output structs.IncidentWrapper
	input := ""
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}
	outputJson, _, err := jsonhttp.CreateAndExecuteJSONRequest("queryUsers", "GET", weburl,
		input, headers, 200, &output)
	log.L.Debugf("Output JSON: %s", outputJson)
	log.L.Debugf("Output JSON: %+v", output)
	SysID := output.Result[0].SysId
	log.L.Debugf("Output sysID: %+v", SysID)
	return SysID, err
}
