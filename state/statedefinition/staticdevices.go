package statedefinition

import (
	"time"

	"github.com/byuoitav/common/nerr"
)

//StaticDevice .
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
	Power           string `json:"power,omitempty"`
	Active          string `json:"active,omitempty"`
	Secure          string `json:"secure,omitempty"`
	HardwareVersion string `json:"hardware-version,omitempty"`
	SoftwareVersion string `json:"software-version,omitempty"`
	FirmwareVersion string `json:"firmware-version"`
	SerialNumber    string `json:"serial-number"`
	IPAddress       string `json:"ip-address"`
	MACAddress      string `json:"mac-address"`
	ModelName       string `json:"model-name"`

	DNSAddress     string `json:"dns-address,omitempty"`
	DefaultGateway string `json:"default-gateway"`

	//Control Processor Specific Fields
	Websocket      string `json:"websocket,omitempty"`
	WebsocketCount *int   `json:"websocket-count,omitempty"`

	//Display Specific Fields
	Blanked      *bool  `json:"blanked,omitempty"`
	Input        string `json:"input,omitempty"`
	LampHours    *int   `json:"lamp-hours,omitempty"`
	Temperature  *int   `json:"temperature,omitempty"`
	ActiveSignal *bool  `json:"active-signal,omitempty"`

	//Audio Device Specific Fields
	Muted  *bool `json:"muted,omitempty"`
	Volume *int  `json:"volume,omitempty"`

	//Fields specific to Microphones
	BatteryChargeBars         *int   `json:"battery-charge-bars,omitempty"`
	BatteryChargeMinutes      *int   `json:"battery-charge-minutes,omitempty"`
	BatteryChargePercentage   *int   `json:"battery-charge-percentage,omitempty"`
	BatteryChargeHoursMinutes string `json:"battery-charge-hours-minutes,omitempty"`
	BatteryCycles             *int   `json:"battery-cycles,omitempty"`
	BatteryType               string `json:"battery-type,omitempty"`
	MicrophoneChannel         string `json:"microphone-channel,omitempty"`
	Interference              string `json:"interference,omitempty"`

	//Fields specific to Vias

	CurrentUserCount *int `json:"current-user-count,omitempty"`
	PresenterCount   *int `json:"presenter-count,omitempty"`

	//meta fields for use in kibana
	Control               string `json:"control,omitempty"`                //the id - used in a URL
	EnableNotifications   string `json:"enable-notifications,omitempty"`   //the id - used in a URL
	SuppressNotifications string `json:"suppress-notifications,omitempty"` //the id - used in a URL
	ViewDashboard         string `json:"ViewDashboard,omitempty"`          //the id - used in a URL

	//Linux Device Information
	CPUUsagePercentage    *float64 `json:"cpu-usage-percent,omitempty"`
	VMemUsage             *float64 `json:"v-mem-used-percent,omitempty"`
	SMemUsage             *float64 `json:"s-mem-used-percent,omitempty"`
	CPUTemp               *float64 `json:"cpu-thermal0-temp,omitempty"`
	DiskWrites            *int     `json:"writes-to-mmcblk0,omitempty"`
	DiskUsagePercentage   *float64 `json:"disk-used-percent,omitempty"`
	AverageProcessesSleep *float64 `json:"avg-procs-u-sleep,omitempty"`

	BroadcomChipTemp *float64 `json:"bcm2835_thermal0-temp,omitempty"`

	//DMPS information
	StatusMessage   string `json:"status-message,omitempty"`
	TransmitRFPower string `json:"transmit-rf-power,omitempty"`

	// HardwareInfo

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
	if new.UpdateTimes["active"].After(base.UpdateTimes["active"]) {
		diff.Active, merged.Active, changes = compareString(base.Active, new.Active, changes)
	}
	if new.UpdateTimes["secure"].After(base.UpdateTimes["secure"]) {
		diff.Secure, merged.Secure, changes = compareString(base.Secure, new.Power, changes)
	}
	if new.UpdateTimes["hardware-version"].After(base.UpdateTimes["hardware-version"]) {
		diff.HardwareVersion, merged.HardwareVersion, changes = compareString(base.HardwareVersion, new.HardwareVersion, changes)
	}
	if new.UpdateTimes["software-version"].After(base.UpdateTimes["software-version"]) {
		diff.SoftwareVersion, merged.SoftwareVersion, changes = compareString(base.SoftwareVersion, new.SoftwareVersion, changes)
	}
	if new.UpdateTimes["firmware-version"].After(base.UpdateTimes["firmware-version"]) {
		diff.FirmwareVersion, merged.FirmwareVersion, changes = compareString(base.FirmwareVersion, new.FirmwareVersion, changes)
	}
	if new.UpdateTimes["serial-number"].After(base.UpdateTimes["serial-number"]) {
		diff.SerialNumber, merged.SerialNumber, changes = compareString(base.SerialNumber, new.SerialNumber, changes)
	}
	if new.UpdateTimes["ip-address"].After(base.UpdateTimes["ip-address"]) {
		diff.IPAddress, merged.IPAddress, changes = compareString(base.IPAddress, new.IPAddress, changes)
	}
	if new.UpdateTimes["mac-address"].After(base.UpdateTimes["mac-address"]) {
		diff.MACAddress, merged.MACAddress, changes = compareString(base.MACAddress, new.MACAddress, changes)
	}
	if new.UpdateTimes["dns-address"].After(base.UpdateTimes["dns-address"]) {
		diff.DNSAddress, merged.DNSAddress, changes = compareString(base.DNSAddress, new.DNSAddress, changes)
	}
	if new.UpdateTimes["default-gateway"].After(base.UpdateTimes["default-gateway"]) {
		diff.DefaultGateway, merged.DefaultGateway, changes = compareString(base.DefaultGateway, new.DefaultGateway, changes)
	}
	if new.UpdateTimes["model-name"].After(base.UpdateTimes["model-name"]) {
		diff.ModelName, merged.ModelName, changes = compareString(base.ModelName, new.ModelName, changes)
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
	if new.UpdateTimes["lamp-hours"].After(base.UpdateTimes["lamp-hours"]) {
		diff.LampHours, merged.LampHours, changes = compareInt(base.LampHours, new.LampHours, changes)
	}
	if new.UpdateTimes["temperature"].After(base.UpdateTimes["temperature"]) {
		diff.Temperature, merged.Temperature, changes = compareInt(base.Temperature, new.Temperature, changes)
	}

	if new.UpdateTimes["active-signal"].After(base.UpdateTimes["active-signal"]) {
		diff.ActiveSignal, merged.ActiveSignal, changes = compareBool(base.ActiveSignal, new.ActiveSignal, changes)
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
		diff.BatteryChargeHoursMinutes, merged.BatteryChargeHoursMinutes, changes = compareString(base.BatteryChargeHoursMinutes, new.BatteryChargeHoursMinutes, changes)
	}
	if new.UpdateTimes["battery-cycles"].After(base.UpdateTimes["battery-cycles"]) {
		diff.BatteryCycles, merged.BatteryCycles, changes = compareInt(base.BatteryCycles, new.BatteryCycles, changes)
	}
	if new.UpdateTimes["battery-type"].After(base.UpdateTimes["battery-type"]) {
		diff.BatteryType, merged.BatteryType, changes = compareString(base.BatteryType, new.BatteryType, changes)
	}
	if new.UpdateTimes["microphone-channel"].After(base.UpdateTimes["microphone-channel"]) {
		diff.MicrophoneChannel, merged.MicrophoneChannel, changes = compareString(base.MicrophoneChannel, new.MicrophoneChannel, changes)
	}
	if new.UpdateTimes["interference"].After(base.UpdateTimes["interference"]) {
		diff.Interference, merged.Interference, changes = compareString(base.Interference, new.Interference, changes)
	}

	//Via specific fields
	if new.UpdateTimes["current-user-count"].After(base.UpdateTimes["current-user-count"]) {
		diff.CurrentUserCount, merged.CurrentUserCount, changes = compareInt(base.CurrentUserCount, new.CurrentUserCount, changes)
	}
	if new.UpdateTimes["presenter-count"].After(base.UpdateTimes["presenter-count"]) {
		diff.PresenterCount, merged.PresenterCount, changes = compareInt(base.PresenterCount, new.PresenterCount, changes)
	}

	//PI Hardware Info fields
	if new.UpdateTimes["cpu-usage-percent"].After(base.UpdateTimes["cpu-usage-percent"]) {
		diff.CPUUsagePercentage, merged.CPUUsagePercentage, changes = compareFloat64(base.CPUUsagePercentage, new.CPUUsagePercentage, changes)
	}
	if new.UpdateTimes["v-mem-used-percent"].After(base.UpdateTimes["v-mem-used-percent"]) {
		diff.VMemUsage, merged.VMemUsage, changes = compareFloat64(base.VMemUsage, new.VMemUsage, changes)
	}
	if new.UpdateTimes["s-mem-used-percent"].After(base.UpdateTimes["s-mem-used-percent"]) {
		diff.SMemUsage, merged.SMemUsage, changes = compareFloat64(base.SMemUsage, new.SMemUsage, changes)
	}
	if new.UpdateTimes["cpu-thermal0-temp"].After(base.UpdateTimes["cpu-thermal0-temp"]) {
		diff.CPUTemp, merged.CPUTemp, changes = compareFloat64(base.CPUTemp, new.CPUTemp, changes)
	}
	if new.UpdateTimes["writes-to-mmcblk0"].After(base.UpdateTimes["writes-to-mmcblk0"]) {
		diff.DiskWrites, merged.DiskWrites, changes = compareInt(base.DiskWrites, new.DiskWrites, changes)
	}
	if new.UpdateTimes["disk-used-percent"].After(base.UpdateTimes["disk-used-percent"]) {
		diff.DiskUsagePercentage, merged.DiskUsagePercentage, changes = compareFloat64(base.DiskUsagePercentage, new.DiskUsagePercentage, changes)
	}
	if new.UpdateTimes["avg-procs-u-sleep"].After(base.UpdateTimes["avg-procs-u-sleep"]) {
		diff.AverageProcessesSleep, merged.AverageProcessesSleep, changes = compareFloat64(base.AverageProcessesSleep, new.AverageProcessesSleep, changes)
	}
	if new.UpdateTimes["bcm2835_thermal0-temp"].After(base.UpdateTimes["bcm2835_thermal0-temp"]) {
		diff.BroadcomChipTemp, merged.BroadcomChipTemp, changes = compareFloat64(base.BroadcomChipTemp, new.BroadcomChipTemp, changes)
	}
	//DMPS fields
	if new.UpdateTimes["status-message"].After(base.UpdateTimes["status-message"]) {
		diff.StatusMessage, merged.StatusMessage, changes = compareString(base.StatusMessage, new.StatusMessage, changes)
	}
	if new.UpdateTimes["transmit-rf-power"].After(base.UpdateTimes["transmit-rf-power"]) {
		diff.TransmitRFPower, merged.TransmitRFPower, changes = compareString(base.TransmitRFPower, new.TransmitRFPower, changes)
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

	for k, v := range new.UpdateTimes {
		if v.After(base.UpdateTimes[k]) {
			merged.UpdateTimes[k] = v
		}
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

func compareFloat64(base, new *float64, changes bool) (*float64, *float64, bool) {
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
