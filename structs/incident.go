package structs

// const (
// 	New               = 10
// 	Assigned          = 20
// 	WorkInProgress    = 30
// 	Pending           = 40
// 	Closed            = 50
// 	ConfirmedComplete = 60
// )

type IncidentRequest struct {
	SysID    string `json:"sys_id,omitempty"`
	Number   string `json:"number,omitempty"`
	Service  string `json:"u_service,omitempty"`
	CallerID string `json:"caller_id,omitempty"`

	AssignmentGroup string `json:"assignment_group,omitempty"`
	State           string `json:"state,omitempty"`

	Room string `json:"u_room,omitempty"`

	ShortDescription string `json:"short_description,omitempty"`

	ContactNumber string `json:"u_work_phone,omitempty"`
	ContactEmail  string `json:"u_email,omitempty"`

	Severity    string `json:"u_severity,omitempty"`
	Reach       string `json:"u_reach,omitempty"`
	WorkStatus  string `json:"u_work_status,omitempty"`
	Sensitivity string `json:"u_sensitivity,omitempty"`

	InternalNotes string `json:"work_notes,omitempty"`
	WorkLog       string `json:"comments,omitempty"`

	//The codes necessary for closing a ticket
	ClosureCode       string `json:"u_closure_code,omitempty"`
	ResolutionService string `json:"u_resolution_service,omitempty"`
	ResolutionAction  string `json:"u_action,omitempty"`
}

type MultiIncidentResponseWrapper struct {
	Result []IncidentResponse `json:"result"`
}

type IncidentResponseWrapper struct {
	Result IncidentResponse `json:"result"`
}

type IncidentResponse struct {
	SysID         string `json:"sys_id,omitempty"`
	Number        string `json:"number,omitempty"`
	InternalNotes string `json:"work_notes,omitempty"`
	WorkLog       string `json:"comments,omitempty"`
}

type QueriedIncidents struct {
	Result []IncidentResponse `json:"result"`
}
