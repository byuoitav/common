package couch

import (
	"fmt"

	"github.com/byuoitav/common/structs"
)

// GetDMJobs returns the device-monitoring's config file if there is one, or the default config file.
func (c *CouchDB) GetDMJobs(deviceID string) (structs.Jobs, error) {
	jobs, err := c.getDMJobs(deviceID)
	return *jobs.Jobs, err
}

func (c *CouchDB) getDMJobs(deviceID string) (jobs, error) {
	var toReturn jobs

	// get the device specific jobs
	err := c.MakeRequest("GET", fmt.Sprintf("%v/%v", deviceMonitoring, deviceID), "", nil, &toReturn)
	if err != nil {
		if _, ok := err.(NotFound); ok {
		} else {
			return toReturn, fmt.Errorf("unable to get device monitoring jobs: %s", err)
		}
	}

	// if that failed in a not-found error, get the default job config
	err = c.MakeRequest("GET", fmt.Sprintf("%v/%v", deviceMonitoring, "default"), "", nil, &toReturn)
	if err != nil {
		return toReturn, fmt.Errorf("unable to get device monitoring jobs: %s", err)
	}

	return toReturn, nil
}
