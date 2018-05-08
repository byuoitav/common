package couch

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/byuoitav/common/structs"
	"github.com/byuoitav/configuration-database-microservice/log"
)

var DeviceValidationRegex *regexp.Regexp

func init() {
	DeviceValidationRegex = regexp.MustCompile(`([A-z,0-9]{2,}-[A-z,0-9]+)-[A-z]+[0-9]+`)
}

func (c *CouchDB) GetAllDevices() ([]structs.Device, error) {
	var toReturn []structs.Device

	// get all devices
	err := MakeRequest("GET", fmt.Sprintf("devices"), "", nil, &toReturn)
	if err != nil {
		msg := fmt.Sprintf("failed to get all devices: %v", err.Error())
		log.L.Error(msg)
		return toReturn, errors.New(msg)
	}

	// get all device types
	types := []structs.DeviceType{}
	err = MakeRequest("GET", fmt.Sprintf("device_types"), "", nil, &types)
	if err != nil {
		msg := fmt.Sprintf("failed to get all device types: %v", err.Error())
		log.L.Error(msg)
		return toReturn, errors.New(msg)
	}

	// create map of typeID -> type
	typeMap := make(map[string]structs.DeviceType)
	for _, t := range types {
		typeMap[t.ID] = t
	}

	// fill types into devices
	for _, d := range toReturn {
		d.Type = typeMap[d.Type.ID]
	}

	return toReturn, nil
}

func (c *CouchDB) GetDeviceByID(ID string) (structs.Device, error) {
	toReturn := structs.Device{}
	err := MakeRequest("GET", fmt.Sprintf("devices/%v", ID), "", nil, &toReturn)
	if err != nil {
		msg := fmt.Sprintf("[couch] Could not get Device %v. %v", ID, err.Error())
		log.L.Warn(msg)
	}

	return toReturn, err
}

func (c *CouchDB) GetDevicesByRoom(roomID string) ([]structs.Device, error) {
	//we query from the - to . (the character after - to get all the elements in the room
	query := IDPrefixQuery{}
	query.Selector.ID.GT = fmt.Sprintf("%v-", roomID)
	query.Selector.ID.LT = fmt.Sprintf("%v.", roomID)
	query.Limit = 1000 //some arbitrarily large number for now.

	b, err := json.Marshal(query)
	if err != nil {
		msg := fmt.Sprintf("There was a problem marshalling the query: %v", err.Error())
		log.L.Warn(msg)
		return []structs.Device{}, errors.New(msg)
	}

	toReturn := structs.DeviceQueryResponse{}
	err = MakeRequest("POST", fmt.Sprintf("devices/_find"), "application/json", b, &toReturn)

	if err != nil {
		msg := fmt.Sprintf("[couch] Could not get room %v. %v", roomID, err.Error())
		log.L.Warn(msg)
	}

	//we need to go through the devices and get their type information.
	//TODO: Cache them so we're not making a thousand requests for duplicate types.
	for i := range toReturn.Docs {
		toReturn.Docs[i].Type, err = c.GetDeviceType(toReturn.Docs[i].Type.ID)
		if err != nil {
			msg := fmt.Sprintf("Problem getting the device type %v. Error: %v", toReturn.Docs[i].Type.ID, err.Error())
			log.L.Warn(msg)
			return toReturn.Docs, errors.New(msg)

		}
	}

	return toReturn.Docs, err
}

