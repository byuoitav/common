package servicenow

import (
	"fmt"
	"strings"

	"github.com/byuoitav/common/jsonhttp"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/structs"
)

const (
	RepairWebURL             = "https://api.byu.edu:443/domains/servicenow/tableapi/v1/table/u_maint_repair"
	RepairModifyWebURL       = "https://api.byu.edu:443/domains/servicenow/apiTable/v1/table/u_maint_repair"
	RepairAssignmentGroup    = "OIT-AV Support"
	RepairRequestOrigination = "On-Site"
	RepairEquipmentReturn    = "On-Site"
	RepairService            = "TEC Room"
	RepairDefaultRequestor   = "AV Metrics Web Service"
	RepairDateNeeded         = "ASAP"
	RepairClosedState        = "Closed Complete"
)

//CreateRepair is to create a new repair ticket from a new room issue
func CreateRepair(input structs.RepairRequest) (structs.RepairResponse, error) {

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	var output structs.RepairResponseWrapper

	outputJSON, _, err := jsonhttp.CreateAndExecuteJSONRequest("CreateRequest", "POST", RepairWebURL,
		input, headers, 20, &output)

	log.L.Debugf("Output JSON: %s", outputJSON)
	log.L.Debugf("Output JSON: %+v", output)

	return output.Result, err
}

//ModifyRepair updates the notes on a repair ticket
func ModifyRepair(input structs.RepairRequest, id string) (structs.RepairResponse, error) {

	weburl := fmt.Sprintf("%s/%s", RepairModifyWebURL, id)

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

//QueryRepairsByRoom gets a list of repairs by room assigned to specified group
func QueryRepairsByRoom(RoomID string) ([]structs.RepairResponse, error) {
	roomIDreplaced := strings.Replace(RoomID, "-", "+", -1)
	GroupName := strings.Replace(RepairAssignmentGroup, " ", "+", -1)

	weburl := fmt.Sprintf("active=true&sysparm_display_value=true&u_room=%s&assignment_group=%s", roomIDreplaced, GroupName)
	weburl = fmt.Sprintf("%s?%s", RepairWebURL, weburl)

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
	weburl := fmt.Sprintf("%s?sysparm_query=number=%s&sysparm_display_value=true", RepairWebURL, RepairNum)

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
