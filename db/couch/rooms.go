package couch

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/byuoitav/common/structs"
)

var roomValidationRegex = regexp.MustCompile(`([A-z,0-9]{2,})-[A-z,0-9]+`)

func (c *CouchDB) GetRoom(id string) (structs.Room, error) {
	toReturn := structs.Room{}
	err := c.MakeRequest("GET", fmt.Sprintf("rooms/%v", id), "", nil, &toReturn)
	if err != nil {
		msg := fmt.Sprintf("[couch] Could not get room %v. %v", id, err.Error())
		c.log.Warn(msg)
	}

	// TODO we need to get the room configuration information
	// TODO we need to get devices

	return toReturn, err
}

func (c *CouchDB) GetAllRooms() ([]structs.Room, error) {
	var toReturn []structs.Room
	var query IDPrefixQuery
	query.Selector.ID.GT = "\x00"
	query.Limit = 1000

	b, err := json.Marshal(query)
	if err != nil {
		c.log.Warnf("There was a problem marshalling the query: %v", err.Error())
		return []structs.Room{}, err
	}

	var resp roomQueryResponse

	err = c.MakeRequest("POST", fmt.Sprintf("rooms/_find"), "application/json", b, &resp)
	if err != nil {
		msg := fmt.Sprintf("[couch] Could not get buildings. %v", err.Error())
		c.log.Warn(msg)
		return []structs.Room{}, errors.New(msg)
	}

	for _, doc := range resp.Docs {
		toReturn = append(toReturn, doc.Room)
	}

	return toReturn, nil
}

func (c *CouchDB) GetRoomsByBuilding(buildingID string) ([]structs.Room, error) {
	//we query from the - to . (the character after - to get all the elements in the room
	query := IDPrefixQuery{}
	query.Selector.ID.GT = fmt.Sprintf("%v-", buildingID)
	query.Selector.ID.LT = fmt.Sprintf("%v.", buildingID)
	query.Limit = 1000 //some arbitrarily large number for now.
	b, err := json.Marshal(query)
	if err != nil {
		c.log.Warnf("There was a problem marshalling the query: %v", err.Error())
		return []structs.Room{}, err
	}
	c.log.Debugf("Getting all rooms for building: %v", buildingID)

	toReturn := structs.RoomQueryResponse{}
	err = c.MakeRequest("POST", fmt.Sprintf("rooms/_find"), "application/json", b, &toReturn)

	if err != nil {
		msg := fmt.Sprintf("[couch] Could not get room %v. %v", buildingID, err.Error())
		c.log.Warn(msg)
	}

	return toReturn.Docs, err
}

/*
CreateRoom creates a room. Required information:
	1. The room must have a valid roomID, that roomID must have a valid BuildingID as a component
	2. The configurationID of the sub configuration item must have at least a valid ID. If the ID doesn't exist currently in the database, the room configuraiton object must meet all requirements to be a valid roomConfiguration.
	3. The room must have a name.
	4. The room must have a designation

	It is important to note that the function will overwrite a room with the same roomID if the Rev field is valid.

	Any devices included in the room will be evaluated for adding, but the room will be evaluated for creation first. If any devices fail creation, this will NOT roll back the creation of the room, or any other devices. All devices wil  be checked for a device ID before moving to creation. If any are lacking, the no cration of ANY device will proceed.
*/

