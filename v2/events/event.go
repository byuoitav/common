package events

import (
	"strings"
	"time"
)


// An Event is generated as a result of something happening in a room and enables other systems to act on it, as well as collect metrics.
type Event struct {
	// GeneratingSystem is the system actually generating the event. i.e. For an API call against a raspberry pi this would be the hostname of the raspberry pi running the AV-API. If the call is against AWS, this would be 'AWS'
	GeneratingSystem string `json:"generating-system"`

	// Timestamp is the time the event took place
	Timestamp time.Time `json:"timestamp"`

	// EventTags is a collection of strings to give more information about what kind of event this is, used in routing and processing events. See the EventTags const delcaration for some common tags.
	EventTags []string `json:"event-tags"`

	// TargetDevice is the device being affected by the event. e.g. a power on event, this would be the device powering on
	TargetDevice BasicDeviceInfo `json:"target-device"`

	// AffectedRoom is the room being affected by the event. e.g. in events arising from an API call this is the room called in the API
	AffectedRoom BasicRoomInfo `json:"affected-room"`

	// Key of the event
	Key string `json:"key"`

	// Value of the event
	Value string `json:"value"`

	// User is the user associated with generating the event
	User string `json:"user"`

	// Data is an optional field to dump data that you wont necessarily want to aggregate on, but you may want to search on
	Data interface{} `json:"data,omitempty"`
}

// BasicRoomInfo contains device information that is easy to aggregate on.
type BasicRoomInfo struct {
	BuildingID string `json:"buildingID,omitempty"`
	RoomID     string `json:"roomID,omitempty"`
}

// BasicDeviceInfo contains device information that is easy to aggregate on.
type BasicDeviceInfo struct {
	BasicRoomInfo
	DeviceID string `json:"deviceID,omitempty"`
}

// GenerateBasicDeviceInfo takes a deviceID and generates a BasicDeviceInfo from it
func GenerateBasicDeviceInfo(deviceID string) BasicDeviceInfo {
	deviceID = strings.ToUpper(deviceID)

	vals := strings.Split(deviceID, "-")
	if len(vals) != 3 {
		return BasicDeviceInfo{DeviceID: deviceID}
	}

	return BasicDeviceInfo{
		BasicRoomInfo: BasicRoomInfo{
			BuildingID: vals[0],
			RoomID:     vals[0] + "-" + vals[1],
		},
		DeviceID: deviceID,
	}
}

// GenerateBasicRoomInfo takes a roomID and generates a BasicRoomInfo from it
func GenerateBasicRoomInfo(roomID string) BasicRoomInfo {
	roomID = strings.ToUpper(roomID)

	vals := strings.Split(roomID, "-")
	if len(vals) != 2 {
		return BasicRoomInfo{
			RoomID: vals[0],
		}
	}

	return BasicRoomInfo{
		BuildingID: vals[0],
		RoomID:     vals[1],
	}
}

// HasTag returns true of the event has the given tag, or false if it doesn't.
func HasTag(e Event, t string) bool {
	for i := range e.EventTags {
		if e.EventTags[i] == t {
			return true
		}
	}

	return false
}
