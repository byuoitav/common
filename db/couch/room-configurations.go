package couch

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/byuoitav/common/structs"
)

func (c *CouchDB) GetRoomConfiguration(id string) (structs.RoomConfiguration, error) {
	rc, err := c.getRoomConfiguration(id)
	switch {
	case err != nil:
		return structs.RoomConfiguration{}, err
	case rc.RoomConfiguration == nil:
		return structs.RoomConfiguration{}, fmt.Errorf("no room configuration %q found", id)
	}

	return *rc.RoomConfiguration, err
}

func (c *CouchDB) getRoomConfiguration(id string) (roomConfiguration, error) {
	var toReturn roomConfiguration

	err := c.MakeRequest("GET", fmt.Sprintf("%v/%v", ROOM_CONFIGURATIONS, id), "", nil, &toReturn)

	if err != nil {
		err = errors.New(fmt.Sprintf("failed to get room configuration %s: %s", id, err))
	}

	return toReturn, err
}

func (c *CouchDB) getRoomConfigurationsByQuery(query IDPrefixQuery) ([]roomConfiguration, error) {
	var toReturn []roomConfiguration

	// marshal query
	b, err := json.Marshal(query)
	if err != nil {
		return toReturn, errors.New(fmt.Sprintf("failed to marshal room configurations query: %s", err))
	}

	// make query for room configs
	var resp roomConfigurationQueryResponse
	err = c.MakeRequest("POST", fmt.Sprintf("%s/_find", ROOM_CONFIGURATIONS), "application/json", b, &resp)
	if err != nil {
		return toReturn, errors.New(fmt.Sprintf("failed to query room configurations: %s", err))
	}

	// return each document
	for _, doc := range resp.Docs {
		toReturn = append(toReturn, doc)
	}

	return toReturn, nil
}

func (c *CouchDB) GetAllRoomConfigurations() ([]structs.RoomConfiguration, error) {
	var toReturn []structs.RoomConfiguration

	// create all device types query
	var query IDPrefixQuery
	query.Selector.ID.GT = "\x00"
	query.Limit = 5000

	// execute query
	configs, err := c.getRoomConfigurationsByQuery(query)
	if err != nil {
		return toReturn, errors.New(fmt.Sprintf("failed to get all room configurations: %s", err))
	}

	// just return the struct
	for _, config := range configs {
		toReturn = append(toReturn, *config.RoomConfiguration)
	}

	return toReturn, nil
}

/*
CreateRoomConfiguraiton adds a room configuration. A valid room configuration must have:
1) an ID
2) a name
3) at least one evaluator.
	An Evaluator must have an ID and a CodeKey.
	Priority will default to 1000.

Note that if the ID is a duplicate, assuming that the `rev` field is valid.
The existing document will get overwritten.
*/
func (c *CouchDB) CreateRoomConfiguration(toAdd structs.RoomConfiguration) (structs.RoomConfiguration, error) {
	var toReturn structs.RoomConfiguration

	// validate room config
	err := toAdd.Validate(true)
	if err != nil {
		return toReturn, err
	}

	// TODO figure out how to check if the evalutaor key is valid

	// marshal room config
	b, err := json.Marshal(toAdd)
	if err != nil {
		return toReturn, errors.New(fmt.Sprintf("failed to marshal room configuration %s: %s", toAdd.ID, err))
	}

	// post room configuration
	var resp CouchUpsertResponse
	err = c.MakeRequest("POST", ROOM_CONFIGURATIONS, "application/json", b, &resp)
	if err != nil {
		if _, ok := err.(*Conflict); ok {
			return toReturn, errors.New(fmt.Sprintf("room configuration already exists; please update this configuration or change id's. error: %s", err))
		}

		return toReturn, errors.New(fmt.Sprintf("failed to create the room configuration %s: %s", toAdd.ID, err))
	}

	// get new room config from db, return it
	toReturn, err = c.GetRoomConfiguration(toAdd.ID)
	if err != nil {
		return toReturn, errors.New(fmt.Sprintf("failed to get room configuration %s after creating it: %s", toAdd.ID, err))
	}

	return toReturn, nil
}

func (c *CouchDB) DeleteRoomConfiguration(id string) error {
	// validate no rooms depend on this type
	rooms, err := c.GetRoomsByRoomConfiguration(id)
	if err != nil {
		return err
	}

	if len(rooms) != 0 {
		return errors.New(fmt.Sprintf("can't delete room configuration %s. %v rooms still depend on it.", id, len(rooms)))
	}

	// get the rev of the room configuration
	config, err := c.getRoomConfiguration(id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get room configuration %s to delete. does it exist? (error: %s)", id, err))
	}

	// delete room config
	err = c.MakeRequest("DELETE", fmt.Sprintf("%v/%v?rev=%v", ROOM_CONFIGURATIONS, id, config.Rev), "", nil, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("unable to delete room configuration %s: %s", config.ID, err))
	}

	return nil
}

func (c *CouchDB) UpdateRoomConfiguration(id string, rc structs.RoomConfiguration) (structs.RoomConfiguration, error) {
	return structs.RoomConfiguration{}, errors.New(fmt.Sprintf("not implemented"))
}
