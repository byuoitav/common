package couch

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/byuoitav/configuration-database-microservice/log"
	"github.com/byuoitav/configuration-database-microservice/structs"
)

//GetBuildingByID gets the company's building with the corresponding ID from the couch database
func (c *CouchDB) GetBuildingByID(id string) (structs.Building, error) {

	toReturn := structs.Building{}
	err := MakeRequest("GET", fmt.Sprintf("buildings/%v", id), "", nil, &toReturn)
	if err != nil {
		msg := fmt.Sprintf("[couch] Could not get building %v. %v", id, err.Error())
		log.L.Warn(msg)
	}

	return toReturn, err
}

//GetAllBuildings returns all buildings for the company specified
func (c *CouchDB) GetAllBuildings() ([]structs.Building, error) {
	var toReturn []structs.Building
	var bulk structs.BulkBuildingResponse

	err := MakeRequest("GET", fmt.Sprintf("buildings/_all_docs?limit=1000&include_docs=true"), "", nil, &bulk)
	if err != nil {
		msg := fmt.Sprintf("[couch] Could not get buildings. %v", err.Error())
		log.L.Warn(msg)
		return toReturn, errors.New(msg)
	}

	for _, row := range bulk.Rows {
		toReturn = append(toReturn, row.Doc)
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
	log.L.Debugf("Starting adding a building: %v", toAdd.Name)

	if len(toAdd.ID) < 2 || len(toAdd.Name) < 2 {
		msg := "Cannot create building, must have at least a name and an ID"
		log.L.Warn(msg)
	}

	//there's not a lot to buildings, so we can just add

	log.L.Debug("Checks passed, creating building.")

	b, err := json.Marshal(toAdd)
	if err != nil {

		msg := fmt.Sprintf("Couldn't marshal building into JSON. Error: %v", err.Error())
		log.L.Error(msg) // this is a slightly bigger deal
		return toAdd, errors.New(msg)
	}

	resp := CouchUpsertResponse{}

	err = MakeRequest("PUT", fmt.Sprintf("buildings/%v", toAdd.ID), "", b, &resp)
	if err != nil {
		log.L.Debugf("%v", err)
		if conflict, ok := err.(*Confict); ok {
			msg := fmt.Sprintf("Error: %v Make changes on an updated version of the configuration.", conflict.Error())
			log.L.Warn(msg)
			return toAdd, errors.New(msg)
		}
		//ther was some other problem
		msg := fmt.Sprintf("unknown problem creating the Building: %v", err.Error())
		log.L.Warn(msg)
		return toAdd, errors.New(msg)
	}

	log.L.Debug("Building created, retriving new configuration from database.")

	//return the created config
	toAdd, err = c.GetBuildingByID(toAdd.ID)
	if err != nil {
		msg := fmt.Sprintf("There was a problem getting the newly created building: %v", err.Error())
		log.L.Warn(msg)
		return toAdd, errors.New(msg)
	}

	log.L.Debug("Done.")
	return toAdd, nil
}

func (c *CouchDB) DeleteBuilding(id string) error {
	log.L.Infof("Starting delete for building: %v", id)
	building, err := c.GetBuildingByID(id)
	if err != nil {
		msg := fmt.Sprintf("There was a problem deleting the building: %v", err.Error())
		log.L.Warn(msg)
		return errors.New(msg)
	}

	log.L.Debugf("Checking for rooms in building %v", id)
	//first we need to check if there are any rooms in the building, if there are, we don't allow a delete at this level
	rms, err := c.GetRoomsByBuilding(id)
	if err != nil {
		msg := fmt.Sprintf("Problem checking the building for rooms: %v", err.Error())
		log.L.Warn(msg)
		return errors.New(msg)
	}

	if len(rms) > 0 {
		log.L.Infof("There were %v rooms in building %v, aborting delete", len(rms), id)
		return errors.New(fmt.Sprintf("Couldn't delete building %v, there are still rooms associated with the building. Remove all rooms from building before deleting.", id))
	}

	log.L.Debugf("No rooms found in building %v. Proceeding with deletion", id)

	err = MakeRequest("DELETE", fmt.Sprintf("buildings/%s?rev=%v", id, building.Rev), "", nil, nil)
	if err != nil {
		msg := fmt.Sprintf("There was a problem deleting the building: %v", err.Error())
		log.L.Warn(msg)
		return errors.New(msg)
	}

	log.L.Debugf("Building %v deleted successfully.", id)
	return nil
}
