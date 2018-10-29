package couch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/byuoitav/common/log"
)

const (
	BUILDINGS           = "buildings"
	ROOMS               = "rooms"
	DEVICES             = "devices"
	DEVICE_TYPES        = "device_types"
	ROOM_CONFIGURATIONS = "room_configurations"
	UI_CONFIGS          = "ui-configuration"
	OPTIONS             = "options"
	ICONS               = "Icons"
	ROLES               = "DeviceRoles"
	ROOM_DESIGNATIONS   = "RoomDesignations"
	TAGS                = "Tags"
	DMPSLIST            = "dmps"
)

type CouchDB struct {
	address  string
	username string
	password string
}

func NewDB(address, username, password string) *CouchDB {

	return &CouchDB{
		address:  strings.Trim(address, "/"),
		username: username,
		password: password,
	}
}

func (c *CouchDB) MakeRequest(method, endpoint, contentType string, body []byte, toFill interface{}) error {
	url := fmt.Sprintf("%v/%v", c.address, endpoint)

	// start building the request
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	// add auth
	if len(c.username) > 0 && len(c.password) > 0 {
		req.SetBasicAuth(c.username, c.password)
	}

	// add headers
	if len(contentType) > 0 {
		req.Header.Add("Content-Type", contentType)
	}
	req.Header.Add("accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode/100 != 2 {
		var ce CouchError
		err = json.Unmarshal(b, &ce)
		if err != nil {
			return errors.New(fmt.Sprintf("received a non-200 response from %v. Body: %s", url, b))
		}

		log.L.Info("Non-200 response: %v", ce)
		return CheckCouchErrors(ce)
	}

	if toFill == nil {
		return nil
	}

	//otherwise we unmarshal
	err = json.Unmarshal(b, toFill)
	if err != nil {
		log.L.Infof("Can't unmarshal %v", err.Error())
		//check to see if it was a known error from couch
		var ce CouchError
		err = json.Unmarshal(b, &ce)
		if err != nil {
			return errors.New(fmt.Sprintf("unknown response from couch: %s", b))
		}

		//it was an error, we can check on error types
		return CheckCouchErrors(ce)
	}

	return nil
}

func (c *CouchDB) ExecuteQuery(query IDPrefixQuery, responseToFill interface{}) error {
	//	var toFill interface{}

	// marshal query
	b, err := json.Marshal(query)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to marshal query: %s", err))
	}

	//var toReturn []interface{}
	var database string
	//	var sliceType reflect.Type

	switch responseToFill.(type) {
	case buildingQueryResponse:
		database = BUILDINGS
		//	sliceType = reflect.TypeOf(responseToFill)
	}

	// execute query
	err = c.MakeRequest("POST", fmt.Sprintf("%s/find", database), "application/json", b, &responseToFill)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to query database %s: %s", database, err))
	}

	//	sliceType = reflect.ValueOf(responseToFill)

	return nil
}

func CheckCouchErrors(ce CouchError) error {
	switch strings.ToLower(ce.Error) {
	case "not_found":
		return &NotFound{fmt.Sprintf("The ID requested was unknown. Message: %v.", ce.Reason)}
	case "conflict":
		return &Conflict{fmt.Sprintf("There was a conflict updating/creating the document: %v", ce.Reason)}
	case "bad_request":
		return &BadRequest{fmt.Sprintf("The request was bad: %v", ce.Reason)}
	default:
		return errors.New(fmt.Sprintf("unknown error type: %v. Message: %v", ce.Error, ce.Reason))
	}
}

type IDPrefixQuery struct {
	Selector struct {
		ID struct {
			GT string `json:"$gt,omitempty"`
			LT string `json:"$lt,omitempty"`
		} `json:"_id"`
	} `json:"selector"`
	Limit int `json:"limit"`
}

type CouchUpsertResponse struct {
	OK  bool   `json:"ok"`
	ID  string `json:"id"`
	Rev string `json:"rev"`
}

type CouchError struct {
	Error  string `json:"error"`
	Reason string `json:"reason"`
}

type NotFound struct {
	msg string
}

func (n NotFound) Error() string {
	return n.msg
}

type Conflict struct {
	msg string
}

func (c Conflict) Error() string {
	return c.msg
}

type BadRequest struct {
	msg string
}

func (br BadRequest) Error() string {
	return br.msg
}
