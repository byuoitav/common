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
	MaintenanceMode *bool    `json:"maintenance-mode,omitempty"`
	Tags            []string `json:"tags,omitempty"`

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

	//bool fields
	if new.UpdateTimes["maintenance-mode"].After(base.UpdateTimes["maintenance-mode"]) {
		diff.MaintenanceMode, merged.MaintenanceMode, changes = compareBool(base.MaintenanceMode, new.MaintenanceMode, changes)
	}

	if new.UpdateTimes["tags"].After(base.UpdateTimes["tags"]) {
		diff.Tags, merged.Tags, changes = compareTags(base.Tags, new.Tags, changes)
	}

	return
}
