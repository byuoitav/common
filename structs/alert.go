package structs

import (
	"time"

	"github.com/byuoitav/common/v2/events"
)

// RoomIssue .
type RoomIssue struct {
	RoomIssueID string `json:"id"`

	events.BasicRoomInfo

	Severity AlertSeverity `json:"severity"`

	RoomTags []string `json:"room-tags"`

	AlertTypes      []AlertType     `json:"alert-types"`
	AlertDevices    []string        `json:"alert-devices"`
	AlertCategories []AlertCategory `json:"alert-categories"`
	AlertCount      int             `json:"alert-count"`

	ActiveAlertTypes      []AlertType     `json:"active-alert-types"`
	ActiveAlertDevices    []string        `json:"active-alert-devices"`
	ActiveAlertCategories []AlertCategory `json:"active-alert-categories"`
	AlertActiveCount      int             `json:"active-alert-count"`

	SystemType string `json:"system-type"`

	Source string `json:"-"`

	Alerts []Alert `json:"alerts"`

	//Editable fields
	IssueTags []string `json:"issue-tags"`

	IncidentID string `json:"incident-id"`

	Notes string `json:"notes"`

	Responders    []string  `json:"responders"`
	HelpSentAt    time.Time `json:"help-sent-at"`
	HelpArrivedAt time.Time `json:"help-arrived-at"`

	//resolution fields
	Resolved       bool           `json:"resolved"`
	ResolutionInfo ResolutionInfo `json:"resolution-info"`

	//notes-log isn't editable
	NotesLog []string `json:"notes-log"`
}

// Alert is a struct that contains the information regarding an alerting event.
type Alert struct {
	events.BasicDeviceInfo

	AlertID string `json:"id,omitempty"`

	Type     AlertType     `json:"type"`
	Category AlertCategory `json:"category"`
	Severity AlertSeverity `json:"severity"`

	Message    string      `json:"message"`
	MessageLog []string    `json:"message-log"`
	Data       interface{} `json:"data,omitempty"`
	SystemType string      `json:"system-type"`

	AlertStartTime      time.Time `json:"start-time"`
	AlertEndTime        time.Time `json:"end-time"`
	AlertLastUpdateTime time.Time `json:"update-time"`

	Active bool `json:"active"`

	AlertTags  []string `json:"alert-tags"`
	DeviceTags []string `json:"device-tags"`
	RoomTags   []string `json:"room-tags"`

	Requester string `json:"requester,omitempty"`

	Source string `json:"-"`

	ManualResolve bool `json:"manual-resolve"`
}

// AlertType is an enum of the different types of alerts
type AlertType string

const (
	Communication AlertType = "communication"
	Heartbeat     AlertType = "heartbeat"
)

// AlertCategory is an enum of the different categories of alerts
type AlertCategory string

// Here is a list of AlertCategory
const (
	System AlertCategory = "system"
	User   AlertCategory = "user"
)

// AlertSeverity is an enum of the different levels of severity for alerts
type AlertSeverity string

// Here is a list of AlertSeverities
const (
	Critical AlertSeverity = "critical"
	Warning  AlertSeverity = "warning"
	Low      AlertSeverity = "low"
)

var AlertSeverities = [...]string{
	"critical",
	"warning",
	"low",
}

// ResolutionInfo is a struct that contains the information about the resolution of the alert
type ResolutionInfo struct {
	Code       string    `json:"resolution-code"`
	Notes      string    `json:"notes"`
	ResolvedAt time.Time `json:"resolved-at"`
}

func ContainsAllTags(tagList []string, tags ...string) bool {
	for i := range tags {
		hasTag := false

		for j := range tagList {
			if tagList[j] == tags[i] {
				hasTag = true
				continue
			}
		}

		if !hasTag {
			return false
		}
	}

	return true
}

func AddToTags(tagList []string, tags ...string) []string {
	for _, t := range tags {
		if !ContainsAllTags(tagList, t) {
			tagList = append(tagList, t)
		}
	}
	return tagList
}

func ContainsAnyTags(tagList []string, tags ...string) bool {
	for i := range tags {
		for j := range tagList {
			if tagList[j] == tags[i] {
				return true
			}
		}
	}

	return false
}

func (r *RoomIssue) CalculateAggregateInfo() {
	r.AlertTypes = []AlertType{}
	r.ActiveAlertTypes = []AlertType{}
	r.AlertCategories = []AlertCategory{}
	r.ActiveAlertCategories = []AlertCategory{}

	activeCount := 0

	for i := range r.Alerts {

		//active alert stuff
		if r.Alerts[i].Active {
			activeCount++

			found := false
			for j := range r.ActiveAlertTypes {
				if r.ActiveAlertTypes[j] == r.Alerts[i].Type {
					found = true
					break
				}
			}
			if !found {
				r.ActiveAlertTypes = append(r.ActiveAlertTypes, r.Alerts[i].Type)
			}

			found = false
			for j := range r.ActiveAlertCategories {
				if r.ActiveAlertCategories[j] == r.Alerts[i].Category {
					found = true
					break
				}

			}
			if !found {
				r.ActiveAlertCategories = append(r.ActiveAlertCategories, r.Alerts[i].Category)
			}

			found = false
			for j := range r.ActiveAlertDevices {
				if r.ActiveAlertDevices[j] == r.Alerts[i].DeviceID {
					found = true
					break
				}

			}
			if !found {
				r.ActiveAlertDevices = append(r.ActiveAlertDevices, r.Alerts[i].DeviceID)
			}
		}

		found := false
		for j := range r.AlertTypes {
			if r.AlertTypes[j] == r.Alerts[i].Type {
				found = true
				break
			}
		}
		if !found {
			r.AlertTypes = append(r.AlertTypes, r.Alerts[i].Type)
		}

		found = false
		for j := range r.AlertCategories {
			if r.AlertCategories[j] == r.Alerts[i].Category {
				found = true
				break
			}

		}
		if !found {
			r.AlertCategories = append(r.AlertCategories, r.Alerts[i].Category)
		}

		found = false
		for j := range r.AlertDevices {
			if r.AlertDevices[j] == r.Alerts[i].DeviceID {
				found = true
				break
			}

		}
		if !found {
			r.AlertDevices = append(r.AlertDevices, r.Alerts[i].DeviceID)
		}
	}
	r.AlertActiveCount = activeCount
	r.AlertCount = len(r.Alerts)
}
