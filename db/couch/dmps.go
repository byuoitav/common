package couch

import (
	"errors"
	"fmt"

	"github.com/byuoitav/common/structs"
)

// GetDMPSList - get the list of DMPSes to pull events from
func (c *CouchDB) GetDMPSList() (structs.DMPSList, error) {
	var toReturn structs.DMPSList

	err := c.MakeRequest("GET", fmt.Sprintf("%v/%v", DMPSLIST, "dmps_list"), "", nil, &toReturn)

	if err != nil {
		err = errors.New(fmt.Sprintf("failed to get DMPSList: %s", err))
	}

	return toReturn, err
}
