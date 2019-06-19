package couch

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"sync"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/common/structs"
)

func (c *CouchDB) GetRoom(id string) (structs.Room, error) {
	room, err := c.getRoom(id)
	if err != nil {
		return structs.Room{}, err
	}
	//if err was nil then room may be.
	return *room.Room, nil
}

func (c *CouchDB) getRoom(id string) (room, error) {
	var toReturn room

	// get the base room
	err := c.MakeRequest("GET", fmt.Sprintf("%s/%v", ROOMS, id), "", nil, &toReturn)
	if err != nil {
		return toReturn, errors.New(fmt.Sprintf("failed to get room %s: %s", id, err))
	}

	// get the devices in room
	devices, err := c.GetDevicesByRoom(id)
	if err != nil {
		return toReturn, errors.New(fmt.Sprintf("failed to get devices in room %s: %s", id, err))
	}

	// fill devices into room
	for _, device := range devices {
		toReturn.Devices = append(toReturn.Devices, device)
	}

	// get room configuration
	toReturn.Configuration, err = c.GetRoomConfiguration(toReturn.Configuration.ID)
	if err != nil {
		return toReturn, errors.New(fmt.Sprintf("failed to get room configuration %s for room %s: %s", toReturn.Configuration.ID, id, err))
	}

	return toReturn, nil
}

func (c *CouchDB) getRoomsByQuery(query IDPrefixQuery) ([]room, error) {
	var toReturn []room
	var resp roomQueryResponse

	err := c.ExecuteQuery(query, resp)
	if err != nil {
		return toReturn, errors.New(fmt.Sprintf("failed to get rooms by query: %s", err))
	}

	for _, doc := range resp.Docs {
		toReturn = append(toReturn, doc)
	}

	return toReturn, nil
}

func (c *CouchDB) GetAllRooms() ([]structs.Room, error) {
	var toReturn []structs.Room

	// create all room query
	var query IDPrefixQuery
	query.Selector.ID.GT = "\x00"
	query.Limit = 1000

	// marshal query
	b, err := json.Marshal(query)
	if err != nil {
		return toReturn, errors.New(fmt.Sprintf("failed to marshal query to get all rooms: %s", err))
	}

	// make request to get rooms
	var resp roomQueryResponse
	err = c.MakeRequest("POST", fmt.Sprintf("%v/_find", ROOMS), "application/json", b, &resp)
	if err != nil {
		return toReturn, errors.New(fmt.Sprintf("failed to get all rooms: %s", err))
	}

	// add each doc to toReturn
	for _, doc := range resp.Docs {
		toReturn = append(toReturn, *doc.Room)
	}

	return toReturn, nil
}

