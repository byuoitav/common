package structs

import (
	"strings"
	"time"
	"unicode"
)

// RoomIssue defines information regarding issues in rooms on campus
type RoomIssue struct {
	ID               string     `json:"id"`
	RoomID           string     `json:"roomID"`
	Start            time.Time  `json:"start"`
	End              time.Time  `json:"end"`
	Severity         int        `json:"severity"`
	OpenAlertCount   int        `json:"openAlertCount"`
	ClosedAlertCount int        `json:"closedAlertCount"`
	Alerts           []Alert    `json:"alerts"`
	Resolution       Resolution `json:"resolution"`
	Events           []Event    `json:"events"`
}

// Resolution - how a room issue was resolved
type Resolution struct {
	Code  int    `json:"code"`
	Notes string `json:"notes"`
}

// EventType - an enum for different types of events
type EventType int

const (
	eventTypeAlertStart EventType = iota + 1
	eventTypeAlertEnd
	eventTypeNote
	eventTypePersonSent
	eventTypePersonArrived
	eventTypeChangedSeverity
	eventTypeAcknowledged
)

// Event - the basic event interface
type Event interface {
	Type() EventType
	At() time.Time
	String() string
}

// EventCommon - the information that all events have in common
type EventCommon struct {
	Type EventType `json:"type"`
	At   time.Time `json:"at"`
}

// EventAlert - information for alert related events
type EventAlert struct {
	EventCommon

	AlertID string `json:"alertID"`
}

// EventNote - info for note related events
type EventNote struct {
	EventCommon

	Note string `json:"note"`
}

// EventPerson - info about human interaction events
type EventPerson struct {
	EventCommon

	PersonID   string `json:"personID"`
	PersonName string `json:"personName"`
	PersonLink string `json:"personLink"`
}

// EventChangedSeverity - info about severity changing events
type EventChangedSeverity struct {
	EventCommon

	From int `json:"from"`
	To   int `json:"to"`
}

// EventAcknowledged - info about acknowledging room issues
type EventAcknowledged struct {
	EventCommon

	PersonID   string `json:"personID"`
	PersonName string `json:"personName"`
	PersonLink string `json:"personLink"`
}

// Alert is a struct that contains the information regarding an alerting event.
type Alert struct {
	ID          string    `json:"id"`
	RoomIssueID string    `json:"roomIssueID"`
	DeviceID    string    `json:"deviceID"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	AlertType   string    `json:"alertType"`
	MessageLog  []string  `json:"messageLog"`
}

// TimeToResolve .
func (a Alert) TimeToResolve() string {
	diff := a.End.Sub(a.Start)
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
