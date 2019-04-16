package structs

import "time"

type RepairRequest struct {
	SysID              string `json:"sys_id,omitempty"`
	Number             string `json:"number,omitempty"`
	Service            string `json:"u_service,omitempty"`
	AssignmentGroup    string `json:"assignment_group,omitempty"`
	RequestOriginator  string `json:"u_requested_by,omitempty"`
	Parent             string `json:"parent,omitempty"`
	RequestOrigination string `json:"u_request_origination,omitempty"`
	RequestDate        string `json:"u_dropped_off_date,omitempty"`
	Building           string `json:"u_building,omitempty"`
	State              string `json:"state,omitempty"`
	Room               string `json:"u_room,omitempty"`
	ShortDescription   string `json:"short_description,omitempty"`
	Description        string `json:"description,omitempty"`

	DateNeeded      string `json:"u_completion_asap_or_date,omitempty"`
	EquipmentReturn string `json:"u_pickup_or_delivery,omitempty"`

	InternalNotes string `json:"work_notes,omitempty"`
	WorkLog       string `json:"comments,omitempty"`
}

type RepairResponse struct {
	SysID         string `json:"sys_id,omitempty"`
	Number        string `json:"number,omitempty"`
	InternalNotes string `json:"work_notes,omitempty"`

	//we are commenting these out because they don't come back consistent (sometimes "", sometimes a struct)
	//and that doesn't play nice with GO
	// Service            ServiceNowLinkValue `json:"u_service,omitempty"`
	// AssignmentGroup    ServiceNowLinkValue `json:"assignment_group,omitempty"`
	// RequestOriginator  string              `json:"u_requested_by,omitempty"`
	// Parent             ServiceNowLinkValue              `json:"parent,omitempty"`
	// RequestOrigination string              `json:"u_request_origination,omitempty"`
	// RequestDate        string              `json:"u_dropped_off_date,omitempty"`
	// Building           ServiceNowLinkValue              `json:"u_building,omitempty"`
	// State              string              `json:"state,omitempty"`
	// Room               ServiceNowLinkValue              `json:"u_room,omitempty"`
	// ShortDescription   string              `json:"short_description,omitempty"`
	// Description        string              `json:"description,omitempty"`

	// DateNeeded      string `json:"u_completion_asap_or_date,omitempty"`
	// EquipmentReturn string `json:"u_pickup_or_delivery,omitempty"`
}

type RepairResponseWrapper struct {
	Result RepairResponse `json:"result"`
}

type MultiRepairResponseWrapper struct {
	Result []RepairResponse `json:"result"`
}

// ClassHalfHourBlock defines the information needed for a 30 min block of a class
type ClassHalfHourBlock struct {
	BlockStart     string    `json:"block-start"`
	ClassName      string    `json:"class-name"`
	ClassTime      string    `json:"class-time"`
	Teacher        Person    `json:"teacher"`
	Days           string    `json:"days"`
	ClassStartTime time.Time `json:"class-start-time"`
	ClassEndTime   time.Time `json:"class-end-time"`
}