func (c *CouchDB) GetRoomsByBuilding(id string) ([]structs.Room, error) {
	var toReturn []structs.Room

	// create query from - to . (the character after - to get all the elements in the room)
	var query IDPrefixQuery
	query.Selector.ID.GT = fmt.Sprintf("%v-", id)
	query.Selector.ID.LT = fmt.Sprintf("%v.", id)
	query.Limit = 1000

	// marshal query
	b, err := json.Marshal(query)
	if err != nil {
		return toReturn, errors.New(fmt.Sprintf("failed to marshal query to get rooms in building %s: %s", id, err))
	}

	// make request to get rooms
	var resp roomQueryResponse
	err = c.MakeRequest("POST", fmt.Sprintf("%v/_find", ROOMS), "application/json", b, &resp)
	if err != nil {
		return toReturn, errors.New(fmt.Sprintf("failed to get rooms in building %s: %s", id, err))
	}

	// add each doc to toReturn
	for _, doc := range resp.Docs {
		toReturn = append(toReturn, *doc.Room)
	}

	return toReturn, nil
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
func (c *CouchDB) CreateRoom(toAdd structs.Room) (structs.Room, error) {
	var toReturn structs.Room

	// validate room struct
	err := toAdd.Validate()
	if err != nil {
		return toReturn, err
	}

	// ensure it's in a real building
	_, err = c.GetBuilding(strings.Split(toAdd.ID, "-")[0])
	if err != nil {
		if _, ok := err.(*NotFound); ok {
			return toReturn, errors.New(fmt.Sprintf("unable to create room %s: building %s doesn't exist.", toAdd.ID, strings.Split(toAdd.ID, "-")[0]))
		}

		return toReturn, errors.New(fmt.Sprintf("unable to validate room %s is in a real building: %s", toAdd.ID, err))
	}

	// ensure room configuration exists (if it doesn't, create it!)
	config, err := c.GetRoomConfiguration(toAdd.Configuration.ID)
	if err != nil {
		if _, ok := err.(*NotFound); ok { // room config wasn't found
			// create new room configuration for this room
			config, err = c.CreateRoomConfiguration(toAdd.Configuration)
			if err != nil {
				return toReturn, errors.New(fmt.Sprintf("unable to create room %s: %s", toAdd.ID, err))
			}
		} else { // some other error looking for room config
			return toReturn, errors.New(fmt.Sprintf("unable to validate if room configuration %s exists or not: %s", toAdd.Configuration.ID, err))
		}
	}

	// we only want to post room config ID up, so replace config with just that
	toAdd.Configuration = structs.RoomConfiguration{ID: config.ID}

	// save the devices to create after creating the room
	var devices []structs.Device
	copy(devices, toAdd.Devices)

	// don't post devices to room table
	toAdd.Devices = []structs.Device{}

	// marshal room
	b, err := json.Marshal(toAdd)
	if err != nil {
		return toReturn, errors.New(fmt.Sprintf("failed to marshal room %s: %s", toAdd.ID, err))
	}

	// post up room!
	var resp CouchUpsertResponse
	err = c.MakeRequest("POST", fmt.Sprintf("%v", ROOMS), "application/json", b, &resp)
	if err != nil {
		if _, ok := err.(*Conflict); ok { // there was a conflict creating room
			return toReturn, errors.New(fmt.Sprintf("unable to create new room, because it already exists. error: %s", err))
		}

		return toReturn, errors.New(fmt.Sprintf("unknown error creating room %s: %s", toAdd.ID, err))
	}

	// get room back from database
	toReturn, err = c.GetRoom(toAdd.ID)
	if err != nil {
		return toReturn, errors.New(fmt.Sprintf("error getting room %s after creating it: %s", toAdd.ID, err))
	}

	// create the devices
	var wg sync.WaitGroup
	for index := range devices {
		wg.Add(1)

		go func(i int) {
			d, err := c.CreateDevice(devices[i])
			if err == nil {
				toReturn.Devices = append(toReturn.Devices, d)
			}
			wg.Done()
		}(index)
	}

	wg.Wait()

	return toReturn, nil
}

func (c *CouchDB) DeleteRoom(id string) error {
	// get the room to delete
	room, err := c.getRoom(id)
	if err != nil {
		return errors.New(fmt.Sprintf("unable to get room %s to delete: %s", id, err))
	}

	// delete each of the devices from the room
	var wg sync.WaitGroup
	for index := range room.Devices {
		wg.Add(1)

		go func(i int) {
			c.DeleteDevice(room.Devices[i].ID)
			wg.Done()
		}(index)
	}
	wg.Wait()

	// delete the room
	err = c.MakeRequest("DELETE", fmt.Sprintf("%v/%v?rev=%v", ROOMS, room.ID, room.Rev), "", nil, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("unable to delete room %s: %s", id, err))
	}

	return nil
}

