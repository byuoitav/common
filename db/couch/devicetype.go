package couch

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/byuoitav/configuration-database-microservice/log"
	"github.com/byuoitav/configuration-database-microservice/structs"
)

func GetDeviceTypeByID(deviceTypeID string) (structs.DeviceType, error) {

	toReturn := structs.DeviceType{}

	err := MakeRequest("GET", fmt.Sprintf("device_type/%v", deviceTypeID), "", nil, &toReturn)

	if err != nil {
		msg := fmt.Sprintf("Could not get deviceType %v. %v", deviceTypeID, err.Error())
		log.L.Warn(msg)
		return toReturn, errors.New(msg)
	}

	return toReturn, err
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
func CreateDeviceType(toAdd structs.DeviceType) (structs.DeviceType, error) {

	log.L.Infof("Starting creation or udpate of device type %v", toAdd.ID)

	if len(toAdd.ID) < 2 || len(toAdd.Class) < 2 {
		msg := fmt.Sprintf("Device types must have a valid ID, name, and class.")
		log.L.Warn(msg)
		return toAdd, errors.New(msg)
	}

	log.L.Debug("Passed basic checks, checking ports.")

	for i := range toAdd.Ports {
		if !validatePort(toAdd.Ports[i]) {
			msg := "Port was malformed, check the name, id, and type fields"
			log.L.Warn(msg)
			return toAdd, errors.New(msg)
		}
	}
	log.L.Debug("Passed port checks, checking Commands")

	for i := range toAdd.Commands {
		if err := validateCommand(toAdd.Commands[i]); err != nil {
			log.L.Warn(err.Error())
			return toAdd, err
		}
	}

	log.L.Debug("Passed command checking. Adding the deviceType")

	b, err := json.Marshal(toAdd)
	if err != nil {
		msg := fmt.Sprintf("Couldn't marshal device type into JSON. Error: %v", err.Error())
		log.L.Error(msg)
		return toAdd, errors.New(msg)
	}

	resp := CouchUpsertResponse{}

	err = MakeRequest("PUT", fmt.Sprintf("device_types/%v", toAdd.ID), "", b, &resp)
	if err != nil {
		if nf, ok := err.(Confict); ok {
			msg := fmt.Sprintf("There was a conflict updating the device type: %v. Make changes on an updated version of the configuration.", nf.Error())
			log.L.Warn(msg)
			return structs.DeviceType{}, errors.New(msg)
		}
		//ther was some other problem
		msg := fmt.Sprintf("unknown problem creating the device type: %v", err.Error())
		log.L.Warn(msg)
		return structs.DeviceType{}, errors.New(msg)
	}

	log.L.Debug("Device Type created, retriving new record from database.")

	//return the created device type
	toAdd, err = GetDeviceTypeByID(toAdd.ID)
	if err != nil {
		msg := fmt.Sprintf("There was a problem getting the newly created device type: %v", err.Error())
		log.L.Warn(msg)
		return toAdd, errors.New(msg)
	}

	return toAdd, nil
}

func validatePort(p structs.Port) bool {
	if len(p.ID) < 3 || len(p.Name) < 3 || len(p.PortType) < 3 {
		return false
	}
	return true
}

func validateCommand(c structs.Command) error {
	if len(c.ID) < 3 || len(c.Name) < 3 {
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
	if len(m.ID) < 3 || len(m.Name) < 3 {
		return errors.New("Invalid Mircroservice. Check Name, and ID")
	}

	//check the address

	if _, err := url.ParseRequestURI(m.Address); err != nil {
		return errors.New("Invalid Microservice Address")
	}

	return nil
}

func checkEndpoint(e structs.Endpoint) error {
	if len(e.ID) < 3 || len(e.Name) < 3 {
		return errors.New("Invalid endpoint. Check Name, and ID")
	}

	if _, err := url.ParseRequestURI(e.Path); err != nil {
		return errors.New("Invalid endpoint path")
	}

	return nil
}
