package structs

type Repair struct {
	Assigned_to        string `json:"assigned_to,omitempty"`
	CallerId           string `json:"caller_id,omitempty"`
	AssignmentGroup    string `json:"assignment_group,omitempty"`
	Description        string `json:"description,omitempty"`
	ShortDescription   string `json:"short_description,omitempty"`
	SysId              string `json:"sys_id,omitempty"`
	Number             string `json:"number,omitempty"`
	RequestOrigination string `json:"u_request_origination,omitempty"`
	DateNeeded         string `json:"u_completion_asap_or_date,omitempty"`
	RequestDate        string `json:"u_dropped_off_date,omitempty"`
	Service            string `json:"u_service,omitempty"`
	EquiptmentReturn   string `json:"u_pickup_or_delivery,omitempty"`
	InternalNotes      string `json:"work_notes,omitempty"`
	WorkLog            string `json:"comments,omitempty"`
	State              string `json:"state,omitempty"`
	Room               string `json:"u_room,omitempty"`
}

type RepairWrapper struct {
	Result []Repair `json:"result"`
}

type RecieveRepair struct {
	Assigned_to        string      `json:"assigned_to,omitempty"`
	CallerId           CallerID    `json:"caller_id,omitempty"`
	AssignmentGroup    AssignGroup `json:"assignment_group,omitempty"`
	Description        string      `json:"description,omitempty"`
	ShortDescription   string      `json:"short_description,omitempty"`
	SysId              string      `json:"sys_id,omitempty"`
	Number             string      `json:"number,omitempty"`
	RequestOrigination string      `json:"u_request_origination,omitempty"`
	DateNeeded         string      `json:"u_completion_asap_or_date,omitempty"`
	RequestDate        string      `json:"u_dropped_off_date,omitempty"`
	Service            string      `json:"u_service,omitempty"`
	EquiptmentReturn   string      `json:"u_pickup_or_delivery,omitempty"`
	InternalNotes      string      `json:"work_notes,omitempty"`
	WorkLog            string      `json:"comments,omitempty"`
	State              string      `json:"state,omitempty"`
	Room               ReturnRoom  `json:"u_room,omitempty"`
}

type ReceiveRepairWrapper struct {
	Result RecieveRepair `json:"result"`
}

type QueriedRepairs struct {
	Result []RecieveRepair `json:"result"`
}
