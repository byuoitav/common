package statedefinition

import (
	"strings"
	"time"

	"github.com/byuoitav/common/nerr"
)

//Designations
const (
	Production = "production"
	Stage      = "stage"
	Test       = "test"
	Dev        = "development"
)

//SystemTypes
const (
	DMPS       = "dmps"
	Pi         = "pi"
	Scheduling = "scheduling"
	Timeclock  = "timeclock"
)

type StaticRoom struct {
	//information fields
	BuildingID string `json:"buildingID,omitempty"`
	RoomID     string `json:"roomID,omitempty"`

	//State fields
	MaintenenceMode        *bool     `json:"maintenence-mode,omitempty"`       //if the system is in maintenence mode.
	MaintenenceModeEndTime time.Time `json:"maintenence-mode-until,omitempty"` //if the system is in maintenence mode, when to put it back in monitoring.
	Monitoring             *bool     `json:"monitoring,omitempty"`             //if the system is in monitoring currently.

	Designation string   `json:"designation,omitempty"`
	SystemType  []string `json:"system-type,omitempty"` //pi, dmps, scheduling, timeclock. If a room has more than one there may be multiple entries into this field.

	Tags []string `json:"tags,omitempty"`

	UpdateTimes map[string]time.Time `json:"update-times"`
}

//Compare rooms takes two rooms and compares them, changes from new to base will only be included if they have a timestamp in UpdateTimes later than that in base for the same field
func CompareRooms(base, new StaticRoom) (diff, merged StaticRoom, changes bool, err *nerr.E) {

	merged = base

	//information fields
	if new.UpdateTimes["building"].After(base.UpdateTimes["building"]) {
		diff.BuildingID, merged.BuildingID, changes = compareString(base.BuildingID, new.BuildingID, changes)
	}
	if new.UpdateTimes["room"].After(base.UpdateTimes["room"]) {
		diff.RoomID, merged.RoomID, changes = compareString(base.RoomID, new.RoomID, changes)
	}

	if new.UpdateTimes["designation"].After(base.UpdateTimes["designation"]) {
		diff.Designation, merged.Designation, changes = compareString(base.Designation, new.Designation, changes)
	}

	//bool fields
	if new.UpdateTimes["maintenence-mode"].After(base.UpdateTimes["maintenence-mode"]) {
		diff.MaintenenceMode, merged.MaintenenceMode, changes = compareBool(base.MaintenenceMode, new.MaintenenceMode, changes)
	}
	if new.UpdateTimes["monitoring"].After(base.UpdateTimes["monitoring"]) {
		diff.Monitoring, merged.Monitoring, changes = compareBool(base.Monitoring, new.Monitoring, changes)
	}

	//time fields
	if new.UpdateTimes["maintenence-mode-until"].After(base.UpdateTimes["maintenence-mode-until"]) {
		diff.MaintenenceModeEndTime, merged.MaintenenceModeEndTime, changes = compareTime(base.MaintenenceModeEndTime, new.MaintenenceModeEndTime, changes)
	}

	if new.UpdateTimes["tags"].After(base.UpdateTimes["tags"]) {
		diff.Tags, merged.Tags, changes = compareTags(base.Tags, new.Tags, changes)
	}
	if new.UpdateTimes["system-type"].After(base.UpdateTimes["system-type"]) {
		diff.SystemType, merged.SystemType, changes = compareTags(base.SystemType, new.SystemType, changes)
	}

	return
}

func (r *StaticRoom) HasSystemType(s string) bool {
	for i := range r.SystemType {
		if r.SystemType[i] == s {
			return true
		}
	}
	return false
}

func IsDefaultSystemType(t string) bool {
	return t == DMPS || t == Pi
}

func IsDeviceOfType(d, t string) bool {
	switch t {
	case DMPS:
		return GetSuffix(d) == "DMPS"
	case Pi:
		return GetSuffix(d) == "CP"
	case Scheduling:
		return GetSuffix(d) == "SP"
	case Timeclock:
		return GetSuffix(d) == "TC"
	}
	return false
}

func GetSuffix(d string) string {
	t := d[strings.LastIndex(d, "-")+1:]
	return strings.TrimRight(t, "0123456789")
}
