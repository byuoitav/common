package structs

// const (
// 	New               = 10
// 	Assigned          = 20
// 	WorkInProgress    = 30
// 	Pending           = 40
// 	Closed            = 50
// 	ConfirmedComplete = 60
// )

type Incident struct {
	Rfc         string `json:"rfc,omitempty"`
	Assigned_to string `json:"assigned_to,omitempty"`
	U_email     string `json:"u_email,omitempty"`
	Due_date    string `json:"due_date,omitempty"`

	State            string `json:"state,omitempty"`
	Room             string `json:"u_room,omitempty"`
	Service          string `json:"u_service,omitempty"`
	Department       string `json:"u_department,omitempty"`
	Severity         string `json:"u_severity,omitempty"`
	Sensitivity      string `json:"u_sensitivity,omitempty"`
	WorkStatus       string `json:"u_work_status,omitempty"`
	Reach            string `json:"u_reach,omitempty"`
	AssignmentGroup  string `json:"assignment_group,omitempty"`
	Description      string `json:"description,omitempty"`
	ShortDescription string `json:"short_description,omitempty"`
	CallerId         string `json:"caller_id,omitempty"`
	SysId            string `json:"sys_id,omitempty"`
	Number           string `json:"number,omitempty"`
	InternalNotes    string `json:"work_notes,omitempty"`
	WorkLog          string `json:"comments,omitempty"`

	//The codes necessary for closing a ticket
	ClosureCode       string `json:"u_closure_code,omitempty"`
	ResolutionService string `json:"u_resolution_service,omitempty"`
	ResolutionAction  string `json:"u_action,omitempty"`
}
type IncidentWrapper struct {
	Result Incident `json:"result"`
}

type ReturnRoom struct {
	RoomName string `json:"display_value"`
}

type ReceiveIncident struct {
	Rfc         string `json:"rfc,omitempty"`
	Assigned_to string `json:"assigned_to,omitempty"`
	U_email     string `json:"u_email,omitempty"`
	Due_date    string `json:"due_date,omitempty"`

	State            string      `json:"state,omitempty"`
	Room             ReturnRoom  `json:"u_room,omitempty"`
	Service          string      `json:"u_service,omitempty"`
	Department       string      `json:"u_department,omitempty"`
	Severity         string      `json:"u_severity,omitempty"`
	Sensitivity      string      `json:"u_sensitivity,omitempty"`
	WorkStatus       string      `json:"u_work_status,omitempty"`
	Reach            string      `json:"u_reach,omitempty"`
	AssignmentGroup  AssignGroup `json:"assignment_group,omitempty"`
	Description      string      `json:"description,omitempty"`
	ShortDescription string      `json:"short_description,omitempty"`
	CallerId         CallerID    `json:"caller_id,omitempty"`
	SysId            string      `json:"sys_id,omitempty"`
	Number           string      `json:"number,omitempty"`
	InternalNotes    string      `json:"work_notes,omitempty"`
	WorkLog          string      `json:"comments,omitempty"`
}
type ReceiveIncidentWrapper struct {
	Result ReceiveIncident `json:"result"`
}
type CallerID struct {
	Caller string `json:"display_value"`
}
type AssignGroup struct {
	AssignGroup string `json:"display_value"`
}

type ResolutionCategories struct {
	Result []Category `json:"result"`
}

type Category struct {
	UAction string `json:"u_action"`
	Hint    string `json:"u_attribute_help_text"`
}

type QueriedIncidents struct {
	Result []ReceiveIncident `json:"result"`
}
