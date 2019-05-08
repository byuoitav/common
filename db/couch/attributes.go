package couch

import (
	"encoding/json"
	"fmt"

	"github.com/byuoitav/common/structs"
)

// GetAttributeGroup returns an attribute group from the database
func (c *CouchDB) GetAttributeGroup(groupID string) (structs.Group, error) {
	resp, err := c.getAttributeGroup(groupID)
	return resp.Group, err
}

func (c *CouchDB) getAttributeGroup(groupID string) (attributeGroup, error) {
	var toReturn attributeGroup

	err := c.MakeRequest("GET", fmt.Sprintf("%v/%v", ATTRIBUTES, groupID), "", nil, toReturn)
	if err != nil {
		err = fmt.Errorf("failed to get attribute group %s: %s", groupID, err)
	}

	return toReturn, err
}

// GetAllAttributeGroups returns a list of all the attribute groups in the database
func (c *CouchDB) GetAllAttributeGroups() ([]structs.Group, error) {
	var toReturn []structs.Group
	var query IDPrefixQuery
	query.Selector.ID.GT = "\x00"
	query.Limit = 1000

	b, err := json.Marshal(query)
	if err != nil {
		return toReturn, fmt.Errorf("failed to marshal query to get all attribute groups: %s", err)
	}

	var resp attributeQueryResponse

	err = c.MakeRequest("POST", fmt.Sprintf("%v/_find", ATTRIBUTES), "application/json", b, &resp)
	if err != nil {
		return toReturn, fmt.Errorf("failed to get all attribute groups: %s", err)
	}

	for _, doc := range resp.Docs {
		toReturn = append(toReturn, doc.Group)
	}

	return toReturn, err
}
