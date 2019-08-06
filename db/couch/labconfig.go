package couch

import (
	"fmt"

	"github.com/byuoitav/common/structs"
)

// GetLabConfig returns the Lab Configuration for the given room
func (c *CouchDB) GetLabConfig(roomID string) (structs.LabConfig, error) {

	var config structs.LabConfig

	err := c.MakeRequest("GET", fmt.Sprintf("%v/%v", LAB_CONFIGS, roomID), "", nil, &config)
	if err != nil {
		return config, fmt.Errorf("Error while getting Lab Config from DB for room %s: %s", roomID, err)
	}

	return config, nil
}
