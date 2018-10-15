package couch

import (
	"encoding/json"
	"fmt"

	sd "github.com/byuoitav/common/state/statedefinition"
)

/*
The big TODO
Alright so here goes
-I need to confirm EXACTLY what data is wanted and what functions they want to be able to call.
	-This includes being able to put a filter onto the data they're getting back
		-This may be satisfied with GetDeviceStatesByQuery
-Also, anything else that could possibly be wanted should be included. We want this thing to be prepped for the full-blown SMEE

*/

func (c *CouchDB) GetDeviceState(id string) (sd.StaticDevice, error) {
	DeviceState, err := c.getDeviceState(id)
	return DeviceState, err
}
func (c *CouchDB) getDeviceState(id string) (sd.StaticDevice, error) {
	var toReturn sd.StaticDevice // get the device state
	err := c.MakeRequest("GET", fmt.Sprintf("%s/%v", DEVICE_STATES, id), "", nil, &toReturn)
	if err != nil {
		return toReturn, fmt.Errorf("failed to get device state for %s: %s", id, err)
	}

	if len(toReturn.DeviceID) == 0 {
		return toReturn, fmt.Errorf("failed to get device state for %s: %s", id, err)
	}
	return toReturn, err
}

//TODO -> includeType may be something that people don't care about
//TODO -> Find what people want
func (c *CouchDB) getDeviceStatesByQuery(query IDPrefixQuery, includeType bool) ([]sd.StaticDevice, error) {
	var toReturn []sd.StaticDevice

	// marshal query
	b, err := json.Marshal(query)
	if err != nil {
		return toReturn, fmt.Errorf("failed to marshal device state query: %s", err)
	}

	// make query for devices
	var resp deviceStateQueryResponse
	err = c.MakeRequest("POST", fmt.Sprintf("%s/_find", DEVICE_STATES), "application/json", b, &resp)
	if err != nil {
		return toReturn, fmt.Errorf("failed to query device state: %s", err)
	}
	/*
		if includeType {
			// get all types
			types, err := c.GetAllDeviceStateTypes()
			if err != nil {
				return toReturn, fmt.Errorf("failed to get devices types for device state query:%s", err)
			}

			// make a map of type.ID -> type
			typesMap := make(map[string]structs.DeviceTypes)
			for _, t := range types {
				typesMap[t.ID] = t
			}

			// fill in device types
			for _, d := range resp.Docs {
				d.Type = typesMap[d.Type.ID]
			}
		}
	*/
	// return each document
	for _, doc := range resp.Docs {
		toReturn = append(toReturn, doc)
	}

	return toReturn, nil
}
func (c *CouchDB) GetAllDeviceStates() ([]sd.StaticDevice, error) {
	var toReturn []sd.StaticDevice

	// create all device state query
	var query IDPrefixQuery
	query.Selector.ID.GT = "\x00"
	query.Limit = 5000

	// query device states
	deviceStates, err := c.getDeviceStatesByQuery(query, false)
	if err != nil {
		return toReturn, fmt.Errorf("failed getting all device states: %s", err)
	}

	// get the struct out of each device
	for _, deviceState := range deviceStates {
		toReturn = append(toReturn, deviceState)
	}

	return toReturn, nil
}

func (c *CouchDB) GetDeviceStatesByRoom(roomID string) ([]sd.StaticDevice, error) {
	var toReturn []sd.StaticDevice

	deviceStates, err := c.getDeviceStatesByRoom(roomID)
	if err != nil {
		return toReturn, nil
	}

	for _, deviceState := range deviceStates {
		toReturn = append(toReturn, deviceState)
	}

	return toReturn, nil
}

func (c *CouchDB) getDeviceStatesByRoom(roomID string) ([]sd.StaticDevice, error) {
	var toReturn []sd.StaticDevice

	// create query
	var query IDPrefixQuery
	query.Selector.ID.Regex = fmt.Sprintf("%v.*", roomID)
	query.Limit = 1000

	// query devices
	toReturn, err := c.getDeviceStatesByQuery(query, true)
	if err != nil {
		return toReturn, fmt.Errorf("failed getting device states in room %s: %s", roomID, err)
	}

	return toReturn, nil
}

func (c *CouchDB) GetDeviceStatesByBuilding(buildingID string) ([]sd.StaticDevice, error) {
	var toReturn []sd.StaticDevice

	deviceStates, err := c.getDeviceStatesByBuilding(buildingID)
	if err != nil {
		return toReturn, nil
	}

	for _, deviceState := range deviceStates {
		toReturn = append(toReturn, deviceState)
	}

	return toReturn, nil
}

func (c *CouchDB) getDeviceStatesByBuilding(buildingID string) ([]sd.StaticDevice, error) {
	var toReturn []sd.StaticDevice

	// create query
	var query IDPrefixQuery
	query.Selector.ID.Regex = fmt.Sprintf("%v.*", buildingID)
	query.Limit = 1000

	// query devices
	toReturn, err := c.getDeviceStatesByQuery(query, true)
	if err != nil {
		return toReturn, fmt.Errorf("failed getting device states in building %s: %s", buildingID, err)
	}

	return toReturn, nil
}
