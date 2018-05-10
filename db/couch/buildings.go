package couch

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/byuoitav/common/structs"
)

//GetBuilding gets the company's building with the corresponding ID from the couch database
func (c *CouchDB) GetBuilding(id string) (structs.Building, error) {
	resp, err := c.getBuilding(id)
	return resp.Building, err
}

func (c *CouchDB) getBuilding(id string) (building, error) {
	var toReturn building

	err := c.MakeRequest("GET", fmt.Sprintf("buildings/%v", id), "", nil, &toReturn)
	if err != nil {
		err = &NotFound{
			msg: fmt.Sprintf("failed to get building %v. %s", id, err),
		}

		c.log.Warn(err)
	}

	return toReturn, err
}

//GetAllBuildings returns all buildings for the company specified
func (c *CouchDB) GetAllBuildings() ([]structs.Building, error) {
	var toReturn []structs.Building
	var query IDPrefixQuery
	query.Selector.ID.GT = "\x00"
	query.Limit = 1000

	b, err := json.Marshal(query)
	if err != nil {
		c.log.Warnf("There was a problem marshalling the query: %v", err.Error())
		return toReturn, err
	}

	var resp buildingQueryResponse

	err = c.MakeRequest("POST", fmt.Sprintf("buildings/_find"), "application/json", b, &resp)
	if err != nil {
		msg := fmt.Sprintf("[couch] Could not get buildings. %v", err.Error())
		c.log.Warn(msg)
		return toReturn, errors.New(msg)
	}

	for _, doc := range resp.Docs {
		toReturn = append(toReturn, doc.Building)
	}

	return toReturn, err
}

/*
AddBuilding adds a building. The building must have at least:
1) ID
2) Name

The function will also overwrite the existing building providing the _rev field is set properly
*/
func (c *CouchDB) CreateBuilding(toAdd structs.Building) (structs.Building, error) {
	c.log.Debugf("Starting adding a building: %v", toAdd.ID)

	if len(toAdd.ID) < 2 {
		msg := "Cannot create building, must have an ID"
		c.log.Warn(msg)
	}

	//there's not a lot to buildings, so we can just add

	c.log.Debug("Checks passed, creating building.")

	b, err := json.Marshal(toAdd)
	if err != nil {

		msg := fmt.Sprintf("Couldn't marshal building into JSON. Error: %v", err.Error())
		c.log.Error(msg) // this is a slightly bigger deal
		return toAdd, errors.New(msg)
	}

	resp := CouchUpsertResponse{}

	err = c.MakeRequest("PUT", fmt.Sprintf("buildings/%v", toAdd.ID), "", b, &resp)
	if err != nil {
		c.log.Debugf("%v", err)
		if conflict, ok := err.(*Conflict); ok {
			msg := fmt.Sprintf("Error: %v Make changes on an updated version of the configuration.", conflict.Error())
			c.log.Warn(msg)
			return toAdd, errors.New(msg)
		}
		//ther was some other problem
		msg := fmt.Sprintf("unknown problem creating the Building: %v", err.Error())
		c.log.Warn(msg)
		return toAdd, errors.New(msg)
	}

	c.log.Debug("Building created, retriving new configuration from database.")

	//return the created config
	toAdd, err = c.GetBuilding(toAdd.ID)
	if err != nil {
		msg := fmt.Sprintf("There was a problem getting the newly created building: %v", err.Error())
		c.log.Warn(msg)
		return toAdd, errors.New(msg)
	}

	c.log.Debug("Done.")
	return toAdd, nil
}

func (c *CouchDB) DeleteBuilding(id string) error {
	c.log.Infof("Starting delete for building: %v", id)
	building, err := c.getBuilding(id)
	if err != nil {
		msg := fmt.Sprintf("There was a problem deleting the building: %v", err.Error())
		c.log.Warn(msg)
		return errors.New(msg)
	}

	c.log.Debugf("Checking for rooms in building %v", id)
	//first we need to check if there are any rooms in the building, if there are, we don't allow a delete at this level
	rms, err := c.GetRoomsByBuilding(id)
	if err != nil {
		msg := fmt.Sprintf("Problem checking the building for rooms: %v", err.Error())
		c.log.Warn(msg)
		return errors.New(msg)
	}

	if len(rms) > 0 {
		c.log.Infof("There were %v rooms in building %v, aborting delete", len(rms), id)
		return errors.New(fmt.Sprintf("Couldn't delete building %v, there are still rooms associated with the building. Remove all rooms from building before deleting.", id))
	}

	c.log.Debugf("No rooms found in building %v. Proceeding with deletion", id)

	err = c.MakeRequest("DELETE", fmt.Sprintf("buildings/%s?rev=%v", id, building.Rev), "", nil, nil)
	if err != nil {
		msg := fmt.Sprintf("There was a problem deleting the building: %v", err.Error())
		c.log.Warn(msg)
		return errors.New(msg)
	}

	c.log.Debugf("Building %v deleted successfully.", id)
	return nil
}

func (c *CouchDB) deleteBuildingWithoutCascade(id string) error {
	c.log.Infof("starting to delete building *without cascading*: %s", id)

	building, err := c.getBuilding(id)
	if err != nil {
		msg := fmt.Sprintf("error getting the building for deletion: %s", err)
		c.log.Warn(msg)
		return errors.New(msg)
	}

	err = c.MakeRequest("DELETE", fmt.Sprintf("buildings/%s?rev=%v", id, building.Rev), "", nil, nil)
	if err != nil {
		msg := fmt.Sprintf("error deleting the building %s: %s", id, err)
		c.log.Warn(msg)
		return errors.New(msg)
	}

	return nil
}

func (c *CouchDB) UpdateBuilding(id string, building structs.Building) (structs.Building, error) {
	var toReturn structs.Building
	var roomsToUpdate []structs.Room

	b, err := json.Marshal(building)
	if err != nil {
		c.log.Warnf("there was a problem marshalling the query: %s", err)
		return toReturn, err
	}

	// TODO reorder this so it is after the PUT request
	// we need to 'change' the id for the new document
	if id != building.ID {
		// delete the old building
		err = c.deleteBuildingWithoutCascade(id)
		if err != nil {
			msg := fmt.Sprintf("failed to delete old building %s: %s", id, err)
			c.log.Warn(msg)
			return toReturn, err
		}

		// save the old rooms to be updated afterwords
		roomsToUpdate, err = c.GetRoomsByBuilding(id)
		if err != nil {
			msg := fmt.Sprintf("failed to get rooms inside the building %s: %s", id, err)
			c.log.Warn(msg)
			return toReturn, err
		}
	}

	// update the building document
	err = c.MakeRequest("PUT", fmt.Sprintf("buildings/%s", building.ID), "application/json", b, &toReturn)
	if err != nil {
		msg := fmt.Sprintf("error updating the building %s: %s", building.ID, err)
		c.log.Warn(msg)
		return toReturn, err
	}

	// update all the rooms that need to be updated if the building id changed
	for _, room := range roomsToUpdate {
		c.log.Debug(room)
		go func() {
			// create the new roomID
			oldID := room.ID
			split := strings.Split(oldID, "-")
			split[0] = building.ID
			room.ID = strings.Join(split, "-")

			c.log.Debugf("Changing room %s to %s, to update building from %s to %s", oldID, room.ID, id, building.ID)

			c.UpdateRoom(oldID, room)
		}()
	}

	return toReturn, nil
}
