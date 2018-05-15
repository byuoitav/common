package couch

import (
	"errors"
	"fmt"

	"github.com/byuoitav/common/structs"
)

func (c *CouchDB) GetUIConfig(roomID string) (structs.UIConfig, error) {
	config, err := c.getUIConfig(roomID)
	return *config.UIConfig, err
}

func (c *CouchDB) getUIConfig(roomID string) (uiconfig, error) {
	var toReturn uiconfig

	err := c.MakeRequest("GET", fmt.Sprintf("%v/%v", UI_CONFIGS, roomID), "", nil, &toReturn)
	if err != nil {
		err = errors.New(fmt.Sprintf("failed to get ui config %s: %s", roomID, err))
	}

	return toReturn, err
}