/*
Create Device. As amazing as it may seem, this fuction creates a device in the databse.

For a device to be created, it must contain the following attributes:

	1. A valid ID
		a. The room portion corresponds to an existing room
	2. A valid name
	3. A valid type
		a. Either the ID corresponds to an existing Type, or all elements are available to create a new type. Note that if the type ID matches, but the current type doesn't match the existing ID, the current type with that ID in the Database will NOT be overwritten.
	4. A valid Class
	5. One or more roles:
		a. A role must have a valid ID and Name

Ports must pass validation - criteria are covered in the CreateDeviceType function.
However in addition, if the port includes devices those devices must be valid devices

If a device is passed into the fuction with a valid 'rev' field, the current device with that ID will be overwritten.
`rev` must be omitted to create a new device.
*/
func (c *CouchDB) CreateDevice(toAdd structs.Device) (structs.Device, error) {
	log.L.Infof("Starting add of Device: %v", toAdd.ID)

	log.L.Debug("Starting checks. Checking name and class.")
	if len(toAdd.Name) < 3 || len(toAdd.Class) < 3 {
		return lde(fmt.Sprintf("Couldn't create device - invalid name or Class"))
	}

	log.L.Debug("Name and class are good. Checking Roles")
	if len(toAdd.Roles) < 1 {
		return lde(fmt.Sprintf("Must include at least one role"))
	}

	for i := range toAdd.Roles {
		if err := checkRole(toAdd.Roles[i]); err != nil {
			return lde(fmt.Sprintf("Couldn't create device: %v", err.Error()))
		}
	}
	log.L.Debug("Roles are all valid. Checking ID")

	vals := DeviceValidationRegex.FindAllStringSubmatch(toAdd.ID, 1)
	if len(vals) == 0 {
		return lde(fmt.Sprintf("Couldn't create Device. Invalid deviceID format %v. Must match `[A-z,0-9]{2,}-[A-z,0-9]+-[A-z]+[0-9]+`", toAdd.ID))
	}

	log.L.Debug("Device ID is well formed, checking for valid room.")

	_, err := c.GetRoomByID(vals[0][1])

	if err != nil {
		if nf, ok := err.(NotFound); ok {
			return lde(fmt.Sprintf("Trying to create a device in a non-existant Room: %v. Create the room before adding the device. message: %v", vals[0][1], nf.Error()))
		}

		return lde(fmt.Sprintf("unknown problem creating the device: %v", err.Error()))
	}
	log.L.Debug("Device has a valid roomID. Checking for a valid type.")

	if len(toAdd.Type.ID) < 1 {
		return lde("Couldn't create a device, a type ID must be included")
	}

	deviceType, err := c.GetDeviceType(toAdd.Type.ID)
	if err != nil {
		if nf, ok := err.(NotFound); ok {
			log.L.Debug("Device Type not found, attempting to create. Message: %v", nf.Error())

			deviceType, err = c.CreateDeviceType(toAdd.Type)
			if err != nil {
				return lde("Trying to create a device with a non-existant device type and not enough information to create the type. Check the included type ID")
			}
			log.L.Debug("Type created")
		} else {
			lde(fmt.Sprintf("Unkown issue creating the device: %v", err.Error()))
		}
	}

	//it should only include the type ID
	toAdd.Type = structs.DeviceType{ID: deviceType.ID}

	log.L.Debug("Type is good. Checking ports.")
	for i := range toAdd.Ports {
		if err := c.checkPort(toAdd.Ports[i]); err != nil {
			return lde(fmt.Sprintf("Couldn't create device: %v", err.Error()))
		}
	}

	log.L.Debug("Ports are good. Checks passed. Creating device.")
	b, err := json.Marshal(toAdd)
	if err != nil {
		return lde(fmt.Sprintf("Couldn't marshal device into JSON. Error: %v", err.Error()))
	}

	resp := CouchUpsertResponse{}

	err = MakeRequest("PUT", fmt.Sprintf("devices/%v", toAdd.ID), "", b, &resp)
	if err != nil {
		if nf, ok := err.(Confict); ok {
			return lde(fmt.Sprintf("There was a conflict updating the device: %v. Make changes on an updated version of the configuration.", nf.Error()))
		}
		//ther was some other problem
		return lde(fmt.Sprintf("unknown problem creating the room: %v", err.Error()))
	}

	log.L.Debug("device created, retriving new device from database.")

	//return the created room
	toAdd, err = c.GetDeviceByID(toAdd.ID)
	if err != nil {
		lde(fmt.Sprintf("There was a problem getting the newly created room: %v", err.Error()))
	}

	log.L.Debug("Done creating device.")
	return toAdd, nil
}

//log device error
//alias to help cut down on cruft
func lde(msg string) (dev structs.Device, err error) {
	log.L.Warn(msg)
	err = errors.New(msg)
	return
}

func checkRole(r structs.Role) error {
	if len(r.ID) < 3 {
		return errors.New("Invalid role, check name and ID.")
	}
	return nil
}

func (c *CouchDB) checkPort(p structs.Port) error {
	if !validatePort(p) {
		return errors.New("Invalid port, check Name, ID, and Port Type")
	}

	//now we need to check the source and destination device
	if len(p.SourceDevice) > 0 {
		if _, err := c.GetDeviceByID(p.SourceDevice); err != nil {
			return errors.New(fmt.Sprintf("Invalid port %v, source device %v doesn't exist. Create it before adding it to a port", p.ID, p.SourceDevice))
		}
	}
	if len(p.DestinationDevice) > 0 {
		if _, err := c.GetDeviceByID(p.DestinationDevice); err != nil {
			return errors.New(fmt.Sprintf("Invalid port %v, destination device %v doesn't exist. Create it before adding it to a port", p.ID, p.DestinationDevice))
		}
	}

	//we're all good
	return nil
}

func (c *CouchDB) GetDevicesByRoomAndRole(roomID, role string) ([]structs.Device, error) {
	toReturn := []structs.Device{}

	devs, err := c.GetDevicesByRoom(roomID)
	if err != nil {
		msg := fmt.Sprintf("Couldn't get devices for filtering: %v", err.Error())
		log.L.Warn(msg)
		return toReturn, errors.New(msg)
	}

	//go through the devices and check if they have the role indicated
	for _, d := range devs {
		if structs.HasRole(d, role) {
			toReturn = append(toReturn, d)
		}
	}

	return toReturn, nil
}

func (c *CouchDB) GetDevicesByRoleAndType(role, dtype string) ([]structs.Device, error) {
	var toReturn []structs.Device

	devs, err := c.GetAllDevices()
	if err != nil {
		msg := fmt.Sprintf("unable to get device list: %v", err.Error())
		log.L.Warn(msg)
		return toReturn, err
	}

	for _, d := range devs {
		if structs.HasRole(d, role) && strings.EqualFold(d.Type.ID, dtype) {
			toReturn = append(toReturn, d)
		}
	}

	return toReturn, nil
}

func (c *CouchDB) DeleteDevice(id string) error {
	log.L.Debugf("[%s] Deleting device", id)

	device, err := c.GetDeviceByID(id)
	if err != nil {
		msg := fmt.Sprintf("[%s] error looking for device to delete: %s", id, err.Error())
		log.L.Warn(msg)
		return errors.New(msg)
	}

	/* TODO get rev
	err = MakeRequest("DELETE", fmt.Sprintf("devices/%s?rev=%v", device.ID, device.Rev), "", nil, nil)
	if err != nil {
		msg := fmt.Sprintf("[%s] error deleting device: %s", id, err.Error())
		log.L.Warn(msg)
		return errors.New(msg)
	}
	*/
	log.L.Debug(device)

	return nil
}
