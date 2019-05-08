package couch

import (
	"fmt"

	"github.com/byuoitav/common/structs"
)

// GetDMPSList - get the list of DMPSes to pull events from
func (c *CouchDB) GetDMPSList() (structs.DMPSList, error) {
	var toReturn structs.DMPSList

	err := c.MakeRequest("GET", fmt.Sprintf("%v/%v", DMPSLIST, "dmps_list"), "", nil, &toReturn)

	if err != nil {
		err = fmt.Errorf("failed to get DMPSList: %s", err)
	}

	return toReturn, err
}

// GetOtherCrestronList - get the list of other crestron devices to monitor
func (c *CouchDB) GetOtherCrestronList() (structs.DMPSList, error) {
	var toReturn structs.DMPSList

	err := c.MakeRequest("GET", fmt.Sprintf("%v/%v", DMPSLIST, "CrstCustom"), "", nil, &toReturn)

	if err != nil {
		err = fmt.Errorf("failed to get CrstCustom: %s", err)
	}

	return toReturn, err
}
