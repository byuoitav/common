package structs

import (
	"time"

	"github.com/byuoitav/common/v2/events"
)

// RoomIssue .
type RoomIssue struct {
	RoomIssueID string `json:"id"`

	events.BasicRoomInfo

	RoomTags []string `json:"room-tags"`

	AlertTypes      []AlertType     `json:"alert-types"`
	AlertDevices    []string        `json:"alert-devices"`
	AlertCategories []AlertCategory `json:"alert-categories"`
	AlertSeverities []AlertSeverity `json:"alert-severities"`
	AlertCount      int             `json:"alert-count"`

	ActiveAlertTypes      []AlertType     `json:"active-alert-types"`
	ActiveAlertDevices    []string        `json:"active-alert-devices"`
	ActiveAlertCategories []AlertCategory `json:"active-alert-categories"`
	ActiveAlertSeverities []AlertSeverity `json:"active-alert-severities"`
	AlertActiveCount      int             `json:"active-alert-count"`

	SystemType string `json:"system-type"`

	Source string `json:"-"`

	Alerts []Alert `json:"alerts"`

	//Editable fields
	IssueTags []string `json:"issue-tags"`

	IncidentID []string `json:"incident-id"`

	Notes string `json:"notes"`

	RoomIssueResponses []RoomIssueResponse `json:"responses"`

	//resolution fields
	Resolved       bool           `json:"resolved"`
	ResolutionInfo ResolutionInfo `json:"resolution-info"`

	//notes-log isn't editable
	NotesLog []string `json:"notes-log"`
}

//RoomIssueResponse represents information about a tech being dispatched on a room issue
type RoomIssueResponse struct {
	Responders    []string  `json:"responders"`
	HelpSentAt    time.Time `json:"help-sent-at"`
	HelpArrivedAt time.Time `json:"help-arrived-at"`
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

var AlertSeverities = []AlertSeverity{
	Critical,
	Warning,
	Low,
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

func AddToSeverity(list []AlertSeverity, toAdd ...AlertSeverity) []AlertSeverity {
	for i := range toAdd {
		found := false
		for j := range list {
			if toAdd[i] == list[j] {
				found = true
				break
			}
		}
		if !found {
			list = append(list, toAdd[i])
		}
	}

	return list
}

func AddToType(list []AlertType, toAdd ...AlertType) []AlertType {
	for i := range toAdd {
		found := false
		for j := range list {
			if toAdd[i] == list[j] {
				found = true
				break
			}
		}
		if !found {
			list = append(list, toAdd[i])
		}
	}

	return list
}

func AddToCategory(list []AlertCategory, toAdd ...AlertCategory) []AlertCategory {
	for i := range toAdd {
		found := false
		for j := range list {
			if toAdd[i] == list[j] {
				found = true
				break
			}
		}
		if !found {
			list = append(list, toAdd[i])
		}
	}

	return list
}

func (r *RoomIssue) CalculateAggregateInfo() {
	r.AlertTypes = []AlertType{}
	r.ActiveAlertTypes = []AlertType{}

	r.AlertCategories = []AlertCategory{}
	r.ActiveAlertCategories = []AlertCategory{}

	r.AlertSeverities = []AlertSeverity{}
	r.ActiveAlertSeverities = []AlertSeverity{}

	r.AlertDevices = []string{}
	r.ActiveAlertDevices = []string{}

	activeCount := 0

	for i := range r.Alerts {

		//active alert stuff
		if r.Alerts[i].Active {
			activeCount++
			r.ActiveAlertDevices = AddToTags(r.ActiveAlertDevices, r.Alerts[i].DeviceID)
			r.ActiveAlertTypes = AddToType(r.ActiveAlertTypes, r.Alerts[i].Type)
			r.ActiveAlertCategories = AddToCategory(r.ActiveAlertCategories, r.Alerts[i].Category)
			r.ActiveAlertSeverities = AddToSeverity(r.ActiveAlertSeverities, r.Alerts[i].Severity)
		}

		r.AlertDevices = AddToTags(r.AlertDevices, r.Alerts[i].DeviceID)
		r.AlertTypes = AddToType(r.AlertTypes, r.Alerts[i].Type)
		r.AlertCategories = AddToCategory(r.AlertCategories, r.Alerts[i].Category)
		r.AlertSeverities = AddToSeverity(r.AlertSeverities, r.Alerts[i].Severity)
	}
	r.AlertActiveCount = activeCount
	r.AlertCount = len(r.Alerts)
}
