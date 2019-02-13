package servicenow

import (
	"fmt"
	"net/http"
	"os"

	"github.com/byuoitav/common"
	"github.com/byuoitav/common/jsonhttp"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/structs"
)

var token = os.Getenv("WSO2_TOKEN")

func main() {
	log.SetLevel("debug")
	port := ":8025"
	router := common.NewRouter()
	//Create incident test
	// TestAlert := structs.Alert{
	// 	BuildingID: "ITB",
	// 	RoomID:     "1108",
	// 	DeviceID:   "ITB-1108-CP5",
	// 	Message:    "There is an issue with the pi, it is not turning on",
	// 	Data:       "Stuff",
	// }
	// CreateIncident(TestAlert)

	//Modify incident test
	// ModifyAlert := structs.Alert{
	// 	HelpSentAt:    time.Now(),
	// 	HelpArrivedAt: time.Now().Add(5),
	// }
	// SysID := "89233ae61bdb674003e68622dd4bcb1b"
	// ModifyIncident(SysID, ModifyAlert)

	//query incident resolution category (for closing tickets) test
	// table := "u_inc_resolution_cat"
	// PrintIncidentResolutionCategory(table)

	//close incident test
	// sysID := "89233ae61bdb674003e68622dd4bcb1b"
	// resolutionaction := "Replaced"
	// notes := "I replaced the pi and the room is working now"
	// CloseIncident(sysID, resolutionaction, notes)

	//test Query all incidents for AV-Support
	// GroupName := "AV-Support"
	// QueryIncidentsByGroup(GroupName)

	//test query by room
	// BuildingID := "ITB"
	// RoomID := "1108"
	// QueryIncidentsByRoom(BuildingID, RoomID)

	//Query all users
	QueryAllUsers()
	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}
	router.StartServer(&server)

}

func CreateIncident(Alert structs.Alert) (structs.Incident, error) {

	weburl := "https://api.byu.edu/domains/servicenow/incident/v1.1/incident"
	room := fmt.Sprintf("%s %s", Alert.BuildingID, Alert.RoomID)
	workStatus := "Very Low"
	sensitivity := "Very Low"
	severity := "Very Low"
	reach := "Very Low"
	assignmentGroup := "AV-Support"
	shortDescription := fmt.Sprintf("%s in room %s-%s has the following alert: %s.", Alert.DeviceID, Alert.BuildingID, Alert.RoomID, Alert.Message)
	description := fmt.Sprintf("%s in room %s-%s has the following alert: %s.", Alert.DeviceID, Alert.BuildingID, Alert.RoomID, Alert.Message)
	internalNotes := fmt.Sprintf("%v", Alert.Data)
	input := structs.Incident{
		Service:          "TEC Room",
		Room:             room,
		WorkStatus:       workStatus,
		Sensitivity:      sensitivity,
		Severity:         severity,
		Reach:            reach,
		AssignmentGroup:  assignmentGroup,
		ShortDescription: shortDescription,
		Description:      description,
		InternalNotes:    internalNotes,
		CallerId:         "mjsmith3",
	}
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	var output structs.IncidentWrapper
	outputJson, _, err := jsonhttp.CreateAndExecuteJSONRequest("CreateRequest", "POST", weburl,
		input, headers, 20, &output)
	log.L.Debugf("Output JSON: %s", outputJson)
	log.L.Debugf("Output JSON: %+v", output)
	return output.Result[0], err
}

//we need to be able to access the sysID of the incident Ticket
//TO DO: takes incident ID and string for internal notes
func ModifyIncident(Alert structs.Alert) (structs.ReceiveIncident, error) {
	IncidentNumber := Alert.IncidentID
	SysID, _ := GetSysID(IncidentNumber)
	weburl := fmt.Sprintf("https://api.byu.edu/domains/servicenow/incident/v1.1/incident/%s?sysparm_display_value=true", SysID)
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

	input := structs.Incident{
		State:         state,
		InternalNotes: internalNotes,
		Description:   "This is a description, want to see what happens",
	}
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	var output structs.ReceiveIncidentWrapper
	outputJson, _, err := jsonhttp.CreateAndExecuteJSONRequest("ModifyIncident", "PUT", weburl,
		input, headers, 20, &output)
	log.L.Debugf("Output JSON: %s", outputJson)
	log.L.Debugf("Output JSON: %+v", output)
	return output.Result, err

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

func CloseIncident(Alert structs.Alert) (structs.ReceiveIncident, error) {
	IncidentID := Alert.IncidentID
	SysID, _ := GetSysID(IncidentID)
	weburl := fmt.Sprintf("https://api.byu.edu/domains/servicenow/incident/v1.1/incident/%s?sysparm_display_value=true", SysID)
	log.L.Debugf("WebURL: %s", weburl)
	state := "Closed"
	closurecode := "Resolved"
	resolutionservice := "TEC Room"

	input := structs.Incident{
		State:             state,
		ClosureCode:       closurecode,
		ResolutionService: resolutionservice,
		ResolutionAction:  Alert.ResolutionInfo.Code,
		WorkLog:           Alert.ResolutionInfo.Notes,
	}
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	var output structs.ReceiveIncidentWrapper
	outputJson, _, err := jsonhttp.CreateAndExecuteJSONRequest("ModifyIncident", "PUT", weburl,
		input, headers, 20, &output)
	log.L.Debugf("Output JSON: %s", outputJson)
	log.L.Debugf("Output JSON: %+v", output)
	return output.Result, err
}

//query all incidents for a given assignment group
func QueryIncidentsByGroup(GroupName string) (structs.QueriedIncidents, error) {
	weburl := fmt.Sprintf("https://api.byu.edu/domains/servicenow/incident/v1.1/incident?active=true&assignment_group=%s&sysparm_display_value=true", GroupName)
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
	return output, err
}

//query all incidents by room number
func QueryIncidentsByRoom(BuildingID string, RoomID string) (structs.QueriedIncidents, error) {
	weburl := fmt.Sprintf("https://api.byu.edu/domains/servicenow/incident/v1.1/incident?active=true&sysparm_display_value=true&u_room=%s+%s", BuildingID, RoomID)
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
	return output, err
}

//TODO query all of the users in the system (Net_id)

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

//Get ticket by INC# return the sysID
func GetSysID(IncidentNumber string) (string, error) {
	weburl := fmt.Sprintf("https://api.byu.edu:443/domains/servicenow/tableapi/v1/table/incident?sysparm_query=number=%s", IncidentNumber)
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
