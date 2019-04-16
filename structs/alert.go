package structs

import (
	"strings"
	"time"
	"unicode"

	"github.com/byuoitav/common/v2/events"
)

// RoomIssue .
type RoomIssue struct {
	RoomIssueID string `json:"id,omitempty"`

	events.BasicRoomInfo

	RoomTags []string `json:"room-tags,omitempty"`

	AlertTypes      []AlertType     `json:"alert-types,omitempty"`
	AlertDevices    []string        `json:"alert-devices,omitempty"`
	AlertCategories []AlertCategory `json:"alert-categories,omitempty"`
	AlertSeverities []AlertSeverity `json:"alert-severities,omitempty"`
	AlertCount      int             `json:"alert-count"`

	ActiveAlertTypes      []AlertType     `json:"active-alert-types,omitempty"`
	ActiveAlertDevices    []string        `json:"active-alert-devices,omitempty"`
	ActiveAlertCategories []AlertCategory `json:"active-alert-categories,omitempty"`
	ActiveAlertSeverities []AlertSeverity `json:"active-alert-severities,omitempty"`
	AlertActiveCount      int             `json:"active-alert-count"`

	SystemType string `json:"system-type,omitempty"`

	Source string `json:"-"`

	Alerts []Alert `json:"alerts,omitempty"`

	//Editable fields
	IssueTags []string `json:"issue-tags,omitempty"`

	IncidentID []string `json:"incident-id,omitempty"`

	Notes string `json:"notes,omitempty"`

	RoomIssueResponses []RoomIssueResponse `json:"responses,omitempty"`

	//resolution fields
	Resolved       bool           `json:"resolved"`
	ResolutionInfo ResolutionInfo `json:"resolution-info,omitempty"`

	//notes-log isn't editable
	NotesLog []string `json:"notes-log,omitempty"`
}

//RoomIssueResponse represents information about a tech being dispatched on a room issue
type RoomIssueResponse struct {
	Responders    []Person  `json:"responders,omitempty"`
	HelpSentAt    time.Time `json:"help-sent-at,omitempty"`
	HelpArrivedAt time.Time `json:"help-arrived-at,omitempty"`
}

// Alert is a struct that contains the information regarding an alerting event.
type Alert struct {
	events.BasicDeviceInfo

	AlertID string `json:"id,omitempty,omitempty"`

	Type     AlertType     `json:"type,omitempty"`
	Category AlertCategory `json:"category,omitempty"`
	Severity AlertSeverity `json:"severity,omitempty"`

	Message    string      `json:"message,omitempty"`
	MessageLog []string    `json:"message-log,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	SystemType string      `json:"system-type,omitempty"`

	AlertStartTime      time.Time `json:"start-time,omitempty"`
	AlertEndTime        time.Time `json:"end-time,omitempty"`
	AlertLastUpdateTime time.Time `json:"update-time,omitempty"`

	Active bool `json:"active"`

	AlertTags  []string `json:"alert-tags,omitempty"`
	DeviceTags []string `json:"device-tags,omitempty"`
	RoomTags   []string `json:"room-tags,omitempty"`

	Requester string `json:"requester,omitempty"`

	Source string `json:"-"`

	ManualResolve bool `json:"manual-resolve"`
}

// TimeToResolve .
func (a Alert) TimeToResolve() string {
	diff := a.AlertEndTime.Sub(a.AlertStartTime)
	str := diff.Truncate(time.Second).String()
	ret := strings.Builder{}

	for _, r := range str {
		ret.WriteRune(r)

		if unicode.IsLetter(r) {
			ret.WriteRune(' ')
		}
	}

	return strings.TrimSpace(ret.String())
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
	Critical AlertSeverity = "Critical"
	Warning  AlertSeverity = "Warning"
	Low      AlertSeverity = "Low"
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
