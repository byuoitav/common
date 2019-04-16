package couch

/*
// GetDMActions returns the device-monitoring's config file if there is one, or the default config file.
func (c *CouchDB) GetDMActions(deviceID string) ([]*actions.Actions, error) {
	actions, err := c.getDMActions(deviceID)
	if err != nil {
		return nil, err
	}

	return actions, err
}

func (c *CouchDB) getDMActions(deviceID string) ([]*Actions, error) {
	// get the device specific jobs
	err := c.MakeRequest("GET", fmt.Sprintf("%v/%v", deviceMonitoring, deviceID), "", nil, &doc)
	if err != nil {
	} else {
		return doc, nil
	}

	// get the device type

	// if that failed, get the default job config for this device type
	/*
		err := c.MakeRequest("GET", fmt.Sprintf("%v/%v", deviceMonitoring, deviceID), "", nil, &toReturn)
		if err != nil {
			if _, ok := err.(*NotFound); ok {
			} else if _, ok := err.(NotFound); ok {
			} else {
				return toReturn, fmt.Errorf("unable to get device monitoring jobs: %s", err)
			}
		} else {
			return toReturn, nil
		}
*/

/*
	// if that failed in a not-found error, get the default job config
	err = c.MakeRequest("GET", fmt.Sprintf("%v/%v", deviceMonitoring, "default"), "", nil, &toReturn)
	if err != nil {
		return toReturn, fmt.Errorf("unable to get default device monitoring jobs: %s", err)
	}
*/

// return toReturn, nil
// }
