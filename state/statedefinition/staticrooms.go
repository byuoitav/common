package statedefinition

import (
	"time"

	"github.com/byuoitav/common/nerr"
)

type StaticRoom struct {
	//information fields
	BuildingID string `json:"buildingID,omitempty"`
	RoomID     string `json:"roomID,omitempty"`

	//State fields
	NotificationsSuppressed *bool `json:"notifications-suppressed,omitempty"`
	Alerting                *bool `json:"alerting,omitempty"`

	LastStateReceived time.Time `json:"last-state-received,omitempty"`
	LastHeartbeat     time.Time `json:"last-heartbeat,omitempty"`
	LastUserInput     time.Time `json:"last-user-input,omitempty"`

	Power string `json:"power,omitempty"`

	//meta fields for Kibana
	ViewDevices           string `json:"view-devices"`
	ViewAlerts            string `json:"view-alerts"`
	EnableNotifications   string `json:"enable-notifications,omitempty"`   //the Hostname - used in a URL
	SuppressNotifications string `json:"suppress-notifications,omitempty"` //the Hostname - used in a URL

	UpdateTimes map[string]time.Time `json:"field-update-times"`
}

//Compare rooms takes two rooms and compares them, changes from new to base will only be included if they have a timestamp in UpdateTimes later than that in base for the same field
func CompareRooms(base, new StaticRoom) (diff, merged StaticRoom, changes bool, err *nerr.E) {

	//information fields
	if new.UpdateTimes["building"].After(base.UpdateTimes["building"]) {
		diff.BuildingID, merged.BuildingID, changes = compareString(base.BuildingID, new.BuildingID, changes)
	}
	if new.UpdateTimes["room"].After(base.UpdateTimes["room"]) {
		diff.RoomID, merged.RoomID, changes = compareString(base.RoomID, new.RoomID, changes)
	}

	//state fields

	//bool fields
	if new.UpdateTimes["notifications-suppressed"].After(base.UpdateTimes["notifications-suppressed"]) {
		diff.NotificationsSuppressed, merged.NotificationsSuppressed, changes = compareBool(base.NotificationsSuppressed, new.NotificationsSuppressed, changes)
	}

	if new.UpdateTimes["alerting"].After(base.UpdateTimes["alerting"]) {
		diff.Alerting, merged.Alerting, changes = compareBool(base.Alerting, new.Alerting, changes)
	}

	//time fields
	if new.UpdateTimes["last-state-receieved"].After(base.UpdateTimes["last-state-receieved"]) {
		diff.LastStateReceived, merged.LastStateReceived, changes = compareTime(base.LastStateReceived, new.LastStateReceived, changes)
	}
	if new.UpdateTimes["last-heartbeat"].After(base.UpdateTimes["last-heartbeat"]) {
		diff.LastHeartbeat, merged.LastHeartbeat, changes = compareTime(base.LastHeartbeat, new.LastHeartbeat, changes)
	}
	if new.UpdateTimes["last-user-input"].After(base.UpdateTimes["last-user-input"]) {
		diff.LastUserInput, merged.LastUserInput, changes = compareTime(base.LastUserInput, new.LastUserInput, changes)
	}

	//string
	if new.UpdateTimes["power"].After(base.UpdateTimes["power"]) {
		diff.Power, merged.Power, changes = compareString(base.Power, new.Power, changes)
	}

	//meta fields
	if new.UpdateTimes["view-devices"].After(base.UpdateTimes["view-devices"]) {
		diff.ViewDevices, merged.ViewDevices, changes = compareString(base.ViewDevices, new.ViewDevices, changes)
	}
	if new.UpdateTimes["view-alerts"].After(base.UpdateTimes["view-alerts"]) {
		diff.ViewAlerts, merged.ViewAlerts, changes = compareString(base.ViewAlerts, new.ViewAlerts, changes)
	}
	if new.UpdateTimes["enable-notifications"].After(base.UpdateTimes["enable-notifications"]) {
		diff.EnableNotifications, merged.EnableNotifications, changes = compareString(base.EnableNotifications, new.EnableNotifications, changes)
	}
	if new.UpdateTimes["suppress-notifications"].After(base.UpdateTimes["suppress-notifications"]) {
		diff.SuppressNotifications, merged.SuppressNotifications, changes = compareString(base.SuppressNotifications, new.SuppressNotifications, changes)
	}

	return
}
