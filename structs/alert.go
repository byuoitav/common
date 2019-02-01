package structs

import (
	"time"

	"github.com/byuoitav/common/v2/events"
)

// Alert is a struct that contains the information regarding an alerting event.
type Alert struct {
	events.BasicDeviceInfo

	AlertID string `json:"_id"`

	Type       AlertType     `json:"type"`
	Category   AlertCategory `json:"category"`
	Severity   AlertSeverity `json:"severity"`
	Message    string        `json:"message"`
	Data       interface{}   `json:"data,omitempty"`
	IncidentID string        `json:"incident-id"`

	AlertStartTime      time.Time `json:"start-time"`
	AlertEndTime        time.Time `json:"end-time"`
	AlertLastUpdateTime time.Time `json:"update-time"`

	Respolved      bool           `json:"resolved"`
	Responders     []string       `json:"responders"`
	HelpSentAt     time.Time      `json:"help-sent-at"`
	HelpArrivedAt  time.Time      `json:"help-arrived-at"`
	ResolutionInfo ResolutionInfo `json:"resolution-info"`

	AlertTags  []string `json:"alert-tags"`
	RoomTags   []string `json:"room-tags"`
	DeviceTags []string `json:"device-tags"`
}

// AlertType is an enum of the different types of alerts
type AlertType string

// AlertCategory is an enum of the different categories of alerts
type AlertCategory string

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