func (c *CouchDB) UpdateRoom(id string, room structs.Room) (structs.Room, error) {
	var toReturn structs.Room

	// validate the room
	err := room.Validate()
	if err != nil {
		return toReturn, err
	}

	// verify the room configuration is real, if it isn't, then create it
	config, err := c.GetRoomConfiguration(room.Configuration.ID)
	if err != nil {
		if _, ok := err.(*NotFound); ok { // room config wasn't found
			// create new room configuration for this room
			config, err = c.CreateRoomConfiguration(room.Configuration)
			if err != nil {
				return toReturn, errors.New(fmt.Sprintf("unable to create room %s: %s", room.ID, err))
			}
		} else { // some other error looking for room config
			return toReturn, errors.New(fmt.Sprintf("unable to validate if room configuration %s exists or not: %s", room.Configuration.ID, err))
		}
	}

	// strip off devices and extra room configuration information
	room.Devices = nil
	room.Configuration = structs.RoomConfiguration{ID: config.ID}

	// get the current room
	r, err := c.getRoom(id)
	if err != nil {
		return toReturn, errors.New(fmt.Sprintf("unable to get room %s to update: %s'", id, err))
	}

	if id == room.ID { // the room ID isn't changing
		// marshal the updated room
		b, err := json.Marshal(room)
		if err != nil {
			return toReturn, errors.New(fmt.Sprintf("failed to unmarshal new room %s: %s", id, err))
		}

		// update the room
		err = c.MakeRequest("PUT", fmt.Sprintf("%v/%s?rev=%v", ROOMS, id, r.Rev), "application/json", b, &toReturn)
		if err != nil {
			return toReturn, errors.New(fmt.Sprintf("failed to update room %s: %s", id, err))
		}
	} else { // the room ID is changing
		// delete the old room
		err = c.DeleteRoom(id)
		if err != nil {
			return toReturn, errors.New(fmt.Sprintf("failed to delete old room %s: %s", id, err))
		}

		// move the old room's devices into the new room's devices, and edit their ID's
		for _, device := range room.Devices {
			device.ID = strings.Replace(device.ID, id, room.ID, 1)
			room.Devices = append(room.Devices, device)
		}

		// create new room
		toReturn, err = c.CreateRoom(room)
		if err != nil {
			return toReturn, errors.New(fmt.Sprintf("failed to update room %s: %s", id, err))
		}
	}

	return toReturn, nil
}

func (c *CouchDB) GetRoomsByDesignation(designation string) ([]structs.Room, *nerr.E) {

	var toReturn []structs.Room

	// get all rooms
	rooms, err := c.GetAllRooms()
	if err != nil {
		return toReturn, nerr.Translate(err).Addf("failed to get rooms by room designation.")
	}

	// filter for ones that have the room configuration
	for _, room := range rooms {
		if strings.EqualFold(room.Designation, designation) {
			toReturn = append(toReturn, room)
		}
	}

	return toReturn, nil
}

// TODO could use a query to be faster
func (c *CouchDB) GetRoomsByRoomConfiguration(configID string) ([]structs.Room, error) {
	var toReturn []structs.Room

	// get all rooms
	rooms, err := c.GetAllRooms()
	if err != nil {
		return toReturn, errors.New(fmt.Sprintf("failed to get rooms by room configuration: %s", err))
	}

	// filter for ones that have the room configuration
	for _, room := range rooms {
		if strings.EqualFold(room.Configuration.ID, configID) {
			toReturn = append(toReturn, room)
		}
	}

	return toReturn, nil
}

// GetRoomAttachments gets the attachments in a room
func (c *CouchDB) GetRoomAttachments(room string) ([]string, error) {
	log.L.Infof(room)
	var roomAttachments roomAttachmentResponse

	//re := regexp.MustCompile(`[A-Za-z|\d|\/\/|:|.|-]*\.jpg`)

	err := c.MakeRequest("GET", fmt.Sprintf("%s/%v", ROOM_ATTACHMENTS, room), "", nil, &roomAttachments)
	if err != nil {
		return []string{}, errors.New(fmt.Sprintf("failed to get room %s: %s", room, err))
	}
	log.L.Infof("what about here\n %v", roomAttachments)
	var toReturn []string
	for k := range roomAttachments.Attachments {
		toReturn = append(toReturn, k)
		log.L.Infof(k)
	}

	// for _, s := range toReturn {
	// 	log.L.Infof(s)
	// }

	sort.Strings(toReturn)
	for i := range toReturn {
		u := &url.URL{}
		err := u.UnmarshalBinary([]byte("https://couchdb-stg.avs.byu.edu:5984/room_attachments/" + room + "/" + toReturn[i]))
		{
			if err != nil {
				log.L.Infof("There was a problem getting the image URLs", err)
			}
		}
		toReturn[i] = u.String()
		//toReturn[i] = "https://couchdb-stg.avs.byu.edu:5984/room_attachments/" + room + "/" + toReturn[i]
	}

	log.L.Infof("This is the data now: ", toReturn)
	return toReturn, nil
}