func (c *CouchDB) CreateRoom(room structs.Room) (structs.Room, error) {
	c.log.Debugf("Starting room creation for %v", room.ID)

	vals := roomValidationRegex.FindAllStringSubmatch(room.ID, 1)
	if len(vals) == 0 {
		msg := fmt.Sprintf("Couldn't create room. Invalid roomID format %v. Must match `(A-z,0-9]{2,}-[A-z,0-9]+`", room.ID)
		c.log.Warn(msg)
		return structs.Room{}, errors.New(msg)
	}

	//we really should check all the other information here, too
	if len(room.Name) < 1 || len(room.Designation) < 1 {
		msg := "Couldn't create room. The room must include a name and a designation."
		c.log.Warn(msg)
		return structs.Room{}, errors.New(msg)
	}
	c.log.Debugf("RoomID and other information is valid, checking for valid buildingID: %v", vals[0][1])

	_, err := c.GetBuilding(vals[0][1])

	if err != nil {
		if nf, ok := err.(*NotFound); ok {
			msg := fmt.Sprintf("Trying to create a room in a non-existant building: %v. Create the building before adding the room. message: %v", vals[0][1], nf.Error())
			c.log.Warn(msg)
			return structs.Room{}, errors.New(msg)
		}

		msg := fmt.Sprintf("unknown problem creating the room: %v", err.Error())
		c.log.Warn(msg)
		return structs.Room{}, errors.New(msg)
	}

	c.log.Debugf("We have a valid buildingID. Checking for a valid room configuration ID")

	if len(room.Configuration.ID) < 1 {
		msg := fmt.Sprintf("Could not create room: A room configuration ID must be included")
		c.log.Warn(msg)
		return structs.Room{}, errors.New(msg)
	}
	//get the configuration and check to see if it's not there. If it isn't there, try to add it. If it can't be addedfor whatever reason (it doesn't meet the rquirements) error out.
	config, err := c.GetRoomConfiguration(room.Configuration.ID)
	if err != nil {
		if _, ok := err.(*NotFound); ok {
			c.log.Debugf("Room configuration %v not found, attempting to create.", room.Configuration.ID)

			//this is where we try to create the configuration
			config, err = c.CreateRoomConfiguration(room.Configuration)
			if err != nil {

				msg := "Trying to create a room with a non-existant room configuration and not enough information to create the configuration. Check the included configuration ID."
				c.log.Warn(msg)
				return structs.Room{}, errors.New(msg)
			}
			c.log.Debugf("Room configuration created.")
		} else {

			msg := fmt.Sprintf("unknown problem creating the room: %v", err.Error())
			c.log.Warn(msg)
			return structs.Room{}, errors.New(msg)
		}
	}

	//the configuration should only have the ID.
	room.Configuration = structs.RoomConfiguration{ID: config.ID}
	c.log.Debug("Room configuration passed. Creating the room.")

	//save the devices for later, if there are any, then remove the frmo the room for putting into the database
	c.log.Debugf("There are %v devices included, saving to be added later.", len(room.Devices))

	devs := []structs.Device{}
	copy(devs, room.Devices)
	room.Devices = []structs.Device{}

	b, err := json.Marshal(room)
	if err != nil {
		msg := fmt.Sprintf("Couldn't marshal room into JSON. Error: %v", err.Error())
		c.log.Error(msg)
		return structs.Room{}, errors.New(msg)
	}

	resp := CouchUpsertResponse{}

	err = c.MakeRequest("PUT", fmt.Sprintf("rooms/%v", room.ID), "", b, &resp)
	if err != nil {
		if nf, ok := err.(*Conflict); ok {
			msg := fmt.Sprintf("There was a conflict updating the room: %v. Make changes on an updated version of the configuration.", nf.Error())
			c.log.Warn(msg)
			return structs.Room{}, errors.New(msg)
		}
		//ther was some other problem
		msg := fmt.Sprintf("unknown problem creating the room: %v", err.Error())
		c.log.Warn(msg)
		return structs.Room{}, errors.New(msg)
	}

	c.log.Debug("room created, retriving new room from database.")

	//return the created room
	room, err = c.GetRoom(room.ID)
	if err != nil {
		msg := fmt.Sprintf("There was a problem getting the newly created room: %v", err.Error())
		c.log.Warn(msg)
		return room, errors.New(msg)
	}

	c.log.Debug("Done creating room, evaluating devices for creation.")

	// Do the devices.
	room.Devices = []structs.Device{}

	for d := range devs {
		dev, err := c.CreateDevice(devs[d])
		if err != nil {
			c.log.Info("Error creating device %v as part of room. Error: %v.", devs[d].ID, err.Error())
			continue
		}
		room.Devices = append(room.Devices, dev)
	}

	return room, nil
}

func (c *CouchDB) DeleteRoom(id string) error {
	c.log.Infof("[%s] Deleting room", id)

	// get the room
	room, err := c.GetRoom(id)
	if err != nil {
		msg := fmt.Sprintf("[%s] error looking for room to delete: %s", id, err.Error())
		c.log.Warn(msg)
		return errors.New(msg)
	}

	/* TODO get rev
	// delete each of the devices from the room
	c.log.Debugf("[%s] Deleting devices from room", id)
	for _, d := range room.Devices {
		c.log.Debugf("[%s] Deleting device %s", id, d.ID)
		err = c.MakeRequest("DELETE", fmt.Sprintf("devices/%s?rev=%v", d.ID, d.Rev), "", nil, nil)
		if err != nil {
			msg := fmt.Sprintf("[%s] error deleting device %s: %s", id, d.ID, err.Error())
			c.log.Warn(msg)
			return errors.New(msg)
		}
	}

	// delete the room
	c.log.Debugf("[%s] Successfully deleted devices from room. Deleting room...", id)
	err = c.MakeRequest("DELETE", fmt.Sprintf("rooms/%s?rev=%v", room.ID, room.Rev), "", nil, nil)
	if err != nil {
		msg := fmt.Sprintf("[%s] error deleting room: %s", id, err.Error())
		c.log.Warn(msg)
		return errors.New(msg)
	}
	*/
	c.log.Debug(room)

	c.log.Infof("[%s] Successfully deleted room", id)
	return nil
}

func (c *CouchDB) UpdateRoom(id string, room structs.Room) (structs.Room, error) {
	return structs.Room{}, nil
}
