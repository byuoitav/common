package statedefinition

import (
	"time"

	"github.com/byuoitav/common/nerr"
)

//*************************
//IMPORTANT - if you add fields to this struct be sure to change the CompareDevices function
//*************************
type StaticDevice struct {
	//common fields
	DeviceID                string           `json:"deviceID,omitempty"`
	Alerting                *bool            `json:"alerting,omitempty"`
	Alerts                  map[string]Alert `json:"alerts,omitempty"`
	NotificationsSuppressed *bool            `json:"notifications-suppressed,omitempty"`
	Building                string           `json:"building,omitempty"`
	Room                    string           `json:"room,omitempty"`
	Hostname                string           `json:"hostname,omitempty"`
	LastStateReceived       time.Time        `json:"last-state-received,omitempty"`
	LastHeartbeat           time.Time        `json:"last-heartbeat,omitempty"`
	LastUserInput           time.Time        `json:"last-user-input,omitempty"`

	DeviceType  string `json:"device-type,omitempty"`
	DeviceClass string `json:"device-class,omitempty"`
	DeviceName  string `json:"device-name,omitempty"`

	Tags []string `json:"tags,omitempty"`

	//semi-common fields LastHeartbeat time.Time `json:"last-heartbeat,omitempty"` LastUserInput time.Time `json:"last-user-input,omitempty"`
	Power string `json:"power,omitempty"`

	//Control Processor Specific Fields
	Websocket      string `json:"websocket,omitempty"`
	WebsocketCount *int   `json:"websocket-count,omitempty"`

	//Display Specific Fields
	Blanked *bool  `json:"blanked,omitempty"`
	Input   string `json:"input,omitempty"`

	//Audio Device Specific Fields
	Muted  *bool `json:"muted,omitempty"`
	Volume *int  `json:"volume,omitempty"`

	//Fields specific to Microphones
	BatteryChargeBars         *int   `json:"battery-charge-bars,omitempty"`
	BatteryChargeMinutes      *int   `json:"battery-charge-minutes,omitempty"`
	BatteryChargePercentage   *int   `json:"battery-charge-percentage,omitempty"`
	BatteryChargeHoursMinutes *int   `json:"battery-charge-hours-minutes,omitempty"`
	BatteryCycles             *int   `json:"battery-cycles,omitempty"`
	BatteryType               string `json:"battery-type,omitempty"`
	Interference              string `json:"intererence,omitempty"`

	//meta fields for use in kibana
	Control               string `json:"control,omitempty"`                //the id - used in a URL
	EnableNotifications   string `json:"enable-notifications,omitempty"`   //the id - used in a URL
	SuppressNotifications string `json:"suppress-notifications,omitempty"` //the id - used in a URL
	ViewDashboard         string `json:"ViewDashboard,omitempty"`          //the id - used in a URL

	UpdateTimes map[string]time.Time `json:"field-state-received"`
}

