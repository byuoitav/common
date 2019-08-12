package couch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/structs"
)

//GetDeploymentInfo returns a config definition for a service
func (c *CouchDB) GetDeploymentInfo(serviceID string) (structs.FullConfig, error) {
	toReturn, err := c.getDeploymentInfo(serviceID)
	return toReturn, err
}

func (c *CouchDB) getDeploymentInfo(serviceID string) (structs.FullConfig, error) {
	var toReturn structs.FullConfig
	err := c.MakeRequest("GET", fmt.Sprintf("%s/%s", DEPLOY, serviceID), "", nil, &toReturn)
	return toReturn, err
}

//GetDeviceDeploymentInfo returns the the necessary elements for a device of a given designation
func (c *CouchDB) GetDeviceDeploymentInfo(deviceType string) (structs.DeviceDeploymentConfig, error) {
	var toReturn structs.DeviceDeploymentConfig
	err := c.MakeRequest("GET", fmt.Sprintf("%s/%s", CAMPUS, deviceType), "", nil, &toReturn)
	return toReturn, err
}

//GetServiceInfo returns a service config definition
func (c *CouchDB) GetServiceInfo(serviceID string) (structs.ServiceConfigWrapper, error) {
	toReturn, err := c.getServiceInfo(serviceID)
	return toReturn, err
}

func (c *CouchDB) getServiceInfo(serviceID string) (structs.ServiceConfigWrapper, error) {
	var toReturn structs.ServiceConfigWrapper
	err := c.MakeRequest("GET", fmt.Sprintf("deployment-information/%s", serviceID), "", nil, &toReturn)
	return toReturn, err
}

// GetServiceAttachment .
func (c *CouchDB) GetServiceAttachment(service, designation string) ([]byte, error) {
	url := fmt.Sprintf("%v/%v/%v/%v", c.address, DEPLOY, service, fmt.Sprintf("%v-%v", service, designation))

	// start building the request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		var ce CouchError
		err = json.Unmarshal(b, &ce)
		if err != nil {
			return nil, fmt.Errorf("received a non-200 response from %v. Body: %s", url, b)
		}

		log.L.Infof("Non-200 response: %v", ce.Error)
		return nil, CheckCouchErrors(ce)
	}

	return b, nil
}

// GetServiceZip .
func (c *CouchDB) GetServiceZip(service, designation string) ([]byte, error) {
	url := fmt.Sprintf("%v/%v/%v/%v", c.address, DEPLOY, service, fmt.Sprintf("%v.tar.gz", designation))

	// start building the request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		var ce CouchError
		err = json.Unmarshal(b, &ce)
		if err != nil {
			return nil, fmt.Errorf("received a non-200 response from %v. Body: %s", url, b)
		}

		log.L.Infof("Non-200 response: %v", ce.Error)
		return nil, CheckCouchErrors(ce)
	}

	return b, nil
}
