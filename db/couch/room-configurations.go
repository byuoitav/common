package couch

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/byuoitav/common/structs"
)

func (c *CouchDB) GetRoomConfiguration(id string) (structs.RoomConfiguration, error) {
	c.log.Debugf("Getting room configuration: %v", id)

	toReturn := structs.RoomConfiguration{}
	err := c.MakeRequest("GET", fmt.Sprintf("room_configurations/%v", id), "", nil, &toReturn)

	if err != nil {
		msg := fmt.Sprintf("Could not get room configuration %v. %v", id, err.Error())
		c.log.Warn(msg)
	}

	return toReturn, err
}

func (c *CouchDB) GetAllRoomConfigurations() ([]structs.RoomConfiguration, error) {
	return []structs.RoomConfiguration{}, nil
}

/*
CreatRoomConfiguraiton adds a room configuration. A valid room configuration must have:
1) an ID
2) a name
3) at least one evaluator.
	An Evaluator must have an ID and a CodeKey.
	Priority will default to 1000.

Note that if the ID is a duplicate, assuming that the `rev` field is valid.
The existing document will get overwritten.
*/
func (c *CouchDB) CreateRoomConfiguration(config structs.RoomConfiguration) (structs.RoomConfiguration, error) {
	c.log.Debugf("Creating a room configuration: %v", config.ID)

	if len(config.ID) == 0 {
		c.log.Warn("Couldn't create configuration, ID is required.")
		return config, errors.New("Couldn't create configuration, ID is required.")
	}

	if len(config.Evaluators) == 0 {
		msg := "Couldn't create configuration, at least on evaluator is required."
		c.log.Warn(msg)
		return config, errors.New(msg)
	}

	//now we need to go through and check each Evaluator.
	//TODO: Figure out some way to check if the evaluator key is valid

	for i := range config.Evaluators {
		if len(config.Evaluators[i].ID) < 1 || len(config.Evaluators[i].CodeKey) < 1 {
			msg := "Couldn't Create configuration. All evaluators must have an ID and a codeKey"
			c.log.Warn(msg)
			return config, errors.New(msg)
		}
		//check if priority is 0, if so, set it to 1000

		if config.Evaluators[i].Priority == 0 {
			config.Evaluators[i].Priority = 1000
		}
	}

	c.log.Debugf("All checks passed. Creating configuration.")

	resp := CouchUpsertResponse{}

	b, err := json.Marshal(config)
	if err != nil {

		msg := fmt.Sprintf("Couldn't marshal configuration into JSON. Error: %v", err.Error())
		c.log.Error(msg) // this is a slightly bigger deal
		return config, errors.New(msg)
	}

	err = c.MakeRequest("PUT", fmt.Sprintf("room_configurations/%v", config.ID), "", b, &resp)
	if err != nil {
		if nf, ok := err.(*Conflict); ok {
			msg := fmt.Sprintf("There was a conflict updating the configuration: %v. Make changes on an updated version of the configuration.", nf.Error())
			c.log.Warn(msg)
			return config, errors.New(msg)
		}
		//ther was some other problem
		msg := fmt.Sprintf("unknown problem creating the configuration: %v", err.Error())
		c.log.Warn(msg)
		return config, errors.New(msg)
	}

	c.log.Debug("Configuration created, retriving new configuration from database.")

	//return the created config
	config, err = c.GetRoomConfiguration(config.ID)
	if err != nil {
		msg := fmt.Sprintf("There was a problem getting the newly created configuration: %v", err.Error())
		c.log.Warn(msg)
		return config, errors.New(msg)
	}

	c.log.Debug("Done.")
	return config, nil
}

func (c *CouchDB) DeleteRoomConfiguration(id string) error {
	return nil
}

func (c *CouchDB) UpdateRoomConfiguration(id string, rc structs.RoomConfiguration) (structs.RoomConfiguration, error) {
	return structs.RoomConfiguration{}, nil
}
