package couch

import (
	"fmt"

	"github.com/byuoitav/common/structs"
)

// GetScheduleConfig returns the Scheduling Panel Configuration for the given room
func (c *CouchDB) GetScheduleConfig(roomID string) (structs.ScheduleConfig, error) {
	var config structs.ScheduleConfig

	err := c.MakeRequest("GET", fmt.Sprintf("%v/%v", SCHEDULING_CONFIGS, roomID), "", nil, &config)
	if err != nil {
		return config, fmt.Errorf("Error while getting Scheduling Config from DB for room %s: %s", roomID, err)
	}

	return config, nil
}
