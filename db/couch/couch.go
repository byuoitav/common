package couch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type CouchDB struct {
	address  string
	username string
	password string
	log      *zap.SugaredLogger
}

func NewDB(address, username, password string, logger *zap.SugaredLogger) *CouchDB {
	return &CouchDB{
		address:  address,
		username: username,
		password: password,
		log:      logger,
	}
}

func (c *CouchDB) MakeRequest(method, endpoint, contentType string, body []byte, toFill interface{}) error {
	url := fmt.Sprintf("%v/%v", c.address, endpoint)
	c.log.Debugf("Making %s request to %v", method, url)

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
		req.Header.Add("content-type", contentType)
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

		c.log.Infof("Got a non-200 response from: %v. Code: %v", endpoint, resp.StatusCode)

		ce := CouchError{}
		err = json.Unmarshal(b, &ce)
		if err != nil {
			msg := fmt.Sprintf("Received a non-200 response from %v. Body: %s", url, b)
			c.log.Warn(msg)
			return errors.New(msg)
		}
		return c.checkCouchErrors(ce)
	}

	if toFill == nil {
		return nil
	}

	//otherwise we unmarshal
	err = json.Unmarshal(b, toFill)
	if err != nil {
		c.log.Infof("Couldn't umarshal response into the provided struct: %v", err.Error())

		//check to see if it was a known error from couch
		ce := CouchError{}
		err = json.Unmarshal(b, &ce)
		if err != nil {
			msg := fmt.Sprintf("Unknown response from couch: %s", b)
			c.log.Warn(msg)
			return errors.New(msg)
		}
		//it was an error, we can check on error types
		return c.checkCouchErrors(ce)
	}

	return nil
}

func (c *CouchDB) checkCouchErrors(ce CouchError) error {
	c.log.Debugf("Checking for couch error type: %v", ce.Error)
	switch strings.ToLower(ce.Error) {
	case "not_found":
		c.log.Debug("Error type found: Not Found.")
		return &NotFound{fmt.Sprintf("The ID requested was unknown. Message: %v.", ce.Reason)}
	case "conflict":
		c.log.Debug("Error type found: Conflict.")
		return &Conflict{fmt.Sprintf("There was a conflict updating/creating the document: %v", ce.Reason)}
	case "bad_request":
		c.log.Debug("Error type found: Bad Request.")
		return &BadRequest{fmt.Sprintf("The request was bad: %v", ce.Reason)}
	default:
		msg := fmt.Sprintf("Unknown error type: %v. Message: %v", ce.Error, ce.Reason)
		c.log.Warn(msg)
		return errors.New(msg)
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