//CompareDevices takes a base devices, and calculates the difference between the two, returning it in the staticDevice return value. Bool denotes if there were any differences
func CompareDevices(base, new StaticDevice) (diff StaticDevice, merged StaticDevice, changes bool, err *nerr.E) {

	//base is our base
	merged = base

	//common fields
	if new.UpdateTimes["deviceID"].After(base.UpdateTimes["deviceID"]) {
		diff.DeviceID, merged.DeviceID, changes = compareString(base.DeviceID, new.DeviceID, changes)
	}
	if new.UpdateTimes["alerting"].After(base.UpdateTimes["alerting"]) {
		diff.Alerting, merged.Alerting, changes = compareBool(base.Alerting, new.Alerting, changes)
	}

	//handle alerts special case - it's alerts.<name>
	diff.Alerts, merged.Alerts, changes = compareAlerts(base.Alerts, new.Alerts, base.UpdateTimes, new.UpdateTimes, changes)

	if new.UpdateTimes["notifications-suppressed"].After(base.UpdateTimes["notifications-suppressed"]) {
		diff.NotificationsSuppressed, merged.NotificationsSuppressed, changes = compareBool(base.NotificationsSuppressed, new.NotificationsSuppressed, changes)
	}
	if new.UpdateTimes["building"].After(base.UpdateTimes["building"]) {
		diff.Building, merged.Building, changes = compareString(base.Building, new.Building, changes)
	}
	if new.UpdateTimes["room"].After(base.UpdateTimes["room"]) {
		diff.Room, merged.Room, changes = compareString(base.Room, new.Room, changes)
	}
	if new.UpdateTimes["hostname"].After(base.UpdateTimes["hostname"]) {
		diff.Hostname, merged.Hostname, changes = compareString(base.Hostname, new.Hostname, changes)
	}

	if new.UpdateTimes["device-name"].After(base.UpdateTimes["device-name"]) {
		diff.DeviceName, merged.DeviceName, changes = compareString(base.DeviceName, new.DeviceName, changes)
	}
	if new.UpdateTimes["device-class"].After(base.UpdateTimes["device-class"]) {
		diff.DeviceClass, merged.DeviceClass, changes = compareString(base.DeviceClass, new.DeviceClass, changes)
	}

	if new.UpdateTimes["device-type"].After(base.UpdateTimes["device-type"]) {
		diff.DeviceType, merged.DeviceType, changes = compareString(base.DeviceType, new.DeviceType, changes)
	}

	if new.UpdateTimes["tags"].After(base.UpdateTimes["tags"]) {
		diff.Tags, merged.Tags, changes = compareTags(base.Tags, new.Tags, changes)

	}

	//semi-common fields
	if new.UpdateTimes["last-heartbeat"].After(base.UpdateTimes["last-heartbeat"]) {
		diff.LastHeartbeat, merged.LastHeartbeat, changes = compareTime(base.LastHeartbeat, new.LastHeartbeat, changes)
	}
	if new.UpdateTimes["last-user-input"].After(base.UpdateTimes["last-user-input"]) {
		diff.LastUserInput, merged.LastUserInput, changes = compareTime(base.LastUserInput, new.LastUserInput, changes)
	}
	if new.UpdateTimes["last-state-received"].After(base.UpdateTimes["last-state-received"]) {
		diff.LastStateReceived, merged.LastStateReceived, changes = compareTime(base.LastStateReceived, new.LastStateReceived, changes)
	}
	if new.UpdateTimes["power"].After(base.UpdateTimes["power"]) {
		diff.Power, merged.Power, changes = compareString(base.Power, new.Power, changes)
	}

	//Conrol processor specific fields
	if new.UpdateTimes["websocket"].After(base.UpdateTimes["websocket"]) {
		diff.Websocket, merged.Websocket, changes = compareString(base.Websocket, new.Websocket, changes)
	}
	if new.UpdateTimes["websocket-count"].After(base.UpdateTimes["websocket-count"]) {
		diff.WebsocketCount, merged.WebsocketCount, changes = compareInt(base.WebsocketCount, new.WebsocketCount, changes)
	}

	//Display specific fields
	if new.UpdateTimes["blanked"].After(base.UpdateTimes["blanked"]) {
		diff.Blanked, merged.Blanked, changes = compareBool(base.Blanked, new.Blanked, changes)
	}
	if new.UpdateTimes["input"].After(base.UpdateTimes["input"]) {
		diff.Input, merged.Input, changes = compareString(base.Input, new.Input, changes)
	}
	//Audio Device specific fields
	if new.UpdateTimes["muted"].After(base.UpdateTimes["muted"]) {
		diff.Muted, merged.Muted, changes = compareBool(base.Muted, new.Muted, changes)
	}
	if new.UpdateTimes["volume"].After(base.UpdateTimes["volume"]) {
		diff.Volume, merged.Volume, changes = compareInt(base.Volume, new.Volume, changes)
	}

	//Microphone specific fields
	if new.UpdateTimes["battery-charge-bars"].After(base.UpdateTimes["battery-charge-bars"]) {
		diff.BatteryChargeBars, merged.BatteryChargeBars, changes = compareInt(base.BatteryChargeBars, new.BatteryChargeBars, changes)
	}
	if new.UpdateTimes["battery-charge-minutes"].After(base.UpdateTimes["battery-charge-minutes"]) {
		diff.BatteryChargeMinutes, merged.BatteryChargeMinutes, changes = compareInt(base.BatteryChargeMinutes, new.BatteryChargeMinutes, changes)
	}
	if new.UpdateTimes["battery-charge-percentage"].After(base.UpdateTimes["battery-charge-percentage"]) {
		diff.BatteryChargePercentage, merged.BatteryChargePercentage, changes = compareInt(base.BatteryChargePercentage, new.BatteryChargePercentage, changes)
	}
	if new.UpdateTimes["battery-chage-hours-minutes"].After(base.UpdateTimes["battery-chage-hours-minutes"]) {
		diff.BatteryChargeHoursMinutes, merged.BatteryChargeHoursMinutes, changes = compareInt(base.BatteryChargeHoursMinutes, new.BatteryChargeHoursMinutes, changes)
	}
	if new.UpdateTimes["battery-cycles"].After(base.UpdateTimes["battery-cycles"]) {
		diff.BatteryCycles, merged.BatteryCycles, changes = compareInt(base.BatteryCycles, new.BatteryCycles, changes)
	}
	if new.UpdateTimes["battery-type"].After(base.UpdateTimes["battery-type"]) {
		diff.BatteryType, merged.BatteryType, changes = compareString(base.BatteryType, new.BatteryType, changes)
	}
	if new.UpdateTimes["interference"].After(base.UpdateTimes["interference"]) {
		diff.Interference, merged.Interference, changes = compareString(base.Interference, new.Interference, changes)
	}

	//meta fields
	if new.UpdateTimes["control"].After(base.UpdateTimes["control"]) {
		diff.Control, merged.Control, changes = compareString(base.Control, new.Control, changes)
	}
	if new.UpdateTimes["enable-notifications"].After(base.UpdateTimes["enable-notifications"]) {
		diff.EnableNotifications, merged.EnableNotifications, changes = compareString(base.EnableNotifications, new.EnableNotifications, changes)
	}
	if new.UpdateTimes["suppress-notifications"].After(base.UpdateTimes["suppress-notifications"]) {
		diff.SuppressNotifications, merged.SuppressNotifications, changes = compareString(base.SuppressNotifications, new.SuppressNotifications, changes)
	}
	if new.UpdateTimes["view-dashboard"].After(base.UpdateTimes["view-dashboard"]) {
		diff.ViewDashboard, merged.ViewDashboard, changes = compareString(base.ViewDashboard, new.ViewDashboard, changes)
	}

	return
}

func compareString(base, new string, changes bool) (string, string, bool) {
	if new != "" {
		if base != new {
			return new, new, true
		}
	}
	return "", base, false || changes
}

func compareBool(base, new *bool, changes bool) (*bool, *bool, bool) {
	if new != nil {
		if base == nil || *base != *new {
			return new, new, true
		}
	}
	return nil, base, false || changes
}

func compareInt(base, new *int, changes bool) (*int, *int, bool) {
	if new != nil {
		if base == nil || *base != *new {
			return new, new, true
		}
	}
	return nil, base, false || changes
}

func compareTime(base, new time.Time, changes bool) (time.Time, time.Time, bool) {
	if !new.IsZero() {
		if !new.Equal(base) {
			return new, new, true
		}
	}
	return time.Time{}, base, false || changes
}

func compareTags(base, new []string, changes bool) ([]string, []string, bool) {
	if new != nil {
		if base == nil || !arraysEqual(base, new) {
			return new, new, true
		}
	}
	return []string{}, base, false || changes
}

//return false if not equal
//this is faster than a map-based compare up to about 150/200 elements, assuming an average of a 7 letter tag.
func arraysEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		found := false
		for j := range b {
			if a[i] == b[j] {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}
