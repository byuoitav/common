package couch

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/byuoitav/common/structs"
)

func (c *CouchDB) GetDeviceType(deviceTypeID string) (structs.DeviceType, error) {

	toReturn := structs.DeviceType{}

	err := c.MakeRequest("GET", fmt.Sprintf("device_type/%v", deviceTypeID), "", nil, &toReturn)

	if err != nil {
		msg := fmt.Sprintf("Could not get deviceType %v. %v", deviceTypeID, err.Error())
		c.log.Warn(msg)
		return toReturn, errors.New(msg)
	}

	return toReturn, err
}

func (c *CouchDB) GetAllDeviceTypes() ([]structs.DeviceType, error) {
	return []structs.DeviceType{}, nil
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

	c.log.Infof("Starting creation or udpate of device type %v", toAdd.ID)

	if len(toAdd.ID) < 2 {
		msg := fmt.Sprintf("Device types must have a valid ID and name")
		c.log.Warn(msg)
		return toAdd, errors.New(msg)
	}

	c.log.Debug("Passed basic checks, checking ports.")

	for i := range toAdd.Ports {
		if !validatePort(toAdd.Ports[i]) {
			msg := "Port was malformed, check the name, id, and type fields"
			c.log.Warn(msg)
			return toAdd, errors.New(msg)
		}
	}
	c.log.Debug("Passed port checks, checking Commands")

	for i := range toAdd.Commands {
		if err := validateCommand(toAdd.Commands[i]); err != nil {
			c.log.Warn(err.Error())
			return toAdd, err
		}
	}

	c.log.Debug("Passed command checking. Adding the deviceType")

	b, err := json.Marshal(toAdd)
	if err != nil {
		msg := fmt.Sprintf("Couldn't marshal device type into JSON. Error: %v", err.Error())
		c.log.Error(msg)
		return toAdd, errors.New(msg)
	}

	resp := CouchUpsertResponse{}

	err = c.MakeRequest("PUT", fmt.Sprintf("device_types/%v", toAdd.ID), "", b, &resp)
	if err != nil {
		if nf, ok := err.(Conflict); ok {
			msg := fmt.Sprintf("There was a conflict updating the device type: %v. Make changes on an updated version of the configuration.", nf.Error())
			c.log.Warn(msg)
			return structs.DeviceType{}, errors.New(msg)
		}
		//ther was some other problem
		msg := fmt.Sprintf("unknown problem creating the device type: %v", err.Error())
		c.log.Warn(msg)
		return structs.DeviceType{}, errors.New(msg)
	}

	c.log.Debug("Device Type created, retriving new record from database.")

	//return the created device type
	toAdd, err = c.GetDeviceType(toAdd.ID)
	if err != nil {
		msg := fmt.Sprintf("There was a problem getting the newly created device type: %v", err.Error())
		c.log.Warn(msg)
		return toAdd, errors.New(msg)
	}

	return toAdd, nil
}

// TODO
func (c *CouchDB) UpdateDeviceType(id string, dt structs.DeviceType) (structs.DeviceType, error) {
	return structs.DeviceType{}, nil
}

// TODO
func (c *CouchDB) DeleteDeviceType(id string) error {
	return nil
}

func validatePort(p structs.Port) bool {
	if len(p.ID) < 3 {
		return false
	}
	return true
}

func validateCommand(c structs.Command) error {
	if len(c.ID) < 3 {
		return errors.New("Invalid base information. Check Name, and ID")
	}

	//check the microservice
	err := checkMicroservice(c.Microservice)
	if err != nil {
		return err
	}
	return checkEndpoint(c.Endpoint)
}

func checkMicroservice(m structs.Microservice) error {
	if len(m.ID) < 3 {
		return errors.New("Invalid Mircroservice. Check Name, and ID")
	}

	//check the address

	if _, err := url.ParseRequestURI(m.Address); err != nil {
		return errors.New("Invalid Microservice Address")
	}

	return nil
}

func checkEndpoint(e structs.Endpoint) error {
	if len(e.ID) < 3 {
		return errors.New("Invalid endpoint. Check Name, and ID")
	}

	if _, err := url.ParseRequestURI(e.Path); err != nil {
		return errors.New("Invalid endpoint path")
	}

	return nil
}
