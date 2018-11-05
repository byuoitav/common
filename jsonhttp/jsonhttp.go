package jsonhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	"github.com/byuoitav/common/log"
)

//CreateRequest Creates an http request with a json body representing the body object passed in
func CreateRequest(method string, url string, body interface{}, headers map[string]string) (*http.Request, error) {

	var bodyBytes []byte
	var ok bool
	var err error

	if bodyBytes, ok = body.([]byte); !ok {
		//marshal
		bodyBytes, err = json.Marshal(body)

		if err != nil {
			return nil, err
		}
	}
	log.L.Debugf("Request body: %s", bodyBytes)

	// start building the request
	requestToReturn, err := http.NewRequest(method, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return requestToReturn, err
	}

	// add headers
	requestToReturn.Header.Add("Content-Type", "application/json")
	requestToReturn.Header.Add("Accept", "application/json")

	for key, value := range headers {
		requestToReturn.Header.Add(key, value)
	}

	return requestToReturn, nil
}

//ExecuteRequest will execute the http request and unmarshal the response into the output object, output interface must be a &pointer
func ExecuteRequest(req *http.Request, output interface{}, timeoutInSeconds int) error {

	if reflect.ValueOf(output).Kind() != reflect.Ptr {
		return fmt.Errorf("output variable must be a pointer")
	}

	client := &http.Client{}
	client.Timeout = time.Duration(timeoutInSeconds) * time.Second
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
		if err != nil {
			return fmt.Errorf("received a non-200 response")
		}
	}

	log.L.Debugf("Response [%v] received", string(b))

	//otherwise we unmarshal
	err = json.Unmarshal(b, &output)
	if err != nil {
		return fmt.Errorf("Can't unmarshal %v. Received: %s", err.Error(), b)
	}

	return nil
}
