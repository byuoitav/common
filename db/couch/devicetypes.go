package couch

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/byuoitav/common/structs"
)

func (c *CouchDB) GetDeviceType(id string) (structs.DeviceType, error) {
	dt, err := c.getDeviceType(id)
	return *dt.DeviceType, err
}

func (c *CouchDB) getDeviceType(id string) (deviceType, error) {
	var toReturn deviceType

	err := c.MakeRequest("GET", fmt.Sprintf("%v/%v", DEVICE_TYPES, id), "", nil, &toReturn)
	if err != nil {
		err = errors.New(fmt.Sprintf("failed to get device type %s: %s", id, err))
	}

	return toReturn, err
}

func (c *CouchDB) getDeviceTypesByQuery(query IDPrefixQuery) ([]deviceType, error) {
	var toReturn []deviceType

	// marshal query
	b, err := json.Marshal(query)
	if err != nil {
		return toReturn, errors.New(fmt.Sprintf("failed to marshal device types query: %s", err))
	}

	// make query for types
	var resp deviceTypeQueryResponse
	err = c.MakeRequest("POST", fmt.Sprintf("%s/_find", DEVICE_TYPES), "application/json", b, &resp)
	if err != nil {
		return toReturn, errors.New(fmt.Sprintf("failed to query device types: %s", err))
	}

	// return each document
	for _, doc := range resp.Docs {
		toReturn = append(toReturn, doc)
	}

	return toReturn, nil
}

func (c *CouchDB) GetAllDeviceTypes() ([]structs.DeviceType, error) {
	var toReturn []structs.DeviceType

	// create all device types query
	var query IDPrefixQuery
	query.Selector.ID.GT = "\x00"
	query.Limit = 5000

	// execute query
	types, err := c.getDeviceTypesByQuery(query)
	if err != nil {
		return toReturn, errors.New(fmt.Sprintf("failed getting all device types: %s", err))
	}

	// return the struct part
	for _, t := range types {
		toReturn = append(toReturn, *t.DeviceType)
	}

	return toReturn, nil
}

/*
CreateDeviceType - now this may come as a shock - creates a device type.
The device type must have the following attributes to be created:
	1. A valid ID (3 or more characters)
	2. A valid Name (3 or more characters) 3. A valid class (3 or more characters) Each command and Port must be validated as well. The criterea for their creation is:
Port:
	1. Must have a valid ID (3 or more characters)
	2. Must have a valid Name (3 or more characters)
	3. Must have a valid PortType (3 or more characters)

Command:
	1. Must have a valid Microservice
		a. Valid ID (3 or more characters)
		b. Valid Name (3 or more characters)
		c. Valid Address (a well formed HTTP or HTTPS host0
	2. Must have a valid Endpoint
		a. Valid ID (3 or more characters)
		b. Valid Name (3 ore more characters)
		c. Valid, well formed URI Path.
	3. Must have a valid command
		a. Valid ID (3 or more characters)
		b. Valid Name (3 ore more characters)

If a device type is submitted with a valid 'rev' field, the device type will be overwritten.
*/
func (c *CouchDB) CreateDeviceType(toAdd structs.DeviceType) (structs.DeviceType, error) {
	var toReturn structs.DeviceType

	// validate device type
	err := toAdd.Validate(true)
	if err != nil {
		return toReturn, err
	}

	// marshal device type
	b, err := json.Marshal(toAdd)
	if err != nil {
		return toReturn, errors.New(fmt.Sprintf("failed to marshal building %s: %s", toAdd.ID, err))
	}

	// post device type
	var resp CouchUpsertResponse
	err = c.MakeRequest("POST", DEVICE_TYPES, "application/json", b, &resp)
	if err != nil {
		if _, ok := err.(*Conflict); ok { // a device type with this id already exists
			return toReturn, errors.New(fmt.Sprintf("device type already exists, please update this type or change id's. error: %s", err))
		}

		return toReturn, errors.New(fmt.Sprintf("failed to create the device type %s: %s", toAdd.ID, err))
	}

	// get new device type from db, and return that
	toReturn, err = c.GetDeviceType(toAdd.ID)
	if err != nil {
		return toReturn, errors.New(fmt.Sprintf("failed to get device type %s after creating it: %s", toAdd.ID, err))
	}

	return toReturn, nil
}

func (c *CouchDB) DeleteDeviceType(id string) error {
	// validate no devices depend on this type
	devices, err := c.GetDevicesByType(id)
	if err != nil {
		return errors.New(fmt.Sprintf("unable to validate no devices depend on this type: %s", err))
	}

	if len(devices) != 0 {
		return errors.New(fmt.Sprintf("can't delete device type %s. %v devices still depend on it.", id, len(devices)))
	}

	// get the rev of the device
	deviceType, err := c.getDeviceType(id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get device type %s to delete. does it exist? (error: %s)", id, err))
	}

	// delete device type
	err = c.MakeRequest("DELETE", fmt.Sprintf("%v/%v?rev=%v", DEVICE_TYPES, deviceType.ID, deviceType.Rev), "", nil, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("unable to delete device type %s: %s", deviceType.ID, err))
	}

	return nil
}

func (c *CouchDB) UpdateDeviceType(id string, dt structs.DeviceType) (structs.DeviceType, error) {
	return structs.DeviceType{}, errors.New("not implemented")
}
