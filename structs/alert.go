package structs

import (
	"time"
)

// Alert is a struct that contains the information regarding an alerting event.
type Alert struct {
	BuildingID     string         `json:"buildingID"`
	RoomID         string         `json:"roomID"`
	DeviceID       string         `json:"deviceID"`
	Type           AlertType      `json:"type"`
	IncidentID     string         `json:"incident-id"`
	Severity       AlertSeverity  `json:"severity"`
	Responders     []string       `json:"responders"`
	HelpSentAt     time.Time      `json:"help-sent-at"`
	HelpArrivedAt  time.Time      `json:"help-arrived-at"`
	ResolutionInfo ResolutionInfo `json:"resolution-info"`
	AlertTags      []string       `json:"alert-tags"`
	RoomTags       []string       `json:"room-tags"`
	DeviceTags     []string       `json:"device-tags"`
}

// AlertType is an enum of the different types of alerts
type AlertType string

// Here is a list of AlertTypes
const (
	System AlertType = "system"
	User   AlertType = "user"
)

// AlertSeverity is an enum of the different levels of severity for alerts
type AlertSeverity int

// Here is a list of AlertSeverities
const (
	Critical AlertSeverity = iota + 1
	Warning
)

// ResolutionInfo is a struct that contains the information about the resolution of the alert
type ResolutionInfo struct {
	Code       string    `json:"resolution-code"`
	Notes      string    `json:"notes"`
	ResolvedAt time.Time `json:"resolved-at"`
}
