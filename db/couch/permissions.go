package couch

import (
	"errors"
	"fmt"

	"github.com/byuoitav/common/structs"
)

func (c *CouchDB) GetAuth() (structs.Auth, error) {
	var toReturn structs.Auth

	err := c.MakeRequest("GET", fmt.Sprintf("%v/%v", "auth", "auth"), "", nil, &toReturn)
	if err != nil {
		err = errors.New(fmt.Sprintf("failed to get permissions: %s", err))
	}

	return toReturn, err
}
