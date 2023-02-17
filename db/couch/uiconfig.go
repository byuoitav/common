package couch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/structs"
)

// GetUIConfig returns a UIConfig file from the database.
func (c *CouchDB) GetUIConfig(roomID string) (structs.UIConfig, error) {
	config, err := c.getUIConfig(roomID)
	switch {
	case err != nil:
		return structs.UIConfig{}, err
	case config.UIConfig == nil:
		return structs.UIConfig{}, errors.New("idk how this happened")
	}

	return *config.UIConfig, nil
}

func (c *CouchDB) getUIConfig(roomID string) (uiconfig, error) {
	var toReturn uiconfig

	err := c.MakeRequest("GET", fmt.Sprintf("%v/%v", UI_CONFIGS, roomID), "", nil, &toReturn)
	if err != nil {
		err = fmt.Errorf("failed to get ui config %s: %s", roomID, err)
	}

	return toReturn, err
}

// CreateUIConfig adds a new UIConfig file to the database.
func (c *CouchDB) CreateUIConfig(roomID string, toAdd structs.UIConfig) (structs.UIConfig, error) {
	var toReturn structs.UIConfig

	b, err := json.Marshal(toAdd)
	if err != nil {
		return toReturn, fmt.Errorf("failed to marshal the config file for %s: %s", roomID, err)
	}

	// Send up the UIConfig
	var resp CouchUpsertResponse
	err = c.MakeRequest("PUT", fmt.Sprintf("%v/%v", UI_CONFIGS, roomID), "application/json", b, &resp)
	if err != nil {
		if _, ok := err.(*Conflict); ok { // UIConfig with same ID already in database
			return toReturn, fmt.Errorf("unable to create ui config, because it already exists. error: %s", err)
		}

		return toReturn, fmt.Errorf("unknown error creating ui config for %s: %s", roomID, err)
	}

	toReturn, err = c.GetUIConfig(roomID)
	if err != nil {
		return toReturn, fmt.Errorf("unable to get ui config for %s : %s", roomID, err)
	}

	return toReturn, nil
}

// DeleteUIConfig removes a UIConfig file from the database.
func (c *CouchDB) DeleteUIConfig(id string) error {
	// Get the UIConfig to delete
	config, err := c.getUIConfig(id)
	if err != nil {
		return fmt.Errorf("failed to get ui config %s to delete: %s", id, err)
	}

	// Delete the UIConfig
	err = c.MakeRequest("DELETE", fmt.Sprintf("%v/%v?rev=%v", UI_CONFIGS, config.ID, config.Rev), "", nil, nil)
	if err != nil {
		return fmt.Errorf("failed to delete ui config for %s: %s", id, err)
	}

	return nil
}

// UpdateUIConfig sends an updated template to the database.
func (c *CouchDB) UpdateUIConfig(id string, update structs.UIConfig) (structs.UIConfig, error) {
	var toReturn structs.UIConfig

	if id == update.ID { // the template ID isn't changing
		// get the rev of the template
		oldConfig, err := c.getUIConfig(id)
		if err != nil {
			return toReturn, fmt.Errorf("unable to get ui config %s to update: %s", id, err)
		}

		// marshal the updated UIConfig
		b, err := json.Marshal(update)
		if err != nil {
			return toReturn, fmt.Errorf("unable to marshal updated ui config for %s : %s", update.ID, err)
		}

		// update the UIConfig
		err = c.MakeRequest("PUT", fmt.Sprintf("%s/%s?rev=%v", UI_CONFIGS, id, oldConfig.Rev), "application/json", b, &toReturn)
		if err != nil {
			return toReturn, fmt.Errorf("failed to update ui config for %s: %s", id, err)
		}
	} else { // the UIConfig ID is changing :|
		// delete the old UIConfig
		err := c.DeleteUIConfig(id)
		if err != nil {
			return toReturn, fmt.Errorf("unable to delete old ui config for %s: %s", id, err)
		}

		// marshal the new UIConfig
		b, err := json.Marshal(update)
		if err != nil {
			return toReturn, fmt.Errorf("unable to marshal new ui config for %s : %s", update.ID, err)
		}

		// post new UIConfig
		var resp CouchUpsertResponse
		err = c.MakeRequest("PUT", fmt.Sprintf("%v/%v", UI_CONFIGS, id), "", b, &resp)
		if err != nil {
			if _, ok := err.(*Conflict); ok { // a UIConfig with the same ID already exists
				return toReturn, fmt.Errorf("ui config already exists, please update this ui config or change IDs. error: %s", err)
			}

			// or an unknown error
			return toReturn, fmt.Errorf("unable to create ui config for %s : %s", id, err)
		}
	}

	toReturn, err := c.GetUIConfig(id)
	if err != nil {
		return structs.UIConfig{}, err
	}

	return toReturn, nil
}

// GetUIAttachment returns the attachment for the given ui if it exists.
// returns the content-type header, the attachment, and an error if the request against couchdb failed.
func (c *CouchDB) GetUIAttachment(ui, attachment string) (string, []byte, error) {
	url := fmt.Sprintf("%v/%v/%v/%v", c.address, UI_CONFIGS, ui, attachment)

	// start building the request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", nil, err
	}

	// add auth
	if len(c.username) > 0 && len(c.password) > 0 {
		req.SetBasicAuth(c.username, c.password)
	}

	// build client with a timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}

	if resp.StatusCode/100 != 2 {
		var ce CouchError
		err = json.Unmarshal(b, &ce)
		if err != nil {
			return "", nil, fmt.Errorf("received a non-200 response from %v. Body: %s", url, b)
		}

		log.L.Infof("Non-200 response: %v", ce.Error)
		return "", nil, CheckCouchErrors(ce)
	}

	return resp.Header.Get("content-type"), b, nil
}

// GetAllUIConfigs returns a list of all the UI Config documents in the database
func (c *CouchDB) GetAllUIConfigs() ([]structs.UIConfig, error) {
	var toReturn []structs.UIConfig
	var query IDPrefixQuery
	query.Selector.ID.GT = "\x00"
	query.Limit = 2048

	b, err := json.Marshal(query)
	if err != nil {
		return toReturn, fmt.Errorf("failed to marshal query to get all UI configs: %s", err)
	}

	var resp uiconfigQueryResponse

	err = c.MakeRequest("POST", fmt.Sprintf("%v/_find", UI_CONFIGS), "application/json", b, &resp)
	if err != nil {
		return toReturn, fmt.Errorf("failed to get all UI configs: %s", err)
	}

	for _, doc := range resp.Docs {
		toReturn = append(toReturn, *doc.UIConfig)
	}

	return toReturn, err
}
